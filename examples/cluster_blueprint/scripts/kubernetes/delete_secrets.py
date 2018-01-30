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
