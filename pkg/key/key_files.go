package key

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
)

func FindAvailableKeys(keyDirPath string) ([]*Key, error) {
	var keys []*Key

	if keyDirPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot get the users home directory: %w", err)
		}

		keyDirPath = homeDir + string(os.PathSeparator) + ".age"
	}

	if _, err := os.Stat(keyDirPath); os.IsNotExist(err) {
		log.Printf("key directory does not exist: %v\n", keyDirPath)
		log.Printf("creating key directory: %v\n", keyDirPath)

		if err := os.Mkdir(keyDirPath, 0700); err != nil {
			return nil, fmt.Errorf("cannot create the key directory '%s': %w", keyDirPath, err)
		}
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
		return nil, fmt.Errorf("reading key files: %w", err)
	}

	return keys, nil
}

func GetAvailableKeys(keyDirPath string) []*Key {
	keys, err := FindAvailableKeys(keyDirPath)
	if err != nil {
		log.Fatalf("reading key files: %v", err)
	}

	if keys == nil {
		if keyDirPath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("cannot get the users home directory: %v", err)
			}
			keyDirPath = homeDir + string(os.PathSeparator) + ".age"
		}
		keyDir := os.DirFS(keyDirPath)
		log.Fatalf("reading key files: No keys were found in the the key dir '%s%s'", keyDir, string(os.PathSeparator))
	}

	return keys
}
