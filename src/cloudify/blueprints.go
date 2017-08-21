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
	var blueprints CloudifyBlueprints

	err := cl.Get("blueprints", &blueprints)
	if err != nil {
		log.Fatal(err)
	}

	return blueprints
}

func (cl *CloudifyClient) DeleteBlueprints(blueprint_id string) CloudifyBlueprintGet {
	var blueprint CloudifyBlueprintGet

	err := cl.Delete("blueprints/"+blueprint_id, &blueprint)
	if err != nil {
		log.Fatal(err)
	}

	return blueprint
}
