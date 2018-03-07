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

    ctx.logger.info('Downloading or building Cfy GO Client.')
    full_install = ctx.node.properties.get('full_install', 'all')

    try:
        download_service("cfy-go")
        if full_install == "all":
            # download cluster provider
            download_service("cfy-kubernetes")

            # download scale tools
            download_service("cfy-autoscale")
    except HttpException:
        if full_install != "all":
            ctx.logger.info('Download cfy-go repo.')
            cwd = '/opt/cloudify-kubernetes-provider/'
            if not os.path.isdir(cwd):
                if execute_command(['sudo', 'mkdir', '-p', cwd]) is False:
                    raise NonRecoverableError("Can't create directory.")
            if execute_command(['sudo', 'chmod', '-R', '777', cwd]) is False:
                raise NonRecoverableError("Can't change owner.")
        else:
            ctx.logger.debug('Download provider repo.')
            cwd = '/opt/'
            if execute_command(['sudo', 'chmod', '-R', '777', cwd]) is False:
                raise NonRecoverableError("Can't change owner.")
            extra_args = {'cwd': cwd}
            command = ['git', 'clone', 'https://github.com/cloudify-incubator/'
                       'cloudify-kubernetes-provider.git', '--depth', '1',
                       '-b', 'testing']
            if execute_command(command, extra_args=extra_args) is False:
                raise NonRecoverableError("Can't download sources.")
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
            if execute_command(['sudo', 'cp', temp_git_file,
                                git_modules_file]) is False:
                raise NonRecoverableError("Can't update links.")
            if execute_command(['git', 'submodule', 'init'],
                               extra_args=extra_args) is False:
                raise NonRecoverableError("Can't init subsources.")
            if execute_command(['git', 'submodule', 'update'],
                               extra_args=extra_args) is False:
                raise NonRecoverableError("Can't update subsources.")
