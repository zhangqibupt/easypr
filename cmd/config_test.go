package cmd

import "testing"

func TestIsValidGitHubRepoURL(t *testing.T) {

	var validUrls = []string{
		"https://github.freewheel.tv/data/common.git/",
		"git@github.freewheel.tv:core/common.git/",
	}

	for _, url := range validUrls {
		if !isValidGitHubRepoURL(url) {
			t.Errorf("isValidGitHubRepoURL(%q) = false, want true", url)
		}
	}

	var invalidUrls = []string{
		"ftp://github.com/user/repo",
		"github.com/user/repo",
		"git@example.com:user/repo",
	}

	for _, url := range invalidUrls {
		if isValidGitHubRepoURL(url) {
			t.Errorf("isValidGitHubRepoURL(%q) = true, want false", url)
		}
	}

}
