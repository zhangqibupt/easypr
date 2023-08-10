package lib

import (
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
)

type Config struct {
	Assignees []string `json:"assignees"`
}

func GetConfigPath() (string, error) {
	gitDirectory, err := TopLevelDirectory()
	if err != nil {
		return "", err
	}
	return gitDirectory + "/.git/.fwpr_config.json", nil
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
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

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		if debug {
			color.Red("Failed to parse config file: %s", err)
		}
		return nil, err
	}
	return &config, nil
}

func SaveConfig(config *Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		if debug {
			color.Red("Failed to marshal config: %s", err)
		}
		return err
	}

	path, err := GetConfigPath()
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
