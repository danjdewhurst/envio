package main

import (
	"os"

	"github.com/danjdewhurst/envio/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
