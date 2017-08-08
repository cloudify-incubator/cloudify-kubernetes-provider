package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/base64"
    "encoding/json"
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
    var host = "localhost"
    var user = "admin"
    var password = "secret"
    var tenant = "default_tenant"
    fmt.Println(GetVersion(host, user, password, tenant))
    fmt.Println(GetStatus(host, user, password, tenant))
}
