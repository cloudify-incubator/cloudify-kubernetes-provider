package cloudifyprovider

import (
	"encoding/json"
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"github.com/golang/glog"
	"io"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	"os"
)

const (
	providerName = "cloudify"
)

// CloudProvider implents Instances, Zones, and LoadBalancer
type CloudProvider struct {
	client    *cloudify.CloudifyClient
	instances *CloudifyIntances
	balancers *CloudifyBalancer
}

// Initialize passes a Kubernetes clientBuilder interface to the cloud provider
func (r *CloudProvider) Initialize(clientBuilder controller.ControllerClientBuilder) {
	glog.Warning("Initialize")
}

// ProviderName returns the cloud provider ID.
func (r *CloudProvider) ProviderName() string {
	return providerName
}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	glog.Warning("LoadBalancer")
	if r.client != nil {
		if r.balancers != nil {
			return r.balancers, true
		} else {
			r.balancers = NewCloudifyBalancer(r.client)
			return r.balancers, true
		}
	}
	return nil, false
}

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) Zones() (cloudprovider.Zones, bool) {
	glog.Warning("Zones")
	return nil, false
}

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) Instances() (cloudprovider.Instances, bool) {
	glog.Warning("Instances")
	if r.client != nil {
		if r.instances != nil {
			return r.instances, true
		} else {
			r.instances = NewCloudifyIntances(r.client)
			return r.instances, true
		}
	}
	return nil, false
}

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) Clusters() (cloudprovider.Clusters, bool) {
	glog.Warning("Clusters")
	return nil, false
}

// Routes returns a routes interface along with whether the interface is supported.
func (r *CloudProvider) Routes() (cloudprovider.Routes, bool) {
	glog.Warning("Routers")
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

type CloudifyProviderConfig struct {
	Host     string `json:"host,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Tenant   string `json:"tenant,omitempty"`
}

func newCloudifyCloud(config io.Reader) (cloudprovider.Interface, error) {
	glog.Warning("New Cloudify client")

	var cloudConfig CloudifyProviderConfig
	cloudConfig.Host = os.Getenv("CFY_HOST")
	cloudConfig.User = os.Getenv("CFY_USER")
	cloudConfig.Password = os.Getenv("CFY_PASSWORD")
	cloudConfig.Tenant = os.Getenv("CFY_TENANT")
	if config != nil {
		err := json.NewDecoder(config).Decode(&cloudConfig)
		if err != nil {
			return nil, err
		}
	}

	glog.Warning("Config %+v", cloudConfig)
	return &CloudProvider{
		client: cloudify.NewClient(
			cloudConfig.Host, cloudConfig.User,
			cloudConfig.Password, cloudConfig.Tenant),
	}, nil
}

func init() {
	glog.Warning("Cloudify init")
	cloudprovider.RegisterCloudProvider(providerName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newCloudifyCloud(config)
	})
}
