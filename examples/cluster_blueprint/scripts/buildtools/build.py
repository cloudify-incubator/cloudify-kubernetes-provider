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
import tempfile
from cloudify import ctx
from cloudify.exceptions import HttpException
from cloudify.state import ctx_parameters as inputs


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
        execute_command(['sudo', 'cp', cfy_binary, service_path])
    # fix file attributes
    execute_command(['sudo', 'chmod', '555', service_path])
    execute_command(['sudo', 'chown', 'root:root', service_path])
    ctx.logger.debug('{} attributes fixed'.format(service_name))


if __name__ == '__main__':

    ctx.logger.info('Downloading or building Cfy GO Client.')
    full_install = inputs.get('full_install', 'all')

    try:
        download_service("cfy-go")
        if full_install == "all":
            # download cluster provider
            download_service("cfy-kubernetes")

            # download scale tools
            download_service("cfy-autoscale")
    except HttpException:
        if full_install != "all":
            ctx.logger.debug('Download cfy-go repo.')
            cwd = '/opt/cloudify-kubernetes-provider/'
            execute_command(['sudo', 'mkdir', '-p', cwd])
            execute_command(['sudo', 'chmod', '-R', '777', cwd])
            ctx.logger.debug('Download cfy-go repo.')
            extra_args = {
                'cwd': cwd,
                'env': {
                    'GOBIN': os.path.join(cwd, 'bin'),
                    'GOPATH': cwd,
                    'PATH': ':'.join(
                        [os.getenv('PATH'), os.path.join(cwd, 'bin')])
                }
            }
            command = ['go', 'get', 'github.com/cloudify-incubator/'
                       'cloudify-rest-go-client/cfy-go']
            execute_command(command, extra_args=extra_args)
        else:
            ctx.logger.debug('Download provider repo.')
            cwd = '/opt/'
            extra_args = {'cwd': cwd}
            command = ['git', 'clone', 'https://github.com/cloudify-incubator/'
                       'cloudify-kubernetes-provider.git', '--depth', '1',
                       '-b', 'testing']
            execute_command(command, extra_args=extra_args)

            cwd = os.path.join(cwd, 'cloudify-kubernetes-provider/')

            extra_args = {
                'cwd': cwd,
                'env': {
                    'GOBIN': os.path.join(cwd, 'bin'),
                    'GOPATH': cwd,
                    'PATH': ':'.join([os.getenv('PATH'),
                                      os.path.join(cwd, 'bin')])
                }
            }
            git_modules_file = os.path.join(cwd, '.gitmodules')
            _, temp_git_file = tempfile.mkstemp()
            with open(git_modules_file, 'r') as infile:
                with open(temp_git_file, 'w') as outfile:
                    for line in infile.readlines():
                        outfile.write(
                            line.replace(
                                'git@github.com:', 'https://github.com/'))

            ctx.logger.debug('Download submodules sources.')
            execute_command(['sudo', 'cp', temp_git_file, git_modules_file])
            execute_command(['git', 'submodule', 'init'],
                            extra_args=extra_args)
            execute_command(['git', 'submodule', 'update'],
                            extra_args=extra_args)
            command = ['go', 'install', 'src/cfy-kubernetes.go']
            execute_command(command, extra_args=extra_args)
            command = ['go', 'build', '-v', '-o', 'bin/cluster-autoscaler',
                       'src/k8s.io/autoscaler/cluster-autoscaler/main.go',
                       'src/k8s.io/autoscaler/cluster-autoscaler/version.go']
            execute_command(command, extra_args=extra_args)

        execute_command(['sudo', 'cp',
                         "/opt/cloudify-kubernetes-provider/bin/*",
                         '/usr/bin/'])
