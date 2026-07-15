package main

import (
	"os"

	"github.com/cak-huel/cloudmanager/v2/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
