#!/usr/bin/env python
import os
from cloudify import ctx
try:
    import yaml
except ImportError:
    ctx.logger.info("Need to install yaml package")
    import pip
    pip.main(['install', 'pyyaml'])
    import yaml

if __name__ == '__main__':

    admin_file_dest = os.path.join(os.path.expanduser('~'), '.kube/config')

    with open(admin_file_dest, 'r') as outfile:
        try:
            kubernetes_master_config = yaml.load(outfile)
        except yaml.YAMLError as e:
            RecoverableError(
                'Unable to read Kubernetes Admin file: {0}: {1}'.format(
                    admin_file_dest, str(e)))

    ctx.instance.runtime_properties['configuration_file_content'] = \
        kubernetes_master_config
