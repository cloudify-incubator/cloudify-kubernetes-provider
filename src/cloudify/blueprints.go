/*
Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

func (cl *CloudifyClient) GetBlueprints() CloudifyBlueprints {
	body := cl.RestCl.Get("http://" + cl.Host + "/api/v3.1/blueprints")

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

func (cl *CloudifyClient) DeleteBlueprints(blueprint_id string) CloudifyBlueprintGet {
	body := cl.RestCl.Delete("http://" + cl.Host + "/api/v3.1/blueprints/" + blueprint_id)

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
