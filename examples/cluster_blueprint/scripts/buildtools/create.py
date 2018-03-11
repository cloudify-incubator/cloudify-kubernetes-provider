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
import platform
import subprocess
from cloudify import ctx
from cloudify.exceptions import (
    NonRecoverableError, HttpException, OperationRetry
)


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
        ctx.logger.info('Installing build requirements.')
        linux_distro = ctx.node.properties.get('linux_distro', 'centos')

        if not linux_distro:
            distro, _, _ = \
                platform.linux_distribution(full_distribution_name=False)
            linux_distro = distro.lower()

        if 'centos' in linux_distro:

            build_output = execute_command(['sudo', 'yum', 'install', '-q',
                                            '-y', 'git'])
            if build_output is False:
                raise OperationRetry('Failed to install git')

            import_gpg_key = execute_command([
                'sudo', 'rpm', '--import',
                'https://mirror.go-repo.io/centos/RPM-GPG-KEY-GO-REPO'])
            if import_gpg_key is False:
                raise OperationRetry('Failed to import Go GPG key')

            go_repo_temp = ctx.download_resource('resources/go.repo')
            execute_command(['sudo', 'mv', go_repo_temp,
                             '/etc/yum.repos.d/go.repo'])

            go_install = execute_command(['sudo', 'yum', '-y', 'install',
                                          'golang'])
            if go_install is False:
                raise OperationRetry('Failed to import Go GPG key')

        elif 'ubuntu' in linux_distro:
            execute_command(['sudo', 'add-apt-repository',
                             'ppa:longsleep/golang-backports'])
            execute_command(['sudo', 'apt-get', 'update'])
            execute_command(['sudo', 'apt-get', 'install', 'golang-go', 'git'])
        else:
            raise NonRecoverableError('Unsupported platform.')
