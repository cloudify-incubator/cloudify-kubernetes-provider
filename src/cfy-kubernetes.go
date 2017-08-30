package main

import (
	_ "cloudifyprovider" // only init from package
	"flag"
	"fmt"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app/options"
	"k8s.io/kubernetes/pkg/cloudprovider"
	_ "k8s.io/kubernetes/pkg/cloudprovider/providers"
	"k8s.io/kubernetes/pkg/util/logs"
	"k8s.io/kubernetes/pkg/version"
	"log"
	"os"
)

var versionShow bool

//super super ugly way!!!! look to src/k8s.io/kubernetes/cmd/cloud-controller-manager/app/options/options.go
func AddNativeFlags(s *options.CloudControllerManagerServer, fs *flag.FlagSet) *flag.FlagSet {
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	fs.StringVar(&s.CloudConfigFile, "cloud-config", s.CloudConfigFile, "The path to the cloud provider configuration file.  Empty string for no configuration file.")
	fs.BoolVar(&versionShow, "version", false, "Path to kubeconfig file with authorization and master location information.")
	return fs
}

func main() {
	s := options.NewCloudControllerManagerServer()
	fs := AddNativeFlags(s, flag.CommandLine)

	logs.InitLogs()
	defer logs.FlushLogs()

	fs.Parse(os.Args[1:])

	if versionShow {
		fmt.Printf("Kubernetes %s\n", version.Get())
		os.Exit(0)
	}

	cloud, err := cloudprovider.InitCloudProvider("cloudify", s.CloudConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(s, cloud); err != nil {
		log.Fatal(err)
	}
}
