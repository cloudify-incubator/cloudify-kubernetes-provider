package cloudifyprovider

import (
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
)

type CloudifyIntances struct {
	client *cloudify.CloudifyClient
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (r *CloudifyIntances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return fmt.Errorf("Not implemented")
}

// CurrentNodeName returns the name of the node we are currently running on
func (r *CloudifyIntances) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func NewCloudifyIntances(client *cloudify.CloudifyClient) *CloudifyIntances{
	return &CloudifyIntances{
		client: client,
	}
}
