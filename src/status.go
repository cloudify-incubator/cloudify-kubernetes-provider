package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type CloudifyBlueprints struct {
	CloudifyBaseMessage
	Metadata CloudifyMetadata `json:"metadata"`
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

type CloudifyDeployments struct {
	CloudifyBaseMessage
	Metadata CloudifyMetadata `json:"metadata"`
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
			fmt.Printf("Version: %v\n", ver.Version)
			fmt.Printf("Edition: %s\n", ver.Edition)
		}
	case "status":
		{
			stat := GetStatus(host, user, password, tenant)
			fmt.Println("Manager status: " + stat.Status)

			for _, service := range stat.Services {
				fmt.Println("- Name: " + service.DisplayName)
				fmt.Println("  Intances: ")
				for _, instance := range service.Instances {
					fmt.Println("  - LoadState: ", instance.LoadState)
					fmt.Println("    Description: ", instance.Description)
					fmt.Println("    State: ", instance.State)
					fmt.Println("    MainPID: ", instance.MainPID)
					fmt.Println("    Id: ", instance.Id)
					fmt.Println("    ActiveState: ", instance.ActiveState)
					fmt.Println("    SubState: ", instance.SubState)
				}
			}
		}
	case "blueprints":
		{
			fmt.Printf("%+v\n", GetBlueprints(host, user, password, tenant))
		}
	case "deployments":
		{
			fmt.Printf("%+v\n", GetDeployments(host, user, password, tenant))
		}
	default:
		fmt.Println("Supported only: status, version")
	}
}
