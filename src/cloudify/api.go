package cloudify

import (
	"cloudify/rest"
	"encoding/json"
	"log"
)

type CloudifyBaseMessage struct {
	Message         string `json:"message,omitempty"`
	ErrorCode       string `json:"error_code,omitempty"`
	ServerTraceback string `json:"server_traceback,omitempty"`
}

type CloudifyVersion struct {
	CloudifyBaseMessage
	Date    string `json:"date"`
	Edition string `json:"edition"`
	Version string `json:"version"`
	Build   string `json:"build"`
	Commit  string `json:"commit"`
}

func GetVersion(host, user, password, tenant string) CloudifyVersion {
	body := rest.Get("http://"+host+"/api/v3.1/version", user, password, tenant)

	var ver CloudifyVersion

	err := json.Unmarshal(body, &ver)
	if err != nil {
		log.Fatal(err)
	}

	if len(ver.ErrorCode) > 0 {
		log.Fatal(ver.Message)
	}
	return ver
}

type CloudifyInstanceStatus struct {
	LoadState   string `json:"LoadState"`
	Description string `json:"Description"`
	State       string `json:"state"`
	MainPID     uint   `json:"MainPID"`
	Id          string `json:"Id"`
	ActiveState string `json:"ActiveState"`
	SubState    string `json:"SubState"`
}

type CloudifyInstanceService struct {
	Instances   []CloudifyInstanceStatus `json:"instances"`
	DisplayName string                   `json:"display_name"`
}

func (s CloudifyInstanceService) Status() string {
	var state string = "unknown"

	for _, instance := range s.Instances {
		if state != "failed" {
			state = instance.State
		}
	}

	return state
}

type CloudifyStatus struct {
	CloudifyBaseMessage
	Status   string                    `json:"status"`
	Services []CloudifyInstanceService `json:"services"`
}

func GetStatus(host, user, password, tenant string) CloudifyStatus {
	body := rest.Get("http://"+host+"/api/v3.1/status", user, password, tenant)

	var stat CloudifyStatus

	err := json.Unmarshal(body, &stat)
	if err != nil {
		log.Fatal(err)
	}

	if len(stat.ErrorCode) > 0 {
		log.Fatal(stat.Message)
	}

	return stat
}

type CloudifyPagination struct {
	Total  uint `json:"total"`
	Offset uint `json:"offset"`
	Size   uint `json:"size"`
}

type CloudifyMetadata struct {
	Pagination CloudifyPagination `json:"pagination"`
}

type CloudifyResource struct {
	Id              string `json:"id"`
	Description     string `json:"description"`
	Tenant          string `json:"tenant_name"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	CreatedBy       string `json:"created_by"`
	PrivateResource bool   `json:"private_resource"`
}
type CloudifyBlueprint struct {
	CloudifyResource
	MainFileName string `json:"main_file_name"`
	// TODO describe "plan" struct
}

type CloudifyBlueprints struct {
	CloudifyBaseMessage
	Metadata CloudifyMetadata    `json:"metadata"`
	Items    []CloudifyBlueprint `json:"items"`
}

func GetBlueprints(host, user, password, tenant string) CloudifyBlueprints {
	body := rest.Get("http://"+host+"/api/v3.1/blueprints", user, password, tenant)

	var blueprints CloudifyBlueprints

	err := json.Unmarshal(body, &blueprints)
	if err != nil {
		log.Fatal(err)
	}

	if len(blueprints.ErrorCode) > 0 {
		log.Fatal(blueprints.Message)
	}

	return blueprints
}

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
	CloudifyBaseMessage
	// have id, owner information
	CloudifyResource
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
	CloudifyBaseMessage
	Metadata CloudifyMetadata     `json:"metadata"`
	Items    []CloudifyDeployment `json:"items"`
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

type CloudifyExecutionPost struct {
	WorkflowId   string `json:"workflow_id"`
	DeploymentId string `json:"deployment_id"`
}

type CloudifyExecution struct {
	// can be response from api
	CloudifyBaseMessage
	// have id, owner information
	CloudifyResource
	// contain information from post
	CloudifyExecutionPost
	IsSystemWorkflow bool   `json:"is_system_workflow"`
	Error            string `json:"error"`
	BlueprintId      string `json:"blueprint_id"`
	Status           string `json:"status"`
	// TODO describe "parameters" struct
}

type CloudifyExecutions struct {
	CloudifyBaseMessage
	Metadata CloudifyMetadata    `json:"metadata"`
	Items    []CloudifyExecution `json:"items"`
}

func GetExecutions(host, user, password, tenant string) CloudifyExecutions {
	body := rest.Get("http://"+host+"/api/v3.1/executions", user, password, tenant)

	var executions CloudifyExecutions

	err := json.Unmarshal(body, &executions)
	if err != nil {
		log.Fatal(err)
	}

	if len(executions.ErrorCode) > 0 {
		log.Fatal(executions.Message)
	}

	return executions
}

func PostExecution(host, user, password, tenant string, exec CloudifyExecutionPost) CloudifyExecution {
	json_data, err := json.Marshal(exec)
	if err != nil {
		log.Fatal(err)
	}

	body := rest.Post("http://"+host+"/api/v3.1/executions", user, password, tenant, json_data)

	var execution CloudifyExecution

	err_post := json.Unmarshal(body, &execution)
	if err_post != nil {
		log.Fatal(err_post)
	}

	if len(execution.ErrorCode) > 0 {
		log.Fatal(execution.Message)
	}

	return execution
}
