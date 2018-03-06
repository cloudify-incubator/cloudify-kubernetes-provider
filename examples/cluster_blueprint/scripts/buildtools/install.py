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

import platform
import subprocess
from cloudify import ctx
from cloudify.state import ctx_parameters as inputs
from cloudify.exceptions import NonRecoverableError, OperationRetry


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

    ctx.logger.info('Installing build requirements.')
    linux_distro = inputs.get('linux_distro', 'centos')
    agent_user = inputs.get('agent_user', 'centos')

    if not linux_distro:
        distro, _, _ = \
            platform.linux_distribution(full_distribution_name=False)
        linux_distro = distro.lower()

    if 'centos' in linux_distro:

        build_output = \
            execute_command(
                'sudo yum install -q -y git')
        if build_output is False:
            raise OperationRetry(
                'Failed to install git')

        import_gpg_key = execute_command(
            'sudo rpm --import https://mirror.go-repo.io/'
            'centos/RPM-GPG-KEY-GO-REPO')
        if import_gpg_key is False:
            raise OperationRetry(
                'Failed to import Go GPG key')

        go_repo_temp = ctx.download_resource('resources/go.repo')
        execute_command(
            'sudo mv {0} /etc/yum.repos.d/go.repo'.format(go_repo_temp))

        go_install = execute_command('sudo yum -y install golang')
        if go_install is False:
            raise OperationRetry(
                'Failed to import Go GPG key')

    elif 'ubuntu' in linux_distro:
        execute_command(
            'sudo add-apt-repository ppa:longsleep/golang-backports')
        execute_command('sudo apt-get update')
        execute_command('sudo apt-get install golang-go git')
    else:
        raise NonRecoverableError('Unsupported platform.')
