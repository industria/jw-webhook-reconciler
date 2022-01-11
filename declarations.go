package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Declaration struct {
	name        string
	description string
	events      []string
	siteIds     []string
	endpoint    string
}

func declarations(specFile string) ([]Declaration, error) {
	if len(specFile) == 0 {
		return nil, fmt.Errorf("missing file name")
	}

	f, err := os.Open(specFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s : %v", specFile, err)
	}
	defer f.Close()

	var declarations map[string]struct {
		Description string   `json:"description"`
		Events      []string `json:"events"`
		Site_ids    []string `json:"site_ids"`
		Endpoint    string   `json:"endpoint"`
	}
	err = json.NewDecoder(f).Decode(&declarations)
	if err != nil {
		return nil, fmt.Errorf("unable to decode %s : %v", specFile, err)
	}
	var result = make([]Declaration, 0, len(declarations))
	for name, decl := range declarations {
		result = append(result, Declaration{name, decl.Description, decl.Events, decl.Site_ids, decl.Endpoint})
	}

	return result, nil
}
