#!/usr/bin/env python
#
# Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import os
import subprocess
from cloudify import ctx
from cloudify.exceptions import NonRecoverableError, HttpException


def execute_command(command, extra_args=None):

    ctx.logger.debug('command: {0}.'.format(repr(command)))

    subprocess_args = {
        'args': command,
        'stdout': subprocess.PIPE,
        'stderr': subprocess.PIPE
    }
    if extra_args is not None and isinstance(extra_args, dict):
        subprocess_args.update(extra_args)

    ctx.logger.debug('subprocess_args {0}.'.format(subprocess_args))

    process = subprocess.Popen(**subprocess_args)
    output, error = process.communicate()

    ctx.logger.debug('command: {0} '.format(repr(command)))
    ctx.logger.debug('output: {0} '.format(output))
    ctx.logger.debug('error: {0} '.format(error))
    ctx.logger.debug('process.returncode: {0} '.format(process.returncode))

    if process.returncode:
        ctx.logger.error('Running `{0}` returns {1} error: {2}.'
                         .format(repr(command), process.returncode,
                                 repr(error)))
        return False

    return output


def download_service(service_name):
    service_path = "/usr/bin/" + service_name
    if not os.path.isfile(service_path):
        cfy_binary = ctx.download_resource('resources/{}'
                                           .format(service_name))
        ctx.logger.debug('{} downloaded.'.format(service_name))
        if execute_command(['sudo', 'cp', cfy_binary, service_path]) is False:
            raise NonRecoverableError("Can't copy {}.".format(service_path))
    # fix file attributes
    if execute_command(['sudo', 'chmod', '555', service_path]) is False:
        raise NonRecoverableError("Can't chmod {}.".format(service_path))
    if execute_command(['sudo', 'chown', 'root:root', service_path]) is False:
        raise NonRecoverableError("Can't chown {}.".format(service_path))
    ctx.logger.debug('{} attributes fixed'.format(service_name))


if __name__ == '__main__':
    full_install = ctx.node.properties.get('full_install', 'all')

    try:
        download_service("cfy-go")
        if full_install == "all":
            # download cluster provider
            download_service("cfy-kubernetes")

            # download scale tools
            download_service("cfy-autoscale")
    except HttpException:
        cwd = '/opt/cloudify-kubernetes-provider/'
        extra_args = {
            'cwd': cwd,
            'env': {
                'GOBIN': os.path.join(cwd, 'bin'),
                'GOPATH': cwd,
                'PATH': ':'.join(
                    [os.getenv('PATH'), os.path.join(cwd, 'bin')])
            }
        }
        if full_install != "all":
            ctx.logger.info('Download cfy-go repo.')
            command = ['go', 'get', 'github.com/cloudify-incubator/'
                       'cloudify-rest-go-client/cfy-go']
            if execute_command(command, extra_args=extra_args) is False:
                raise NonRecoverableError("Can't build cfy-go.")
        else:
            ctx.logger.info('Build cfy-kubernetes')
            command = ['go', 'install', 'src/cfy-kubernetes.go']
            if execute_command(command, extra_args=extra_args) is False:
                raise NonRecoverableError("Can't build cfy-kubernetes.")

            ctx.logger.info('Build cfy-autoscale')
            command = ['go', 'build', '-v', '-o', 'bin/cfy-autoscale',
                       'src/k8s.io/autoscaler/cluster-autoscaler/main.go',
                       'src/k8s.io/autoscaler/cluster-autoscaler/version.go']
            if execute_command(command, extra_args=extra_args) is False:
                raise NonRecoverableError("Can't build cfy-autoscale.")

            ctx.logger.info('Build cfy-go')
            command = ['go', 'get', 'github.com/cloudify-incubator/'
                       'cloudify-rest-go-client/cfy-go']
            if execute_command(command, extra_args=extra_args) is False:
                raise NonRecoverableError("Can't build cfy-go.")

        for_copy = ["/opt/cloudify-kubernetes-provider/bin/cfy-go"]
        if full_install == "all":
            for_copy += [
                "/opt/cloudify-kubernetes-provider/bin/cfy-kubernetes",
                "/opt/cloudify-kubernetes-provider/bin/cfy-autoscale"]

        for path in for_copy:
            if execute_command(['sudo', 'cp', path, '/usr/bin/']) is False:
                raise NonRecoverableError("Can't copy {} tool.".format(path))
