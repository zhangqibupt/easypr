package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.freewheel.tv/qzhang/fwpr/lib"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "set the default config for pull request, such as assignees",
}

var assigneeCmd = &cobra.Command{
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

		color.Cyan("default assignees for current repo: %s", config.Assignees)
	},
}

func init() {
	configCmd.AddCommand(assigneeCmd)
	configCmd.AddCommand(listCmd)
}
