package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
)

func GetRequest(url, user, password, tenant, method string, body io.Reader) *http.Request {
	log.Printf("Use: %v:%v@%v#%s\n", user, password, url, tenant)

	var auth_string string
	auth_string = user + ":" + password
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth_string)))
	if len(tenant) > 0 {
		req.Header.Add("Tenant", tenant)
	}

	return req
}

func Get(url string, user string, password string, tenant string) []byte {
	req := GetRequest(url, user, password, tenant, "GET", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response %s\n", string(body))
	return body
}

func Delete(url string, user string, password string, tenant string) []byte {
	req := GetRequest(url, user, password, tenant, "DELETE", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response %s\n", string(body))
	return body
}

func Post(url, user, password, tenant string, data []byte) []byte {
	req := GetRequest(url, user, password, tenant, "POST", bytes.NewBuffer(data))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response %s\n", string(body))
	return body
}

func Put(url, user, password, tenant string, data []byte) []byte {
	req := GetRequest(url, user, password, tenant, "PUT", bytes.NewBuffer(data))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response %s\n", string(body))
	return body
}

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
	body := Get("http://"+host+"/api/v3.1/version", user, password, tenant)

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
	body := Get("http://"+host+"/api/v3.1/status", user, password, tenant)

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
	body := Get("http://"+host+"/api/v3.1/blueprints", user, password, tenant)

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
	body := Get("http://"+host+"/api/v3.1/deployments", user, password, tenant)

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
	body := Delete("http://"+host+"/api/v3.1/deployments/"+deployment_id, user, password, tenant)

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

	body := Put("http://"+host+"/api/v3.1/deployments/"+deployment_id, user, password, tenant, json_data)

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
	body := Get("http://"+host+"/api/v3.1/executions", user, password, tenant)

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

	body := Post("http://"+host+"/api/v3.1/executions", user, password, tenant, json_data)

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

func PrintBottomLine(columnSizes []int) {
	fmt.Printf("+")
	for _, size := range columnSizes {
		fmt.Print(strings.Repeat("-", size+2))
		fmt.Printf("+")
	}
	fmt.Printf("\n")
}

func PrintLine(columnSizes []int, lines []string) {
	fmt.Printf("|")
	for col, size := range columnSizes {
		fmt.Print(" " + lines[col] + " ")
		fmt.Print(strings.Repeat(" ", size-utf8.RuneCountInString(lines[col])))
		fmt.Printf("|")
	}
	fmt.Printf("\n")
}

func PrintTable(titles []string, lines [][]string) {
	columnSizes := make([]int, len(titles))

	// column title sizes
	for col, name := range titles {
		if columnSizes[col] < utf8.RuneCountInString(name) {
			columnSizes[col] = utf8.RuneCountInString(name)
		}
	}

	// column value sizes
	for _, values := range lines {
		for col, name := range values {
			if columnSizes[col] < utf8.RuneCountInString(name) {
				columnSizes[col] = utf8.RuneCountInString(name)
			}
		}
	}

	PrintBottomLine(columnSizes)
	// titles
	PrintLine(columnSizes, titles)
	PrintBottomLine(columnSizes)
	// lines
	for _, values := range lines {
		PrintLine(columnSizes, values)
	}
	PrintBottomLine(columnSizes)
}

var host string
var user string
var password string
var tenant string

func basicOptions(name string) *flag.FlagSet {
	var commonFlagSet *flag.FlagSet
	commonFlagSet = flag.NewFlagSet("name", flag.ExitOnError)
	commonFlagSet.StringVar(&host, "host", "localhost", "Manager host name")
	commonFlagSet.StringVar(&user, "user", "admin", "Manager user name")
	commonFlagSet.StringVar(&password, "password", "secret", "Manager user password")
	commonFlagSet.StringVar(&tenant, "tenant", "default_tenant", "Manager tenant")
	return commonFlagSet
}

func infoOptions() int {
	defaultError := "state/version subcommand is required"

	if len(os.Args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	operFlagSet := basicOptions("status")

	operFlagSet.Parse(os.Args[3:])

	switch os.Args[2] {
	case "state":
		{
			stat := GetStatus(host, user, password, tenant)

			fmt.Printf("Retrieving manager services status... [ip=%v]\n", host)
			fmt.Printf("Manager status: %v\n", stat.Status)
			fmt.Printf("Services:\n")
			var lines [][]string = make([][]string, len(stat.Services))
			for pos, service := range stat.Services {
				lines[pos] = make([]string, 2)
				lines[pos][0] = service.DisplayName
				lines[pos][1] = service.Status()
			}
			PrintTable([]string{"service", "status"}, lines)
		}
	case "version":
		{
			ver := GetVersion(host, user, password, tenant)
			fmt.Printf("Retrieving manager services version... [ip=%v]\n", host)
			PrintTable([]string{"Version", "Edition"}, [][]string{{ver.Version, ver.Edition}})
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

func blueprintsOptions() int {
	defaultError := "list subcommand is required"

	if len(os.Args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	operFlagSet := basicOptions("blueprints")

	operFlagSet.Parse(os.Args[3:])

	switch os.Args[2] {
	case "list":
		{
			blueprints := GetBlueprints(host, user, password, tenant)
			var lines [][]string = make([][]string, len(blueprints.Items))
			for pos, blueprint := range blueprints.Items {
				lines[pos] = make([]string, 7)
				lines[pos][0] = blueprint.Id
				lines[pos][1] = blueprint.Description
				lines[pos][2] = blueprint.MainFileName
				lines[pos][3] = blueprint.CreatedAt
				lines[pos][4] = blueprint.UpdatedAt
				lines[pos][5] = blueprint.Tenant
				lines[pos][6] = blueprint.CreatedBy
			}
			PrintTable([]string{"id", "description", "main_file_name", "created_at", "updated_at", "tenant_name", "created_by"}, lines)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

func deploymentsOptions() int {
	defaultError := "list/create/delete subcommand is required"

	if len(os.Args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	operFlagSet := basicOptions("deployments")

	switch os.Args[2] {
	case "list":
		{
			operFlagSet.Parse(os.Args[3:])
			deployments := GetDeployments(host, user, password, tenant)
			var lines [][]string = make([][]string, len(deployments.Items))
			for pos, deployment := range deployments.Items {
				lines[pos] = make([]string, 6)
				lines[pos][0] = deployment.Id
				lines[pos][1] = deployment.BlueprintId
				lines[pos][2] = deployment.CreatedAt
				lines[pos][3] = deployment.UpdatedAt
				lines[pos][4] = deployment.Tenant
				lines[pos][5] = deployment.CreatedBy
			}
			PrintTable([]string{"id", "blueprint_id", "created_at", "updated_at", "tenant_name", "created_by"}, lines)
		}
	case "create":
		{
			if len(os.Args) < 4 {
				fmt.Println("Deployment Id requered")
				return 1
			}

			var blueprint string
			operFlagSet.StringVar(&blueprint, "blueprint", "", "The unique identifier for the blueprint")

			operFlagSet.Parse(os.Args[4:])

			var depl CloudifyDeploymentPost
			depl.BlueprintId = blueprint
			depl.Inputs = map[string]string{}
			deployment := CreateDeployments(host, user, password, tenant, os.Args[3], depl)

			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 6)
			lines[0][0] = deployment.Id
			lines[0][1] = deployment.BlueprintId
			lines[0][2] = deployment.CreatedAt
			lines[0][3] = deployment.UpdatedAt
			lines[0][4] = deployment.Tenant
			lines[0][5] = deployment.CreatedBy
			PrintTable([]string{"id", "blueprint_id", "created_at", "updated_at", "tenant_name", "created_by"}, lines)
		}
	case "delete":
		{
			if len(os.Args) < 4 {
				fmt.Println("Deployment Id requered")
				return 1
			}

			operFlagSet.Parse(os.Args[4:])
			deployment := DeleteDeployments(host, user, password, tenant, os.Args[3])
			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 6)
			lines[0][0] = deployment.Id
			lines[0][1] = deployment.BlueprintId
			lines[0][2] = deployment.CreatedAt
			lines[0][3] = deployment.UpdatedAt
			lines[0][4] = deployment.Tenant
			lines[0][5] = deployment.CreatedBy
			PrintTable([]string{"id", "blueprint_id", "created_at", "updated_at", "tenant_name", "created_by"}, lines)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

func executionsOptions() int {
	defaultError := "list/start subcommand is required"

	if len(os.Args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	operFlagSet := basicOptions("executions")

	switch os.Args[2] {
	case "list":
		{
			operFlagSet.Parse(os.Args[3:])
			executions := GetExecutions(host, user, password, tenant)
			var lines [][]string = make([][]string, len(executions.Items))
			for pos, execution := range executions.Items {
				lines[pos] = make([]string, 8)
				lines[pos][0] = execution.Id
				lines[pos][1] = execution.WorkflowId
				lines[pos][2] = execution.Status
				lines[pos][3] = execution.DeploymentId
				lines[pos][4] = execution.CreatedAt
				lines[pos][5] = execution.Error
				lines[pos][6] = execution.Tenant
				lines[pos][7] = execution.CreatedBy
			}
			PrintTable([]string{"id", "workflow_id", "status", "deployment_id", "created_at", "error", "tenant_name", "created_by"}, lines)
		}
	case "start":
		{

			if len(os.Args) < 4 {
				fmt.Println("Workflow Id requered")
				return 1
			}

			var deployment string
			operFlagSet.StringVar(&deployment, "deployment", "", "The unique identifier for the deployment")
			operFlagSet.Parse(os.Args[4:])

			var exec CloudifyExecutionPost
			exec.WorkflowId = os.Args[3]
			exec.DeploymentId = deployment

			execution := PostExecution(host, user, password, tenant, exec)

			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 8)
			lines[0][0] = execution.Id
			lines[0][1] = execution.WorkflowId
			lines[0][2] = execution.Status
			lines[0][3] = execution.DeploymentId
			lines[0][4] = execution.CreatedAt
			lines[0][5] = execution.Error
			lines[0][6] = execution.Tenant
			lines[0][7] = execution.CreatedBy
			PrintTable([]string{"id", "workflow_id", "status", "deployment_id", "created_at", "error", "tenant_name", "created_by"}, lines)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

func main() {
	defaultError := "Supported only: status, version, blueprints, deployments, executions, executions-install"
	if len(os.Args) < 2 {
		fmt.Println(defaultError)
		return
	}

	switch os.Args[1] {
	case "status":
		{
			os.Exit(infoOptions())
		}
	case "blueprints":
		{
			os.Exit(blueprintsOptions())
		}
	case "deployments":
		{
			os.Exit(deploymentsOptions())
		}
	case "executions":
		{
			os.Exit(executionsOptions())
		}
	default:
		{
			fmt.Println(defaultError)
			os.Exit(1)
		}
	}
}
