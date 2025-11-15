package cmd

import "github.com/spf13/cobra"

func completionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "Generate the autocompletion script for the specified shell",
	}
}
