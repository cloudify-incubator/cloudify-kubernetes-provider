from cloudify import ctx
from cloudify.state import ctx_parameters as inputs
import os

ctx.logger.info('Inputs : {0}'.format(repr(inputs)))

path = inputs['path']
os.system("umount '{}' 2>/dev/null > /dev/null".format(path))
