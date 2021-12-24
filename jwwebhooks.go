package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type WebhookMetadata struct {
	Description string   `json:"description"`
	Events      []string `json:"events"`
	Name        string   `json:"name"`
	SiteIds     []string `json:"site_ids"`
	Url         string   `json:"webhook_url"`
}

type WebhookDefinition struct {
	Created      string          `json:"created"`
	Id           string          `json:"id"`
	LastModified string          `json:"last_modified"`
	MetaData     WebhookMetadata `json:"metadata"`
	// relationship Ignored
	// schema Ignored
	Type string `json:"type"`
}

type WebhookResponse struct {
	Page       int                 `json:"page"`
	PageLength int                 `json:"page_length"`
	Total      int                 `json:"total"`
	Webhooks   []WebhookDefinition `json:"webhooks"`
}

type WebhookCreatePatch struct {
	Metadata WebhookMetadata `json:"metadata"`
}

var secret = ""

const service_url = "https://api.jwplayer.com/v2/webhooks/"

var jwClient = &http.Client{Timeout: time.Second * 10}

func createRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if nil != err {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+secret)
	return req, err
}

func Setup(apiSecret string) {
	secret = apiSecret
}

func WebhooksDefinitions() ([]WebhookDefinition, error) {
	req, err := createRequest("GET", service_url, nil)
	if nil != err {
		return []WebhookDefinition{}, err
	}
	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("page_length", "250")
	req.URL.RawQuery = q.Encode()

	res, err := jwClient.Do(req)
	if nil != err {
		return []WebhookDefinition{}, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return []WebhookDefinition{}, err
	}

	var webhooks WebhookResponse
	err = json.Unmarshal(b, &webhooks)
	if nil != err {
		fmt.Println("Spec read")
		//fmt.Println(definitions)
		return []WebhookDefinition{}, err
	}

	return webhooks.Webhooks, nil
}

func CreateWebhook(declaration Declaration) error {

	metadata := WebhookMetadata{declaration.description, declaration.events, declaration.name, declaration.siteIds, declaration.endpoint}
	create := WebhookCreatePatch{metadata}

	b, err := json.Marshal(create)
	if err != nil {
		return err
	}

	req, err := createRequest("POST", service_url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := jwClient.Do(req)
	if nil != err {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(ioutil.Discard, res.Body) // Throw away the body which must be read to EOF
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		return fmt.Errorf("webhook declaration %v not created, service returend statuscode %d", declaration, res.StatusCode)
	}

	return nil
}

func UpdateWebhook(id string, declaration Declaration) error {

	metadata := WebhookMetadata{declaration.description, declaration.events, declaration.name, declaration.siteIds, declaration.endpoint}
	update := WebhookCreatePatch{metadata}

	b, err := json.Marshal(update)
	if err != nil {
		return err
	}

	patchUrl := service_url + id + "/"
	req, err := createRequest("PATCH", patchUrl, bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := jwClient.Do(req)
	if nil != err {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(ioutil.Discard, res.Body) // Throw away the body which must be read to EOF
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("webhook declaration %v not updated, service returend statuscode %d", declaration, res.StatusCode)
	}
	return nil
}

func DeleteWebhook(id string) error {
	deleteUrl := service_url + id + "/"
	req, err := createRequest("DELETE", deleteUrl, nil)
	if err != nil {
		return err
	}
	res, err := jwClient.Do(req)
	if nil != err {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(ioutil.Discard, res.Body) // Throw away the body which must be read to EOF
	if err != nil {
		return err
	}

	if res.StatusCode != 204 {
		return fmt.Errorf("webhook id %s not deleted, service returend statuscode %d", id, res.StatusCode)
	}
	return nil
}
