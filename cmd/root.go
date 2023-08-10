package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhangqibuptse/easypr/git"
)

var (
	rootCmd = &cobra.Command{
		Use:   "easypr",
		Short: "A tool to make Pull Requests and Cherry-picking easy.",
		Long:  `easypr is a tool to create multiple Pull Requests based on current branch to make Cherry-picking easy.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				git.EnableDebug()
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

	completion := completionCommand()
	// mark completion hidden
	completion.Hidden = true

	rootCmd.AddCommand(completion)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.Flags().SortFlags = false

	return rootCmd.Execute()
}
