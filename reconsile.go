package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

// Hand the JSON used uppercase the suffixes would not be needed
type WebhookDefinition struct {
	Description string   `json:"description"`
	Events      []string `json:"events"`
	Site_ids    []string `json:"site_ids"`
	Endpoint    string   `json:"endpoint"`
}

func usage() {
	fmt.Printf("Usage: %s arguments <command> \n", os.Args[0])
	fmt.Printf("  command : list, diff or apply\n\n")
	flag.PrintDefaults()
}

func validCommand(cmd string) bool {
	cmds := []string{"list", "diff", "apply"}
	for _, s := range cmds {
		if s == cmd {
			return true
		}
	}
	return false
}

func main() {
	var secret = flag.String("secret", "", "API secret to use for communicating with JW")
	var showIds = flag.Bool("id", false, "Include JW webhook id in out from list.")
	var spec = flag.String("spec", "", "Path to the specification file.")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	cmd := flag.Args()[0]
	if !validCommand(cmd) {
		fmt.Printf("Unknown command %s \n", cmd)
		flag.Usage()
		os.Exit(1)
	}

	var _, err = os.Stat(*spec)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("File %s not found\n", *spec)
		os.Exit(1)
	}

	f, err := ioutil.ReadFile(*spec)
	if err != nil {
		fmt.Printf("Failed to read %s\n", *spec)
		os.Exit(1)
	}

	var definitions map[string]WebhookDefinition
	json.Unmarshal([]byte(f), &definitions)

	fmt.Println(definitions)

	fmt.Println("Specification file:", *spec)
	fmt.Println("Show ids:", *showIds)
	fmt.Println("Secret:", *secret)

}
