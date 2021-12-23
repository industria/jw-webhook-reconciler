package main

import (
	"encoding/json"
	"fmt"
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

func createRequest(method string) (*http.Request, error) {
	req, err := http.NewRequest(method, service_url, nil)
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
	req, err := createRequest("GET")
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
	// TODO: b is basically what we nned to post

	fmt.Println(string(b))
	return nil
}

func UpdateWebhook(declaration Declaration) error {

	metadata := WebhookMetadata{declaration.description, declaration.events, declaration.name, declaration.siteIds, declaration.endpoint}
	update := WebhookCreatePatch{metadata}

	b, err := json.Marshal(update)
	if err != nil {
		return err
	}
	// TODO: b is basically what we nned to Patch

	fmt.Println(string(b))
	return nil
}

func DeleteWebhook(declaration Declaration) error {

	// TODO delete to url
	return nil
}
