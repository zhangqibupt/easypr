package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os/exec"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade myapp to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		// Run upgrade logic
		performUpgrade()
	},
}

func performUpgrade() {
	cmd := exec.Command("go", "install", "github.com/zhangqibuptse/easypr@latest")
	color.Cyan("Upgrading easypr to the latest version by 'go install github.com/zhangqibuptse/easypr@latest'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("Failed to upgrade due to:%s", output)
		return
	}
	color.Green("Upgrade successfully")
}
