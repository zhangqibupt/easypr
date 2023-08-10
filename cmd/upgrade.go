package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os/exec"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade fwpr to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		performUpgrade()
	},
}

func performUpgrade() {
	cmd := exec.Command("go", "install", "github.freewheel.tv/qzhang/fwpr@latest")
	color.Cyan("Upgrading fwpr to the latest version by 'go install github.freewheel.tv/qzhang/fwpr@latest'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Failed to upgrade due to:%s", output)
		return
	}
	color.Green("Upgrade successfully")
}
