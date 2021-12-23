package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Hand the JSON used uppercase the suffixes would not be needed
type webhookDeclaration struct {
	Description string   `json:"description"`
	Events      []string `json:"events"`
	Site_ids    []string `json:"site_ids"`
	Endpoint    string   `json:"endpoint"`
}

type Declaration struct {
	name        string
	description string
	events      []string
	siteIds     []string
	endpoint    string
}

func newDeclaration(name string, decl webhookDeclaration) *Declaration {
	return &Declaration{name, decl.Description, decl.Events, decl.Site_ids, decl.Endpoint}
}

func Declarations(declarationFile string) ([]Declaration, error) {
	var result []Declaration

	f, err := ioutil.ReadFile(declarationFile)
	if err != nil {
		return result, fmt.Errorf("failed to read %s", declarationFile)
	}

	var declarations map[string]webhookDeclaration
	err = json.Unmarshal([]byte(f), &declarations)
	if err != nil {
		return result, fmt.Errorf("failed to parse %s", declarationFile)
	}

	for name, decl := range declarations {
		newDeclaration(name, decl)
		result = append(result, *newDeclaration(name, decl))
	}

	return result, nil
}
