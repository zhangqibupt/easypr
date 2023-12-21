package lib

import (
	"fmt"
	"net/url"
	"os/exec"
	"sort"
	"strings"

	"github.com/fatih/color"
)

var debug = false

func execGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if debug {
		fmt.Printf("[git] git %s\n", strings.Join(args, " "))
		fmt.Printf("[git] output: %s\n", output)
	}

	return string(output), err
}

func CurrentBranch() (string, error) {
	color.Cyan("Getting current branch")
	out, err := execGit("symbolic-ref", "--short", "HEAD")
	if err != nil {
		color.Red("Get current branch error: %s", out)
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func RemoteBranches() ([]string, error) {
	out, err := execGit("branch", "-r")
	if err != nil {
		color.Red("Failed to fetch remote branches due to:%s", out)
		return nil, err
	}
	var branches []string

	lines := strings.Split(out, "\n")

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "origin/") {
			branches = append(branches, strings.TrimSpace(line))
		}
	}
	var priorityBranches, otherBranches, topBranches []string
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

func OtherLocalBranches() ([]string, error) {
	out, err := execGit("branch")
	if err != nil {
		color.Red("Failed to fetch local branches due to:%s", out)
		return nil, err
	}
	var branches []string

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "* ") {
			continue
		}

		branches = append(branches, strings.TrimSpace(line))
	}
	return branches, nil
}

func HasUncommittedChanges() (bool, error) {
	out, err := execGit("status", "-s")
	if err != nil {
		color.Red("Failed to check git status due to:%s", out)
		return false, err
	}

	if strings.TrimSpace(out) != "" {
		return true, nil
	}
	return false, nil
}

func RemoteURL(remote string) (string, error) {
	out, err := execGit("remote", "get-url", remote)
	if err != nil {
		color.Red("Get remote url error: %s", out)
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func CreatePRLink(source, target string) (string, error) {
	color.Cyan("Fetching %s from remote", target)
	out, err := execGit("fetch", "origin")
	if err != nil {
		color.Red("Fetch origin error: %s", out)
		return "", err
	}

	color.Cyan("Rebasing %s to %s", source, target)
	out, err = execGit("rebase", target)
	if err != nil {
		color.Red("Rebase error: %s", out)
		return "", err
	}

	return generatePRLink(source, target, "")
}

func CreateCherryPickPRLink(source, target string, commits []string, originalPRLink string) (string, error) {
	newBranchName := generateNewBranchName(source, target)

	if exist := branchExists(newBranchName); exist {
		color.Cyan("Branch %s already exists, trying to recreate it", newBranchName)
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
			color.Red("cherry pick %s error: %s", commit, out)
			_, _ = execGit("cherry-pick", "--abort")
			return "", err
		}
	}

	return generatePRLink(newBranchName, target, originalPRLink)
}

func RecreateCPBranch(localCPBranch string, commits []string) error {
	cpRemoteBaseBranch := fmt.Sprintf("origin/%s", strings.Split(localCPBranch, CherryPickPlaceholder)[1])

	color.Cyan("Recreating branch %s based on %s", localCPBranch, cpRemoteBaseBranch)
	if err := DeleteBranch(localCPBranch); err != nil {
		return err
	}

	if err := CreateBranch(localCPBranch, cpRemoteBaseBranch); err != nil {
		return err
	}

	for _, commit := range commits {
		color.Cyan("Cherry picking %s to branch %s", commit, localCPBranch)
		out, err := execGit("cherry-pick", commit)
		if err != nil {
			color.Red("cherry pick %s error: %s", commit, out)
			_, _ = execGit("cherry-pick", "--abort")
			return err
		}
	}

	return ForcePush(localCPBranch)
}

func ForcePush(branch string) error {
	color.Cyan("Force Pushing %s to remote", branch)
	out, err := execGit("push", "origin", "-f", fmt.Sprintf("%s:%s", branch, branch))
	if err != nil {
		color.Red("Force Push %s to remote error: %s", branch, out)
		return err
	}
	return nil
}

func generatePRLink(source, target string, orignalPRLink string) (string, error) {
	isCherryPick := len(orignalPRLink) > 0

	color.Cyan("Pushing %s to remote", source)
	if err := ForcePush(source); err != nil {
		return "", err
	}

	originalRepo, err := RemoteURL("origin")
	if err != nil {
		return "", err
	}

	targetRepo := originalRepo
	forked := false
	config, _ := LoadRepoConfig()
	if config != nil && config.Upstream != "" {
		forked = true
		targetRepo = config.Upstream
	}

	color.Cyan("Generating pull request link")

	if strings.HasPrefix(targetRepo, "git@") {
		targetRepo = strings.Replace(targetRepo, ":", "/", 1)
		targetRepo = strings.Replace(targetRepo, "git@", "https://", 1)
	}
	baseURL := strings.Replace(targetRepo, ".git", "", 1)

	labels := []string{generateLabel(target)}
	if isCherryPick {
		labels = append(labels, "cherry-pick")
	}

	if forked {
		space := extractRepoNameFromURL(originalRepo)
		source = fmt.Sprintf("%s:%s", space, source)
	}

	fullURL := fmt.Sprintf("%s/compare/%s...%s?quick_pull=1&labels=%s", baseURL, shortBranchName(target), source, strings.Join(labels, ","))
	if isCherryPick {
		description := url.QueryEscape(fmt.Sprintf("Auto Generated by [fwpr](https://github.freewheel.tv/qzhang/fwpr)\nOrignal PR [Link](%s)\n", orignalPRLink))
		fullURL = fmt.Sprintf("%s&title=%s&body=%s", fullURL, url.QueryEscape(source+"#CP"), description)
	}

	fullURL = fillInAssignees(fullURL)

	return fullURL, nil
}

func extractRepoNameFromURL(url string) string {
	// Remove ".git" extension
	url = strings.TrimSuffix(url, ".git")
	// Split the URL by slashes
	parts := strings.Split(url, "/")
	// The repository name is the last part after the last slash
	repoName := parts[len(parts)-2]
	return repoName
}

func fillInAssignees(fullURL string) string {
	c, _ := LoadRepoConfig()
	if c == nil || len(c.Assignees) == 0 {
		return fullURL
	}

	return fmt.Sprintf("%s&assignees=%s", fullURL, strings.Join(c.Assignees, ","))
}

func generateLabel(target string) string {
	branchName := shortBranchName(target)
	if strings.HasPrefix(branchName, "V_") {
		label := strings.Replace(branchName, "V_", "", 1)
		label = strings.ReplaceAll(label, "_", ".")
		return label
	}
	return branchName
}

func CommitsBetween(source, target string) ([]string, error) {
	output, err := execGit("log", fmt.Sprintf("%s..%s", target, source), "--oneline")
	if err != nil {
		color.Red("Get commits between %s and %s error: %s", source, target, output)
		return nil, err
	}

	var commits []string
	for _, line := range strings.Split(output, "\n") {
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

func shortBranchName(branch string) string {
	return strings.Replace(branch, "origin/", "", 1)
}

func generateNewBranchName(source, target string) string {
	prefix := CPBranchPrefix(source)
	branchName := fmt.Sprintf("%s%s", prefix, shortBranchName(target))
	if len(branchName) > 200 {
		branchName = branchName[:200]
	}
	return branchName
}

var CherryPickPlaceholder = "_to_"

func CPBranchPrefix(source string) string {
	return fmt.Sprintf("%s%s", source, CherryPickPlaceholder)
}

func Checkout(branch string) error {
	out, err := execGit("checkout", shortBranchName(branch))
	if err != nil {
		color.Red("checkout %s error: %s", branch, out)
		return err
	}
	return nil
}

func DeleteBranch(branch string) error {
	out, err := execGit("branch", "-D", branch)
	if err != nil {
		color.Red("failed to delete branch %s error: %s", branch, out)
		return err
	}

	return nil
}

func branchExists(branch string) bool {
	_, err := execGit("rev-parse", "--verify", branch)
	if err != nil {
		return false
	}
	return true
}

func CreateBranch(newBranch, baseBranch string) error {
	out, err := execGit("checkout", "-b", newBranch, baseBranch)
	if err != nil {
		color.Red("failed to create branch %s error: %s", newBranch, out)
		return err
	}
	return nil
}

func EnableDebug() {
	debug = true
	color.Yellow("Debug mode enabled")
}

func TopLevelDirectory() (string, error) {
	director, err := execGit("rev-parse", "--show-toplevel")
	if err != nil {
		color.Red("failed to get git repo information: %s", director)
		return "", err
	}
	return strings.TrimSpace(director), nil

}
