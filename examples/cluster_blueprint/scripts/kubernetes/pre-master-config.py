#!/usr/bin/env python
#
# Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import subprocess
import time
from cloudify import ctx
from cloudify.exceptions import OperationRetry
import re

JOIN_COMMAND_REGEX = '^kubeadm join[\sA-Za-z0-9\.\:\-\_]*'
BOOTSTRAP_TOKEN_REGEX = '[a-z0-9]{6}.[a-z0-9]{16}'
BOOTSTRAP_HASH_REGEX = '^sha256:[a-z0-9]{64}'
IP_PORT_REGEX = '[0-9]+(?:\.[0-9]+){3}:[0-9]+'
JCRE_COMPILED = re.compile(JOIN_COMMAND_REGEX)
BTRE_COMPILED = re.compile(BOOTSTRAP_TOKEN_REGEX)
BHRE_COMPILED = re.compile(BOOTSTRAP_HASH_REGEX)
IPRE_COMPILED = re.compile(IP_PORT_REGEX)


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


def setup_secrets(_split_master_port, _bootstrap_token, _bootstrap_hash):
    master_ip = split_master_port[0]
    master_port = split_master_port[1]
    ctx.instance.runtime_properties['master_ip'] = _split_master_port[0]
    ctx.instance.runtime_properties['master_port'] = _split_master_port[1]
    ctx.instance.runtime_properties['bootstrap_token'] = _bootstrap_token
    ctx.instance.runtime_properties['bootstrap_hash'] = _bootstrap_hash
    from cloudify import manager
    cfy_client = manager.get_rest_client()

    _secret_key = 'kubernetes_master_ip'
    if cfy_client and not len(cfy_client.secrets.list(key=_secret_key)) == 1:
        cfy_client.secrets.create(key=_secret_key, value=master_ip)
    else:
        cfy_client.secrets.update(key=_secret_key, value=master_ip)
    ctx.logger.info('Set secret: {0}.'.format(_secret_key))

    _secret_key = 'kubernetes_master_port'
    if cfy_client and not len(cfy_client.secrets.list(key=_secret_key)) == 1:
        cfy_client.secrets.create(key=_secret_key, value=master_port)
    else:
        cfy_client.secrets.update(key=_secret_key, value=master_port)
    ctx.logger.info('Set secret: {0}.'.format(_secret_key))

    _secret_key = 'bootstrap_token'
    if cfy_client and not len(cfy_client.secrets.list(key=_secret_key)) == 1:
        cfy_client.secrets.create(key=_secret_key, value=_bootstrap_token)
    else:
        cfy_client.secrets.update(key=_secret_key, value=_bootstrap_token)
    ctx.logger.info('Set secret: {0}.'.format(_secret_key))

    _secret_key = 'bootstrap_hash'
    if cfy_client and not len(cfy_client.secrets.list(key=_secret_key)) == 1:
        cfy_client.secrets.create(key=_secret_key, value=_bootstrap_hash)
    else:
        cfy_client.secrets.update(key=_secret_key, value=_bootstrap_hash)
    ctx.logger.info('Set secret: {0}.'.format(_secret_key))


def cleanup_and_retry():
    reset_cluster_command = 'sudo kubeadm reset'
    output = execute_command(reset_cluster_command)
    ctx.logger.info('reset_cluster_command {1}'
                    .format(reset_cluster_command, output))
    raise OperationRetry('Restarting kubernetes because of a problem.')


if __name__ == '__main__':
    ctx.logger.info("Reload kubeadm")
    status = execute_command('sudo systemctl daemon-reload')
    if status is False:
        raise OperationRetry('Failed daemon-reload')

    restart_service = execute_command('sudo systemctl stop kubelet')
    if restart_service is False:
        raise OperationRetry('Failed to stop kubelet')

    time.sleep(5)

    restart_service = execute_command('sudo systemctl start kubelet')

    for retry_count in range(10):
        proc = subprocess.Popen(
            ["sudo systemctl status kubelet | "
             "grep 'Active:'| awk '{print $2}'"],
            stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True
        )
        (out, err) = proc.communicate()
        ctx.logger.info("#{}: Kubelet state: {}".format(retry_count, out))
        if out.strip() in ['active']:
            break
        ctx.logger.info("Wait little more.")
        time.sleep(5)
    else:
        raise OperationRetry("Error: Service kubelet inactive.")

    ctx.logger.info("Init kubeadm")
    # echo 1 | sudo tee /proc/sys/net/bridge/bridge-nf-call-iptables
    status = execute_command(
        "sudo sysctl net.bridge.bridge-nf-call-iptables=1")
    if status is False:
        raise OperationRetry('Failed to set bridge-nf-call-iptables')

    status = execute_command('sudo kubeadm reset')
    if status is False:
        raise OperationRetry('sudo kubeadm reset failed')

    start_output = execute_command(
        'sudo kubeadm init --pod-network-cidr 10.244.0.0/16 --token-ttl 0'
    )
    if start_output is False:
        raise OperationRetry('kubeadm init failed')

    # Check if start succeeded.
    if start_output is False or not isinstance(start_output, basestring):
        ctx.logger.error('Kubernetes master failed to start.')
        cleanup_and_retry()
    ctx.logger.info('Kubernetes master started successfully.')

    # Slice and dice the start_master_command start_output.
    ctx.logger.info('Attempting to retrieve Kubernetes cluster information.')
    split_start_output = \
        [line.strip() for line in start_output.split('\n') if line.strip()]
    del line

    ctx.logger.debug(
        'Kubernetes master start output, split and stripped: {0}'.format(
            split_start_output))
    split_join_command = ''
    for li in split_start_output:
        ctx.logger.debug('li in split_start_output: {0}'.format(li))
        if re.match(JCRE_COMPILED, li):
            split_join_command = re.split('\s', li)
    del li
    ctx.logger.info('split_join_command: {0}'.format(split_join_command))

    if not split_join_command:
        ctx.logger.error('No join command in split_start_output: {0}'
                         .format(split_join_command))
        cleanup_and_retry()

    for li in split_join_command:
        ctx.logger.info('Sorting bits and pieces: li: {0}'.format(li))
        if re.match(BHRE_COMPILED, li):
            bootstrap_hash = li
        elif re.match(BTRE_COMPILED, li):
            bootstrap_token = li
        elif re.match(IPRE_COMPILED, li):
            split_master_port = li.split(':')
    setup_secrets(split_master_port, bootstrap_token, bootstrap_hash)

    ctx.logger.info("Reload kubeadm")
    status = execute_command(
        'sudo sed -i s|admission-control=Initializers,NamespaceLifecycle,'
        'LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,'
        'DefaultTolerationSeconds,NodeRestriction,ResourceQuota|'
        'admission-control=Initializers,NamespaceLifecycle,LimitRanger,'
        'ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,'
        'NodeRestriction,ResourceQuota|g '
        '/etc/kubernetes/manifests/kube-apiserver.yaml'
    )

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

    for retry_count in range(10):
        proc = subprocess.Popen(
            ["sudo systemctl status kubelet | "
             "grep 'Active:'| awk '{print $2}'"],
            stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True
        )
        (out, err) = proc.communicate()
        ctx.logger.info("#{}: Kubelet state: {}".format(retry_count, out))
        if out.strip() in ['active']:
            break
        ctx.logger.info("Wait little more.")
        time.sleep(5)
    else:
        raise OperationRetry("Error: Service kubelet inactive.")
