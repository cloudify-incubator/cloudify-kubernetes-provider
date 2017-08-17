package cloudifyprovider

import (
	"cloudify"
	"io"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	"log"
	"os"
)

const (
	providerName = "cloudify"
)

// CloudProvider implents Instances, Zones, and LoadBalancer
type CloudProvider struct {
	client *cloudify.CloudifyClient
}

// Initialize passes a Kubernetes clientBuilder interface to the cloud provider
func (r *CloudProvider) Initialize(clientBuilder controller.ControllerClientBuilder) {
	log.Println("Initialize")
}

// ProviderName returns the cloud provider ID.
func (r *CloudProvider) ProviderName() string {
	return providerName
}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	log.Println("LoadBalancer")
	return nil, false
}

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) Zones() (cloudprovider.Zones, bool) {
	log.Println("Zones")
	return nil, false
}

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) Instances() (cloudprovider.Instances, bool) {
	log.Println("Instances")
	return nil, false
}

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) Clusters() (cloudprovider.Clusters, bool) {
	log.Println("Clusters")
	return nil, false
}

// Routes returns a routes interface along with whether the interface is supported.
func (r *CloudProvider) Routes() (cloudprovider.Routes, bool) {
	log.Println("Routers")
	return nil, false
}

// HasClusterID returns true if a ClusterID is required and set
func (r *CloudProvider) HasClusterID() bool {
	return false
}

// ScrubDNS provides an opportunity for cloud-provider-specific code to process DNS settings for pods.
func (r *CloudProvider) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nameservers, searches
}

func newCloudifyCloud(config io.Reader) (cloudprovider.Interface, error) {
	log.Println("New Cloudify client")
	return &CloudProvider{
		client: cloudify.NewClient(os.Getenv("CFY_HOST"), os.Getenv("CFY_USER"), os.Getenv("CFY_PASSWORD"), os.Getenv("CFY_TENANT")),
	}, nil
}

func init() {
	log.Println("Cloudify init")
	cloudprovider.RegisterCloudProvider(providerName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newCloudifyCloud(config)
	})
}
