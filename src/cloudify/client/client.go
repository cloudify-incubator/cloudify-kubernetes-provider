package client

type CloudifyClient struct {
}

func GetClient() *CloudifyClient {
	var cliCloudify CloudifyClient
	return &cliCloudify
}
