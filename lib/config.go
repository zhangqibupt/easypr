package lib

import (
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path/filepath"
)

type GlobalConfig struct {
	GithubAccessToken  string `json:"github_access_token"`
	AWSAccessKeyID     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
}

func LoadGlobalConfig() (*GlobalConfig, error) {
	path, err := getGlobalConfigPath()
	if err != nil {
		return nil, err
	}

	// check if file exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		if debug {
			color.Red("Failed to read config file from %s: %s", path, err)
		}
		return nil, err
	}

	var config GlobalConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		if debug {
			color.Red("Failed to parse config file: %s", err)
		}
		return nil, err
	}
	return &config, nil
}

func getGlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".easypr_config.json"), nil
}

func SaveGlobalConfig(config *GlobalConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		if debug {
			color.Red("Failed to marshal config: %s", err)
		}
		return err
	}

	path, err := getGlobalConfigPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		if debug {
			color.Red("Failed to write config file: %s", err)
		}
		return err
	}
	return nil
}
