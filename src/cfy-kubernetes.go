package main

import (
	_ "cloudifyprovider" // only init from package
	"fmt"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app/options"
	"k8s.io/kubernetes/pkg/cloudprovider"
	_ "k8s.io/kubernetes/pkg/cloudprovider/providers"
	"k8s.io/kubernetes/pkg/version/verflag"
	"log"
)

func main() {
	s := options.NewCloudControllerManagerServer()
	fmt.Printf("%+v\n", s)

	verflag.PrintAndExitIfRequested()

	cloud, err := cloudprovider.InitCloudProvider("cloudify", s.CloudConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(s, cloud); err != nil {
		log.Fatal(err)
	}
}
