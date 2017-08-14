package cloudify

import (
	"cloudify/rest"
	"encoding/json"
	"log"
)

type CloudifyVersion struct {
	rest.CloudifyBaseMessage
	Date    string `json:"date"`
	Edition string `json:"edition"`
	Version string `json:"version"`
	Build   string `json:"build"`
	Commit  string `json:"commit"`
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
	rest.CloudifyBaseMessage
	Status   string                    `json:"status"`
	Services []CloudifyInstanceService `json:"services"`
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
