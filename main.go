package main

import (
	"os"
	"flag"
	"fmt"
)

var Version = "development"


func main() {
	version := flag.Bool("version", false, "version")
	flag.Parse()

	if *version {
		fmt.Println("Version:", Version)
		os.Exit(0)
	}

}
