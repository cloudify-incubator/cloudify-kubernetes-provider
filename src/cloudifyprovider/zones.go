package cloudifyprovider

import (
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type CloudifyZones struct {
	client *cloudify.CloudifyClient
}

// GetZone is an implementation of Zones.GetZone
func (r *CloudifyZones) GetZone() (cloudprovider.Zone, error) {
	glog.Infof("GetZone")
	return cloudprovider.Zone{
		FailureDomain: "FailureDomain",
		Region:        "Region",
	}, nil
}

func NewCloudifyZones(client *cloudify.CloudifyClient) *CloudifyZones {
	return &CloudifyZones{
		client: client,
	}
}
