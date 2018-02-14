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

    ctx.logger.info('Add kubernetes repository')
    linux_distro = inputs.get('linux_distro', 'centos')
    agent_user = inputs.get('agent_user', 'centos')

    if not linux_distro:
        distro, _, _ = \
            platform.linux_distribution(full_distribution_name=False)
        linux_distro = distro.tolower()

    if 'centos' in linux_distro:

        kubernetes_repo_temp = ctx.download_resource(
            'resources/kubernetes.repo'
        )
        copy_repo = execute_command(
            'sudo mv {0} /etc/yum.repos.d/kubernetes.repo'.format(
                kubernetes_repo_temp))
        if copy_repo is False:
            raise OperationRetry(
                'Failed to copy repository description')

        disable_selinux = execute_command('sudo setenforce 0')
        if disable_selinux is False:
            raise OperationRetry(
                'Failed to disable selinux')

        kubernetes_install = execute_command(
            'sudo yum -y install kubeadm-1.8.5-0 '
            'kubelet-1.8.5-0 kubectl-1.8.5-0'
        )
        if kubernetes_install is False:
            raise OperationRetry(
                'Failed to install kubenrnetes')

    elif 'ubuntu' in linux_distro:
        execute_command('apt-get update && apt-get install -y ' +
                        'apt-transport-https curl')
        execute_command(
            'curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg |' +
            ' apt-key add -')

        kubernetes_repo_temp = ctx.download_resource(
            'resources/kubernetes.list')
        copy_repo = execute_command(
            'sudo mv {0} /etc/apt/sources.list.d/kubernetes.list'.format(
                kubernetes_repo_temp))
        if copy_repo is False:
            raise OperationRetry(
                'Failed to copy repository description')

        execute_command('sudo apt-get update')
        execute_command('sudo apt-get install -y kubelet=1.8.5-00 '
                        'kubeadm=1.8.5-00 kubectl=1.8.5-00')
    else:
        raise NonRecoverableError('Unsupported platform.')
