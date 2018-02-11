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
from cloudify.exceptions import HttpException
from cloudify.state import ctx_parameters as inputs


def execute_command(_command, extra_args=None):

    ctx.logger.debug('_command {0}.'.format(_command))

    subprocess_args = {
        'args': _command.split(),
        'stdout': subprocess.PIPE,
        'stderr': subprocess.PIPE
    }
    if extra_args is not None and isinstance(extra_args, dict):
        subprocess_args.update(extra_args)

    ctx.logger.debug('subprocess_args {0}.'.format(subprocess_args))

    process = subprocess.Popen(**subprocess_args)
    output, error = process.communicate()

    ctx.logger.debug('command: {0} '.format(_command))
    ctx.logger.debug('output: {0} '.format(output))
    ctx.logger.debug('error: {0} '.format(error))
    ctx.logger.debug('process.returncode: {0} '.format(process.returncode))

    if process.returncode:
        ctx.logger.error('Running `{0}` returns error.'.format(_command))
        return False

    return output


if __name__ == '__main__':

    ctx.logger.info('Downloading or building cfy-kubernetes.')

    cfy_kubernetes_binary_path = \
        inputs.get(
            'cfy_kubernetes_binary_path',
            '/usr/bin/cfy-kubernetes')
    execute_command('sudo mkdir -p /opt/cloudify-kubernetes-provider/bin')

    if not os.path.isfile(cfy_kubernetes_binary_path):
        try:
            cfy_kubernetes_binary = \
                ctx.download_resource('resources/cfy-kubernetes')
        except HttpException:
            ctx.logger.debug('Build cfy-kubernetes.')
            _cwd = '/opt/cloudify-kubernetes-provider'
            _extra_args = {
                'cwd': _cwd,
                'env': {
                    'GOBIN': os.path.join(_cwd, 'bin'),
                    'GOPATH': _cwd,
                    'PATH': ':'.join(
                        [os.getenv('PATH'), os.path.join(_cwd, 'bin')])
                }
            }
            _command = \
                'go install src/cfy-kubernetes.go'
            execute_command(_command, extra_args=_extra_args)
            cfy_kubernetes_binary = \
                '/opt/cloudify-kubernetes-provider/bin/cfy-kubernetes'
        else:
            ctx.logger.debug('cfy-kubernetes already built/downloaded.')

            execute_command('sudo chmod -R 755 /opt/')
            execute_command('sudo mkdir -p /opt/bin')

            # TODO: Understand how to make sure we don't need this.
            execute_command(
                'sudo cp {0} {1}'.format(
                    cfy_kubernetes_binary,
                    '/opt/cloudify-kubernetes-provider/bin/cfy-kubernetes'))
            execute_command(
                'sudo cp {0} {1}'.format(
                    cfy_kubernetes_binary, cfy_kubernetes_binary_path))

        ctx.logger.debug('cfy-kubernetes built/downloaded.')
        execute_command(
            'sudo cp {0} {1}'.format(
                cfy_kubernetes_binary, cfy_kubernetes_binary_path))
