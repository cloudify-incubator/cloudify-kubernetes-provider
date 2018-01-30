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

    ctx.logger.info('Downloading or building Cfy Provider.')

    provider_repo = \
        inputs.get('cfy_go_binary_path', '/usr/bin/cfy-go')
    agent_user = \
        inputs.get('agent_user', 'centos')
    group = \
        inputs.get('group', 'wheel')

    execute_command('sudo CGO_ENABLED=0 go install -a -installsuffix cgo std')
    execute_command('sudo mkdir -p /opt/cloudify-kubernetes-provider')
    execute_command('sudo chmod -R 777 /opt/')
    _extra_args = {
        'env': {
            'USER': agent_user
        }
    }
    _command = \
        'sudo chown {0}:{1} /opt/cloudify-' \
        'kubernetes-provider'.format(agent_user, group)
    execute_command(_command)

    try:
        ctx.download_resource('resources/cfy-go')
    except HttpException:
        ctx.logger.debug('Download provider repo.')
        _cwd = '/opt/'
        _extra_args = {'cwd': _cwd}
        _command = \
            'git clone https://github.com/cloudify-incubator/' \
            'cloudify-kubernetes-provider.git --depth 1 -b master'
        execute_command(_command, extra_args=_extra_args)
        _cwd = \
            os.path.join(_cwd, 'cloudify-kubernetes-provider/')

        git_modules_file = os.path.join(_cwd, '.gitmodules')
        _, temp_git_file = tempfile.mkstemp()
        with open(git_modules_file, 'r') as infile:
            with open(temp_git_file, 'w') as outfile:
                for line in infile.readlines():
                    outfile.write(
                        line.replace(
                            'git@github.com:', 'https://github.com/'))

        ctx.logger.debug('Download submodules sources.')
        execute_command(
            'sudo cp {0} {1}'.format(temp_git_file, git_modules_file))
        _extra_args = {'cwd': _cwd}
        execute_command('git submodule init', extra_args=_extra_args)
        execute_command('git submodule update', extra_args=_extra_args)
