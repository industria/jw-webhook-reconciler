package main

import (
	"flag"
	"fmt"
)

var secret = flag.String("secret", "", "API secret to use for communicating with JW")

func main() {

	//args := os.Args
	fmt.Println(secret)
}
