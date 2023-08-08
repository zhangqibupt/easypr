package git

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/fatih/color"
)

var debug = false

func execGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()

	if debug {
		fmt.Printf("[git] git %s\n", strings.Join(args, " "))
		fmt.Printf("[git] output: %s %s", output, err)
	}

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func CurrentBranch() (string, error) {
	color.Cyan("Getting current branch")
	out, err := execGit("symbolic-ref", "--short", "HEAD")
	if err != nil {
		color.Red("Get current branch error: %s, %s", out, err)
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func RemoteBranches() ([]string, error) {
	out, err := execGit("branch", "-r")
	if err != nil {
		color.Red("Failed to fetch remote branches due to:%s, %s", out, err)
		return nil, err
	}
	branches := []string{}

	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "origin/") {
			branches = append(branches, strings.TrimSpace(line))
		}
	}
	priorityBranches := []string{}
	otherBranches := []string{}
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

	sort.Slice(priorityBranches, func(i, j int) bool {
		return strings.Compare(priorityBranches[i], priorityBranches[j]) > 0
	})

	branches = append(topBranches, priorityBranches...)
	branches = append(branches, otherBranches...)

	return branches, nil
}

func HasUncommittedChanges() (bool, error) {
	out, err := execGit("status", "-s")
	if err != nil {
		color.Red("Failed to check git status due to:%s, %s", out, err)
		return false, err
	}

	if strings.TrimSpace(out) != "" {
		return true, nil
	}
	return false, nil
}

// stash 未提交的修改
func Stash() error {
	_, err := execGit("stash")
	return err
}

func RemoteURL(remote string) (string, error) {
	out, err := execGit("remote", "get-url", remote)
	if err != nil {
		color.Red("Get remote url error: %s, %s", out, err)
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func CreatePRLink(source, target string) (string, error) {
	color.Cyan("Fetching %s from remote", target)
	out, err := execGit("fetch", "origin")
	if err != nil {
		color.Red("Fetch origin error: %s %s", out, err)
		return "", err
	}

	color.Cyan("Rebasing %s to %s", source, target)
	out, err = execGit("rebase", target)
	if err != nil {
		color.Red("Rebase error: %s %s", out, err)
		return "", err
	}

	return generatePRLink(source, target)
}

func CreateCherryPickPRLink(source, target string, commits []string) (string, error) {
	newBranchName := generateNewBranchName(source, target)

	if exist := branchExists(newBranchName); exist {
		color.Yellow("Branch %s already exists, trying to recreate it", newBranchName)
		if err := DeleteBranch(newBranchName); err != nil {
			return "", err
		}
	}

	color.Cyan("Creating branch %s based on %s", newBranchName, target)
	if err := CreateBranch(newBranchName, target); err != nil {
		return "", err
	}

	for _, commit := range commits {
		color.Cyan("Cherry picking %s to branch %s", commit, newBranchName)
		out, err := execGit("cherry-pick", commit)
		if err != nil {
			color.Red("cherry pick %s error: %s %s", commit, out, err)
			execGit("cherry-pick", "--abort")
			return "", err
		}
	}

	return generatePRLink(newBranchName, target)
}

func generatePRLink(source, target string) (string, error) {
	color.Cyan("pushing %s to remote", source)
	out, err := execGit("push", "origin", "-f", fmt.Sprintf("%s:%s", source, source))
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

	var commits []string
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		commits = append(commits, strings.Fields(line)[0])
	}
	sortedCommits := make([]string, len(commits))
	for i, commit := range commits {
		sortedCommits[len(commits)-1-i] = commit
	}
	color.Cyan("Found %d commits between %s and %s", len(sortedCommits), source, target)

	return sortedCommits, nil
}

func shortName(branch string) string {
	return strings.Replace(branch, "origin/", "", 1)
}

func generateNewBranchName(source, target string) string {
	branchName := fmt.Sprintf("%s_to_%s", source, shortName(target))
	if len(branchName) > 200 {
		branchName = branchName[:200]
	}
	return branchName
}

func Checkout(branch string) error {
	out, err := execGit("checkout", shortName(branch))
	if err != nil {
		color.Red("checkout %s error: %s %s", branch, out, err)
		return err
	}
	return nil
}

func DeleteBranch(branch string) error {
	out, err := execGit("branch", "-D", branch)
	if err != nil {
		color.Red("failed to delete branch %s error: %s %s", branch, out, err)
		return err
	}

	return nil
}

func branchExists(branch string) bool {
	out, _ := execGit("rev-parse", "--verify", branch)
	return strings.TrimSpace(out) != ""
}

func CreateBranch(newBranch, baseBranch string) error {
	out, err := execGit("checkout", "-b", newBranch, baseBranch)
	if err != nil {
		color.Red("failed to create branch %s error: %s %s", newBranch, out, err)
		return err
	}
	return nil
}

func EnableDebug() {
	debug = true
}
