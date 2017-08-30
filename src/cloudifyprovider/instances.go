package cloudifyprovider

import (
	"cloudify"
	"fmt"
)

func ListInstances(client *cloudify.CloudifyClient) {
	fmt.Printf("%+v", client)
}
