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
from cloudify.plugins import lifecycle


@workflow
def delete(ctx,
           scalable_entity_name,
           delta,
           scale_compute,
           removed_ids_exclude_hint,
           removed_ids_include_hint,
           ignore_failure=False,
           **kwargs):
    """Scales in/out the subgraph of node_or_group_name.

    If a node name is passed, and `scale_compute` is set to false, the
    subgraph will consist of all the nodes that are contained in the node and
    the node itself.
    If a node name is passed, and `scale_compute` is set to true, the subgraph
    will consist of all nodes that are contained in the compute node that
    contains the node and the compute node itself.
    If a group name or a node that is not contained in a compute
    node, is passed, this property is ignored.

    `delta` is used to specify the scale factor.
    For `delta > 0`: If current number of instances is `N`, scale out to
    `N + delta`.
    For `delta < 0`: If current number of instances is `N`, scale in to
    `N - |delta|`.

    :param ctx: cloudify context
    :param scalable_entity_name: the node or group name to scale
    :param delta: scale in/out factor
    :param scale_compute: should scale apply on compute node containing
                          the specified node
    :param ignore_failure: ignore operations failures in uninstall workflow
    """
    if isinstance(delta, basestring):
        try:
            delta = int(delta)
        except ValueError:
            raise ValueError('The delta parameter must be a number. Got: {0}'
                             .format(delta))

    ctx.logger.info(('Delete with scalable_entity_name: {}, delta: {}, ' +
                     'removed_ids_exclude_hint: {}, ' +
                     'removed_ids_include_hint: {}').format(
                         repr(scalable_entity_name),
                         repr(delta),
                         repr(removed_ids_exclude_hint),
                         repr(removed_ids_include_hint)
                     ))
    if delta == 0:
        ctx.logger.info('delta parameter is 0, so no scaling will take place.')
        return

    scaling_group = ctx.deployment.scaling_groups.get(scalable_entity_name)
    if scaling_group:
        curr_num_instances = scaling_group['properties']['current_instances']
        planned_num_instances = curr_num_instances + delta
        scale_id = scalable_entity_name
    else:
        node = ctx.get_node(scalable_entity_name)
        if not node:
            raise ValueError("No scalable entity named {0} was found".format(
                scalable_entity_name))
        host_node = node.host_node
        scaled_node = host_node if (scale_compute and host_node) else node
        curr_num_instances = scaled_node.number_of_instances
        planned_num_instances = curr_num_instances + delta
        scale_id = scaled_node.id

    if planned_num_instances < 0:
        raise ValueError('Provided delta: {0} is illegal. current number of '
                         'instances of entity {1} is {2}'
                         .format(delta,
                                 scalable_entity_name,
                                 curr_num_instances))
    modification = ctx.deployment.start_modification({
        scale_id: {
            'instances': planned_num_instances,

            # These following parameters are not exposed at the moment,
            # but should be used to control which node instances get scaled in
            # (when scaling in).
            # They are mentioned here, because currently, the modification API
            # is not very documented.
            # Special care should be taken because if `scale_compute == True`
            # (which is the default), then these ids should be the compute node
            # instance ids which are not necessarily instances of the node
            # specified by `scalable_entity_name`.

            'removed_ids_exclude_hint': removed_ids_exclude_hint,
            # Node instances denoted by these instance ids should be *kept* if
            # possible.
            # 'removed_ids_exclude_hint': [],

            'removed_ids_include_hint': removed_ids_include_hint,
            # Node instances denoted by these instance ids should be *removed*
            # if possible.
            # 'removed_ids_include_hint': []
        }
    })
    graph = ctx.graph_mode()
    try:
        ctx.logger.info('Deployment modification started. '
                        '[modification_id={0}]'.format(modification.id))
        if delta > 0:
            added_and_related = set(modification.added.node_instances)
            added = set(i for i in added_and_related
                        if i.modification == 'added')
            related = added_and_related - added
            try:
                lifecycle.install_node_instances(
                    graph=graph,
                    node_instances=added,
                    related_nodes=related)
            except Exception as ex:
                ctx.logger.error('Scale out failed, scaling back in. {}'
                                 .format(repr(ex)))
                for task in graph.tasks_iter():
                    graph.remove_task(task)
                lifecycle.uninstall_node_instances(
                    graph=graph,
                    node_instances=added,
                    ignore_failure=ignore_failure,
                    related_nodes=related)
                raise ex
        else:
            removed_and_related = set(modification.removed.node_instances)
            removed = set(i for i in removed_and_related
                          if i.modification == 'removed')
            related = removed_and_related - removed
            lifecycle.uninstall_node_instances(
                graph=graph,
                node_instances=removed,
                ignore_failure=ignore_failure,
                related_nodes=related)
    except Exception as ex:
        ctx.logger.warn('Rolling back deployment modification. '
                        '[modification_id={0}]: {}'
                        .format(modification.id, repr(ex)))
        try:
            modification.rollback()
        except Exception as ex:
            ctx.logger.warn('Deployment modification rollback failed. The '
                            'deployment model is most likely in some corrupted'
                            ' state. [modification_id={0}]: {}'
                            .format(modification.id, repr(ex)))
            raise ex
        raise ex
    else:
        try:
            modification.finish()
        except Exception as ex:
            ctx.logger.warn('Deployment modification finish failed. The '
                            'deployment model is most likely in some corrupted'
                            ' state.[modification_id={0}]: {}'
                            .format(modification.id, repr(ex)))
            raise ex
