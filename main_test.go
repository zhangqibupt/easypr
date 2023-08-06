package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/zhangqibuptse/easypr/git"
)

func Test(t *testing.T) {
	// 获取当前分支
	// set current
	// set current execution folder
	err := os.Chdir("/Users/qzhang/workspace/freewheel/common/src/go/src/order_service")
	if err != nil {
		log.Fatal(err)
	}
	currentBranch, err := git.CurrentBranch()
	if err != nil {
		log.Fatal(err)
	}

	// 获取远程分支列表
	branches, err := git.RemoteBranches()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(branches)
	fmt.Println(currentBranch)

	// use cobra to let user select target branch
	selectedBranches := []string{}

	err := createCmd.Flags().StringSliceVarP(&selectedBranches, "branches", "b", nil, "Target branches")

	if err != nil {
		return err
	}

	if len(selectedBranches) == 0 {

		fmt.Println("Select target branches:")

		for i, branch := range branches {
			fmt.Printf("%d: %s\n", i+1, branch)
		}

		// 交互选择分支
		selectedBranches, err = interact.MultiSelect(branches)

		if err != nil {
			return err
		}

	}

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

}
