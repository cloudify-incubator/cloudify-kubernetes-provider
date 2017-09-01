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
	return fmt.Errorf("Not implemented")
}

// GetLoadBalancer is an implementation of LoadBalancer.GetLoadBalancer
func (r *CloudifyBalancer) GetLoadBalancer(clusterName string, service *api.Service) (status *api.LoadBalancerStatus, exists bool, retErr error) {
	glog.Infof("GetLoadBalancer [%s]", clusterName)
	return &api.LoadBalancerStatus{}, false, fmt.Errorf("Not implemented")
}

// EnsureLoadBalancerDeleted is an implementation of LoadBalancer.EnsureLoadBalancerDeleted.
func (r *CloudifyBalancer) EnsureLoadBalancerDeleted(clusterName string, service *api.Service) error {
	glog.Infof("EnsureLoadBalancerDeleted [%s]", clusterName)
	return fmt.Errorf("Not implemented")
}

// EnsureLoadBalancer is an implementation of LoadBalancer.EnsureLoadBalancer.
func (r *CloudifyBalancer) EnsureLoadBalancer(clusterName string, service *api.Service, nodes []*api.Node) (*api.LoadBalancerStatus, error) {
	glog.Infof("EnsureLoadBalancer [%s]", clusterName)
	return &api.LoadBalancerStatus{}, fmt.Errorf("Not implemented")
}

func NewCloudifyBalancer(client *cloudify.CloudifyClient) *CloudifyBalancer {
	return &CloudifyBalancer{
		client: client,
	}
}
