package main

import (
	"os"

	"gambim.com/blockchain/client"
)

func main() {
	defer os.Exit(0)
	cli := client.CommandLine{}
	cli.Run()
}
