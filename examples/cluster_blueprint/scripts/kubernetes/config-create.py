#!/usr/bin/env python

import platform
import ssl
import subprocess
import tempfile
import os
from cloudify import ctx
from cloudify.state import ctx_parameters as inputs
from cloudify.exceptions import NonRecoverableError

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


if __name__ == '__main__':

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
    with open(os.path.expanduser('~') + "/cfy.json", 'w') as outfile:
        outfile.write("{" + CONFIG.format(
            cfy_deployment,
            cfy_instance,
            cfy_tenant,
            cfy_pass,
            cfy_user,
            cfy_host if not cfy_ssl else "https://" + cfy_host) + "}")

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
            cert_file.write(ssl.get_server_certificate((
                cfy_host, 443)))

        if 'centos' in linux_distro:
            execute_command([
                'sudo', 'bash', '-c',
                'cat {} >> /etc/pki/tls/certs/ca-bundle.crt'
                .format(temp_cert_file)
            ])
        else:
            raise NonRecoverableError('Unsupported platform.')
