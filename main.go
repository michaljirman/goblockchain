package main

import (
	"os"

	"github.com/michaljirman/goblockchain/cli"
)

// Main entry point of the program.
// Examples:
//DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go createblockchain -address "John"
//DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go printchain
//DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go getbalance -address "John"
//DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go send -from "John" -to "fred" -amount 50
//DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go getbalance -address "fred"
//DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go getbalance -address "John"

//DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go createwallet
//DEVELOPMENT_LOGGER=TRUE DEBUG=TRUE go run main.go listaddresses
func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}
