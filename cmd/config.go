package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"easypr/lib"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var setAssigneeCmd = &cobra.Command{
	Use:   "set-assignees [name1 name2...]",
	Short: "Set assignees for pull request",
	Run: func(cmd *cobra.Command, args []string) {
		config := lib.RepoConfig{
			Assignees: args,
		}
		err := lib.SaveRepoConfig(&config)
		if err != nil {
			color.Red("Failed to save config: %s", err)
			return
		}
		color.Green("Successfully set assignees for current repo")
	},
}

var setUpstreamCmd = &cobra.Command{
	Use:   "set-upstream [repo]",
	Short: "Specify the target repo of the pull request, it is useful when the repo is forked and you want to create pull request to the original repo",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			color.Red("Please specify the upstream repo")
			return
		}

		if !isValidGitHubRepoURL(args[0]) {
			color.Red("Please specify a valid repo url")
			return
		}

		config, _ := lib.LoadRepoConfig()
		if config == nil {
			config = &lib.RepoConfig{}
		}
		config.Upstream = args[0]

		err := lib.SaveRepoConfig(config)
		if err != nil {
			color.Red("Failed to save config: %s", err)
			return
		}
		color.Green("Successfully set upstream for current repo")
	},
}

var setAccessTokenCmd = &cobra.Command{
	Use:   "set-access-token [token]",
	Short: "Set the access token for authenticating with GitHub API",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			color.Red("Please specify the access token")
			return
		}

		// Add additional validation for the access token if needed

		config, err := lib.LoadGlobalConfig()
		if err != nil {
			color.Red("Failed to load config: %s", err)
			return
		}

		if config == nil {
			config = &lib.GlobalConfig{}
		}

		config.GithubAccessToken = args[0]

		err = lib.SaveGlobalConfig(config)
		if err != nil {
			color.Red("Failed to save config: %s", err)
			return
		}

		color.Green("Successfully set access token")
	},
}

func isValidGitHubRepoURL(url string) bool {
	repoRegex := regexp.MustCompile(`((http|git|ssh|http(s)|file|/?)|(git@[\w.]+))(:(//)?)([\w.@:/\-~]+)(\.git)(/)?`)

	return repoRegex.MatchString(url)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List config for pull request",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := lib.LoadRepoConfig()
		if err != nil {
			color.Red("Failed to load config: %s", err)
			return
		}

		if config == nil {
			color.Yellow("You haven't set config for this repo yet.")
			return
		}

		color.Cyan("pr assignees for current repo: %s", config.Assignees)
		color.Cyan("upstream for current repo: %s", config.Upstream)
	},
}

func init() {
	configCmd.AddCommand(setAssigneeCmd)
	configCmd.AddCommand(setUpstreamCmd)
	configCmd.AddCommand(setAccessTokenCmd)
	configCmd.AddCommand(listCmd)
	configCmd.AddCommand(copilotCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configUseCmd)
}

// Config struct with the settings
type Config struct {
	Model        string `json:"model"`
	CopilotToken string `json:"copilot_token"`
}

const modelCopilot = "copilot"
const modelTwinkle = "twinkle"

// configCmd is the root command for configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure settings, Set the default config for pull request, such as assignees. Setup the model to be used or init Copilot token",
}

// configUseCmd sets the model to be used for test case generation
var configUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Set the model to be used for test case generation (copilot or twinkle), by default it is set to twinkle.",
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		// Validate and set model
		model := args[0]
		if model != modelCopilot && model != modelTwinkle {
			log.Errorf("Invalid model: %s, must be one of %v", model, []string{
				modelCopilot,
				modelTwinkle,
			})
			return
		}

		// Load config, set model, and save it
		config, err := loadConfig()
		if err != nil {
			log.Errorf("Failed to load config: %v", err)
			return
		}

		config.Model = model
		if err := saveConfig(config); err != nil {
			log.Errorf("Failed to save config: %v", err)
			return
		}

		fmt.Printf("Config updated: model set to %s\n", model)
		if config.Model == modelCopilot && config.CopilotToken == "" {
			log.Warn("Copilot token is not set. You must init the Copilot token manually before using it. Please run `easypr config copilot init-token` and follow the instructions.")
		}
	},
}

// configShowCmd displays the current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the current configuration settings",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the configuration
		config, err := loadConfig()
		if err != nil {
			log.Errorf("Failed to load config: %v", err)
			return
		}

		// Display the current settings
		log.Println("Current Configuration:")
		log.Printf("\tModel: %s\n", config.Model)
		if config.CopilotToken != "" {
			log.Printf("\tCopilot Token: %s\n", config.CopilotToken)
		} else {
			log.Println("\tCopilot Token: Not set")
		}
	},
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		return ""
	}
	return filepath.Join(homeDir, ".easypr", "config.json")
}

func saveConfig(config *Config) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(getConfigPath())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Serialize the config into JSON
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write the config to the file
	if err := ioutil.WriteFile(getConfigPath(), configData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	return nil
}

func loadConfig() (*Config, error) {
	configPath := getConfigPath()

	// If the file doesn't exist, create a default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &Config{Model: modelTwinkle} // Set default model to "twinkle"
		if err := saveConfig(defaultConfig); err != nil {
			return nil, err
		}
		return defaultConfig, nil
	}

	// Read the config file
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Deserialize JSON into Config struct
	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

var globalConfig *Config

func getGlobalConfig() *Config {
	if globalConfig == nil {
		var err error
		globalConfig, err = loadConfig()
		if err != nil {
			log.Fatal("Failed to load config: %v", err)
		}
	}
	return globalConfig
}
