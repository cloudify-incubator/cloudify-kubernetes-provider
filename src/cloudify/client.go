package cloudify

import (
	"archive/zip"
	"bytes"
	"cloudify/rest"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

func binaryPut(cl *CloudifyClient, url string, input []byte, input_type string, output rest.CloudifyMessageInterface) error {
	body := cl.RestCl.Put(url, input_type, input)

	err_post := json.Unmarshal(body, output)
	if err_post != nil {
		return err_post
	}

	if len(output.ErrorCode()) > 0 {
		return output
	}
	return nil
}

func dirZipArchive(parentDir string) ([]byte, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	log.Printf("Looking into %s", parentDir)
	err_walk := filepath.Walk(parentDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			f, err_create := w.Create("parent/" + path[len(parentDir):])
			if err_create != nil {
				return err_create
			}

			content, err_read := ioutil.ReadFile(path)
			if err_read != nil {
				return err_read
			}

			_, err_write := f.Write(content)
			if err_write != nil {
				return err_write
			}
			log.Printf("Attached: %s", path[len(parentDir):])
		}
		return nil
	})

	if err_walk != nil {
		return nil, err_walk
	}

	// Make sure to check the error on Close.
	err_zip := w.Close()
	if err_zip != nil {
		return nil, err_zip
	}
	return buf.Bytes(), nil
}

func (cl *CloudifyClient) PutBinary(url string, data []byte, output rest.CloudifyMessageInterface) error {
	return binaryPut(cl, url, data, rest.DataContentType, output)
}

func (cl *CloudifyClient) PutZip(url, path string, output rest.CloudifyMessageInterface) error {
	data, err := dirZipArchive(path)
	if err != nil {
		return err
	}

	return binaryPut(cl, url, data, rest.DataContentType, output)
}

func (cl *CloudifyClient) Put(url string, input interface{}, output rest.CloudifyMessageInterface) error {
	json_data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return binaryPut(cl, url, json_data, rest.JsonContentType, output)
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
