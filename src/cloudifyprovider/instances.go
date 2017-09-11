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
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/types"
	api "k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type CloudifyIntances struct {
	deployment string
	client     *cloudify.CloudifyClient
}

// NodeAddresses returns the addresses of the specified instance.
// This implementation only returns the address of the calling instance. This is ok
// because the gce implementation makes that assumption and the comment for the interface
// states it as a todo to clarify that it is only for the current host
func (r *CloudifyIntances) NodeAddresses(nodeName types.NodeName) ([]api.NodeAddress, error) {
	name := string(nodeName)
	glog.Infof(">NodeAddresses [%s]", name)

	var params = map[string]string{}
	nodeInstances := r.client.GetNodeInstances(params)

	addresses := []api.NodeAddress{}

	for _, nodeInstance := range nodeInstances.Items {
		// skip different deployments
		if nodeInstance.DeploymentId != r.deployment {
			continue
		}

		// skip nodes without ip's
		if nodeInstance.NodeId != "kubeinstance" && nodeInstance.NodeId != "kubemanager" {
			continue
		}

		// check runtime properties
		if nodeInstance.RuntimeProperties != nil {
			if v, ok := nodeInstance.RuntimeProperties["name"]; ok == true {
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
		glog.Infof("InstanceNotFound: %+v", name)
		return nil, cloudprovider.InstanceNotFound
	} else {
		glog.Infof("Addresses: %+v", addresses)
	}
	return addresses, nil
}

// NodeAddressesByProviderID returns the node addresses of an instances with the specified unique providerID
// This method will not be called from the node that is requesting this ID. i.e. metadata service
// and other local methods cannot be used here
func (r *CloudifyIntances) NodeAddressesByProviderID(providerID string) ([]api.NodeAddress, error) {
	glog.Infof(">NodeAddressesByProviderID [%s]", providerID)

	var params = map[string]string{}
	nodeInstances := r.client.GetNodeInstances(params)

	addresses := []api.NodeAddress{}

	for _, nodeInstance := range nodeInstances.Items {
		// skip different deployments
		if nodeInstance.DeploymentId != r.deployment {
			continue
		}

		// skip nodes without ip's
		if nodeInstance.NodeId != "kubeinstance" && nodeInstance.NodeId != "kubemanager" {
			continue
		}

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

	glog.Infof("Addresses: %+v", addresses)
	return addresses, nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (r *CloudifyIntances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	glog.Infof("?AddSSHKeyToAllInstances [%s]", user)
	return fmt.Errorf("Not implemented:AddSSHKeyToAllInstances")
}

// CurrentNodeName returns the name of the node we are currently running on
func (r *CloudifyIntances) CurrentNodeName(hostname string) (types.NodeName, error) {
	glog.Infof("?CurrentNodeName [%s]", hostname)
	return types.NodeName(hostname), nil
}

// ExternalID returns the cloud provider ID of the specified instance (deprecated).
func (r *CloudifyIntances) ExternalID(nodeName types.NodeName) (string, error) {
	name := string(nodeName)
	glog.Infof("?ExternalID [%s]", name)
	return r.InstanceID(nodeName)
}

const fakeuuid = "fakeuuid:"

// InstanceID returns the cloud provider ID of the specified instance.
func (r *CloudifyIntances) InstanceID(nodeName types.NodeName) (string, error) {
	name := string(nodeName)
	glog.Infof("InstanceID [%s]", name)

	var params = map[string]string{}
	nodeInstances := r.client.GetNodeInstances(params)

	for _, nodeInstance := range nodeInstances.Items {
		// skip different deployments
		if nodeInstance.DeploymentId != r.deployment {
			continue
		}

		// skip nodes without ip's
		if nodeInstance.NodeId != "kubeinstance" && nodeInstance.NodeId != "kubemanager" {
			continue
		}

		// check runtime properties
		if nodeInstance.RuntimeProperties != nil {
			if v, ok := nodeInstance.RuntimeProperties["name"]; ok == true {
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
			if nodeInstance.State == "started" {
				glog.Infof("Node is alive %+v", nodeInstance)
				return fakeuuid + name, nil
			}
		}
	}

	glog.Infof("Node died: %+v", name)

	return "", cloudprovider.InstanceNotFound
}

// InstanceType returns the type of the specified instance.
// Note that if the instance does not exist or is no longer running, we must return ("", cloudprovider.InstanceNotFound)
func (r *CloudifyIntances) InstanceType(nodeName types.NodeName) (string, error) {
	_, err := r.InstanceID(nodeName)
	if err != nil {
		return "", err
	}
	return providerName, nil
}

// InstanceTypeByProviderID returns the cloudprovider instance type of the node with the specified unique providerID
// This method will not be called from the node that is requesting this ID. i.e. metadata service
// and other local methods cannot be used here
func (r *CloudifyIntances) InstanceTypeByProviderID(providerID string) (string, error) {
	glog.Infof("?InstanceTypeByProviderID [%s]", providerID)
	return "", fmt.Errorf("Not implemented:InstanceTypeByProviderID")
}

func NewCloudifyIntances(client *cloudify.CloudifyClient, deployment string) *CloudifyIntances {
	return &CloudifyIntances{
		client:     client,
		deployment: deployment,
	}
}
