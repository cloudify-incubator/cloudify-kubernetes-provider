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
    req.Header.Add("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(auth_string)))
    req.Header.Add("Tenant", tenant)
    if err != nil {
        log.Fatal(err)
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

type Version struct {
    Date        string `json:"date"`
    Edition     string `json:"edition"`
    Community   string `json:"community"`
    Version     string `json:"version"`
    Build       string `json:"build"`
    Commit      string `json:"commit"`
}

func GetVersion(host string, user string, password string, tenant string) Version {
    body := Get("http://" + host + "/api/v3.1/version", "admin", "secret", "default_tenant")

    var ver Version

    err := json.Unmarshal(body, &ver)
    if err != nil {
        log.Fatal(err)
    }

    return ver
}

func main() {
    var host = "localhost"
    var user = "admin"
    var password = "secret"
    var tenant = "default_tenant"
    fmt.Println(GetVersion(host, user, password, tenant))
}
