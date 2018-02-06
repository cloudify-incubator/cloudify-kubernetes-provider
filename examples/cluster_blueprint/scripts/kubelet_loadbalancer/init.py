#!/usr/bin/env python
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
from cloudify import ctx
from cloudify.exceptions import NonRecoverableError
from cloudify.state import ctx_parameters as inputs
from tempfile import NamedTemporaryFile

HAPROXY_HEADER = """
global
    log 127.0.0.1 local0 notice
    user haproxy
    group haproxy
defaults
    log global
    retries 2
    timeout connect 3000
    timeout server 5000
    timeout client 5000
listen stats 0.0.0.0:9000
    mode http
    balance
    timeout client 5000
    timeout connect 4000
    timeout server 30000
    stats uri /haproxy_stats
    stats realm HAProxy\ Statistics
    stats auth admin:password
    stats admin if TRUE
"""


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
    ctx.logger.info('Inputs : {0}'.format(repr(inputs)))
    runtime_properties = ctx.instance.runtime_properties

    runtime_properties["proxy_ports"] = inputs.get("ports")
    runtime_properties["proxy_nodes"] = inputs.get("nodes")
    runtime_properties["proxy_cluster"] = inputs.get("cluster")
    runtime_properties["proxy_name"] = inputs.get("name")
    runtime_properties["proxy_namespace"] = inputs.get("namespace")

    execute_command("sudo systemctl stop haproxy")

    config_text = HAPROXY_HEADER
    for port in inputs.get("ports", []):
        in_port = int(port['port'])
        out_port = int(port['nodePort'])

        config_text = (config_text + "\nfrontend {}\n\toption forceclose\n\t"
                       "bind *:{}\n\tdefault_backend servers_{}_{}\n"
                       "backend servers_{}_{}\n\toption forceclose"
                       .format(ctx.instance.id, in_port,
                               in_port, out_port, in_port, out_port))

        for node in inputs.get("nodes", []):
            config_text = (config_text + "\n\tserver {} {}:{} maxconn 32"
                           .format(node['hostname'], node['ip'], out_port))

    # Get the HAProxy config file path.
    haproxy_cfg_path = inputs.get('haproxy.cfg', '/etc/haproxy/haproxy.cfg')
    # Render the template and write the rendered file to a temporary file.
    with NamedTemporaryFile(delete=False) as temp_config:
        temp_config.write(config_text)
    # Test the temporary file.
    out = execute_command('sudo /usr/sbin/haproxy -f {0} -c'
                          .format(temp_config.name))
    if not out:
        raise NonRecoverableError('Invalid config.')
    # Replace the HAProxy configuration file with the temporary file.
    execute_command('sudo cp {0} {1}'
                    .format(temp_config.name, haproxy_cfg_path))
    execute_command('sudo chmod {0} {1}'
                    .format('0600', haproxy_cfg_path))
    # Reload the HAProxy process.
    execute_command("sudo systemctl start haproxy")
