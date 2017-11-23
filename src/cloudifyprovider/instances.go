/*
Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloudifyprovider

import (
	"fmt"
	cloudify "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify"
	utils "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/utils"
	"github.com/golang/glog"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// Instances - struct with connection settings
type Instances struct {
	deployment string
	client     *cloudify.Client
}

func (r *Instances) getInstances(params map[string]string) []cloudify.NodeInstance {
	// Add filter by deployment
	params["deployment_id"] = r.deployment

	nodeInstances, err := r.client.GetNodeInstancesWithType(
		params, "cloudify.nodes.ApplicationServer.kubernetes.Node")
	if err != nil {
		glog.Infof("Not found instances: %+v", err)
		return []cloudify.NodeInstance{}
	}

	// starting only because we restart kubelet after join
	aliveStates := []string{
		// "initializing", "creating", // workflow started for instance
		// "created", "configuring", // create action, had ip
		"configured", "starting", // configure action, joined to cluster
		"started", // everything done
	}
	instances := []cloudify.NodeInstance{}
	for _, instance := range nodeInstances.Items {
		if utils.InList(aliveStates, instance.State) {
			instances = append(instances, instance)
		}
	}
	return instances
}

// NodeAddresses returns the addresses of the specified instance.
// This implementation only returns the address of the calling instance. This is ok
// because the gce implementation makes that assumption and the comment for the interface
// states it as a todo to clarify that it is only for the current host
func (r *Instances) NodeAddresses(nodeName types.NodeName) ([]api.NodeAddress, error) {
	name := string(nodeName)
	glog.Infof(">NodeAddresses [%s]", name)

	var params = map[string]string{}
	nodeInstances := r.getInstances(params)

	addresses := []api.NodeAddress{}

	for _, nodeInstance := range nodeInstances {
		// check runtime properties
		if nodeInstance.RuntimeProperties != nil {
			if v, ok := nodeInstance.RuntimeProperties["hostname"]; ok == true {
				switch v.(type) {
				case string:
					{
						if v.(string) != name {
							// node with different name
							continue
						} else {
							addresses = append(addresses, api.NodeAddress{
								Type:    api.NodeHostName,
								Address: v.(string),
							})
						}
					}
				}
			} else {
				// node without name
				continue
			}

			if v, ok := nodeInstance.RuntimeProperties["ip"]; ok == true {
				switch v.(type) {
				case string:
					{
						addresses = append(addresses, api.NodeAddress{
							Type:    api.NodeInternalIP,
							Address: v.(string),
						})
					}
				}
			}

			if v, ok := nodeInstance.RuntimeProperties["public_ip"]; ok == true {
				switch v.(type) {
				case string:
					{
						addresses = append(addresses, api.NodeAddress{
							Type:    api.NodeExternalIP,
							Address: v.(string),
						})
					}
				}
			}
		}
	}

	if len(addresses) == 0 {
		glog.Infof("NodeAddresses: InstanceNotFound: %+v", name)
		return nil, cloudprovider.InstanceNotFound
	}

	glog.Infof("NodeAddresses: %+v", addresses)
	return addresses, nil
}

// NodeAddressesByProviderID returns the node addresses of an instances with the specified unique providerID
// This method will not be called from the node that is requesting this ID. i.e. metadata service
// and other local methods cannot be used here
func (r *Instances) NodeAddressesByProviderID(providerID string) ([]api.NodeAddress, error) {
	glog.Infof(">NodeAddressesByProviderID [%s]", providerID)

	var params = map[string]string{}
	nodeInstances := r.getInstances(params)

	addresses := []api.NodeAddress{}

	// hostname
	for _, nodeInstance := range nodeInstances {
		// check runtime properties
		if nodeInstance.RuntimeProperties != nil {
			if v, ok := nodeInstance.RuntimeProperties["hostname"]; ok == true {
				switch v.(type) {
				case string:
					{
						addresses = append(addresses, api.NodeAddress{
							Type:    api.NodeHostName,
							Address: v.(string),
						})
					}
				}
			}
		}
	}

	// internal ip
	for _, nodeInstance := range nodeInstances {
		// check runtime properties
		if nodeInstance.RuntimeProperties != nil {
			if v, ok := nodeInstance.RuntimeProperties["ip"]; ok == true {
				switch v.(type) {
				case string:
					{
						addresses = append(addresses, api.NodeAddress{
							Type:    api.NodeInternalIP,
							Address: v.(string),
						})
					}
				}
			}
		}
	}

	// external ip
	for _, nodeInstance := range nodeInstances {
		// check runtime properties
		if nodeInstance.RuntimeProperties != nil {
			if v, ok := nodeInstance.RuntimeProperties["public_ip"]; ok == true {
				switch v.(type) {
				case string:
					{
						addresses = append(addresses, api.NodeAddress{
							Type:    api.NodeExternalIP,
							Address: v.(string),
						})
					}
				}
			}
		}
	}

	glog.Infof("NodeAddressesByProviderID: %+v", addresses)
	return addresses, nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (r *Instances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	glog.Infof("?AddSSHKeyToAllInstances [%s]", user)
	return fmt.Errorf("Not implemented:AddSSHKeyToAllInstances")
}

// CurrentNodeName returns the name of the node we are currently running on
func (r *Instances) CurrentNodeName(hostname string) (types.NodeName, error) {
	glog.Infof("?CurrentNodeName [%s]", hostname)
	return types.NodeName(hostname), nil
}

// ExternalID returns the cloud provider ID of the specified instance (deprecated).
func (r *Instances) ExternalID(nodeName types.NodeName) (string, error) {
	name := string(nodeName)
	glog.Infof("?ExternalID [%s]", name)
	return r.InstanceID(nodeName)
}

// InstanceID returns the cloud provider ID of the specified instance.
func (r *Instances) InstanceID(nodeName types.NodeName) (string, error) {
	name := string(nodeName)
	glog.Infof("InstanceID [%s]", name)

	var params = map[string]string{}
	nodeInstances := r.getInstances(params)

	for _, nodeInstance := range nodeInstances {
		// check runtime properties
		if nodeInstance.RuntimeProperties != nil {
			if v, ok := nodeInstance.RuntimeProperties["hostname"]; ok == true {
				switch v.(type) {
				case string:
					{
						if v.(string) != name {
							// node with different name
							continue
						}
					}
				}
			} else {
				// node without name
				continue
			}

			glog.Infof("Node is alive: %+v", name)
			return name, nil
		}
	}

	glog.Infof("Node died: %+v", name)

	return "", cloudprovider.InstanceNotFound
}

// InstanceType returns the type of the specified instance.
// Note that if the instance does not exist or is no longer running, we must return ("", cloudprovider.InstanceNotFound)
func (r *Instances) InstanceType(nodeName types.NodeName) (string, error) {
	_, err := r.InstanceID(nodeName)
	if err != nil {
		return "", err
	}
	return providerName, nil
}

// InstanceTypeByProviderID returns the cloudprovider instance type of the node with the specified unique providerID
// This method will not be called from the node that is requesting this ID. i.e. metadata service
// and other local methods cannot be used here
func (r *Instances) InstanceTypeByProviderID(providerID string) (string, error) {
	glog.Infof("?InstanceTypeByProviderID [%s]", providerID)
	return "", fmt.Errorf("Not implemented:InstanceTypeByProviderID")
}

// InstanceExistsByProviderID returns true if the instance with the given provider id still exists and is running.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
func (r *Instances) InstanceExistsByProviderID(providerID string) (bool, error) {
	return false, fmt.Errorf("Not implemented:InstanceExistsByProviderID")
}

// NewInstances - create instance with support kubernetes intances interface.
func NewInstances(client *cloudify.Client, deployment string) *Instances {
	return &Instances{
		client:     client,
		deployment: deployment,
	}
}
