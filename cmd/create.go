package cmd

import (
	"os"
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
	Run:   CreateRun(),
}

func CreateRun() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		os.Chdir("/Users/qzhang/workspace/github/test_easypr")
		color.Cyan("Checking uncommitted changes...")
		if has, _ := git.HasUncommittedChanges(); has {
			color.Red("There are uncommitted changes, please stash them first.")
			return
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

		i, primaryBranch, err := (&promptui.Select{Label: "Select target branch", Items: branches, Searcher: searcher(branches), Size: 10}).Run()
		if err != nil {
			color.Red("Failed to select current branch due to error: %s", err)
			return
		}

		source, err := git.CurrentBranch()
		if err != nil {
			return
		}
		defer git.Checkout(source)

		commits, err := git.CommitsBetween(source, primaryBranch)
		if err != nil {
			return
		}

		if len(commits) == 0 {
			color.Yellow("No commits between %s and %s, existing", source, primaryBranch)
			return
		}

		branches = append(branches[:i], branches[i+1:]...)
		var allItems []*MultipleSelectItem
		for _, item := range branches {
			allItems = append(allItems, &MultipleSelectItem{
				ID: item,
			})
		}

		selected, err := multipleSelect(0, allItems)
		if err != nil {
			color.Red("Failed to select cherry-pick branches due to error: %s", err)
			return
		}

		color.Cyan("Creating PR for primary branch %s...", primaryBranch)
		link, err := git.CreatePRLink(source, primaryBranch)
		if err != nil {
			return
		}
		color.Green("Primary PR link: %s", link)

		if len(selected) == 0 {
			return
		}

		for _, target := range selected {
			color.Cyan("Creating PR for branch %s...", target)
			link, err := git.CreateCherryPickPRLink(source, target, commits)
			if err != nil {
				return
			}
			color.Green("PR link to %s: %s", target, link)
		}
	}
}

func searcher(branches []string) func(input string, index int) bool {
	return func(input string, index int) bool {
		branch := branches[index]
		return fuzzy.Match(strings.ToLower(input), strings.ToLower(branch))
	}
}

var template = &promptui.SelectTemplates{
	Label:    `? Select cherry-pick branch:`,
	Active:   `→ {{if .IsSelected}}{{ .ID | green }} {{ "✔" | green }} {{else}}{{ .ID }}{{end}}`,
	Inactive: `  {{if .IsSelected}}{{ .ID | green }} {{ "✔" | green }} {{else}}{{ .ID }}{{end}}`,
}

type MultipleSelectItem struct {
	ID         string
	IsSelected bool
}

// multipleSelect() prompts user to select one or more items in the given slice
func multipleSelect(selectedPos int, items []*MultipleSelectItem) ([]string, error) {
	// Always prepend a "Done" item to the slice if it doesn't
	// already exist.
	const doneID = "Done"
	if len(items) > 0 && items[0].ID != doneID {
		doneItem := &MultipleSelectItem{
			ID: doneID,
		}
		items = append([]*MultipleSelectItem{doneItem}, items...)
	}

	prompt := promptui.Select{
		Label:        "? Select cherry-pick branch:",
		Items:        items,
		Templates:    template,
		Size:         10,
		CursorPos:    selectedPos,
		HideSelected: true,
		Searcher: func(input string, index int) bool {
			item := items[index]
			if item.ID == doneID {
				return true
			}
			return fuzzy.Match(strings.ToLower(input), strings.ToLower(item.ID))
		},
	}

	selectionIdx, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	chosenItem := items[selectionIdx]

	if chosenItem.ID != doneID {
		chosenItem.IsSelected = !chosenItem.IsSelected
		return multipleSelect(selectionIdx, items)
	}

	// If the user selected the "Done" item, return
	// all selected items.
	var selectedItems []string
	for _, i := range items {
		if i.IsSelected {
			selectedItems = append(selectedItems, i.ID)
		}
	}
	return selectedItems, nil
}

func init() {
	rootCmd.AddCommand(createCmd)
}
