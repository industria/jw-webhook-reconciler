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

type Webhooks struct {
	secret     string
	serviceURL string
	httpClient *http.Client
}

func newWebhooks(secret string) *Webhooks {
	return &Webhooks{
		secret:     secret,
		serviceURL: "https://api.jwplayer.com/v2/webhooks/",
		httpClient: &http.Client{Timeout: time.Second * 10},
	}
}

func (w *Webhooks) request(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+w.secret)
	return req, err
}

func (w *Webhooks) definitions() ([]WebhookDefinition, error) {
	req, err := w.request("GET", w.serviceURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("page_length", "250")
	req.URL.RawQuery = q.Encode()

	res, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var webhooks WebhookResponse
	err = json.Unmarshal(b, &webhooks)
	if err != nil {
		return nil, err
	}

	return webhooks.Webhooks, nil
}

func (w *Webhooks) create(declaration Declaration) error {
	metadata := WebhookMetadata{declaration.description, declaration.events, declaration.name, declaration.siteIds, declaration.endpoint}
	create := WebhookCreatePatch{metadata}

	b, err := json.Marshal(create)
	if err != nil {
		return err
	}

	req, err := w.request("POST", w.serviceURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := w.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(ioutil.Discard, res.Body) // Throw away the body which must be read to EOF
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("webhook declaration %v not created, service returend statuscode %d", declaration, res.StatusCode)
	}

	return nil
}

func (w *Webhooks) update(id string, declaration Declaration) error {
	metadata := WebhookMetadata{declaration.description, declaration.events, declaration.name, declaration.siteIds, declaration.endpoint}
	update := WebhookCreatePatch{metadata}

	b, err := json.Marshal(update)
	if err != nil {
		return err
	}

	patchUrl := w.serviceURL + id + "/"
	req, err := w.request("PATCH", patchUrl, bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := w.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(ioutil.Discard, res.Body) // Throw away the body which must be read to EOF
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook declaration %v not updated, service returend statuscode %d", declaration, res.StatusCode)
	}
	return nil
}

func (w *Webhooks) delete(id string) error {
	deleteUrl := w.serviceURL + id + "/"
	req, err := w.request("DELETE", deleteUrl, nil)
	if err != nil {
		return err
	}
	res, err := w.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(ioutil.Discard, res.Body) // Throw away the body which must be read to EOF
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("webhook id %s not deleted, service returend statuscode %d", id, res.StatusCode)
	}
	return nil
}
