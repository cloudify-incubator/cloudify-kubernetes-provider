#!/usr/bin/env python
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


def install_centos(_docker_version, _agent_user, _update_yum):
    """Install Docker _docker_version on Centos
    """

    ctx.logger.info(
        'Installing Docker {0} on Centos.'.format(_docker_version))

    if not os.path.exists('/etc/yum.repos.d/docker.repo'):
        docker_repo_temp = ctx.download_resource('resources/docker.repo')
        save_repo_file = execute_command(
            'sudo mv {0} /etc/yum.repos.d/docker.repo'.format(
                docker_repo_temp))
        if save_repo_file is False:
            return False

    execute_command('sudo yum install deltarpm epel-release unzip -q -y')

    if _update_yum is True:
        execute_command('sudo yum update -y -q')

    execute_command('sudo groupadd docker')
    execute_command('sudo usermod -aG docker {0}'.format(_agent_user))
    return execute_command(
        'sudo yum install docker-engine-{0} -y -q'.format(_docker_version))


def install_ubuntu():
    raise NonRecoverableError('Not supported.')


if __name__ == '__main__':

    ctx.logger.info(
        'Verifying that Docker is installed on the system.')

    # Allow user overrides
    docker_binary_location = \
        inputs.get('docker_binary_location', '/usr/bin/docker')
    docker_version = inputs.get('docker_version', '1.12.6')
    update_package_manager = inputs.get('update_package_manager', False)
    linux_distro = inputs.get('linux_distro', 'centos')
    agent_user = inputs.get('agent_user', 'centos')

    if not linux_distro:
        distro, _, _ = \
            platform.linux_distribution(full_distribution_name=False)
        linux_distro = distro.tolower()

    if os.path.exists(docker_binary_location) \
            and docker_version in execute_command('docker version'):
        ctx.logger.debug(
            'Docker {0} already installed.'.format(docker_version))
        if ctx.operation.retry_number == 0:
            # No cleanup.
            ctx.instance.runtime_properties['existing_docker_install'] = \
                True
        docker_installed = True
    elif 'centos' in linux_distro:
        docker_installed = \
            install_centos(docker_version, agent_user, update_package_manager)
    else:
        docker_installed = False

    if docker_installed is False:
        raise NonRecoverableError('Failed to install Docker see logs.')

    enabled = execute_command('sudo systemctl enable docker.service')
    if not enabled:
        OperationRetry('Failed to enable docker.service')
    started = execute_command('sudo systemctl start docker')
    if not started:
        OperationRetry('Failed to start docker.service')
