from cloudify import ctx
from cloudify.state import ctx_parameters as inputs
import os

ctx.logger.info('Inputs : {0}'.format(repr(inputs)))

path = inputs['path']
os.system("mkdir -p '{}' 2>/dev/null > /dev/null".format(path))
os.system("dd if=/dev/zero of='{}.img' count=204800 2>/dev/null > /dev/null".format(path))
os.system("mkfs.ext4 -F '{}.img'  2>/dev/null > /dev/null".format(path))
os.system("mount -o loop '{}.img' '{}' 2>/dev/null > /dev/null".format(path, path))
