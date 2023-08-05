package main

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "easypr",
	Short: "A tool for creating pull requests",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	// 获取当前分支
	ref, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}

	currentBranch := ref.Name().Short()
}
