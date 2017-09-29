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
	"encoding/json"
	"fmt"
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"io/ioutil"
	"os"
)

type baseResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type capabilitiesResponse struct {
	Attach bool `json:"attach"`
}

type initResponse struct {
	baseResponse
	Capabilities capabilitiesResponse `json:"capabilities,omitempty"`
}

type mountResponse struct {
	baseResponse
	Attached bool `json:"attached"`
}

type CloudifyConfig struct {
	Host       string `json:"host"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Tenant     string `json:"tenant"`
	Deployment string `json:"deployment"`
	Instance   string `json:"intance"`
}

func getConfig() (*CloudifyConfig, error) {
	var config CloudifyConfig
	var configFile = os.Getenv("CFY_CONFIG")
	if configFile == "" {
		configFile = "/etc/cloudify/mount.json"
	}
	configContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err_marshal := json.Unmarshal([]byte(configContent), &config)
	return &config, err_marshal
}

func initFunction() error {
	var response initResponse
	response.Status = "Success"
	response.Capabilities.Attach = false
	json_data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println(string(json_data))
	return nil
}

func mountFunction(config *CloudifyConfig, path, config_json string) error {
	var in_data_parsed map[string]interface{}
	err := json.Unmarshal([]byte(config_json), &in_data_parsed)
	if err != nil {
		return err
	}

	cl := cloudify.NewClient(config.Host, config.User, config.Password, config.Tenant)

	var exec cloudify.CloudifyExecutionPost
	exec.WorkflowId = "execute_operation"
	exec.DeploymentId = config.Deployment
	exec.Parameters = map[string]interface{}{}
	exec.Parameters["operation"] = "maintenance.mount"
	exec.Parameters["node_ids"] = []string{}
	exec.Parameters["type_names"] = []string{}
	exec.Parameters["run_by_dependency_order"] = false
	exec.Parameters["allow_kwargs_override"] = nil
	exec.Parameters["node_instance_ids"] = []string{config.Instance}
	exec.Parameters["operation_kwargs"] = map[string]interface{}{
		"path":   path,
		"params": in_data_parsed}
	execution := cl.PostExecution(exec)

	fmt.Printf("%+v\n", execution)

	var response mountResponse
	response.Status = "Success"
	response.Attached = true
	json_data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println(string(json_data))
	return nil
}

func unMountFunction(path string) error {
	fmt.Println(path)
	var response mountResponse
	response.Status = "Success"
	response.Attached = false
	json_data, err := json.Marshal(response)
	if err != nil {
		return err
	} else {
		fmt.Println(string(json_data))
	}
	return nil
}

func main() {
	var message string = "Unknown"
	config, config_err := getConfig()
	if config_err != nil {
		message = config_err.Error()
	} else if len(os.Args) > 1 {
		command := os.Args[1]
		if len(os.Args) == 2 && command == "init" {
			err := initFunction()
			if err != nil {
				message = err.Error()
			} else {
				os.Exit(0)
			}
		}
		if len(os.Args) == 4 && command == "mount" {
			err := mountFunction(config, os.Args[2], os.Args[3])
			if err != nil {
				message = err.Error()
			} else {
				os.Exit(0)
			}
		}
		if len(os.Args) == 3 && command == "unmount" {
			err := unMountFunction(os.Args[2])
			if err != nil {
				message = err.Error()
			} else {
				os.Exit(0)
			}
		}
	}
	var response baseResponse
	response.Status = "Not supported"
	response.Message = message
	json_data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(json_data))
	os.Exit(0)
}
