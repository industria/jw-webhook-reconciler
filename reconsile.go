package main

import (
	"flag"
	"fmt"
)

func main() {
	var secret = flag.String("secret", "", "API secret to use for communicating with JW")
	var showIds = flag.Bool("id", false, "Include JW webhook id in out from list.")
	var spec = flag.String("spec", "", "Path to the specification file.")
	flag.Parse()

	//flag.PrintDefaults()

	fmt.Println("Arguments left", flag.NArg())

	//args := os.Args
	fmt.Println("Specification file:", *spec)
	fmt.Println("Show ids:", *showIds)
	fmt.Println("Secret:", *secret)

}
