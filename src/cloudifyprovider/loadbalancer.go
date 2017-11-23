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
	api "k8s.io/api/core/v1"
)

// Balancer - struct with connection settings
type Balancer struct {
	client *cloudify.Client
}

// UpdateLoadBalancer is an implementation of LoadBalancer.UpdateLoadBalancer.
func (r *Balancer) UpdateLoadBalancer(clusterName string, service *api.Service, nodes []*api.Node) error {
	glog.Errorf("?UpdateLoadBalancer [%s]", clusterName)
	return fmt.Errorf("Not implemented:UpdateLoadBalancer")
}

func (r *Balancer) toLBStatus(serviceID string) (*api.LoadBalancerStatus, bool, error) {
	glog.Errorf("?toLBStatus [%s]", serviceID)
	ingress := []api.LoadBalancerIngress{}

	// TODO: show real id
	ingress = append(ingress, api.LoadBalancerIngress{IP: "127.0.0.1"})

	return &api.LoadBalancerStatus{ingress}, true, nil
}

// GetLoadBalancer is an implementation of LoadBalancer.GetLoadBalancer
func (r *Balancer) GetLoadBalancer(clusterName string, service *api.Service) (status *api.LoadBalancerStatus, exists bool, retErr error) {
	glog.Errorf("?GetLoadBalancer [%s]", clusterName)
	return r.toLBStatus(clusterName)
}

// EnsureLoadBalancerDeleted is an implementation of LoadBalancer.EnsureLoadBalancerDeleted.
func (r *Balancer) EnsureLoadBalancerDeleted(clusterName string, service *api.Service) error {
	glog.Errorf("?EnsureLoadBalancerDeleted [%s]", clusterName)

	// TODO: We can delete anything from unexisted services :-)
	return nil
}

// EnsureLoadBalancer is an implementation of LoadBalancer.EnsureLoadBalancer.
func (r *Balancer) EnsureLoadBalancer(clusterName string, service *api.Service, nodes []*api.Node) (*api.LoadBalancerStatus, error) {
	glog.Errorf("?EnsureLoadBalancer [%s]", clusterName)
	status, _, err := r.toLBStatus(clusterName)
	if err != nil {
		return nil, err
	}

	return status, nil
}

// NewBalancer - create instance with support kubernetes balancer interface.
func NewBalancer(client *cloudify.Client) *Balancer {
	return &Balancer{
		client: client,
	}
}
