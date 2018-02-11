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
import ssl
import subprocess
import tempfile
import os
import os.path
from cloudify import ctx
from cloudify.state import ctx_parameters as inputs
from cloudify.exceptions import (
    HttpException,
    NonRecoverableError,
    OperationRetry
)

CONFIG = ('"deployment": "{0}",' +
          '"instance": "{1}",' +
          '"tenant": "{2}",' +
          '"password": "{3}",' +
          '"user": "{4}",' +
          '"host": "{5}"')


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
        try:
            cfy_binary = ctx.download_resource(
                'resources/{}'.format(service_name))
        except HttpException:
            raise NonRecoverableError(
                '{} binary not in resources.'.format(service_name))
        ctx.logger.debug('{} downloaded.'.format(service_name))
        execute_command(['sudo', 'cp', cfy_binary, service_path])
    # fix file attributes
    execute_command(['sudo', 'chmod', '555', service_path])
    execute_command(['sudo', 'chown', 'root:root', service_path])
    ctx.logger.debug('{} attributes fixed'.format(service_name))


def create_service(service_name):
    service_path = '/etc/systemd/system/{}.service'.format(service_name)
    if not os.path.isfile(service_path):
        try:
            _tv = {'home_dir': os.path.expanduser('~')}
            cfy_service = \
                ctx.download_resource_and_render(
                    'resources/{}.service'.format(service_name),
                    template_variables=_tv)
        except HttpException:
            raise NonRecoverableError(
                '{}.service not in resources.'.format(service_name))
        else:
            execute_command(['sudo', 'cp', cfy_service, service_path])
            execute_command(['sudo', 'cp',
                             '/etc/systemd/system/{}.service'
                             .format(service_name),
                             '/etc/systemd/system/multi-user.target.wants/'])

        execute_command(['sudo', 'systemctl', 'daemon-reload'])
        execute_command(['sudo', 'systemctl', 'enable',
                         '{}.service'.format(service_name)])
        execute_command(['sudo', 'systemctl', 'start',
                         '{}.service'.format(service_name)])


def start_check(service_name):
    status = ''
    systemctl_status = execute_command(['sudo', 'systemctl', 'status',
                                        '{}.service'.format(service_name)])
    if not isinstance(systemctl_status, basestring):
        raise OperationRetry(
            'check sudo systemctl status {}.service'.format(service_name))
    for line in systemctl_status.split('\n'):
        if 'Active:' in line:
            status = line.strip()
    zstatus = status.split(' ')
    ctx.logger.info('{} status: {}'.format(zstatus, service_name))
    if not len(zstatus) > 1 and 'active' not in zstatus[1]:
        raise OperationRetry('Wait a little more.')


if __name__ == '__main__':

    # create global config
    config_file = os.path.expanduser('~') + "/cfy.json"
    if not os.path.isfile(config_file):
        ctx.logger.info("Create config {} file".format(config_file))

        linux_distro = inputs.get('linux_distro', 'centos')

        cfy_deployment = \
            inputs.get('cfy_deployment', ctx.deployment.id)

        cfy_instance = \
            inputs.get('cfy_instance', ctx.instance.id)

        cfy_user = \
            inputs.get('cfy_user', 'admin')

        cfy_pass = \
            inputs.get('cfy_password', 'admin')

        cfy_tenant = \
            inputs.get('cfy_tenant', 'default_tenant')

        cfy_host = \
            inputs.get('cfy_host', 'localhost')

        cfy_ssl = \
            inputs.get('cfy_ssl', False)

        ctx.logger.info("create cloudify manager config")

        # services config
        with open(config_file, 'w') as outfile:
            outfile.write("{" + CONFIG.format(
                cfy_deployment,
                cfy_instance,
                cfy_tenant,
                cfy_pass,
                cfy_user,
                cfy_host if not cfy_ssl else "https://" + cfy_host) + "}")

    if ctx.operation.retry_number == 0:
        if not linux_distro:
            distro, _, _ = \
                platform.linux_distribution(full_distribution_name=False)
            linux_distro = distro.tolower()

        if cfy_ssl:
            ctx.logger.info("Set certificate as trusted")

            # cert config
            _, temp_cert_file = tempfile.mkstemp()

            with open(temp_cert_file, 'w') as cert_file:
                cert_file.write("# cloudify certificate\n")
                try:
                    cert_file.write(ssl.get_server_certificate((
                        cfy_host, 443)))
                except Exception as ex:
                    ctx.logger.error("Check https connection to manager {}."
                                     .format(str(ex)))

            if 'centos' in linux_distro:
                execute_command([
                    'sudo', 'bash', '-c',
                    'cat {} >> /etc/pki/tls/certs/ca-bundle.crt'
                    .format(temp_cert_file)
                ])
            else:
                raise NonRecoverableError('Unsupported platform.')

    full_install = inputs.get('full_install', 'all')

    # download mount tools
    download_service("cfy-go")

    if full_install == "all":
        # download scale tools
        download_service("cfy-autoscale")
        create_service("cfy-autoscale")

        # download cluster provider
        download_service("cfy-kubernetes")
        create_service("cfy-kubernetes")

        start_check("cfy-autoscale")
        start_check("cfy-kubernetes")
