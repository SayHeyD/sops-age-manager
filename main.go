//go:build main
// +build main

package main

import (
	"embed"
	"log"

	"github.com/SayHeyD/sops-age-manager/cmd"
)

//go:embed version.txt
var versionFile embed.FS

//go:embed Logo.png
var logoFile []byte

func main() {
	version, err := versionFile.ReadFile("version.txt")
	if err != nil {
		log.Fatalf("error reading version file: %v", err)
	}

	cmd.Execute(string(version), logoFile)
}
