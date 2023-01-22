package key

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
)

func GetAvailableKeys() []*Key {
	var keys []*Key

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot get the users home directory: %v", err))
	}

	keyDirPath := homeDir + string(os.PathSeparator) + ".age"

	keyDir := os.DirFS(keyDirPath)
	err = fs.WalkDir(keyDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileSuffix := ".txt"

		fullPath := keyDirPath + string(os.PathSeparator) + path

		fmt.Println(fullPath)

		if !strings.HasSuffix(path, fileSuffix) {
			return nil
		}
		keyName := strings.TrimSuffix(path, fileSuffix)

		keyFileContent, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		key := NewKey(keyName, string(keyFileContent))
		keys = append(keys, key)

		return nil
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("readKeyFiles: %v", err))
	}

	return keys
}