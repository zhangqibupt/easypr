package lib

import (
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
)

type RepoConfig struct {
	Assignees []string `json:"assignees"`
	Upstream  string   `json:"upstream"`
}

func GetRepoConfigPath() (string, error) {
	gitDirectory, err := TopLevelDirectory()
	if err != nil {
		return "", err
	}
	return gitDirectory + "/.git/.fwpr_config.json", nil
}

func LoadRepoConfig() (*RepoConfig, error) {
	path, err := GetRepoConfigPath()
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

	var config RepoConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		if debug {
			color.Red("Failed to parse config file: %s", err)
		}
		return nil, err
	}
	return &config, nil
}

func SaveRepoConfig(config *RepoConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		if debug {
			color.Red("Failed to marshal config: %s", err)
		}
		return err
	}

	path, err := GetRepoConfigPath()
	if err != nil {
		return err
	}

	_, err = os.Stat(path)
	fileExists := !os.IsNotExist(err)

	// If the file doesn't exist, create it
	if !fileExists {
		file, err := os.Create(path)
		if err != nil {
			if debug {
				color.Red("Failed to create config file: %s", err)
			}
			return err
		}
		defer file.Close()
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
