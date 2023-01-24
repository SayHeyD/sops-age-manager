package key

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
)

func GetAvailableKeys(keyDirPath string) []*Key {
	var keys []*Key

	if keyDirPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("cannot get the users home directory: %v", err)
		}

		keyDirPath = homeDir + string(os.PathSeparator) + ".age"
	}

	keyDir := os.DirFS(keyDirPath)
	err := fs.WalkDir(keyDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileSuffix := ".txt"

		fullPath := keyDirPath + string(os.PathSeparator) + path

		if !strings.HasSuffix(path, fileSuffix) {
			return nil
		}
		keyName := strings.TrimSuffix(path, fileSuffix)

		keyFileContent, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		key := NewKey(keyName, fullPath, string(keyFileContent))

		for _, processedKey := range keys {
			if key.Name == processedKey.Name {
				return fmt.Errorf("multiple keys with the name \"%s\" were detected", key.Name)
			}
		}

		keys = append(keys, key)

		return nil
	})
	if err != nil {
		log.Fatalf("readKeyFiles: %v", err)
	}

	return keys
}
