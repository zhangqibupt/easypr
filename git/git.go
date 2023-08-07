package git

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// 获取当前分支
func CurrentBranch() (string, error) {
	color.Cyan("getting current branch")
	out, err := execGit("symbolic-ref", "--short", "HEAD")
	if err != nil {
		color.Red("get current branch error: %s, %s", out, err)
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// 获取远程分支列表
func RemoteBranches() ([]string, error) {
	out, err := execGit("branch", "-r")
	if err != nil {
		return nil, err
	}
	branches := []string{}

	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		// 只包含origin/开头的分支
		if strings.HasPrefix(strings.TrimSpace(line), "origin/") {
			branches = append(branches, strings.TrimSpace(line))
		}
	}
	// 构造优先级分支slice
	priorityBranches := []string{}
	// 构造非优先级分支slice
	otherBranches := []string{}
	// primary branch
	topBranches := []string{}

	for _, branch := range branches {
		if strings.HasPrefix(branch, "origin/V_") {
			priorityBranches = append(priorityBranches, branch)
		} else if branch == "origin/master" || branch == "origin/main" {
			topBranches = append(topBranches, branch)
		} else {
			otherBranches = append(otherBranches, branch)
		}
	}

	// 对优先级分支排序
	sort.Slice(priorityBranches, func(i, j int) bool {
		return strings.Compare(priorityBranches[i], priorityBranches[j]) > 0
	})

	// 合并
	branches = append(topBranches, priorityBranches...)
	branches = append(branches, otherBranches...)

	return branches, nil

}

// 执行git命令
func execGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// 检查是否有未提交的修改
func HasUncommittedChanges() (bool, string) {
	out, _ := execGit("status", "-s")
	if out != "" {
		return true, out
	}
	return false, ""
}

// stash 未提交的修改
func Stash() error {
	_, err := execGit("stash")
	return err
}

func RemoteURL(remote string) (string, error) {
	out, err := execGit("remote", "get-url", remote)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func CreatePRLink(source, target string) (string, error) {
	color.Cyan("fetching %s from remote", target)
	out, err := execGit("fetch", "origin")
	if err != nil {
		color.Red("fetch origin error: %s %s", out, err)
		return "", err
	}

	color.Cyan("rebasing to %s", target)
	out, err = execGit("rebase", target)
	if err != nil {
		color.Red("rebase error: %s %s", out, err)
		return "", err
	}

	return generatePRLink(source, target)
}

func CreateCherryPickPRLink(source, target string, commits []string) (string, error) {
	newBranchName := generateNewBranchName(source, target)

	color.Cyan("creating branch [%s] based on [%s]", source, target)
	out, err := execGit("checkout", "-b", newBranchName, target)
	if err != nil {
		color.Red("checkout -b %s error: %s %s", newBranchName, out, err)
		return "", err
	}
	color.Cyan("cherry picking %d commits to %s", len(commits), newBranchName)
	for _, commit := range commits {
		out, err = execGit("cherry-pick", commit)
		if err != nil {
			color.Red("cherry pick %s error: %s %s", commit, out, err)
			return "", err
		}
	}

	return generatePRLink(source, target)
}

func generatePRLink(source, target string) (string, error) {
	// push to remote
	color.Cyan("pushing %s to remote", source)
	out, err := execGit("push", "origin", fmt.Sprintf("%s:%s", source, source))
	if err != nil {
		color.Red("push %s to remote error: %s %s", source, out, err)
		return "", err
	}

	color.Cyan("generating pull request link")
	out, err = RemoteURL("origin")
	if err != nil {
		color.Red("get remote url error: %s, %s", out, err)
		return "", err
	}

	if strings.HasPrefix(out, "git@") {
		out = strings.Replace(out, "git@", "https://", 1)
		out = strings.Replace(out, ":", "/", 1)
	}
	out = strings.Replace(out, ".git", "", 1)
	return fmt.Sprintf("%s/compare/%s...%s", out, shortName(target), source), nil
}

func CommitsBetween(source, target string) ([]string, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("%s..%s", target, source), "--oneline")
	output, err := cmd.Output()
	if err != nil {
		color.Red("get commits between %s and %s error: %s %s", source, target, output, err)
		return nil, err
	}
	color.Red("get commits between %s and %s output: %s", source, target, output)

	// 解析输出结果
	commits := strings.Split(string(output), "\n")
	sortedCommits := make([]string, len(commits))
	for i, commit := range commits {
		sortedCommits[len(commits)-1-i] = strings.Fields(commit)[0]
	}
	color.Cyan("found %d commits between %s and %s", len(sortedCommits), source, target)
	color.Cyan("commits: %v", sortedCommits)

	return sortedCommits, nil
}

func shortName(branch string) string {
	return strings.Replace(branch, "origin/", "", 1)
}

func generateNewBranchName(source, target string) string {
	baseBranchName := fmt.Sprintf("%s_to_%s", source, shortName(target))
	branchName := baseBranchName
	for suffix := 1; ; suffix++ {
		_, err := execGit("rev-parse", "--verify", branchName)
		if err != nil {
			return branchName
		}
		branchName = fmt.Sprintf("%s_%d", baseBranchName, suffix)
	}
}
