#!/usr/bin/env python

import subprocess
import socket
import time
from cloudify import ctx
from cloudify.exceptions import OperationRetry


def check_command(command):

    ctx.logger.debug('command {0}.'.format(repr(command)))

    subprocess_args = {
        'args': command,
        'stdout': subprocess.PIPE,
        'stderr': subprocess.PIPE
    }

    ctx.logger.debug('subprocess_args {0}.'.format(repr(subprocess_args)))

    process = subprocess.Popen(**subprocess_args)
    output, error = process.communicate()

    ctx.logger.debug('output: {0} '.format(repr(output)))
    ctx.logger.debug('error: {0} '.format(repr(error)))
    ctx.logger.debug('process.returncode: {0} '.format(repr(process.returncode)))

    if process.returncode:
        ctx.logger.error('Running `{0}` returns error.'.format(command))
        return False

    return True


def execute_command(command):

    ctx.logger.debug('command {0}.'.format(repr(command)))

    subprocess_args = {
        'args': command,
        'stdout': subprocess.PIPE,
        'stderr': subprocess.PIPE
    }

    ctx.logger.debug('subprocess_args {0}.'.format(repr(subprocess_args)))

    process = subprocess.Popen(**subprocess_args)
    output, error = process.communicate()

    ctx.logger.debug('output: {0} '.format(repr(output)))
    ctx.logger.debug('error: {0} '.format(repr(error)))
    ctx.logger.debug('process.returncode: {0} '.format(repr(process.returncode)))

    if process.returncode:
        ctx.logger.error('Running `{0}` returns error.'.format(repr(command)))
        return False

    return output


if __name__ == '__main__':

    # Check if Docker PS works
    docker = check_command(['sudo', 'docker', 'ps'])
    if not docker:
            raise OperationRetry(
                'Docker is not present on the system.')
    ctx.logger.info('Docker is present on the system.')

    # Next check if Cloud Init is running.
    finished = False
    ps = execute_command(['ps', '-ef'])
    for line in ps.split('\n'):
        if '/usr/bin/python /usr/bin/cloud-init modules' in line:
            raise OperationRetry(
                'You provided a Cloud-init Cloud Config to configure '
                'instances. Waiting for Cloud-init to complete.')
    ctx.logger.info('Cloud-init finished.')

    execute_command(["sudo", "sed", "-i", "s|cgroup-driver=systemd|"
                     "cgroup-driver=systemd --provider-id='{}'|g"
                     .format(socket.gethostname()),
                     "/etc/systemd/system/kubelet.service.d/10-kubeadm.conf"])

    ctx.logger.info("Reload kubeadm")
    status = execute_command(["sudo", "systemctl", "daemon-reload"])
    if status is False:
        raise OperationRetry('Failed daemon-reload')

    restart_service = execute_command(["sudo", "systemctl", "stop", "kubelet"])
    if restart_service is False:
        raise OperationRetry('Failed to stop kubelet')

    time.sleep(5)

    restart_service = execute_command(
       ["sudo", "systemctl", "start", "kubelet"])
    if restart_service is False:
        raise OperationRetry('Failed to start kubelet')

    time.sleep(5)
