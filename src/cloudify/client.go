package cloudify

import (
	"cloudify/rest"
)

type CloudifyClient struct {
	RestCl rest.CloudifyRestClient
	Host   string
}

func NewClient(host, user, password, tenant string) *CloudifyClient {
	var cliCl CloudifyClient
	cliCl.Host = host
	cliCl.RestCl.User = user
	cliCl.RestCl.Password = password
	cliCl.RestCl.Tenant = tenant
	return &cliCl
}
