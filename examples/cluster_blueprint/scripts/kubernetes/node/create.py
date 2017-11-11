#!/usr/bin/env python

import socket
from cloudify import ctx
from cloudify.exceptions import NonRecoverableError
from cloudify.state import ctx_parameters as inputs

TYPE = 'cloudify.nodes.Compute'
REL_TYPE = 'cloudify.relationships.contained_in'


def get_instance_by_target_type(_target_id, _target_type):
    """Get the node instance object by node target type
    and ID for extra granularity.
    """

    ctx.logger.debug(
        'Attempting to resolve host with id {0} or type {1}'.format(
            _target_id, _target_type))

    for rel in ctx.instance.relationships:
        # This is the most exact.
        if rel.target.instance.id == _target_id:
            return rel.target.instance
        # Equally exact.
        elif rel.type == REL_TYPE:
            return rel.target.instance
        # Equally exact.
        elif REL_TYPE in rel.type_hierarchy:
            return rel.target.instance
        # This is a possibility, though problematic.
        elif _target_type in rel.target.node.type_hierarchy:
            return rel.target.instance
        # Get out of here.
        else:
            raise NonRecoverableError(
                'Ambiguous host resolution data.')


if __name__ == '__main__':

    ctx.logger.info(
        'Setting initial Kubernetes node data')

    # Allow user to provide specific values.
    target_id = inputs.get('target_id')
    target_type = inputs.get('target_type', TYPE)
    hostname = inputs.get('hostname', socket.gethostname())
    fqdn = inputs.get('fqdn', socket.getfqdn())
    ip = inputs.get('ip', ctx.instance.host_ip)
    public_ip = inputs.get('public_ip')

    if not public_ip:
        host_instance = \
            get_instance_by_target_type(
                target_id, target_type)
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
