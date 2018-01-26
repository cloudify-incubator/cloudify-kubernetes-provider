#!/usr/bin/env python

import os
from os.path import expanduser
import re
import subprocess
import sys
import time
from cloudify import ctx
from cloudify.state import ctx_parameters as inputs
from cloudify.exceptions import NonRecoverableError, OperationRetry


def execute_command(_command, extra_args=None):

    ctx.logger.debug('_command {0}.'.format(_command))

    subprocess_args = {
        'args': _command.split(),
        'stdout': subprocess.PIPE,
        'stderr': subprocess.PIPE,
        'shell': False
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


ctx.logger.info("Reload kubeadm")
status = execute_command('sudo systemctl daemon-reload')
if status is False:
    raise OperationRetry('Failed daemon-reload')


restart_service = execute_command('sudo systemctl stop kubelet')
if restart_service is False:
    raise OperationRetry('Failed to stop kubelet')

time.sleep(5)

restart_service = execute_command('sudo systemctl start kubelet')

for retry_count in range (0, 9):
    proc = subprocess.Popen(["sudo systemctl status kubelet | grep 'Active:'| awk '{print $2}'"], stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    (out, err) = proc.communicate()
    ctx.logger.info ("#{}: Kubelet state: {}".format(retry_count, out))
    if out.strip() in ['active']:
        break
    elif retry_count < 10:
        ctx.logger.info("Wait little more.")
        time.sleep(3)
    else:
        raise OperationRetry("Error: Service kubelet inactive.")
    
ctx.logger.info("Init kubeadm")
status = execute_command("sudo sysctl net.bridge.bridge-nf-call-iptables=1")   #("echo 1 | sudo tee /proc/sys/net/bridge/bridge-nf-call-iptables")
if status is False:
    raise OperationRetry('Failed to set bridge-nf-call-iptables')

status = execute_command('sudo kubeadm reset')
if status is False:
    raise OperationRetry('sudo kubeadm reset failed')

status = execute_command('sudo kubeadm init --pod-network-cidr 10.244.0.0/16 --token-ttl 0')
if status is False:
    raise OperationRetry('kubeadm init failed')

ctx.logger.info("Reload kubeadm")
status = execute_command('sudo sed -i s|admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota|admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota|g /etc/kubernetes/manifests/kube-apiserver.yaml')

status = execute_command("sudo systemctl daemon-reload")
if status is False:
    raise OperationRetry('daemon-reload failed')

status = execute_command('sudo systemctl stop kubelet')
if status is False:
    raise OperationRetry('kubelet stop failed')

time.sleep(5)

status = execute_command('sudo systemctl start kubelet')
if status is False:
    raise OperationRetry('kubelet start failed')
    

for retry_count in range (0, 9):
    proc = subprocess.Popen(["sudo systemctl status kubelet | grep 'Active:'| awk '{print $2}'"], stdout=subprocess.PIPE, stderr=subprocess.PIPE,  shell=True)
    (out, err) = proc.communicate()
    ctx.logger.info ("#{}: Kubelet state: {}".format(retry_count, out))
    if out.strip() in ['active']:
        break
    elif retry_count < 10:
        ctx.logger.info("Wait little more.")
        time.sleep(3)
    else:
        raise OperationRetry("Error: Service kubelet inactive.")
