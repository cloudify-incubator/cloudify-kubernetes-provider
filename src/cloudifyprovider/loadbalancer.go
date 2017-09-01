package cloudifyprovider

import (
	"fmt"
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"github.com/golang/glog"
	api "k8s.io/kubernetes/pkg/api/v1"
)

type CloudifyBalancer struct {
	client *cloudify.CloudifyClient
}

// UpdateLoadBalancer is an implementation of LoadBalancer.UpdateLoadBalancer.
func (r *CloudifyBalancer) UpdateLoadBalancer(clusterName string, service *api.Service, nodes []*api.Node) error {
	glog.Infof("UpdateLoadBalancer [%s]", clusterName)
	return fmt.Errorf("Not implemented:UpdateLoadBalancer")
}

func (r *CloudifyBalancer) toLBStatus(service_id string) (*api.LoadBalancerStatus, bool, error) {
	ingress := []api.LoadBalancerIngress{}

	// TODO: show real id
	ingress = append(ingress, api.LoadBalancerIngress{IP: "127.0.0.1"})

	return &api.LoadBalancerStatus{ingress}, true, nil
}

// GetLoadBalancer is an implementation of LoadBalancer.GetLoadBalancer
func (r *CloudifyBalancer) GetLoadBalancer(clusterName string, service *api.Service) (status *api.LoadBalancerStatus, exists bool, retErr error) {
	glog.Infof("GetLoadBalancer [%s]", clusterName)
	return r.toLBStatus(clusterName)
}

// EnsureLoadBalancerDeleted is an implementation of LoadBalancer.EnsureLoadBalancerDeleted.
func (r *CloudifyBalancer) EnsureLoadBalancerDeleted(clusterName string, service *api.Service) error {
	glog.Infof("EnsureLoadBalancerDeleted [%s]", clusterName)

	// TODO: We can delete anything from unexisted services :-)
	return nil
}

// EnsureLoadBalancer is an implementation of LoadBalancer.EnsureLoadBalancer.
func (r *CloudifyBalancer) EnsureLoadBalancer(clusterName string, service *api.Service, nodes []*api.Node) (*api.LoadBalancerStatus, error) {
	glog.Infof("EnsureLoadBalancer [%s]", clusterName)
	status, _, err := r.toLBStatus(clusterName)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func NewCloudifyBalancer(client *cloudify.CloudifyClient) *CloudifyBalancer {
	return &CloudifyBalancer{
		client: client,
	}
}
