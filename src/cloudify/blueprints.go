package cloudify

import (
	"cloudify/rest"
	"encoding/json"
	"log"
)

type CloudifyBlueprint struct {
	// have id, owner information
	rest.CloudifyResource
	MainFileName string `json:"main_file_name"`
	// TODO describe "plan" struct
}

type CloudifyBlueprintGet struct {
	// can be response from api
	rest.CloudifyBaseMessage
	CloudifyBlueprint
}

type CloudifyBlueprints struct {
	rest.CloudifyBaseMessage
	Metadata rest.CloudifyMetadata `json:"metadata"`
	Items    []CloudifyBlueprint   `json:"items"`
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

func DeleteBlueprints(host, user, password, tenant, blueprint_id string) CloudifyBlueprintGet {
	body := rest.Delete("http://"+host+"/api/v3.1/blueprints/"+blueprint_id, user, password, tenant)

	var blueprint CloudifyBlueprintGet

	err := json.Unmarshal(body, &blueprint)
	if err != nil {
		log.Fatal(err)
	}

	if len(blueprint.ErrorCode) > 0 {
		log.Fatal(blueprint.Message)
	}

	return blueprint
}
