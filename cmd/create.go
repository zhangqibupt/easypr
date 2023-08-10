package cmd

import (
	"github.com/spf13/cobra"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/fatih/color"

	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.freewheel.tv/qzhang/fwpr/git"
)

var debug = false
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate PRs across branches",
	Run:   CreateRun(),
}

func CreateRun() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// Check for uncommitted changes
		color.Cyan("Checking for uncommitted changes")
		hasChanges, err := git.HasUncommittedChanges()
		if err != nil {
			return
		}
		if hasChanges {
			color.Red("Uncommitted changes found. Please commit or stash before continuing.")
			return
		}

		// Fetch remote branch list
		color.Cyan("Fetching remote branch list")
		originURL, err := git.RemoteURL("origin")
		if err != nil {
			return
		}
		if strings.TrimSpace(originURL) == "" {
			color.Red("Failed to fetch remote origin. Check if 'origin' is set using 'git remote -v'.")
			return
		}

		// Fetch remote branches
		branches, err := git.RemoteBranches()
		if err != nil {
			return
		}

		if len(branches) == 0 {
			color.Red("No remote branches found!")
			return
		}

		// Select target branch
		i, targetBranch, err := (&promptui.Select{Label: "Select target branch", Items: branches, Searcher: searcher(branches), Size: 10, HideSelected: true}).Run()
		if err != nil {
			color.Red("Failed to select target branch: %s", err)
			return
		}

		sourceBranch, err := git.CurrentBranch()
		if err != nil {
			return
		}
		defer func() {
			_ = git.Checkout(sourceBranch)
		}()

		commits, err := git.CommitsBetween(sourceBranch, targetBranch)
		if err != nil {
			return
		}

		if len(commits) == 0 {
			color.Yellow("No difference found between %s and %s, exiting...", sourceBranch, targetBranch)
			return
		}

		// Select cherry-pick branches
		branches = append(branches[:i], branches[i+1:]...)
		var allItems []*MultipleSelectItem
		for _, item := range branches {
			allItems = append(allItems, &MultipleSelectItem{
				ID: item,
			})
		}

		selected, err := multipleSelect(1, allItems, "? Select remote target branches for cherry-picking:")
		if err != nil {
			return
		}

		color.Yellow("Creating Pull Request to %s...", targetBranch)
		targetLink, err := git.CreatePRLink(sourceBranch, targetBranch)
		if err != nil {
			return
		}

		if len(selected) > 0 {
			for _, target := range selected {
				color.Yellow("Creating cherry-pick Pull Request to %s...", target)
				link, err := git.CreateCherryPickPRLink(sourceBranch, target, commits)
				if err != nil {
					color.Red("Failed to create cherry-pick Pull Request for branch %s, skipping...", target)
					continue
				}
				defer color.Cyan("Cherry-Pick PR to %s: %s", target, color.GreenString(link))
			}
		}
		color.Cyan("PR to %s: %s", targetBranch, color.GreenString(targetLink))
	}
}

func searcher(branches []string) func(input string, index int) bool {
	return func(input string, index int) bool {
		branch := branches[index]
		return fuzzy.Match(strings.ToLower(input), strings.ToLower(branch))
	}
}

type MultipleSelectItem struct {
	ID         string
	IsSelected bool
}

// multipleSelect() prompts user to select one or more items in the given slice
func multipleSelect(selectedPos int, items []*MultipleSelectItem, label string) ([]string, error) {
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
		//Label: label,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    label,
			Active:   `→ {{if .IsSelected}}{{ .ID | green }} {{ "✔" | green }} {{else}}{{ .ID }}{{end}}`,
			Inactive: `  {{if .IsSelected}}{{ .ID | green }} {{ "✔" | green }} {{else}}{{ .ID }}{{end}}`,
		},
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
		color.Red("Failed to select branches for cherry-picking due to %v", err)
		return nil, err
	}

	chosenItem := items[selectionIdx]

	if chosenItem.ID != doneID {
		chosenItem.IsSelected = !chosenItem.IsSelected
		return multipleSelect(selectionIdx, items, label)
	}

	// If the user selected the "Done" item, return all selected branches.
	var selectedBranches []string
	for _, i := range items {
		if i.IsSelected {
			selectedBranches = append(selectedBranches, i.ID)
		}
	}
	return selectedBranches, nil
}
