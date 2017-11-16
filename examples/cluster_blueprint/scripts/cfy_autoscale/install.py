#!/usr/bin/env python

import os
import subprocess
from cloudify import ctx
from cloudify.exceptions import (
    HttpException,
    NonRecoverableError,
    OperationRetry
)
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

    ctx.logger.info('Downloading or building cfy-autoscale.')

    cfy_autoscale_binary_path = \
        inputs.get(
            'cfy_autoscale_binary_path',
            '/usr/bin/cfy-autoscale')

    execute_command(
        'sudo mkdir -p /opt/cloudify-kubernetes-provider/src/'
        'k8s.io/autoscaler/cluster-autoscaler')

    if ctx.operation.retry_number == 1:

        try:
            cfy_autoscale_binary = \
                ctx.download_resource('resources/cfy-autoscale')
        except HttpException:
            ctx.logger.debug('Build cfy-autoscale.')
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
                'go build -v -o bin/cluster-autoscaler ' \
                'src/k8s.io/autoscaler/cluster-autoscaler/main.go ' \
                'src/k8s.io/autoscaler/cluster-autoscaler/version.go'
            execute_command(_command, extra_args=_extra_args)
            cfy_autoscale_binary = \
                '/opt/cloudify-kubernetes-provider/bin/cluster-autoscaler'
        ctx.logger.debug('cfy-autoscale built/downloaded.')
        execute_command(
            'sudo cp {0} {1}'.format(
                cfy_autoscale_binary, cfy_autoscale_binary_path))
        execute_command(
            'sudo chmod 555 {0}'.format(cfy_autoscale_binary_path))
        execute_command(
            'sudo chown root:root {0}'.format(cfy_autoscale_binary_path))

        try:
            _tv = {'home_dir': os.path.expanduser('~')}
            cfy_autoscale_service = \
                ctx.download_resource_and_render(
                    'resources/cfy-autoscale.service',
                    template_variables=_tv)
        except HttpException:
            raise NonRecoverableError(
                'cfy-autoscale.service not in resources.')
        else:
            execute_command(
                'sudo cp {0} {1}'.format(
                    cfy_autoscale_service,
                    '/etc/systemd/system/cfy-autoscale.service'))
            execute_command('sudo systemctl daemon-reload')
            execute_command('sudo systemctl enable cfy-autoscale.service')
            execute_command('sudo systemctl start cfy-autoscale.service')

    status = ''
    systemctl_status = \
        execute_command('sudo systemctl status cfy-autoscale.service')
    if not isinstance(systemctl_status, basestring):
        raise OperationRetry(
            'check sudo systemctl status cfy-autoscale.service')
    for line in systemctl_status.split('\n'):
        if 'Active:' in line:
            status = line.strip()
    zstatus = status.split(' ')
    ctx.logger.info('cfy-autoscale status: {0}'.format(zstatus))
    if not len(zstatus) > 1 and 'active' not in zstatus[1]:
        raise OperationRetry('Wait a little more.')
