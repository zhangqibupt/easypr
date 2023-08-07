package main

import "github.com/spf13/cobra"
import "github.com/zhangqibuptse/easypr/cmd"

func main() {
	cmd.CreateRun()(&cobra.Command{}, []string{})
}
