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

var debug = false
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate PRs across branches",
	Run:   CreateRun(),
}

func CreateRun() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if debug {
			git.EnableDebug()
			color.Green("Debug mode enabled")
		}

		os.Chdir("/Users/qzhang/workspace/github/test_easypr")

		// check if there is uncommitted changes
		color.Cyan("Checking for uncommitted changes...")
		hasChanage, err := git.HasUncommittedChanges()
		if err != nil {
			return
		}
		if hasChanage {
			color.Red("You have uncommitted changes, please commit or stash them before continue.")
			return
		}

		// check if remote origin is set
		color.Cyan("Fetching remote branche list...")
		originUrl, err := git.RemoteURL("origin")
		if err != nil {
			return
		}
		if strings.TrimSpace(originUrl) == "" {
			color.Red("Failed to get remote origin, run `git remote -v` to check if the 'origin' has been set.")
			return
		}

		// get remote branches
		branches, err := git.RemoteBranches()
		if err != nil {
			return
		}

		if len(branches) == 0 {
			color.Red("No remote branches found!")
			return
		}

		// select target branch
		i, targetBranch, err := (&promptui.Select{Label: "Select target branch", Items: branches, Searcher: searcher(branches), Size: 10}).Run()
		if err != nil {
			color.Red("Failed to select current branch due to error: %s", err)
			return
		}

		sourceBranch, err := git.CurrentBranch()
		if err != nil {
			return
		}
		defer git.Checkout(sourceBranch)

		commits, err := git.CommitsBetween(sourceBranch, targetBranch)
		if err != nil {
			return
		}

		if len(commits) == 0 {
			color.Yellow("No diff commits between %s and %s, existing", sourceBranch, targetBranch)
			return
		}

		// select cherry-pick branches
		branches = append(branches[:i], branches[i+1:]...)
		var allItems []*MultipleSelectItem
		for _, item := range branches {
			allItems = append(allItems, &MultipleSelectItem{
				ID: item,
			})
		}

		selected, err := multipleSelect(0, allItems)
		if err != nil {
			return
		}

		color.Cyan("Creating PR for target branch %s...", targetBranch)
		green := color.New(color.FgGreen).SprintFunc()
		targetLink, err := git.CreatePRLink(sourceBranch, targetBranch)
		if err != nil {
			return
		}
		defer color.Cyan("PR to %s: %s", targetBranch, green(targetLink))

		if len(selected) == 0 {
			return
		}

		for _, target := range selected {
			color.Cyan("Creating cherry-pick PR for branch %s...", target)
			link, err := git.CreateCherryPickPRLink(sourceBranch, target, commits)
			if err != nil {
				color.Red("Failed to create cherry-pick PR for branch %s due to error: %s", target)
				continue
			}
			defer color.Cyan("Cherry-Pick PR to %s: %s", target, green(link))
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
		color.Red("Failed to select cherry-pick branches due to %v", err)
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
	createCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Print debug messages")
	rootCmd.AddCommand(createCmd)
}
