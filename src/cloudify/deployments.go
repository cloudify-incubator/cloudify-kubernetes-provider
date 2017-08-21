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
	var deployments CloudifyDeployments

	err := cl.Get("deployments", &deployments)
	if err != nil {
		log.Fatal(err)
	}

	return deployments
}

func (cl *CloudifyClient) DeleteDeployments(deployment_id string) CloudifyDeploymentGet {
	var deployment CloudifyDeploymentGet

	err := cl.Delete("deployments/"+deployment_id, &deployment)
	if err != nil {
		log.Fatal(err)
	}

	return deployment
}

func (cl *CloudifyClient) CreateDeployments(deployment_id string, depl CloudifyDeploymentPost) CloudifyDeploymentGet {
	var deployment CloudifyDeploymentGet

	err := cl.Put("deployments/"+deployment_id, depl, &deployment)
	if err != nil {
		log.Fatal(err)
	}

	return deployment
}
