package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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

	repo, err := findGitRepo("/Users/qzhang/workspace/freewheel/common/src/go/src/order_service")
	if err != nil {
		log.Fatal(err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		log.Fatal(err)
	}

	//var publicKey *ssh.PublicKeys
	sshPath := os.Getenv("HOME") + "/.ssh/id_ed25519"
	//sshKey, _ := ioutil.ReadFile(sshPath)
	//publicKey, keyError := ssh.NewPublicKeys("git", []byte(sshKey), "")
	//if keyError != nil {
	//	fmt.Println(keyError)
	//}

	keys, err := ssh.NewPublicKeysFromFile("git", sshPath, "")
	if err != nil {
		fmt.Println(err)
	}

	remoteBranchNames := []string{}
	for _, remote := range remotes {
		refs, err := remote.List(&git.ListOptions{
			Auth: keys,
		})
		if err != nil {
			log.Fatal(err)
		}

		branches := getBranches(refs)

		for _, b := range branches {
			remoteBranchNames = append(remoteBranchNames, b.Name().Short())
		}
	}

	// 显示分支列表给用户选择
	fmt.Println("Select target branches:")
	for i, name := range remoteBranchNames {
		fmt.Printf("%d: %s\n", i+1, name)
	}

	//// 用户输入选择
	//var selections []int
	//fmt.Print("Enter your selections (separated by commas): ")
	//fmt.Scanln(&selections)
	//
	//// 获取选择的分支名
	//targetBranches := []string{}
	//for _, sel := range selections {
	//	targetBranches = append(targetBranches, remoteBranchNames[sel-1])
	//}
	//print(targetBranches)

}

func findGitRepo(cwd string) (repo *git.Repository, err error) {
	for {
		repo, err := git.PlainOpen(cwd)
		if err == nil {
			return repo, nil
		}

		// 到达根目录还未找到,返回错误
		if cwd == "/" {
			return nil, err
		}

		// 向上递归一个目录
		cwd = filepath.Dir(cwd)
	}
}

func getBranches(refs []*plumbing.Reference) []*plumbing.Reference {
	var branches []*plumbing.Reference

	for _, ref := range refs {
		if strings.HasPrefix(ref.Name().String(), "refs/heads/") {
			branches = append(branches, ref)
		}
	}

	return branches
}
