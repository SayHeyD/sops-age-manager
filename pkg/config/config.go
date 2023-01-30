package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const (
	defaultConfig = `encryption-key: ""
decryption-key: ""
key-dir: ""`

	configFileEnv = "SOPS_AGE_MANAGER_CONFIG_DIR"
)

type Config struct {
	EncryptionKeyName string `yaml:"encryption-key"`
	DecryptionKeyName string `yaml:"decryption-key"`
	KeyDir            string `yaml:"key-dir"`
}

func GetConfigFilePath() string {
	return os.Getenv(configFileEnv)
}

func NewConfig(encryptionKeyName string, decryptionKeyName string, keyDir string) *Config {
	return &Config{
		EncryptionKeyName: encryptionKeyName,
		DecryptionKeyName: decryptionKeyName,
		KeyDir:            keyDir,
	}
}

func NewConfigFromFile(path string) (*Config, error) {
	contentBytes, err := getConfigFileContents(path)
	if err != nil {
		return nil, fmt.Errorf("trying generate a new config from a file: %v", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(contentBytes, config); err != nil {
		return nil, fmt.Errorf("cannot unmarshal the config file: %v", err)
	}

	return config, nil
}

func (c *Config) Raw(path string) (string, error) {
	contentBytes, err := getConfigFileContents(path)
	if err != nil {
		return "", fmt.Errorf("trying to dump raw config file: %v", err)
	}

	return string(contentBytes), nil
}

func (c *Config) Write(path string) error {

	configPath, err := getConfigDirPath()
	if err != nil {
		return fmt.Errorf("could not get config dir path: %v", err)
	}

	if path != "" {
		configPath = path
	}

	configFile, err := os.Create(configPath + string(os.PathSeparator) + "config.yaml")
	if err != nil {
		return fmt.Errorf("could not create the config file: %v", err)
	}
	defer configFile.Close()

	configFileContentBytes, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("could not marshal config: %v", err)
	}

	trimmedConfigFileContentBytes := strings.Trim(string(configFileContentBytes), "\t\n ")

	_, err = configFile.WriteString(trimmedConfigFileContentBytes)
	if err != nil {
		return fmt.Errorf("could not write to config: %v", err)
	}

	return nil
}

func getConfigDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get the users home directory: %v", err)
	}

	samConfigDir := homeDir + string(os.PathSeparator) + ".sops-age-manager"

	if _, err := os.Stat(samConfigDir); os.IsNotExist(err) {
		if err = os.Mkdir(samConfigDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("cannot create the sops-age-manager config directory: %v", err)
		}
	}

	return samConfigDir, nil
}

func getConfigFileContents(path string) ([]byte, error) {
	configDir, err := getConfigDirPath()
	if err != nil {
		return nil, fmt.Errorf("trying to read config dir: %v", err)
	}
	configFilePath := configDir + string(os.PathSeparator) + "config.yaml"

	if path != "" {
		configFilePath = path
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		configFile, err := os.Create(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("trying to create the config file: %v", err)
		}
		defer configFile.Close()

		_, err = configFile.WriteString(defaultConfig)
		if err != nil {
			return nil, fmt.Errorf("trying write to the config file: %v", err)
		}

		return []byte(defaultConfig), nil
	}

	contentBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("trying to read the config file: %v", err)
	}

	return contentBytes, nil
}
