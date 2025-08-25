// Package main demonstrates merge operations and conflict resolution
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/instruqt/git-exec/pkg/git"
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
	gitInstance, err := git.NewGit()
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Init(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	gitInstance.SetWorkingDirectory(repoPath)

	// Configure user
	err = gitInstance.SetConfig("user.name", "Merge Example")
	if err != nil {
		log.Fatal(err)
	}
	err = gitInstance.SetConfig("user.email", "merge@example.com")
	if err != nil {
		log.Fatal(err)
	}

	// Create initial commit
	mainFile := filepath.Join(repoPath, "main.txt")
	err = os.WriteFile(mainFile, []byte("line 1\nline 2\nline 3\n"), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Add([]string{"main.txt"})
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Commit("Initial commit")
	if err != nil {
		log.Fatal(err)
	}

	// Successful merge example
	err = gitInstance.CreateBranch("feature")
	if err != nil {
		log.Fatal(err)
	}

	_, err = gitInstance.Checkout(git.CheckoutWithBranch("feature"))
	if err != nil {
		log.Fatal(err)
	}

	featureFile := filepath.Join(repoPath, "feature.txt")
	err = os.WriteFile(featureFile, []byte("feature content\n"), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Add([]string{"feature.txt"})
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Commit("Add feature file")
	if err != nil {
		log.Fatal(err)
	}

	_, err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
	if err != nil {
		log.Fatal(err)
	}

	// Successful merge
	result, err := gitInstance.Merge(git.MergeWithBranch("feature"))
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		fmt.Printf("Merge successful: %s (fast-forward: %t)\n", result.MergedBranch, result.FastForward)
	}

	// Conflicting merge example
	err = gitInstance.CreateBranch("conflicting")
	if err != nil {
		log.Fatal(err)
	}

	_, err = gitInstance.Checkout(git.CheckoutWithBranch("conflicting"))
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

	err = gitInstance.Add([]string{"main.txt"})
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Commit("Conflicting changes")
	if err != nil {
		log.Fatal(err)
	}

	// Switch to main and make different changes
	_, err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
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

	err = gitInstance.Add([]string{"main.txt"})
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Commit("Main changes")
	if err != nil {
		log.Fatal(err)
	}

	// Attempt merge - this will create conflicts
	conflictResult, err := gitInstance.Merge(git.MergeWithBranch("conflicting"))
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

		err = gitInstance.ResolveConflicts(resolutions)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Conflicts resolved using 'ours' strategy")

		// Continue the merge
		err = gitInstance.MergeContinue()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Merge completed successfully")
	}

	// Show final commit history
	logs, err := gitInstance.Log(git.LogWithMaxCount("5"))
	if err == nil {
		fmt.Println("\nCommit history:")
		for i, logEntry := range logs {
			fmt.Printf("  %d. %s\n", i+1, logEntry.Message)
		}
	}

	fmt.Printf("\nRepository created at: %s\n", repoPath)
}