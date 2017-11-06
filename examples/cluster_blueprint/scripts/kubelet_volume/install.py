#!/usr/bin/env python

import os
import subprocess
import tempfile
from cloudify import ctx
from cloudify.state import ctx_parameters as inputs

MOUNT = '#!/bin/bash' \
        'echo \$@ >> /var/log/mount-calls.log' \
        '/usr/bin/{0} kubernetes \$1 \$2 \$3 -deployment "{1}" '\
        ' -instance "{2}" -tenant "{3}" '\
        '-password "{4}" -user "{5}" -host "{6}'


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

    ctx.logger.info('Configuring Kubelet Volume Plugin.')

    cfy_go_binary_path = \
        inputs.get('cfy_go_binary_path', '/usr/bin/cfy-go')

    plugin_directory = \
        inputs.get(
            'plugin_directory',
            '/usr/libexec/kubernetes/'
            'kubelet-plugins/volume/exec/cloudify~mount/')

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

    if os.path.exists('/usr/bin/cfy-go'):
        ctx.logger.debug(
            'Cfy Go Binary already at {0}'.format(cfy_go_binary_path))
    else:
        ctx.logger.debug(
            'Copying Cfy Go Binary to {0}'.format(cfy_go_binary_path))
        execute_command(
            'sudo cp /opt/bin/cfy-go {0}'.format(cfy_go_binary_path))

    execute_command('sudo chmod 555 {0}'.format(cfy_go_binary_path))
    execute_command('sudo chown root:root {0}'.format(cfy_go_binary_path))

    _, temp_mount_file = tempfile.mkstemp()

    with open(temp_mount_file, 'w') as outfile:
        outfile.write(MOUNT.format(
            cfy_go_binary_path,
            cfy_deployment,
            cfy_instance,
            cfy_tenant,
            cfy_pass,
            cfy_user,
            cfy_host))

    execute_command('sudo mkdir -p {0}'.format(
        plugin_directory))
    execute_command('sudo cp {0} {1}'.format(
        temp_mount_file,
        os.path.join(plugin_directory, 'mount')))
    execute_command('sudo chmod 555 {0}'.format(
        os.path.join(plugin_directory, 'mount')))
    execute_command('sudo chown root:root {0}'.format(
        os.path.join(plugin_directory, 'mount')))
