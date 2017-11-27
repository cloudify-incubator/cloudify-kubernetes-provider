#!/usr/bin/env python
import subprocess
from cloudify import ctx
from cloudify.exceptions import OperationRetry


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
    copy_repo = execute_command(
        "kubectl create -f https://raw.githubusercontent.com/kubernetes/" +
        "heapster/release-1.4/deploy/kube-config/rbac/heapster-rbac.yaml")
    if copy_repo is False:
        raise OperationRetry(
            'Failed to apply heapster_rbac')
