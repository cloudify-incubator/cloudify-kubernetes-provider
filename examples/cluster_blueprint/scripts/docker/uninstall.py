#!/usr/bin/env python

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


if __name__ == '__main__':

    ctx.logger.info(
        'Verifying that Docker is not installed on the system.')

    existing_install = \
        ctx.instance.runtime_properties.get('existing_docker_install', False)

    if not existing_install:
        stopped = execute_command('sudo systemctl stop docker')
        if not stopped:
            OperationRetry('Failed to stop docker.service')
        uninstalled = execute_command('sudo yum remove -y -q docker-engine')
        if not uninstalled:
            OperationRetry('Failed to uninstall Docker')
        execute_command('sudo rm /etc/yum.repos.d/docker.repo')
