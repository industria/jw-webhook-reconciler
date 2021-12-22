package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

type match struct {
	declaration Declaration
	definition  WebhookDefinition
}

type changeset struct {
	create []Declaration
	modify []match
	delete []WebhookDefinition
}

func declarationInDefinitions(declaration Declaration, definitions []WebhookDefinition) bool {
	for _, definition := range definitions {
		if declaration.name == definition.MetaData.Name {
			return true
		}
	}
	return false
}

var DefinitionNotFound = errors.New("DefinitionNotFound")

func findDefinitionFromDeclaration(declaration Declaration, definitions []WebhookDefinition) (match, error) {
	for _, definition := range definitions {
		if declaration.name == definition.MetaData.Name {
			return match{declaration, definition}, nil
		}
	}
	return match{}, DefinitionNotFound
}

func definitionInDeclarations(definition WebhookDefinition, declarations []Declaration) bool {
	for _, declaration := range declarations {
		if definition.MetaData.Name == declaration.name {
			return true
		}
	}
	return false
}

func different(match match) bool {
	// TODO: compare
	return false
}

func changeSet(declarations []Declaration, definitions []WebhookDefinition) *changeset {
	var create []Declaration
	var modify []match
	var delete []WebhookDefinition

	// Create declarations not found in definitions
	for _, declaration := range declarations {
		if !declarationInDefinitions(declaration, definitions) {
			create = append(create, declaration)
		}
	}

	// Modify if declaration is found in definitions and the attribute values differs
	for _, declaration := range declarations {
		match, err := findDefinitionFromDeclaration(declaration, definitions)
		if nil == err && different(match) {
			modify = append(modify, match)
		}
	}

	// Delete definitions not found in declarations
	for _, definition := range definitions {
		if !definitionInDeclarations(definition, declarations) {
			delete = append(delete, definition)
		}
	}
	return &changeset{create, modify, delete}
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

	declarations, err := Declarations(*spec)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Declarations file:", declarations)

	fmt.Println("Specification file:", *spec)
	fmt.Println("Secret:", *secret)

	Setup(*secret)
	definitions, err := WebhooksDefinitions()
	if err != nil {
		fmt.Printf("Failed to get the webhooks from JW service \n")
		os.Exit(1)
	}

	if cmd == "list" {
		const padding = 3
		w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
		fmt.Fprintf(w, "Id\tName\tURL\tSites\tEvents  \n")

		for _, definition := range definitions {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t \n", definition.Id, definition.MetaData.Name, definition.MetaData.Url, definition.MetaData.SiteIds, definition.MetaData.Events)
		}
		w.Flush()
	} else if cmd == "diff" {

		changeset := changeSet(declarations, definitions)
		fmt.Println("Create:")
		fmt.Println(changeset.create)

		fmt.Println("Modify:")
		fmt.Println(changeset.modify)

		fmt.Println("Delete:")
		fmt.Println(changeset.delete)

	} else {
		fmt.Printf("Unknown command %s\n", cmd)
	}

}
