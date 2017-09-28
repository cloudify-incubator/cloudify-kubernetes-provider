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

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type baseResponse struct {
	Status string `json:"status,omitempty"`
}

type capabilitiesResponse struct {
	Attach bool `json:"attach"`
}

type initResponse struct {
	baseResponse
	Capabilities capabilitiesResponse `json:"capabilities,omitempty"`
}

type mountResponse struct {
	Message    string `json:"message,omitempty"`
	Device     string `json:"device,omitempty"`
	VolumeName string `json:"volumeName,omitempty"`
	Attached   bool   `json:"attached,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		var response baseResponse
		response.Status = "Not supported"
		json_data, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(json_data))
		os.Exit(0)
	}
	command := os.Args[1]

	if command == "init" {
		var response initResponse
		response.Status = "Success"
		response.Capabilities.Attach = false
		json_data, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(json_data))
		os.Exit(0)
	} else {
		fmt.Println(command)
	}
}
