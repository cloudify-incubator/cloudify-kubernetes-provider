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
        'stderr': subprocess.PIPE
    }
    if extra_args is not None and isinstance(extra_args, dict):
        subprocess_args.update(extra_args)

    ctx.logger.debug('subprocess_args {0}.'.format(subprocess_args))
    output = ""
    try:
        process = subprocess.Popen(**subprocess_args)
        output, error = process.communicate()

        ctx.logger.debug('command: {0} '.format(_command))
        ctx.logger.debug('output: {0} '.format(output))
        ctx.logger.debug('error: {0} '.format(error))
        ctx.logger.debug('process.returncode: {0} '.format(process.returncode))

        if process.returncode:
            ctx.logger.error('Running `{0}` returns error.'.format(_command))
            return False
    except Exception as e:
        ctx.logger.error(repr(e))
        raise e
    return output

    

    
    
    
ctx.logger.info("Copy config")
home = expanduser("~")
if not os.path.exists(home + "/.kube"):
    os.makedirs(home + "/.kube")

status = execute_command("sudo cp -v /etc/kubernetes/admin.conf {0}/.kube/config".format(home))
if status is False:
    raise OperationRetry("cp .kube/config failed")

userid = os.getuid()
groupid = os.getgid()

status = execute_command("sudo chown {0}:{1} {2}/.kube/config".format(userid, groupid, home))
if status is False:
    raise OperationRetry('chown operation failed')


ctx.logger.info("Apply network")
for retry_count in range (0, 9):
    proc = subprocess.Popen(["kubectl apply -f https://git.io/weave-kube-1.6"], stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    (out, err) = proc.communicate()
    if proc.returncode != 0:
        ctx.logger.info("#{}: Init network configuration failed?".format(retry_count))
        time.sleep(3)
    else:
        break

ctx.logger.info ("Install cfy-kubernetes provider")
status = execute_command("sudo cp -v /opt/cloudify-kubernetes-provider/bin/cfy-kubernetes /usr/bin/cfy-kubernetes")
if status is False:
    raise OperationRetry('Copy cfy-jubernetes failed')
    
status = execute_command("sudo chmod 555 /usr/bin/cfy-kubernetes")
if status is False:
    raise OperationRetry('chmod 555 cfy-kubernetes failed')

status = execute_command("sudo chown root:root /usr/bin/cfy-kubernetes")
if status is False:
    raise OperationRetry('chmod cfy-kubernetes failed')

ctx.logger.info("Create service")
cfy_kubernetes_temp = ctx.download_resource('resources/cfy-kubernetes.service')

status = execute_command('sudo mv {0} /usr/lib/systemd/system/cfy-kubernetes.service'.format(
                cfy_kubernetes_temp))
if status is False:
    raise OperationRetry('Failed to move cfy-kubernetes.service')

status = execute_command('sudo sed -i s|$HOME|{0}|g /usr/lib/systemd/system/cfy-kubernetes.service'.format(home))
    
ctx.logger.info("Start service")
status = execute_command("sudo systemctl daemon-reload")
if status is False:
    raise OperationRetry('daemon-reload failed')

status = execute_command("sudo systemctl enable cfy-kubernetes.service")
if status is False:
    raise OperationRetry('enable cfy-kubernetes.service failed')

status = execute_command("sudo systemctl start cfy-kubernetes.service")
if status is False:
    raise OperationRetry('start cfy-kubernetes.service failed')

for retry_count in range (0, 9):
    proc = subprocess.Popen(["sudo systemctl status cfy-kubernetes.service | grep 'Active:'| awk '{print $2}'"], stdout=subprocess.PIPE, shell=True)
    (out, err) = proc.communicate()
    ctx.logger.info ("#{}: Kubernetes state: {}".format(retry_count, out.strip()))
    if out.strip() in ['active']:
        break
    elif retry_count < 10:
        ctx.logger.info("Wait little more.")
        time.sleep(3)
    else:
        raise OperationRetry("Error: Service kubelet inactive.")

ctx.logger.info ("Get token")
try:
    proc = subprocess.Popen(["sudo kubeadm token list | grep authentication,signing | awk '{print $1}' | base64"], stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    (out, err) = proc.communicate()
    ctx.instance.runtime_properties['token'] = out
    ctx.logger.info("Token {}".format(out))
except Exception as e:
    ctx.logger.error(repr(e))
    raise e

