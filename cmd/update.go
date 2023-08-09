package cmd

import (
	"github.com/spf13/cobra"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/fatih/color"

	"github.com/zhangqibuptse/easypr/git"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync new commits to cherry-pick branches.",
	Long:  "Sync new commits to cherry-pick branches. \nIt is used when you have created cherry-pick Pull Requests through 'create' command, you made some new commits and you want to sync these commits to these cherry-pick branches.",
	Run:   SyncRun(),
}

func SyncRun() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if debug {
			git.EnableDebug()
			color.Yellow("Debug mode enabled")
		}

		// Check for uncommitted changes
		color.Cyan("Checking for uncommitted changes")
		hasChanges, err := git.HasUncommittedChanges()
		if err != nil {
			return
		}
		if hasChanges {
			color.Red("Uncommitted changes found. Please commit or stash before continuing.")
			return
		}

		// Fetch remote branch list
		color.Cyan("Fetching remote branch list")
		originURL, err := git.RemoteURL("origin")
		if err != nil {
			return
		}
		if strings.TrimSpace(originURL) == "" {
			color.Red("Failed to fetch remote origin. Check if 'origin' is set using 'git remote -v'.")
			return
		}

		// Fetch remote branches
		branches, err := git.RemoteBranches()
		if err != nil {
			return
		}

		if len(branches) == 0 {
			color.Red("No remote branches found!")
			return
		}

		// Select target branch
		_, targetBranch, err := (&promptui.Select{
			//Label:        "Which remote branch you want to merge current branch into? It should be the target branch of the current branch's Pull Request. Usually, this should be master, main, or a release branch.",
			Label:        "Which remote branch you want to merge current branch into?",
			Items:        branches,
			Searcher:     searcher(branches),
			Size:         10,
			HideSelected: true}).Run()
		if err != nil {
			color.Red("Failed to select target branch: %s", err)
			return
		}

		sourceBranch, err := git.CurrentBranch()
		if err != nil {
			return
		}
		defer func() {
			_ = git.Checkout(sourceBranch)
		}()

		commits, err := git.CommitsBetween(sourceBranch, targetBranch)
		if err != nil {
			return
		}

		if len(commits) == 0 {
			color.Yellow("No difference found between %s and %s, exiting...", sourceBranch, targetBranch)
			return
		}

		localBranches, err := git.OtherLocalBranches()
		if err != nil {
			return
		}

		// Find the cherry-pick branches
		var cherryPickBranches []string
		var cherryPickBranchPrefix = git.CPBranchPrefix(sourceBranch)
		color.Cyan("Checking for local cherry-pick branches with prefix '%s'", cherryPickBranchPrefix)
		for _, branch := range localBranches {
			if strings.HasPrefix(branch, cherryPickBranchPrefix) {
				cherryPickBranches = append(cherryPickBranches, branch)
			}
		}

		if len(cherryPickBranches) == 0 {
			color.Yellow("No local cherry-pick branches found, make sure you have created cherry-pick Pull Requests through 'easypr create' command, exiting...")
			return
		}

		// Select cherry-pick branches
		var allItems []*MultipleSelectItem
		for _, item := range cherryPickBranches {
			allItems = append(allItems, &MultipleSelectItem{ID: item, IsSelected: true})
		}

		selected, err := multipleSelect(0, allItems, "? Select local branches you want to sync to:")
		if err != nil {
			return
		}

		if len(selected) == 0 {
			color.Yellow("No local branches selected, exiting...")
			return
		}

		for _, target := range selected {
			color.Yellow("Syncing new commits to %s...", target)
			err := git.RecreateCPBranch(target, commits)
			if err != nil {
				color.Red("Failed to create cherry-pick Pull Request for branch %s, skipping...", target)
				continue
			}
			defer color.Green("%s has been synced successfully", target)
		}
	}
}

func init() {
	syncCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Print debug messages")
	rootCmd.AddCommand(syncCmd)
}
