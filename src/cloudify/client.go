package cloudify

import (
	"cloudify/rest"
)

type CloudifyClient struct {
	RestCl rest.CloudifyRestClient
}

func NewClient(host, user, password, tenant string) *CloudifyClient {
	var cliCl CloudifyClient
	cliCl.RestCl.RestURL = "http://" + host + "/api/v3.1/"
	cliCl.RestCl.User = user
	cliCl.RestCl.Password = password
	cliCl.RestCl.Tenant = tenant
	return &cliCl
}
