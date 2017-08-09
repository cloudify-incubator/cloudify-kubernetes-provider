package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

func Get(url string, user string, password string, tenant string) []byte {
	log.Printf("Use: %v:%v@%v#%s\n", user, password, url, tenant)

	client := &http.Client{}

	var auth_string string
	auth_string = user + ":" + password
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth_string)))
	if len(tenant) > 0 {
		req.Header.Add("Tenant", tenant)
	}

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
	Message         string `json:"message"`
	ErrorCode       string `json:"error_code"`
	ServerTraceback string `json:"server_traceback"`
}

type CloudifyVersion struct {
	CloudifyBaseMessage
	Date    string `json:"date"`
	Edition string `json:"edition"`
	Version string `json:"version"`
	Build   string `json:"build"`
	Commit  string `json:"commit"`
}

func GetVersion(host string, user string, password string, tenant string) CloudifyVersion {
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

func GetStatus(host string, user string, password string, tenant string) CloudifyStatus {
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

func GetBlueprints(host string, user string, password string, tenant string) CloudifyBlueprints {
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

type CloudifyDeployment struct {
	CloudifyResource
	Permalink   string             `json:"permalink"`
	BlueprintId string             `json:"blueprint_id"`
	Workflows   []CloudifyWorkflow `json:"workflows"`
	// TODO describe "inputs" srtuct
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

func GetDeployments(host string, user string, password string, tenant string) CloudifyDeployments {
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

func main() {
	var host string
	var user string
	var password string
	var tenant string
	var command string

	flag.StringVar(&host, "host", "localhost", "Manager host name")
	flag.StringVar(&user, "user", "admin", "Manager user name")
	flag.StringVar(&password, "password", "secret", "Manager user password")
	flag.StringVar(&tenant, "tenant", "default_tenant", "Manager tenant")
	flag.StringVar(&command, "command", "version", "Command for run")
	flag.Parse()

	switch command {
	case "version":
		{
			ver := GetVersion(host, user, password, tenant)
			fmt.Printf("Retrieving manager services version... [ip=%v]\n", host)
			PrintTable([]string{"Version", "Edition"}, [][]string{{ver.Version, ver.Edition}})
		}
	case "status":
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
	case "blueprints":
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
	case "deployments":
		{
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
	default:
		fmt.Println("Supported only: status, version")
	}
}
