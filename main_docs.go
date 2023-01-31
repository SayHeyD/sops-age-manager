//go:build docs
// +build docs

package main

import (
	"github.com/SayHeyD/sops-age-manager/cmd"
	"github.com/spf13/cobra/doc"
	"log"
	"os"
)

func main() {
	if _, err := os.Stat("./docs"); os.IsNotExist(err) {
		if err := os.Mkdir("./docs", os.ModePerm); err != nil {
			log.Fatalf("could not create docs directory: %v", err)
		}
	}

	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs")
	if err != nil {
		log.Fatalf("could not create markdown tree: %v", err)
	}
}
