from cloudify import ctx
from cloudify.state import ctx_parameters as inputs

ctx.logger.info('Inputs : {0}'.format(repr(inputs)))
