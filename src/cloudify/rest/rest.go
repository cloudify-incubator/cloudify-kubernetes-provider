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

func Get(url string, user string, password string, tenant string) []byte {
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

func Delete(url string, user string, password string, tenant string) []byte {
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
