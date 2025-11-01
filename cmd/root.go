package cmd

import (
	"easypr/internal/logger"
	"easypr/lib"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "easypr",
		Short: "A tool to make Pull Requests and Cherry-picking easy.",
		Long: `easypr is a tool to create multiple Pull Requests based on current branch to make Cherry-picking easy.

Get more details from https://easypr
Join our slack channel to submit issues and get supported #kof-user-group https://freewheel.slack.com/archives/C05M8MQNMQQ
`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				lib.EnableDebug()
			}
		},
	}
	log = logger.GetLogger() // Global logger
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
}

// Execute executes the root command.
func Execute() error {
	cobra.EnableCommandSorting = false
	rootCmd.Flags().SortFlags = false

	completion := completionCommand()
	completion.Hidden = true

	rootCmd.AddCommand(completion)

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(reviewCmd)

	// add command to show PR links

	return rootCmd.Execute()
}
