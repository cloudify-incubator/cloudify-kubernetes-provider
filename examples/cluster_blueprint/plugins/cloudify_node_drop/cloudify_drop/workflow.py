# Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
from cloudify.decorators import workflow
from cloudify.plugins.lifecycle import uninstall_node_instances


@workflow
def delete(ctx, **kwargs):
    """Default uninstall workflow"""

    instance_ids = kwargs.get("instance_ids", [])
    ignore_failure = kwargs.get("ignore_failure", False)

    if not instance_ids:
        ctx.logger.info("No instnaces for delete, skip")
        return

    ctx.logger.info("Delete {} instnaces with {} flag"
                    .format(repr(instance_ids), repr(ignore_failure)))
    uninstall_node_instances(
        graph=ctx.graph_mode(),
        node_instances=set([
            instance for instance in ctx.node_instances if (
                instance.id in instance_ids or (
                    instance._node_instance.host_id and
                    instance._node_instance.host_id in instance_ids
                )
            )
        ]),
        ignore_failure=ignore_failure)
