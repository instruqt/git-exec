// Package main demonstrates merge operations and conflict resolution
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/instruqt/git-exec/pkg/git/commands"
	"github.com/instruqt/git-exec/pkg/git/types"
)

func main() {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "git-exec-merge-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	repoPath := filepath.Join(tempDir, "merge-repo")

	// Initialize repository
	git, err := commands.NewGit()
	if err != nil {
		log.Fatal(err)
	}

	err = git.Init(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	git.SetWorkingDirectory(repoPath)

	// Configure user
	err = git.Config("user.name", "Merge Example")
	if err != nil {
		log.Fatal(err)
	}
	err = git.Config("user.email", "merge@example.com")
	if err != nil {
		log.Fatal(err)
	}

	// Create initial commit
	mainFile := filepath.Join(repoPath, "main.txt")
	err = os.WriteFile(mainFile, []byte("line 1\nline 2\nline 3\n"), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = git.Add([]string{"main.txt"})
	if err != nil {
		log.Fatal(err)
	}

	err = git.Commit("Initial commit")
	if err != nil {
		log.Fatal(err)
	}

	// Successful merge example
	err = git.CreateBranch("feature")
	if err != nil {
		log.Fatal(err)
	}

	err = git.Checkout(commands.CheckoutWithBranch("feature"))
	if err != nil {
		log.Fatal(err)
	}

	featureFile := filepath.Join(repoPath, "feature.txt")
	err = os.WriteFile(featureFile, []byte("feature content\n"), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = git.Add([]string{"feature.txt"})
	if err != nil {
		log.Fatal(err)
	}

	err = git.Commit("Add feature file")
	if err != nil {
		log.Fatal(err)
	}

	err = git.Checkout(commands.CheckoutWithBranch("main"))
	if err != nil {
		log.Fatal(err)
	}

	// Successful merge
	result, err := git.Merge(commands.MergeWithBranch("feature"))
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		fmt.Printf("Merge successful: %s (fast-forward: %t)\n", result.MergedBranch, result.FastForward)
	}

	// Conflicting merge example
	err = git.CreateBranch("conflicting")
	if err != nil {
		log.Fatal(err)
	}

	err = git.Checkout(commands.CheckoutWithBranch("conflicting"))
	if err != nil {
		log.Fatal(err)
	}

	// Modify same line differently
	conflictContent := `line 1
conflicting line 2
line 3
`
	err = os.WriteFile(mainFile, []byte(conflictContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = git.Add([]string{"main.txt"})
	if err != nil {
		log.Fatal(err)
	}

	err = git.Commit("Conflicting changes")
	if err != nil {
		log.Fatal(err)
	}

	// Switch to main and make different changes
	err = git.Checkout(commands.CheckoutWithBranch("main"))
	if err != nil {
		log.Fatal(err)
	}

	mainConflictContent := `line 1
main line 2
line 3
`
	err = os.WriteFile(mainFile, []byte(mainConflictContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = git.Add([]string{"main.txt"})
	if err != nil {
		log.Fatal(err)
	}

	err = git.Commit("Main changes")
	if err != nil {
		log.Fatal(err)
	}

	// Attempt merge - this will create conflicts
	conflictResult, err := git.Merge(commands.MergeWithBranch("conflicting"))
	if err != nil {
		log.Fatal(err)
	}

	if !conflictResult.Success {
		fmt.Printf("Merge conflicts detected in %d files\n", len(conflictResult.ConflictedFiles))

		// Show conflict details
		for _, conflict := range conflictResult.Conflicts {
			fmt.Printf("Conflict in %s:\n", conflict.Path)
			for i, section := range conflict.Sections {
				fmt.Printf("  Section %d: Our=%q, Their=%q\n", i+1, section.OurContent, section.TheirContent)
			}
		}

		// Resolve conflicts using "ours" strategy
		resolutions := []types.ConflictResolution{
			{
				FilePath: "main.txt",
				UseOurs:  true,
			},
		}

		err = git.ResolveConflicts(resolutions)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Conflicts resolved using 'ours' strategy")

		// Continue the merge
		err = git.MergeContinue()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Merge completed successfully")
	}

	// Show final commit history
	logs, err := git.Log(commands.LogWithMaxCount("5"))
	if err == nil {
		fmt.Println("\nCommit history:")
		for i, logEntry := range logs {
			fmt.Printf("  %d. %s\n", i+1, logEntry.Message)
		}
	}

	fmt.Printf("\nRepository created at: %s\n", repoPath)
}