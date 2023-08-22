package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.freewheel.tv/qzhang/fwpr/lib"
	"regexp"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "set the default config for pull request, such as assignees",
}

var setAssigneeCmd = &cobra.Command{
	Use:   "set-assignees [name1 name2...]",
	Short: "set assignees for pull request",
	Run: func(cmd *cobra.Command, args []string) {
		config := lib.Config{
			Assignees: args,
		}
		err := lib.SaveConfig(&config)
		if err != nil {
			color.Red("Failed to save config: %s", err)
			return
		}
		color.Green("Successfully set assignees for current repo")
	},
}

var setUpstreamCmd = &cobra.Command{
	Use:   "set-upstream [repo]",
	Short: "specify the target repo of the pull request, it is useful when the repo is forked and you want to create pull request to the original repo",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			color.Red("Please specify the upstream repo")
			return
		}

		if !isValidGitHubRepoURL(args[0]) {
			color.Red("Please specify a valid repo url")
			return
		}

		config, _ := lib.LoadConfig()
		if config == nil {
			config = &lib.Config{}
		}
		config.Upstream = args[0]

		err := lib.SaveConfig(config)
		if err != nil {
			color.Red("Failed to save config: %s", err)
			return
		}
		color.Green("Successfully set upstream for current repo")
	},
}

func isValidGitHubRepoURL(url string) bool {
	repoRegex := regexp.MustCompile(`((http|git|ssh|http(s)|file|/?)|(git@[\w.]+))(:(//)?)([\w.@:/\-~]+)(\.git)(/)?`)

	return repoRegex.MatchString(url)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list config for pull request",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := lib.LoadConfig()
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
	configCmd.AddCommand(listCmd)
}
