package cloudify

import (
	"cloudify/rest"
	"encoding/json"
	"log"
)

type CloudifyWorkflow struct {
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	// TODO describe "parameters" srtuct
}

type CloudifyDeploymentPost struct {
	BlueprintId string            `json:"blueprint_id"`
	Inputs      map[string]string `json:"inputs"`
}

type CloudifyDeployment struct {
	// can be response from api
	rest.CloudifyBaseMessage
	// have id, owner information
	rest.CloudifyResource
	// contain information from post
	CloudifyDeploymentPost
	Permalink string             `json:"permalink"`
	Workflows []CloudifyWorkflow `json:"workflows"`
	// TODO describe "policy_types" struct
	// TODO describe "policy_triggers" struct
	// TODO describe "groups" struct
	// TODO describe "scaling_groups" struct
	// TODO describe "outputs" struct
}

type CloudifyDeployments struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []CloudifyDeployment  `json:"items"`
}

func GetDeployments(host, user, password, tenant string) CloudifyDeployments {
	body := rest.Get("http://"+host+"/api/v3.1/deployments", user, password, tenant)

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

func DeleteDeployments(host, user, password, tenant, deployment_id string) CloudifyDeployment {
	body := rest.Delete("http://"+host+"/api/v3.1/deployments/"+deployment_id, user, password, tenant)

	var deployment CloudifyDeployment

	err := json.Unmarshal(body, &deployment)
	if err != nil {
		log.Fatal(err)
	}

	if len(deployment.ErrorCode) > 0 {
		log.Fatal(deployment.Message)
	}

	return deployment
}

func CreateDeployments(host, user, password, tenant, deployment_id string, depl CloudifyDeploymentPost) CloudifyDeployment {
	json_data, err := json.Marshal(depl)
	if err != nil {
		log.Fatal(err)
	}

	body := rest.Put("http://"+host+"/api/v3.1/deployments/"+deployment_id, user, password, tenant, json_data)

	var deployment CloudifyDeployment

	err_post := json.Unmarshal(body, &deployment)
	if err_post != nil {
		log.Fatal(err_post)
	}

	if len(deployment.ErrorCode) > 0 {
		log.Fatal(deployment.Message)
	}

	return deployment
}
