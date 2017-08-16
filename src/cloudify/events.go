package cloudify

import (
	"cloudify/rest"
	"encoding/json"
	"log"
	"net/url"
)

type CloudifyEvent struct {
	NodeInstanceId    string `json:"node_instance_id"`
	EventType         string `json:"event_type"`
	Operation         string `json:"operation"`
	BlueprintId       string `json:"blueprint_id"`
	NodeName          string `json:"node_name"`
	WorkflowId        string `json:"workflow_id"`
	ErrorCauses       string `json:"error_causes"`
	ReportedTimestamp string `json:"reported_timestamp"`
	DeploymentId      string `json:"deployment_id"`
	Type              string `json:"type"`
	ExecutionId       string `json:"execution_id"`
	Timestamp         string `json:"timestamp"`
	Message           string `json:"message"`
}

type CloudifyEvents struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []CloudifyEvent       `json:"items"`
}

func GetEvents(host, user, password, tenant string, params map[string]string) CloudifyEvents {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	body := rest.Get("http://"+host+"/api/v3.1/events?"+values.Encode(), user, password, tenant)

	var events CloudifyEvents

	err := json.Unmarshal(body, &events)
	if err != nil {
		log.Fatal(err)
	}

	if len(events.ErrorCode) > 0 {
		log.Fatal(events.Message)
	}

	return events
}
