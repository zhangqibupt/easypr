package git

import (
	"os/exec"
	"sort"
	"strings"
)

// 获取当前分支
func CurrentBranch() (string, error) {
	out, err := execGit("symbolic-ref", "--short", "HEAD")
	if err != nil {
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
			topBranches = append(otherBranches, branch)
		} else {
			otherBranches = append(otherBranches, branch)
		}
	}

	// 对优先级分支排序
	sort.Slice(priorityBranches, func(i, j int) bool {
		return strings.Compare(priorityBranches[i], priorityBranches[j]) < 0
	})

	// 合并
	branches = append(topBranches, priorityBranches...)
	branches = append(branches, otherBranches...)

	return branches, nil

}

// 检查目标分支是否包含提交
func CherryTargetBranch(source, target string) ([]string, error) {
	out, err := execGit("cherry", target, source)
	if err != nil {
		return nil, err
	}

	var commits []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+") {
			commits = append(commits, strings.TrimSpace(line[1:]))
		}
	}

	return commits, nil

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
	out, _ := execGit("status", "--porcelain")
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
