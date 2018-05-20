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

import sys
import os
import os.path
import json
import platform
import socket
import ssl
import subprocess
import tempfile

from cloudify import ctx
from cloudify.state import ctx_parameters as inputs
from cloudify import manager
from cloudify_rest_client.exceptions import CloudifyClientError
from cloudify.utils import exception_to_error_cause
from cloudify.exceptions import (
    HttpException,
    NonRecoverableError,
    OperationRetry
)

CONFIG = ('"deployment": "{0}",' +
          '"instance": "{1}",' +
          '"tenant": "{2}",' +
          '"password": "{3}",' +
          '"user": "{4}",' +
          '"host": "{5}",' +
          '"agent": "{6}"')


def generate_traceback_exception():
    _, exc_value, exc_traceback = sys.exc_info()
    response = exception_to_error_cause(exc_value, exc_traceback)
    return response


def execute_command(command, extra_args=None):

    ctx.logger.debug('command: {0}.'.format(repr(command)))

    subprocess_args = {
        'args': command,
        'stdout': subprocess.PIPE,
        'stderr': subprocess.PIPE
    }
    if extra_args is not None and isinstance(extra_args, dict):
        subprocess_args.update(extra_args)

    ctx.logger.debug('subprocess_args {0}.'.format(subprocess_args))

    process = subprocess.Popen(**subprocess_args)
    output, error = process.communicate()

    ctx.logger.debug('command: {0} '.format(repr(command)))
    ctx.logger.debug('output: {0} '.format(output))
    ctx.logger.debug('error: {0} '.format(error))
    ctx.logger.debug('process.returncode: {0} '.format(process.returncode))

    if process.returncode:
        ctx.logger.error('Running `{0}` returns {1} error: {2}.'
                         .format(repr(command), process.returncode,
                                 repr(error)))
        return False

    return output


def download_service(service_name):
    service_path = "/usr/bin/" + service_name
    if not os.path.isfile(service_path):
        try:
            cfy_binary = ctx.download_resource(
                'resources/{}'.format(service_name))
        except HttpException:
            raise NonRecoverableError(
                '{} binary not in resources.'.format(service_name))
        ctx.logger.debug('{} downloaded.'.format(service_name))
        if execute_command(['sudo', 'cp', cfy_binary, service_path]) is False:
            raise NonRecoverableError("Can't copy {}.".format(service_path))
    # fix file attributes
    if execute_command(['sudo', 'chmod', '555', service_path]) is False:
        raise NonRecoverableError("Can't chmod {}.".format(service_path))
    if execute_command(['sudo', 'chown', 'root:root', service_path]) is False:
        raise NonRecoverableError("Can't chown {}.".format(service_path))
    ctx.logger.debug('{} attributes fixed'.format(service_name))


def create_service(service_name):
    service_path = '/etc/systemd/system/{}.service'.format(service_name)
    if not os.path.isfile(service_path):
        try:
            _tv = {'home_dir': os.path.expanduser('~')}
            cfy_service = \
                ctx.download_resource_and_render(
                    'resources/{}.service'.format(service_name),
                    template_variables=_tv)
        except HttpException:
            raise NonRecoverableError(
                '{}.service not in resources.'.format(service_name))

        execute_command(['sudo', 'cp', cfy_service, service_path])
        execute_command(['sudo', 'cp', service_path,
                         '/etc/systemd/system/multi-user.target.wants/'])

        execute_command(['sudo', 'systemctl', 'daemon-reload'])
        execute_command(['sudo', 'systemctl', 'enable',
                         '{}.service'.format(service_name)])
        execute_command(['sudo', 'systemctl', 'start',
                         '{}.service'.format(service_name)])


def start_check(service_name):
    status_string = ''
    systemctl_status = execute_command(['sudo', 'systemctl', 'status',
                                        '{}.service'.format(service_name)])
    if not isinstance(systemctl_status, basestring):
        raise OperationRetry(
            'check sudo systemctl status {}.service'.format(service_name))
    for line in systemctl_status.split('\n'):
        if 'Active:' in line:
            status = line.strip()
            zstatus = status.split(' ')
            ctx.logger.debug('{} status line: {}'
                             .format(service_name, repr(zstatus)))
            if len(zstatus) > 1:
                status_string = zstatus[1]

    ctx.logger.info('{} status: {}'.format(service_name, repr(status_string)))
    if 'active' != status_string:
        raise OperationRetry('Wait a little more.')
    else:
        ctx.logger.info('Service {} is started.'.format(service_name))


def get_instance_host(relationships, rel_type, target_type):
    for rel in relationships:
        if rel.type == rel_type or rel_type in rel.type_hierarchy:
            if target_type in rel.target.node.type_hierarchy:
                return rel.target.instance
            instance = get_instance_host(rel.target.instance.relationships,
                                         rel_type, target_type)
            if instance:
                return instance
    return None


def update_host_address(host_instance, hostname, fqdn, ip, public_ip):
    ctx.logger.info('Setting initial Kubernetes node data')

    if not public_ip:
        public_ip_prop = host_instance.runtime_properties.get(
            'public_ip')
        public_ip_address_prop = host_instance.runtime_properties.get(
            'public_ip_address')
        public_ip = public_ip_prop or public_ip_address_prop or ip

    new_runtime_properties = {
        'name': ctx.instance.id,
        'hostname': hostname,
        'fqdn': fqdn,
        'ip': ip,
        'public_ip': public_ip
    }

    for key, value in new_runtime_properties.items():
        ctx.instance.runtime_properties[key] = value

    ctx.logger.info(
        'Finished setting initial Kubernetes node data.')


def get_kubernetes_deployment_info(relationships, rel_type, target_type):
    deployments = []
    for rel in relationships:
        if rel.type == rel_type or rel_type in rel.type_hierarchy:
            ctx.logger.info(
                'type_hierarchy: {0}'.format(rel.target.node.type_hierarchy))

            if target_type in rel.target.node.type_hierarchy:
                # Get deployment id from deployment proxy node
                dep_id = rel.target.\
                    node.properties['resource_config']['deployment']['id']

                ctx.logger.info('Deployment id: {0}'.format(dep_id))

                # Get deployment type from deployment proxy node outputs
                dep_outputs = rel.target.\
                    instance.runtime_properties['deployment']['outputs']

                ctx.logger.info('Deployment Info: {0}'.format(dep_outputs))

                if type(dep_outputs) is dict:
                    ctx.logger.info('Deployment Dict: {0}'.format(dep_outputs))
                    for key, value in dep_outputs.items():
                        ctx.logger.info(
                            'Key: {0}, Value: {1}'.format(key, value))

                # Deployment type (Load Balancer or Node)
                dep_type = dep_outputs['deployment-type']
                # Node Data Type (kubernetes Node or kubernetes load balancer)
                node_data_type = dep_outputs['deployment-node-data-type']

                ctx.logger.info(
                    'Deployment node data type: {0}'.format(node_data_type))

                ctx.logger.info(
                    'Deployment type: {0}'.format(dep_type))

                deployment_info = dict()
                deployment_info['id'] = dep_id
                deployment_info['deployment_type'] = dep_type
                deployment_info['node_data_type'] = node_data_type
                deployments.append(deployment_info)

    return deployments


def generate_deployment_file():
    data = dict()
    # Get the deployments info for  Nodes && Load balancers
    deployments = get_kubernetes_deployment_info(
        ctx.instance.relationships,
        'cloudify.relationships.depends_on',
        'cloudify.nodes.DeploymentProxy')

    if deployments:
        deployment_file = os.path.expanduser('~') + "/deployment.json"
        data['deployments'] = deployments

        if not os.path.isfile(deployment_file):
            ctx.logger.info("Create deployment {}"
                            " file".format(deployment_file))
            with open(deployment_file, 'w') as json_file:
                json.dump(data, json_file)
        return

    raise NonRecoverableError('Unable to generate deployment file, '
                              'deployment ids are empty !!!.'
                              'Please check your kubernetes.yaml')


def setup_kubernetes_node_data_type():
    ctx.logger.debug(
        'Setup kubernetes node data '
        'type for deployment id {0}'.format(ctx.deployment.id))

    cfy_client = manager.get_rest_client()
    try:
        response = cfy_client.deployments.outputs.get(ctx.deployment.id)

    except CloudifyClientError as ex:
        ctx.logger.debug(
            'Unable to parse outputs for deployment'
            ' {0}'.format(ctx.deployment.id))

        raise OperationRetry('Re-try getting deployment outputs again.')

    except Exception:
        response = generate_traceback_exception()

        ctx.logger.error(
            'Error traceback {0} with message {1}'.format(
                response['traceback'], response['message']))

        raise NonRecoverableError("Failed to get outputs")

    else:
        dep_outputs = response.get('outputs')
        ctx.logger.debug('Deployment outputs: {0}'.format(dep_outputs))
        node_data_type = dep_outputs.get('deployment-node-data-type')

        if node_data_type:
            os.environ['CFY_K8S_NODE_TYPE'] = node_data_type

        else:
            os.environ['CFY_K8S_NODE_TYPE'] =\
                'cloudify.nodes.ApplicationServer.kubernetes.Node'


if __name__ == '__main__':

    # create global config
    config_file = os.path.expanduser('~') + "/cfy.json"

    host_instance = get_instance_host(ctx.instance.relationships,
                                      'cloudify.relationships.contained_in',
                                      'cloudify.nodes.Compute')

    if not host_instance:
        raise NonRecoverableError('Ambiguous host resolution data.')

    cloudify_agent = host_instance.runtime_properties.get('cloudify_agent', {})

    linux_distro = cloudify_agent.get('distro')
    cfy_host = cloudify_agent.get('broker_ip')
    cfy_ssl_port = cloudify_agent.get('rest_port')
    agent_name = cloudify_agent.get('name')

    cfy_user = inputs.get('cfy_user', 'admin')
    cfy_pass = inputs.get('cfy_password', 'admin')
    cfy_tenant = inputs.get('cfy_tenant', 'default_tenant')
    agent_user = inputs.get('agent_user', 'centos')
    full_install = inputs.get('full_install', 'all')

    # Generate deployment file to be used as a reference inside the cfy.json
    generate_deployment_file()

    if not os.path.isfile(config_file):
        ctx.logger.info("Create config {} file".format(config_file))
        deployment_file = os.path.expanduser('~') + "/deployment.json"

        # services config
        with open(config_file, 'w') as outfile:
            agent_file = "/root" if agent_user == "root" else (
                "/home/" + agent_user
            )
            cfy_host_full = cfy_host if not cfy_ssl_port else (
                    "https://" + cfy_host + ":" + str(cfy_ssl_port)
            )
            outfile.write("{" + CONFIG.format(
                deployment_file,
                ctx.instance.id,
                cfy_tenant,
                cfy_pass,
                cfy_user,
                cfy_host_full,
                "{}/.cfy-agent/{}.json".format(agent_file, agent_name)
            ) + "}")

    if ctx.operation.retry_number == 0:
        # Allow user to provide specific values.
        update_host_address(
            host_instance=host_instance,
            hostname=inputs.get('hostname', socket.gethostname()),
            fqdn=inputs.get('fqdn', socket.getfqdn()),
            ip=inputs.get('ip', ctx.instance.host_ip),
            public_ip=inputs.get('public_ip'))

        # certificate logic
        if not linux_distro:
            distro, _, _ = \
                platform.linux_distribution(full_distribution_name=False)
            linux_distro = distro.lower()

        ctx.logger.info("Set certificate as trusted")

        # cert config
        _, temp_cert_file = tempfile.mkstemp()

        with open(temp_cert_file, 'w') as cert_file:
            cert_file.write("# cloudify certificate\n")
            try:
                cert_file.write(ssl.get_server_certificate((
                    cfy_host, cfy_ssl_port)))
            except Exception as ex:
                ctx.logger.error("Check https connection to manager {}."
                                 .format(str(ex)))

        if 'centos' in linux_distro:
            execute_command([
                'sudo', 'cp', temp_cert_file,
                '/etc/pki/ca-trust/source/anchors/cloudify.crt'
            ])
            execute_command([
                'sudo', 'update-ca-trust', 'extract'
            ])
            execute_command([
                'sudo', 'bash', '-c',
                'cat {} >> /etc/pki/tls/certs/ca-bundle.crt'
                .format(temp_cert_file)
            ])
        else:
            raise NonRecoverableError('Unsupported platform.')

    # download cfy-go tools
    if full_install != "loadbalancer":
        # Update the os environment variable to be used by the cfy-go diag
        setup_kubernetes_node_data_type()

        # Download cfy-go service
        download_service("cfy-go")

        # Run the diag command with option ``-node`` which check the
        # kubernetes nodes
        try:
            output = execute_command([
                '/usr/bin/cfy-go', 'status', 'diag',
                '-tenant', cfy_tenant, '-password', cfy_pass,
                '-user', cfy_user, '-host', cfy_host_full,
                '-agent-file', "{}/.cfy-agent/{}.json"
                .format(agent_file, agent_name)])
            ctx.logger.info("Diagnostic: {}".format(output))

        except Exception:
            response = generate_traceback_exception()

            ctx.logger.error(
                'Error traceback {0} with message {1}'.format(
                    response['traceback'], response['message']))

            raise NonRecoverableError("Failed to run daig command")

    if full_install == "all":
        # download cluster provider
        download_service("cfy-kubernetes")
        create_service("cfy-kubernetes")

        # download scale tools
        download_service("cfy-autoscale")
        create_service("cfy-autoscale")

        start_check("cfy-kubernetes")
        start_check("cfy-autoscale")