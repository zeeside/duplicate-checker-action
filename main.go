// This Github Action
package main

import (
	"os"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		path, _ := os.Getwd()
		githubactions.Fatalf("%s directory: %s", err.Error(), path)
	}
	log := NewActionLogger(config)
	fp := NewFileProcessor(config, log.logger)
	fp.Run()
	fp.PrintOutput()
}
