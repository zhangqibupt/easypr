package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display fwpr version",
	Run: func(cmd *cobra.Command, args []string) {
		displayVersion()
	},
}

func displayVersion() {
	color.Cyan("fwpr version: %s", "v1.0.3")
}
