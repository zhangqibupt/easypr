package lib

import (
	"regexp"
	"strings"
)

var ChangeTypeAdd = "add"
var ChangeTypeDelete = "delete"
var ChangeTypeModify = "modify"

type FileChange struct {
	ChangeType      string
	FileName        string
	OldContent      string
	NewContent      string
	RawDiff         string
	LineNumber      int
	OriginalContent string
}

func ParseUnifiedFormat(unifiedDiff string, extensions []string) ([]FileChange, error) {
	var fileChanges []FileChange

	// diff --git a/func/configuration/config/env.prd.us-east-1.yml b/func/configuration/config/env.prd.us-east-1.yml
	// Define regular expressions to match file paths and content
	filePathRegex := regexp.MustCompile(`^diff --git a/(.+)$`)
	//filePathRegex := regexp.MustCompile(`^--- a/(.+)$`)
	contentRegex := regexp.MustCompile(`^([-+ ])(.*)$`)

	var currentFileChange *FileChange

	// Iterate over lines in the Unified Format string
	for _, line := range strings.Split(unifiedDiff, "\n") {
		if match := filePathRegex.FindStringSubmatch(line); match != nil {
			// Start of a new file change
			if currentFileChange != nil {
				fileChanges = append(fileChanges, *currentFileChange)
			}
			fileName := match[1]
			//fileName := extractFileName(match[1])
			currentFileChange = &FileChange{FileName: fileName}
			continue
		}

		if currentFileChange != nil {
			currentFileChange.RawDiff += line + "\n"
		}

		if match := contentRegex.FindStringSubmatch(line); match != nil {
			currentFileChange.RawDiff += line + "\n"
			// Line contains old or new content
			switch match[1] {
			case "-":
				currentFileChange.OldContent += match[2] + "\n"
			case "+":
				currentFileChange.NewContent += match[2] + "\n"
			case " ":
				// Unchanged line, include in both old and new content
				//currentFileChange.OldContent += match[2] + "\n"
				//currentFileChange.NewContent += match[2] + "\n"
			}
		}
	}

	// Add the last file change to the list
	if currentFileChange != nil {
		fileChanges = append(fileChanges, *currentFileChange)
	}

	if len(extensions) > 0 {
		// filter out by extension
		var whiteList []FileChange
		for _, c := range fileChanges {
			for _, ext := range extensions {
				if strings.HasSuffix(c.FileName, ext) {
					whiteList = append(whiteList, c)
				}
			}
		}
		return whiteList, nil
	}

	return fileChanges, nil
}
