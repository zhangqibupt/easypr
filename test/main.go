package main

import (
	"fmt"
	"regexp"
	"strings"
)

type FileChange struct {
	FileName   string
	OldContent string
	NewContent string
}

func parseUnifiedFormat(unifiedDiff string) ([]FileChange, error) {
	var fileChanges []FileChange

	// Define regular expressions to match file paths and content
	filePathRegex := regexp.MustCompile(`^--- ([^\t]+)\t.*$`)
	contentRegex := regexp.MustCompile(`^([-+ ])(.*)$`)

	var currentFileChange *FileChange

	// Iterate over lines in the Unified Format string
	for _, line := range strings.Split(unifiedDiff, "\n") {
		if match := filePathRegex.FindStringSubmatch(line); match != nil {
			// Start of a new file change
			if currentFileChange != nil {
				fileChanges = append(fileChanges, *currentFileChange)
			}
			fileName := extractFileName(match[1])
			currentFileChange = &FileChange{FileName: fileName}
		} else if match := contentRegex.FindStringSubmatch(line); match != nil {
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

	return fileChanges, nil
}

func extractFileName(fullPath string) string {
	// Extract just the file name from the full path
	elements := strings.Split(fullPath, "/")
	return elements[len(elements)-1]
}

func main() {
	unifiedDiff := `--- lao	2002-02-21 23:30:39.942229878 -0800
+++ tzu	2002-02-21 23:30:50.442260588 -0800
@@ -1,7 +1,6 @@
-The Way that can be told of is not the eternal Way;
-The name that can be named is not the eternal name.
 The Nameless is the origin of Heaven and Earth;
-The Named is the mother of all things.
+The named is the mother of all things.
+
 Therefore let there always be non-being,
   so we may see their subtlety,
 And let there always be being,
@@ -9,3 +8,6 @@
 The two are the same,
 But after they are produced,
   they have different names.
+They both may be called deep and profound.
+Deeper and more profound,
+The door of all subtleties!
--- lao	2002-02-21 23:30:39.942229878 -0800
+++ tzu	2002-02-21 23:30:50.442260588 -0800
@@ -1,7 +1,6 @@
-The Way that can be told of is not the eternal Way;
-The name that can be named is not the eternal name.
 The Nameless is the origin of Heaven and Earth;
-The Named is the mother of all things.
+The named is the mother of all things.
+
 Therefore let there always be non-being,
   so we may see their subtlety,
 And let there always be being,
@@ -9,3 +8,6 @@
 The two are the same,
 But after they are produced,
   they have different names.
+They both may be called deep and profound.
+Deeper and more profound,
+The door of all subtleties!`

	fileChanges, err := parseUnifiedFormat(unifiedDiff)
	if err != nil {
		fmt.Println("Error parsing Unified Format:", err)
		return
	}

	for _, change := range fileChanges {
		fmt.Printf("File Name: %s\n", change.FileName)
		fmt.Printf("Old Content:\n%s", change.OldContent)
		fmt.Printf("New Content:\n%s", change.NewContent)
		fmt.Println(strings.Repeat("-", 50))
	}
}
