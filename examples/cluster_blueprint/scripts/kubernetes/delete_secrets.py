#!/usr/bin/env python

from cloudify import ctx
from cloudify import manager

if __name__ == '__main__':
    cfy_client = manager.get_rest_client()
    secrets_keys = ['kubernetes-admin_client_certificate_data',
                    'kubernetes-admin_client_key_data',
                    'kubernetes_certificate_authority_data',
                    'kubernetes_master_ip',
                    'kubernetes_master_port']
    for _secret_key in secrets_keys:
        cfy_client.secrets.delete(key=_secret_key)
        ctx.logger.info('Unset secret: {0}.'.format(_secret_key))
