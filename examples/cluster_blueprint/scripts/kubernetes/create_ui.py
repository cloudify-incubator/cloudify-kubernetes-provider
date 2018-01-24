#!/usr/bin/env python

import subprocess

from cloudify import ctx
from cloudify import manager


def execute_command(_command):

    ctx.logger.debug('_command {0}.'.format(_command))

    subprocess_args = {
        'args': _command.split(),
        'stdout': subprocess.PIPE,
        'stderr': subprocess.PIPE
    }

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

    # Get Cfy Manager Python Rest Client.
    cfy_client = manager.get_rest_client()

    # Try to call kubectl and deploy UI to the k8s
    ui_deploy = execute_command('kubectl create -f https://raw.githubusercontent.com/kubernetes/'
                                'dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml')

    ctx.logger.debug('K8S UI deploy output: {0} '.format(ui_deploy))
