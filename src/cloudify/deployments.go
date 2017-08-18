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

package cloudify

import (
	"cloudify/rest"
	"encoding/json"
	"log"
)

// Check https://blog.golang.org/json-and-go for more info about json marshaling.
type CloudifyWorkflow struct {
	CreatedAt  string                 `json:"created_at"`
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

type CloudifyDeploymentPost struct {
	BlueprintId string                 `json:"blueprint_id"`
	Inputs      map[string]interface{} `json:"inputs"`
}

type CloudifyDeployment struct {
	// have id, owner information
	rest.CloudifyResource
	// contain information from post
	CloudifyDeploymentPost
	Permalink string                 `json:"permalink"`
	Workflows []CloudifyWorkflow     `json:"workflows"`
	Outputs   map[string]interface{} `json:"outputs"`
	// TODO describe "policy_types" struct
	// TODO describe "policy_triggers" struct
	// TODO describe "groups" struct
	// TODO describe "scaling_groups" struct
}

type CloudifyDeploymentGet struct {
	// can be response from api
	rest.CloudifyBaseMessage
	CloudifyDeployment
}

type CloudifyDeployments struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []CloudifyDeployment  `json:"items"`
}

func (cl *CloudifyClient) GetDeployments() CloudifyDeployments {
	body := cl.RestCl.Get("http://" + cl.Host + "/api/v3.1/deployments")

	var deployments CloudifyDeployments

	err := json.Unmarshal(body, &deployments)
	if err != nil {
		log.Fatal(err)
	}

	if len(deployments.ErrorCode) > 0 {
		log.Fatal(deployments.Message)
	}

	return deployments
}

func (cl *CloudifyClient) DeleteDeployments(deployment_id string) CloudifyDeploymentGet {
	body := cl.RestCl.Delete("http://" + cl.Host + "/api/v3.1/deployments/" + deployment_id)

	var deployment CloudifyDeploymentGet

	err := json.Unmarshal(body, &deployment)
	if err != nil {
		log.Fatal(err)
	}

	if len(deployment.ErrorCode) > 0 {
		log.Fatal(deployment.Message)
	}

	return deployment
}

func (cl *CloudifyClient) CreateDeployments(deployment_id string, depl CloudifyDeploymentPost) CloudifyDeploymentGet {
	json_data, err := json.Marshal(depl)
	if err != nil {
		log.Fatal(err)
	}

	body := cl.RestCl.Put("http://"+cl.Host+"/api/v3.1/deployments/"+deployment_id, json_data)

	var deployment CloudifyDeploymentGet

	err_post := json.Unmarshal(body, &deployment)
	if err_post != nil {
		log.Fatal(err_post)
	}

	if len(deployment.ErrorCode) > 0 {
		log.Fatal(deployment.Message)
	}

	return deployment
}
