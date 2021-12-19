package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("Usage: %s <command> arguments.. \n", os.Args[0])
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

	if !validCommand(flag.Arg(0)) {
		fmt.Printf("Unknown command %s", flag.Arg(0))
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Arguments left", flag.NArg())

	//args := os.Args
	fmt.Println("Specification file:", *spec)
	fmt.Println("Show ids:", *showIds)
	fmt.Println("Secret:", *secret)

}
