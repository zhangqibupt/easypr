package git

import (
	"os"
	"testing"
)

func TestCommitsBetween(t *testing.T) {
	os.Chdir("/Users/qzhang/workspace/github/test_easypr")
	got, err := CommitsBetween("feature1", "master")
	t.Logf("CommitsBetween() got = %v", got)
	if err != nil {
		t.Errorf("CommitsBetween() error = %v", err)
		return
	}
}
