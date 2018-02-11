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

import os
from os.path import expanduser
import subprocess
import time
from cloudify import ctx
from cloudify.exceptions import OperationRetry


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
    output = ""
    try:
        process = subprocess.Popen(**subprocess_args)
        output, error = process.communicate()

        ctx.logger.debug('command: {0} '.format(_command))
        ctx.logger.debug('output: {0} '.format(output))
        ctx.logger.debug('error: {0} '.format(error))
        ctx.logger.debug('process.returncode: {0} '.format(process.returncode))

        if process.returncode:
            ctx.logger.error('Running `{0}` returns error.'.format(_command))
            return False
    except Exception as e:
        ctx.logger.error(repr(e))
        raise e
    return output


ctx.logger.info("Copy config")
home = expanduser("~")
if not os.path.exists(home + "/.kube"):
    os.makedirs(home + "/.kube")

status = execute_command(
    "sudo cp -v /etc/kubernetes/admin.conf {0}/.kube/config".format(home)
)
if status is False:
    raise OperationRetry("cp .kube/config failed")

userid = os.getuid()
groupid = os.getgid()

status = execute_command(
    "sudo chown {0}:{1} {2}/.kube/config".format(userid, groupid, home)
)
if status is False:
    raise OperationRetry('chown operation failed')


ctx.logger.info("Apply network")
for retry_count in range(10):
    proc = subprocess.Popen(
        ["kubectl apply -f https://git.io/weave-kube-1.6"],
        stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True
    )
    (out, err) = proc.communicate()
    if proc.returncode != 0:
        ctx.logger.info(
            "#{}: Init network configuration failed?".format(retry_count)
        )
        time.sleep(5)
    else:
        break

ctx.logger.info("Get token")
try:
    proc = subprocess.Popen(
        [
            "sudo kubeadm token list | grep authentication,signing | "
            "awk '{print $1}' | base64"
        ],
        stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    (out, err) = proc.communicate()
    ctx.instance.runtime_properties['token'] = out
    ctx.logger.info("Token {}".format(out))
except Exception as e:
    ctx.logger.error(repr(e))
    raise e
