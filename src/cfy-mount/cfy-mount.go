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
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type capabilitiesResponse struct {
	Attach bool `json:"attach"`
}

type initResponse struct {
	baseResponse
	Capabilities capabilitiesResponse `json:"capabilities,omitempty"`
}

type mountResponse struct {
	baseResponse
	Attached bool `json:"attached"`
}

/*
type otherResponse struct {

	Device     string `json:"device,omitempty"`
	VolumeName string `json:"volumeName,omitempty"`
}
*/

func main() {
	var message string = "Unknown"
	if len(os.Args) > 1 {
		command := os.Args[1]

		if len(os.Args) == 2 && command == "init" {
			var response initResponse
			response.Status = "Success"
			response.Capabilities.Attach = false
			json_data, err := json.Marshal(response)
			if err != nil {
				message = err.Error()
			} else {
				fmt.Println(string(json_data))
				os.Exit(0)
			}
		}
		if len(os.Args) == 4 && command == "mount" {
			path := os.Args[2]
			in_data_unparsed := os.Args[3]
			fmt.Println(path)
			var in_data_parsed map[string]interface{}
			err := json.Unmarshal([]byte(in_data_unparsed), &in_data_parsed)
			if err != nil {
				message = err.Error()
			} else {
				fmt.Printf("%+v", in_data_parsed)
				var response mountResponse
				response.Status = "Success"
				response.Attached = true
				json_data, err := json.Marshal(response)
				if err != nil {
					message = err.Error()
				} else {
					fmt.Println(string(json_data))
					os.Exit(0)
				}
			}
		}
		if len(os.Args) == 3 && command == "unmount" {
			path := os.Args[2]
			fmt.Println(path)
			var response mountResponse
			response.Status = "Success"
			response.Attached = false
			json_data, err := json.Marshal(response)
			if err != nil {
				message = err.Error()
			} else {
				fmt.Println(string(json_data))
				os.Exit(0)
			}
		}
	}
	var response baseResponse
	response.Status = "Not supported"
	response.Message = message
	json_data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(json_data))
	os.Exit(0)

}
