package cmd

import (
	"github.com/spf13/cobra"
	"github.freewheel.tv/qzhang/fwpr/lib"
)

var (
	rootCmd = &cobra.Command{
		Use:   "fwpr",
		Short: "A tool to make Pull Requests and Cherry-picking easy.",
		Long: `fwpr is a tool to create multiple Pull Requests based on current branch to make Cherry-picking easy.

Get more details from https://github.freewheel.tv/qzhang/fwpr
Join our slack channel to submit issues and get supported #fw-useful-tools https://freewheel.slack.com/archives/C05NA6TPM2R
`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				lib.EnableDebug()
			}
		},
	}
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

	return rootCmd.Execute()
}
