//go:build main
// +build main

package main

import (
	"embed"
	"github.com/SayHeyD/sops-age-manager/cmd"
	"log"
)

//go:embed version.txt
var versionFile embed.FS

func main() {
	version, err := versionFile.ReadFile("version.txt")
	if err != nil {
		log.Fatalf("error reading version file: %v", err)
	}

	cmd.Execute(string(version))
}
