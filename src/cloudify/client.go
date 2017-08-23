package cloudify

import (
	"cloudify/rest"
	"encoding/json"
	"io/ioutil"
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

func (cl *CloudifyClient) Get(url string, output rest.CloudifyMessageInterface) error {
	body := cl.RestCl.Get(url, rest.JsonContentType)

	err := json.Unmarshal(body, output)
	if err != nil {
		return err
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}

func (cl *CloudifyClient) GetBinary(url, output_path string) error {
	body := cl.RestCl.Get(url, rest.DataContentType)

	err := ioutil.WriteFile(output_path, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (cl *CloudifyClient) Put(url string, input interface{}, output rest.CloudifyMessageInterface) error {
	json_data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	body := cl.RestCl.Put(url, json_data)

	err_post := json.Unmarshal(body, output)
	if err_post != nil {
		return err_post
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}

func (cl *CloudifyClient) Post(url string, input interface{}, output rest.CloudifyMessageInterface) error {
	json_data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	body := cl.RestCl.Post(url, json_data)

	err_post := json.Unmarshal(body, output)
	if err_post != nil {
		return err_post
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}

func (cl *CloudifyClient) Delete(url string, output rest.CloudifyMessageInterface) error {
	body := cl.RestCl.Delete(url)

	err_post := json.Unmarshal(body, output)
	if err_post != nil {
		return err_post
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}
