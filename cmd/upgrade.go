package cmd

import (
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Short:   "Upgrade fwpr to the latest version",
	Aliases: []string{"u"},
	Run: func(cmd *cobra.Command, args []string) {
		performUpgrade()
	},
}

func performUpgrade() {
	cmd := exec.Command("go", "install", "fwpr@latest")
	color.Cyan("Upgrading fwpr to the latest version by 'go install fwpr@latest'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Failed to upgrade due to:%s", output)
		return
	}
	color.Green("Upgrade successfully")
}
