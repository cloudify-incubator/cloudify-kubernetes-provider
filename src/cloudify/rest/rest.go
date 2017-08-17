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

package rest

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func GetRequest(url, user, password, tenant, method string, body io.Reader) *http.Request {
	log.Printf("Use: %v:%v@%v#%s\n", user, password, url, tenant)

	var auth_string string
	auth_string = user + ":" + password
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth_string)))
	if len(tenant) > 0 {
		req.Header.Add("Tenant", tenant)
	}

	return req
}

func Get(url, user, password, tenant string) []byte {
	req := GetRequest(url, user, password, tenant, "GET", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response %s\n", string(body))
	return body
}

func Delete(url, user, password, tenant string) []byte {
	req := GetRequest(url, user, password, tenant, "DELETE", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response %s\n", string(body))
	return body
}

func Post(url, user, password, tenant string, data []byte) []byte {
	req := GetRequest(url, user, password, tenant, "POST", bytes.NewBuffer(data))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response %s\n", string(body))
	return body
}

func Put(url, user, password, tenant string, data []byte) []byte {
	req := GetRequest(url, user, password, tenant, "PUT", bytes.NewBuffer(data))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response %s\n", string(body))
	return body
}
