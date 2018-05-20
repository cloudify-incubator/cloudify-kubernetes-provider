#!/usr/bin/env python

import pwd
import grp
import os
import re
import getpass
import subprocess
import pip

try:
    import yaml
except ImportError:
    pip.main(['install', 'pyyaml'])
    import yaml

from cloudify import manager
from cloudify import ctx
from cloudify.state import ctx_parameters as inputs
from cloudify.exceptions import OperationRetry
from cloudify.exceptions import RecoverableError


JOIN_COMMAND_REGEX = '^kubeadm join[\sA-Za-z0-9\.\:\-\_]*'
BOOTSTRAP_TOKEN_REGEX = '[a-z0-9]{6}.[a-z0-9]{16}'
BOOTSTRAP_HASH_REGEX = '^sha256:[a-z0-9]{64}'
IP_PORT_REGEX = '[0-9]+(?:\.[0-9]+){3}:[0-9]+'
JCRE_COMPILED = re.compile(JOIN_COMMAND_REGEX)
BTRE_COMPILED = re.compile(BOOTSTRAP_TOKEN_REGEX)
BHRE_COMPILED = re.compile(BOOTSTRAP_HASH_REGEX)
IPRE_COMPILED = re.compile(IP_PORT_REGEX)


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


def cleanup_and_retry():
    reset_cluster_command = 'sudo kubeadm reset'
    output = execute_command(reset_cluster_command)
    ctx.logger.info('reset_cluster_command {1}'.format(reset_cluster_command, output))
    raise OperationRetry('Restarting kubernetes because of a problem.')


def configure_admin_conf():
    # Add the kubeadmin config to environment
    agent_user = getpass.getuser()
    uid = pwd.getpwnam(agent_user).pw_uid
    gid = grp.getgrnam('docker').gr_gid
    admin_file_dest = os.path.join(os.path.expanduser('~'), 'admin.conf')

    execute_command('sudo cp {0} {1}'.format('/etc/kubernetes/admin.conf', admin_file_dest))
    execute_command('sudo chown {0}:{1} {2}'.format(uid, gid, admin_file_dest))

    with open(os.path.join(os.path.expanduser('~'), '.bashrc'), 'a') as outfile:
        outfile.write('export KUBECONFIG=$HOME/admin.conf')
    os.environ['KUBECONFIG'] = admin_file_dest


def store_kubernetes_master_config():
    cfy_client = manager.get_rest_client()
    admin_file_dest = os.path.join(os.path.expanduser('~'), 'admin.conf')

    # Storing the K master configuration.
    kubernetes_master_config = {}
    with open(admin_file_dest, 'r') as outfile:
        try:
            kubernetes_master_config = yaml.load(outfile)
        except yaml.YAMLError as e:
            RecoverableError(
                'Unable to read Kubernetes Admin file: {0}: {1}'.format(
                    admin_file_dest, str(e)))
    ctx.instance.runtime_properties['configuration_file_content'] = \
        kubernetes_master_config

    clusters = kubernetes_master_config.get('clusters')
    _clusters = {}
    for cluster in clusters:
        __name = cluster.get('name')
        _cluster = cluster.get('cluster', {})
        _secret_key = '%s_certificate_authority_data' % __name
        if cfy_client and not len(
                cfy_client.secrets.list(key=_secret_key)) == 1:
            cfy_client.secrets.create(key=_secret_key, value=_cluster.get(
                'certificate-authority-data'))
            ctx.logger.info('Set secret: {0}.'.format(_secret_key))
        else:
            cfy_client.secrets.update(key=_secret_key, value=_cluster.get(
                'certificate-authority-data'))
        ctx.instance.runtime_properties[
            '%s_certificate_authority_data' % __name] = _cluster.get(
            'certificate-authority-data')
        _clusters[__name] = _cluster
    del __name

    contexts = kubernetes_master_config.get('contexts')
    _contexts = {}
    for context in contexts:
        __name = context.get('name')
        _context = context.get('context', {})
        _contexts[__name] = _context
    del __name

    users = kubernetes_master_config.get('users')
    _users = {}
    for user in users:
        __name = user.get('name')
        _user = user.get('user', {})
        _secret_key = '%s_client_certificate_data' % __name
        if cfy_client and not len(
                cfy_client.secrets.list(key=_secret_key)) == 1:
            cfy_client.secrets.create(key=_secret_key, value=_user.get(
                'client-certificate-data'))
            ctx.logger.info('Set secret: {0}.'.format(_secret_key))
        else:
            cfy_client.secrets.update(key=_secret_key, value=_user.get(
                'client-certificate-data'))
        _secret_key = '%s_client_key_data' % __name
        if cfy_client and not len(
                cfy_client.secrets.list(key=_secret_key)) == 1:
            cfy_client.secrets.create(key=_secret_key,
                                      value=_user.get('client-key-data'))
            ctx.logger.info('Set secret: {0}.'.format(_secret_key))
        else:
            cfy_client.secrets.update(key=_secret_key,
                                      value=_user.get('client-key-data'))
        ctx.instance.runtime_properties[
            '%s_client_certificate_data' % __name] = _user.get(
            'client-certificate-data')
        ctx.instance.runtime_properties[
            '%s_client_key_data' % __name] = _user.get('client-key-data')
        _users[__name] = _user
    del __name

    ctx.instance.runtime_properties['kubernetes'] = {
        'clusters': _clusters,
        'contexts': _contexts,
        'users': _users
    }


def setup_kubernetes_bootstrap_data(start_output):
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
        ctx.logger.error('No join command in split_start_output: {0}'.format(
            split_join_command))
        cleanup_and_retry()

    for li in split_join_command:
        ctx.logger.info('Sorting bits and pieces: li: {0}'.format(li))
        if re.match(BHRE_COMPILED, li):
            bootstrap_hash = li
        elif re.match(BTRE_COMPILED, li):
            bootstrap_token = li
        elif re.match(IPRE_COMPILED, li):
            split_master_port = li.split(':')

    # setup as cloudify secrets
    setup_secrets(split_master_port, bootstrap_token, bootstrap_hash)


def start_kubernetes_master():
    # Start Kubernetes Master
    ctx.logger.info('Attempting to start Kubernetes master.')
    init_command = 'sudo kubeadm init'
    cni_provider = inputs.get('cni-provider-blueprint', 'weave.yaml')

    # Each cni provider work with specific pod-cidr
    if cni_provider == 'flannel.yaml':
        init_command = '{0} {1}'.format(init_command,
                                        '--pod-network-cidr=10.244.0.0/16')

    start_output = execute_command(init_command)
    ctx.logger.debug('start_master_command output: {0}'.format(start_output))
    # Check if start succeeded.
    if start_output is False or not isinstance(start_output, basestring):
        ctx.logger.error('Kubernetes master failed to start.')
        cleanup_and_retry()
    ctx.logger.info('Kubernetes master started successfully.')
    return start_output


def setup_secrets(_split_master_port, _bootstrap_token, _bootstrap_hash):
    master_ip = _split_master_port[0]
    master_port = _split_master_port[1]

    ctx.instance.runtime_properties['master_ip'] = master_ip
    ctx.instance.runtime_properties['master_port'] = master_port
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


def set_iptables():
    # echo 1 | sudo tee /proc/sys/net/bridge/bridge-nf-call-iptables
    status = execute_command(
        "sudo sysctl net.bridge.bridge-nf-call-iptables=1")
    if status is False:
        raise OperationRetry('Failed to set bridge-nf-call-iptables')


if __name__ == '__main__':

    ctx.instance.runtime_properties['KUBERNETES_MASTER'] = True

    # Set bridge-nf-call-iptables
    set_iptables()

    # Start kubernetes master
    start_output = start_kubernetes_master()

    # Setup kubernetes master bootstrap data
    setup_kubernetes_bootstrap_data(start_output)

    # configure kubernetes admin config
    configure_admin_conf()

    # Storing Kubernetes master configuration
    store_kubernetes_master_config()
