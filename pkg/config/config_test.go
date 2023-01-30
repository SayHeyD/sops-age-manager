package config

import (
	"fmt"
	"github.com/SayHeyD/sops-age-manager/test"
	"os"
	"strings"
	"testing"
)

func getExpectedDecryptionKeyName() string {
	return "some_key_name"
}

func getExpectedEncryptionKeyName() string {
	return "some_other_key_name"
}

func getExpectedKeyDir() string {
	tempDir := os.TempDir() + "key-dir"
	return tempDir
}

func getExpectedFileContent() string {
	defaultConfigTemplateString := `encryption-key: %s
decryption-key: %s
key-dir: %s`
	return fmt.Sprintf(defaultConfigTemplateString, getExpectedEncryptionKeyName(),
		getExpectedDecryptionKeyName(), getExpectedKeyDir())
}

func TestNewConfig(t *testing.T) {
	t.Parallel()

	expectedEncryptionKeyName := getExpectedEncryptionKeyName()
	expectedDecryptionKeyName := getExpectedDecryptionKeyName()
	expectedKeyDir := getExpectedKeyDir()

	config := NewConfig(expectedEncryptionKeyName, expectedDecryptionKeyName, expectedKeyDir)

	if config.EncryptionKeyName != expectedEncryptionKeyName {
		t.Fatalf("The EncryptionKeyName \"%s\" did not match the expected value \"%s\"", config.EncryptionKeyName, expectedEncryptionKeyName)
	}

	if config.DecryptionKeyName != expectedDecryptionKeyName {
		t.Fatalf("The DecryptionKeyName \"%s\" did not match the expected value \"%s\"", config.DecryptionKeyName, expectedDecryptionKeyName)
	}

	if config.KeyDir != expectedKeyDir {
		t.Fatalf("The KeyDir \"%s\" did not match the expected value \"%s\"", config.KeyDir, expectedKeyDir)
	}
}

func TestNewConfigFromFileShouldReturnANonNilValue(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	testConfigFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"
	if err := os.Setenv(configFileEnv, testConfigFilePath); err != nil {
		t.Fatalf("could not set \"%s\" to \"%s\"", configFileEnv, testConfigFilePath)
	}

	var config *Config

	config, err := NewConfigFromFile()
	if err != nil {
		t.Fatalf("Error creating config from file: %v", err)
	}

	if config == nil {
		t.Fatalf("Returned object is nil")
	}
}

func TestNewConfigFromFileShouldReturnAConfigWithTheCorrectValues(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	// defer testDir.CleanTestDir(t)
	expectedEncryptionKeyName := getExpectedEncryptionKeyName()
	expectedDecryptionKeyName := getExpectedDecryptionKeyName()
	expectedKeyDir := getExpectedKeyDir()
	expectedFileContent := getExpectedFileContent()

	testConfigFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"
	if err := os.Setenv(configFileEnv, testConfigFilePath); err != nil {
		t.Fatalf("could not set \"%s\" to \"%s\"", configFileEnv, testConfigFilePath)
	}

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

	config, err = NewConfigFromFile()
	if err != nil {
		t.Fatalf("Error creating config from file: %v", err)
	}

	if config.EncryptionKeyName != expectedEncryptionKeyName {
		t.Fatalf("The encryption key name \"%s\" differs from whats expected \"%s\"", config.EncryptionKeyName, expectedEncryptionKeyName)
	}

	if config.DecryptionKeyName != expectedDecryptionKeyName {
		t.Fatalf("The decryption key name \"%s\" differs from whats expected \"%s\"", config.DecryptionKeyName, expectedEncryptionKeyName)
	}

	if config.KeyDir != expectedKeyDir {
		t.Fatalf("The key dir \"%s\" differs from whats expected \"%s\"", config.KeyDir, expectedKeyDir)
	}
}

func TestConfigWriteGeneratesNewFileWhenNotExists(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	configFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"

	if err := os.Setenv(configFileEnv, configFilePath); err != nil {
		t.Fatalf("could not set \"%s\" to \"%s\"", configFileEnv, configFilePath)
	}

	expectedEncryptionKeyName := getExpectedEncryptionKeyName()
	expectedDecryptionKeyName := getExpectedDecryptionKeyName()
	expectedKeyDir := getExpectedKeyDir()
	expectedFileContent := getExpectedFileContent()

	if err := NewConfig(expectedEncryptionKeyName, expectedDecryptionKeyName, expectedKeyDir).Write(); err != nil {
		t.Fatalf("could not write config file in folder \"%s\": %v", testDir.Path, err)
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

func TestGetConfigDirPathReturnsCorrectDirectory(t *testing.T) {
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

func TestGetConfigFileContentsShouldReturnTheDefaultConfigIfNoFileExists(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	testConfigFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"

	fileContent, err := getConfigFileContents(testConfigFilePath)
	fileContentString := strings.Trim(string(fileContent), "\n\t ")
	if err != nil {
		t.Fatalf("Error creating config from file: %v", err)
	}

	if fileContentString != defaultConfig {
		t.Fatalf("The file content \"%s\" does not match with the expected value \"%s\"", fileContentString, defaultConfig)
	}
}

func TestGetConfigFileContentsShouldCreateAFileIfNotExist(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	testConfigFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"

	_, err := getConfigFileContents(testConfigFilePath)
	if err != nil {
		t.Fatalf("Error creating config from file: %v", err)
	}

	if _, err := os.Stat(testConfigFilePath); os.IsNotExist(err) {
		t.Fatalf("File does not exist after getConfigFileContents() is called")
	}
}

func TestGetConfigFileContentsShouldNotCreateNewFileIfOneAlreadyExists(t *testing.T) {
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

	_, err = getConfigFileContents(testConfigFilePath)
	if err != nil {
		t.Fatalf("Error reading config from file: %v", err)
	}

	fileContent, err := os.ReadFile(testConfigFilePath)
	if err != nil {
		t.Fatalf("Could not read the test config file: %v", err)
	}

	fileContentString := string(fileContent)

	if fileContentString != expectedFileContent {
		t.Fatalf("The file content \"%s\" differs from whats expected \"%s\"", fileContentString, expectedFileContent)
	}
}

func TestGetConfigFileContentsShouldReturnAExpectedFileContents(t *testing.T) {
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

	configContent, err := getConfigFileContents(testConfigFilePath)
	if err != nil {
		t.Fatalf("Error getting config from file: %v", err)
	}

	if string(configContent) != expectedFileContent {
		t.Fatalf("The fetched content \"%s\" differs from whats expected \"%s\"",
			string(configContent), expectedFileContent)
	}
}

func TestRawShouldReturnANonEmptyString(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	testConfigFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"
	if err := os.Setenv(configFileEnv, testConfigFilePath); err != nil {
		t.Fatalf("could not set \"%s\" to \"%s\"", configFileEnv, testConfigFilePath)
	}

	var config *Config

	config, err := NewConfigFromFile()
	if err != nil {
		t.Fatalf("Error creating config from file: %v", err)
	}

	configContent, err := config.Raw()
	if err != nil {
		t.Fatalf("Could not get the config from a file: %v", err)
	}

	if configContent == "" {
		t.Fatalf("Returned string is empty")
	}
}

func TestRawShouldReturnTheExpectedFileContent(t *testing.T) {
	t.Parallel()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	expectedFileContent := getExpectedFileContent()

	testConfigFilePath := testDir.Path + string(os.PathSeparator) + "config.yaml"
	if err := os.Setenv(configFileEnv, testConfigFilePath); err != nil {
		t.Fatalf("could not set \"%s\" to \"%s\"", configFileEnv, testConfigFilePath)
	}

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

	config, err = NewConfigFromFile()
	if err != nil {
		t.Fatalf("Could not get the config from a file: %v", err)
	}

	configContent, err := config.Raw()
	if err != nil {
		t.Fatalf("Could not get the content of the config: %v", err)
	}

	if configContent != expectedFileContent {
		t.Fatalf("raw config content \"%s\" does not match with the expected config \"%s\"", configContent, expectedFileContent)
	}
}

func TestGetConfigFilePathReturnsEmptyStringIfVarIsUnset(t *testing.T) {
	t.Parallel()

	err := os.Unsetenv(configFileEnv)
	if err != nil {
		t.Fatalf("could not unset env \"%s\"", configFileEnv)
	}

	configFilePath := getConfigFilePath()
	if configFilePath != "" {
		t.Fatalf("\"configFileEnv\" renturned value \"%s\" expected was \"\"", configFilePath)
	}
}

func TestGetConfigFilePathReturnsCorrectString(t *testing.T) {
	t.Parallel()

	expectedFilePath := "/some/random/directory/config.yaml"

	err := os.Setenv(configFileEnv, expectedFilePath)
	if err != nil {
		t.Fatalf("could not set env \"%s\" to \"%s\"", configFileEnv, expectedFilePath)
	}

	configFilePath := getConfigFilePath()
	if configFilePath != expectedFilePath {
		t.Fatalf("\"configFilePath\" renturned value \"%s\" expected was \"%s\"", configFilePath, expectedFilePath)
	}
}
