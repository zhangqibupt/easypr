package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/manifoldco/promptui"

	"github.com/fatih/color"

	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/zhangqibuptse/easypr/git"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate PRs across branches",
	Run: func(cmd *cobra.Command, args []string) {
		//os.Chdir("/Users/qzhang/workspace/freewheel/common/src/go/src/order_service")
		color.Cyan("Checking uncommitted changes...")
		if has, _ := git.HasUncommittedChanges(); has {
			color.Yellow("There are uncommitted changes, stashing them...")
			err := git.Stash()
			if err != nil {
				color.Red("Failed to stash uncommitted changes due to error: %s", err)
				return
			}
			// TODO should pop stash after all
		}

		color.Cyan("Fetching remote branche list...")
		// check if origin is set
		_, err := git.RemoteURL("origin")
		if err != nil {
			color.Red("Failed to get remote origin, run `git remote -v` to check if the origin has been set.")
			return
		}
		branches, err := git.RemoteBranches()
		if err != nil {
			color.Red("Failed to get remote branches due to error: %s", err)
			return
		}

		i, currBranch, err := (&promptui.Select{
			Label: "Select current branch",
			Items: branches,
			Searcher: func(input string, index int) bool {
				branch := branches[index]
				return fuzzy.Match(strings.ToLower(input), strings.ToLower(branch))
			},
			Size: 20,
		}).Run()

		if err != nil {
			color.Red("Failed to select current branch due to error: %s", err)
		}
		println(i, currBranch)
		// 选择其他分支
		//targets := promptui.MultiSelect{
		//	Label: "Select target branches",
		//	Items: branches,
		//}.Run()
		// 过滤出目标分支
		//var targetBranches []string
		// ...过滤逻辑

		// 对每个目标分支
		//for _, target := range targetBranches {
		//
		//	// 检查目标分支是否已包含部分提交
		//	existedCommits, err := git.CherryTargetBranch(currentBranch, target)
		//
		//	// 对当前分支进行cherry-pick
		//	err = git.CherryPick(currentBranch, target, existedCommits)
		//	if err != nil {
		//		// 处理错误
		//	}
		//
		//	// 提交并推送目标分支
		//	err = git.Push(target)
		//	if err != nil {
		//		// 处理错误
		//	}
		//
		//	// 创建PR
		//	pr, err := github.CreatePR(currentBranch, target)
		//	if err != nil {
		//		// 处理错误
		//	}
		//
		//	// 显示创建的PR
		//	fmt.Println("PR created:", pr.URL)
		//}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
