#!/usr/bin/env python
# -*- coding: utf-8 -*-

#
# Copyright (c) 2016 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import logging

import ovirtsdk4 as sdk
import ovirtsdk4.types as types

logging.basicConfig(level=logging.DEBUG, filename='example.log')

# This example will connect to the server and print the names and identifiers of all the virtual machines:

# Create the connection to the server:
connection = sdk.Connection(
    url='https://master.tmax.dom/ovirt-engine/api',
    username='admin@internal',
    password='asdfasdf',
    ca_file='../ca.crt',
    debug=True,
    log=logging.getLogger(),
)

if 

# Get the reference to the "vms" service:
vms_service = connection.system_service().vms_service()

# Use the "list" method of the "vms" service to list all the virtual machines of the system:
vms = vms_service.list()

# Print the virtual machine names and identifiers:
for vm in vms:
  print("%s: %s" % (vm.name, vm.id))

# Close the connection to the server:
connection.close()