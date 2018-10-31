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

from cloudify import ctx
from cloudify.state import ctx_parameters as inputs

ctx.logger.info('Inputs : {0}'.format(repr(inputs)))

path = inputs['path']
os.system("mkdir -p '{}' 2>/dev/null > /dev/null".format(path))
os.system("dd if=/dev/zero of='{}.img' count=204800 2>/dev/null > /dev/null"
          .format(path))
os.system("mkfs.ext4 -F '{}.img'  2>/dev/null > /dev/null".format(path))
os.system("mount -o loop '{}.img' '{}' 2>/dev/null > /dev/null"
          .format(path, path))
