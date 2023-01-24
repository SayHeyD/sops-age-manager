package config

import (
	"fmt"
	"github.com/SayHeyD/sops-age-manager/test"
	"os"
	"testing"
)

func getExpectedKeyName() string {
	return "some_key_name"
}

func getExpectedFileContent() string {
	return fmt.Sprintf("key: %s", getExpectedKeyName())
}

func TestNewConfig(t *testing.T) {
	t.Parallel()

	expectedKeyName := getExpectedKeyName()

	config := NewConfig(expectedKeyName)

	if config.KeyName != expectedKeyName {
		t.Fatalf("The KeyName \"%s\" did not match the expected value \"%s\"", config.KeyName, expectedKeyName)
	}
}

func TestNewConfigFromFileShouldReturnANonNilValue(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	testConfigFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"

	var config *Config

	config, err := NewConfigFromFile(testConfigFilePath)
	if err != nil {
		t.Fatalf("Error creating config from file: %v", err)
	}

	if config == nil {
		t.Fatalf("Returned object is nil")
	}
}

func TestNewConfigFromFileShouldCreateAFileIfNotExist(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	testConfigFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"

	_, err := NewConfigFromFile(testConfigFilePath)
	if err != nil {
		t.Fatalf("Error creating config from file: %v", err)
	}

	if _, err := os.Stat(testConfigFilePath); os.IsNotExist(err) {
		t.Fatalf("File does not exist after NewConfigFromFile() is called")
	}
}

func TestNewConfigFromFileShouldNotCreateNewFileIfOneAlreadyExists(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	expectedFileContent := getExpectedFileContent()

	testConfigFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"

	configFile, err := os.Create(testConfigFilePath)
	if err != nil {
		t.Fatalf("Error creating config for testing: %v", err)
	}

	if _, err = configFile.WriteString(expectedFileContent); err != nil {
		t.Fatalf("Error writing to config for testing: %v", err)
	}

	if err = configFile.Close(); err != nil {
		t.Fatalf("Error closing the config for testing: %v", err)
	}

	_, err = NewConfigFromFile(testConfigFilePath)
	if err != nil {
		t.Fatalf("Error creating config from file: %v", err)
	}

	fileContent, err := os.ReadFile(testConfigFilePath)
	fileContentString := string(fileContent)

	if fileContentString != expectedFileContent {
		t.Fatalf("The file content \"%s\" differs from whats expected \"%s\"", fileContentString, expectedFileContent)
	}
}

func TestNewConfigFromFileShouldReturnAConfigWithTheCorrectValues(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)
	expectedKeyName := getExpectedKeyName()
	expectedFileContent := getExpectedFileContent()

	testConfigFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"

	configFile, err := os.Create(testConfigFilePath)
	if err != nil {
		t.Fatalf("Error creating config for testing: %v", err)
	}

	if _, err = configFile.WriteString(expectedFileContent); err != nil {
		t.Fatalf("Error writing to config for testing: %v", err)
	}

	if err = configFile.Close(); err != nil {
		t.Fatalf("Error closing the config for testing: %v", err)
	}

	var config *Config

	config, err = NewConfigFromFile(testConfigFilePath)
	if err != nil {
		t.Fatalf("Error creating config from file: %v", err)
	}

	if config.KeyName != expectedKeyName {
		t.Fatalf("The key name \"%s\" differs from whats expected \"%s\"", config.KeyName, expectedKeyName)
	}
}

func TestConfigWriteGeneratesNewFileWhenNotExists(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)
	expectedKeyName := getExpectedKeyName()
	expectedFileContent := getExpectedFileContent()

	configFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"

	if err := NewConfig(expectedKeyName).Write(configFilePath); err != nil {
		t.Fatalf("could not write config file \"%s\": %v", configFilePath, err)
	}

	configFileContentBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		t.Fatalf("could not read config file \"%s\": %v", configFilePath, err)
	}

	configFileContent := string(configFileContentBytes)

	if configFileContent != expectedFileContent {
		t.Fatalf("The file content \"%s\" does not match with the expected content \"%s\"", configFileContent, expectedFileContent)
	}
}

func TestGetConfigDirPath(t *testing.T) {
	t.Parallel()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("An error occured while executing getting the users homedir: %v", err)
	}

	expectedConfigDir := homeDir + string(os.PathSeparator) + ".sops-age-manager"

	configDir, err := getConfigDirPath()
	if err != nil {
		t.Fatalf("An error occured while executing getConfigDirPath(): %v", err)
	}

	if configDir != expectedConfigDir {
		t.Fatalf("The path returned by getConfigDirPath() \"%s\" did not match the expected value \"%s\"", configDir, expectedConfigDir)
	}
}
