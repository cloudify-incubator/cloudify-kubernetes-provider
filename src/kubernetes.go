
package main

import (
	"fmt"
    "k8s.io/kubernetes/cmd/cloud-controller-manager/app/options"
)

func main() {
	s := options.NewCloudControllerManagerServer()
    fmt.Printf("%+v\n", s)
}
