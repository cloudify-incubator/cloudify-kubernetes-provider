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

package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app/options"
	"k8s.io/kubernetes/pkg/cloudprovider"
	_ "k8s.io/kubernetes/pkg/cloudprovider/providers"
	_ "k8s.io/kubernetes/pkg/cloudprovider/providers/cloudifyprovider" // only init from package
	"k8s.io/kubernetes/pkg/kubectl/util/logs"
	"k8s.io/kubernetes/pkg/version"
	"os"
)

var versionShow bool

/*
addNativeFlags - add supported flags.
Note: super super ugly way!!!!
look to src/k8s.io/kubernetes/cmd/cloud-controller-manager/app/options/options.go
*/
func addNativeFlags(s *options.CloudControllerManagerServer, fs *flag.FlagSet) *flag.FlagSet {
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	fs.StringVar(&s.CloudConfigFile, "cloud-config", s.CloudConfigFile, "The path to the cloud provider configuration file.  Empty string for no configuration file.")
	fs.BoolVar(&versionShow, "version", false, "Path to kubeconfig file with authorization and master location information.")
	return fs
}

var versionString = "0.1"

func main() {
	s := options.NewCloudControllerManagerServer()
	fs := addNativeFlags(s, flag.CommandLine)

	logs.InitLogs()
	defer logs.FlushLogs()

	fs.Parse(os.Args[1:])

	if versionShow {
		fmt.Printf("Kubernetes %s\n", version.Get())
		fmt.Printf("CFY Go client: %s\n", versionString)
		os.Exit(0)
	}

	cloud, err := cloudprovider.InitCloudProvider("cloudify", s.CloudConfigFile)
	if err != nil {
		glog.Fatal(err)
	}

	if err := app.Run(s, cloud); err != nil {
		glog.Fatal(err)
	}
}
