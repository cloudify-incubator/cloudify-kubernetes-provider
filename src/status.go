package main

import (
    "fmt"
    "log"
    "flag"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "encoding/base64"
)

func Get(url string, user string, password string, tenant string) []byte {
    client := &http.Client{}

    var auth_string string
    auth_string = user + ":" + password
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }

    req.Header.Add("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(auth_string)))
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
    return body
}

type CloudifyVersion struct {
    Date        string `json:"date"`
    Edition     string `json:"edition"`
    Community   string `json:"community"`
    Version     string `json:"version"`
    Build       string `json:"build"`
    Commit      string `json:"commit"`
}

func GetVersion(host string, user string, password string, tenant string) CloudifyVersion {
    body := Get("http://" + host + "/api/v3.1/version", "admin", "secret", "default_tenant")

    var ver CloudifyVersion

    err := json.Unmarshal(body, &ver)
    if err != nil {
        log.Fatal(err)
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
    Status      string                      `json:"status"`
    Services    []CloudifyInstanceService   `json:"services"`
}

func GetStatus(host string, user string, password string, tenant string) CloudifyStatus {
    body := Get("http://" + host + "/api/v3.1/status", "admin", "secret", "default_tenant")

    var stat CloudifyStatus

    err := json.Unmarshal(body, &stat)
    if err != nil {
        log.Fatal(err)
    }

    return stat
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
        case "version": fmt.Println(GetVersion(host, user, password, tenant))
        case "status": {
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
        default: fmt.Println("Supported only: status, version")
    }
}
