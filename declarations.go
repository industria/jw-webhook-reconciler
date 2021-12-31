package main

import (
	"encoding/json"
	"fmt"
	"os"
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

func declarations(specFile string) ([]Declaration, error) {
	f, err := os.Open(specFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s : %v", specFile, err)
	}
	defer f.Close()
	var declarations map[string]webhookDeclaration
	err = json.NewDecoder(f).Decode(&declarations)
	if err != nil {
		return nil, fmt.Errorf("unable to decode %s : %v", specFile, err)
	}
	var result = make([]Declaration, 0, len(declarations))
	for name, decl := range declarations {
		newDeclaration(name, decl)
		result = append(result, *newDeclaration(name, decl))
	}

	return result, nil
}
