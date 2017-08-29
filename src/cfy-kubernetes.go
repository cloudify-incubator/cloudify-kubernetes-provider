package main

import (
	_ "cloudifyprovider" // only init from package
	"flag"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app/options"
	"k8s.io/kubernetes/pkg/util/logs"
	"k8s.io/kubernetes/pkg/cloudprovider"
	_ "k8s.io/kubernetes/pkg/cloudprovider/providers"
	"k8s.io/kubernetes/pkg/version/verflag"
	"log"
	"os"
)

//super super ugly way!!!! look to src/k8s.io/kubernetes/cmd/cloud-controller-manager/app/options/options.go
func AddNativeFlags(s *options.CloudControllerManagerServer, fs *flag.FlagSet) *flag.FlagSet {
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	return fs
}

func main() {
	s := options.NewCloudControllerManagerServer()
	fs := AddNativeFlags(s, flag.CommandLine)

	logs.InitLogs()
	defer logs.FlushLogs()
	verflag.PrintAndExitIfRequested()

	fs.Parse(os.Args[1:])
	cloud, err := cloudprovider.InitCloudProvider("cloudify", s.CloudConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(s, cloud); err != nil {
		log.Fatal(err)
	}
}
