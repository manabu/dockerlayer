package main

import (
	"fmt"
	"os"

	"github.com/manabu/dockerlayer/cmd"
	"github.com/manabu/dockerlayer/config"
)

const version = "0.1.3-dev"

func init() {
	config.VersionString = version
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
