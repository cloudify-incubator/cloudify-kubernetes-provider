#!/usr/bin/env python

import base64
import subprocess
import tempfile
from cloudify import ctx
from cloudify.state import ctx_parameters as inputs
from cloudify.exceptions import OperationRetry, NonRecoverableError

JOIN = 'sudo kubeadm join --token {0} {1}:6443 --skip-preflight-checks'
IP_TABLES_PATH = '/proc/sys/net/bridge/bridge-nf-call-iptables'


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

    token = inputs.get('token')
    ip = inputs.get('ip')
    ctx.logger.info('Try to join to {0} by {1}'.format(ip, token))

    if ctx.operation.retry_number == 1:

        _, temp_mount_file = tempfile.mkstemp()
        with open(temp_mount_file, 'w') as outfile:
            outfile.write('1')

        execute_command(
            'sudo cp {0} {1}'.format(
                temp_mount_file, IP_TABLES_PATH))

        token_decoded = base64.b64decode(token)
        execute_command(JOIN.format(token_decoded, ip))

    status = ''
    systemctl_status = execute_command('sudo systemctl status kubelet')
    if not isinstance(systemctl_status, basestring):
        raise OperationRetry('check sudo systemctl status kubelet')
    for line in systemctl_status.split('\n'):
        if 'Active:' in line:
            status = line.strip()
    zstatus = status.split(' ')
    ctx.logger.info('Kublet status: {0}'.format(zstatus))
    if not len(zstatus) > 1 and 'active' not in zstatus[1]:
        raise OperationRetry('Wait a little more.')
