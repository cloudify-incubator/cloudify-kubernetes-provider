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
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/types"
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

// GetZoneByProviderID implements Zones.GetZoneByProviderID
// This is particularly useful in external cloud providers where the kubelet
// does not initialize node data.
func (r *CloudifyZones) GetZoneByProviderID(providerID string) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, fmt.Errorf("GetZoneByProviderID not implemented")
}

// GetZoneByNodeName implements Zones.GetZoneByNodeName
// This is particularly useful in external cloud providers where the kubelet
// does not initialize node data.
func (r *CloudifyZones) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, fmt.Errorf("GetZoneByNodeName not imeplemented")
}

func NewCloudifyZones(client *cloudify.CloudifyClient) *CloudifyZones {
	return &CloudifyZones{
		client: client,
	}
}
