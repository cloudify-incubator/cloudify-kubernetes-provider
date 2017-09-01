package cloudifyprovider

import (
	"fmt"
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/types"
	api "k8s.io/kubernetes/pkg/api/v1"
)

type CloudifyIntances struct {
	client *cloudify.CloudifyClient
}

// NodeAddresses returns the addresses of the specified instance.
// This implementation only returns the address of the calling instance. This is ok
// because the gce implementation makes that assumption and the comment for the interface
// states it as a todo to clarify that it is only for the current host
func (r *CloudifyIntances) NodeAddresses(nodeName types.NodeName) ([]api.NodeAddress, error) {
	name := string(nodeName)
	glog.Infof("NodeAddresses [%s]", name)
	addresses := []api.NodeAddress{}

	// TODO: Use real ip's
	hostIps := []string{"127.0.0.1"}

	for _, ip := range hostIps {
		addresses = append(addresses, api.NodeAddress{
			Type:    api.NodeExternalIP,
			Address: ip,
		})
		addresses = append(addresses, api.NodeAddress{
			Type:    api.NodeInternalIP,
			Address: ip,
		})
	}

	return addresses, nil
}

// NodeAddressesByProviderID returns the node addresses of an instances with the specified unique providerID
// This method will not be called from the node that is requesting this ID. i.e. metadata service
// and other local methods cannot be used here
func (r *CloudifyIntances) NodeAddressesByProviderID(providerID string) ([]api.NodeAddress, error) {
	glog.Infof("NodeAddressesByProviderID [%s]", providerID)
	return []api.NodeAddress{}, fmt.Errorf("Not implemented:NodeAddressesByProviderID")
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (r *CloudifyIntances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return fmt.Errorf("Not implemented:AddSSHKeyToAllInstances")
}

// CurrentNodeName returns the name of the node we are currently running on
func (r *CloudifyIntances) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

// ExternalID returns the cloud provider ID of the specified instance (deprecated).
func (r *CloudifyIntances) ExternalID(nodeName types.NodeName) (string, error) {
	name := string(nodeName)
	glog.Infof("ExternalID [%s]", name)
	return r.InstanceID(nodeName)
}

const fakeuuid = "fakeuuid:"

// ExternalID returns the cloud provider ID of the specified instance (deprecated).
func (r *CloudifyIntances) InstanceID(nodeName types.NodeName) (string, error) {
	name := string(nodeName)
	glog.Infof("InstanceID [%s]", name)

	// TODO: return real uuid?
	return fakeuuid + name, nil
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
	glog.Infof("InstanceTypeByProviderID [%s]", providerID)
	return "", fmt.Errorf("Not implemented:InstanceTypeByProviderID")
}

func NewCloudifyIntances(client *cloudify.CloudifyClient) *CloudifyIntances {
	return &CloudifyIntances{
		client: client,
	}
}
