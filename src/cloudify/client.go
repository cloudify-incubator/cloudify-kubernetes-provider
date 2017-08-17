package cloudify

type CloudifyClient struct {
	Host     string
	User     string
	Password string
	Tenant   string
}

func NewClient(host, user, password, tenant string) *CloudifyClient {
	var cliCl CloudifyClient
	cliCl.Host = host
	cliCl.User = user
	cliCl.Password = password
	cliCl.Tenant = tenant
	return &cliCl
}
