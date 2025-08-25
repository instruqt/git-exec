// Package main demonstrates bare repository operations and reference management
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
	tempDir, err := os.MkdirTemp("", "git-exec-bare-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a normal repository first
	sourceRepo := filepath.Join(tempDir, "source-repo")
	gitInstance, err := git.NewGit()
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Init(sourceRepo)
	if err != nil {
		log.Fatal(err)
	}

	gitInstance.SetWorkingDirectory(sourceRepo)

	// Configure user
	err = gitInstance.SetConfig("user.name", "Bare Example")
	if err != nil {
		log.Fatal(err)
	}
	err = gitInstance.SetConfig("user.email", "bare@example.com")
	if err != nil {
		log.Fatal(err)
	}

	// Create some commits
	for i := 1; i <= 3; i++ {
		filename := fmt.Sprintf("file%d.txt", i)
		filepath := filepath.Join(sourceRepo, filename)
		content := fmt.Sprintf("Content for file %d", i)
		
		err = os.WriteFile(filepath, []byte(content), 0644)
		if err != nil {
			log.Fatal(err)
		}
		
		err = gitInstance.Add([]string{filename})
		if err != nil {
			log.Fatal(err)
		}
		
		err = gitInstance.Commit(fmt.Sprintf("Add file %d", i))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Check if source is bare (should be false)
	isBare, err := gitInstance.IsBareRepository()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Source repository is bare: %v\n", isBare)

	// Create a bare repository by cloning
	bareRepo := filepath.Join(tempDir, "server.git")
	fmt.Printf("\nCreating bare repository at: %s\n", bareRepo)
	
	err = gitInstance.Clone(sourceRepo, bareRepo, git.CloneWithBare())
	if err != nil {
		log.Fatal(err)
	}

	// Work with the bare repository
	gitInstance.SetWorkingDirectory(bareRepo)

	// Check if it's bare
	isBare, err = gitInstance.IsBareRepository()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server repository is bare: %v\n", isBare)

	// List all references
	refs, err := gitInstance.ListRefs()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nReferences in bare repository:\n")
	for _, ref := range refs {
		fmt.Printf("  %s (%s) -> %s\n", ref.Name, ref.Type, ref.Commit[:8])
	}

	// Get the first commit hash for demonstration
	logs, err := gitInstance.Log(git.LogWithMaxCount("3"))
	if err != nil {
		log.Fatal(err)
	}

	if len(logs) >= 2 {
		firstCommit := logs[len(logs)-1].Commit // Oldest commit
		secondCommit := logs[len(logs)-2].Commit

		// Create a new reference pointing to the first commit
		fmt.Printf("\nCreating new reference 'refs/heads/v1.0' pointing to first commit\n")
		err = gitInstance.UpdateRef("refs/heads/v1.0", firstCommit)
		if err != nil {
			log.Fatal(err)
		}

		// Create a tag reference
		fmt.Printf("Creating tag 'refs/tags/release-1.0' pointing to second commit\n")
		err = gitInstance.UpdateRef("refs/tags/release-1.0", secondCommit)
		if err != nil {
			log.Fatal(err)
		}

		// List references again
		refs, err = gitInstance.ListRefs()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nUpdated references:\n")
		for _, ref := range refs {
			switch ref.Type {
			case types.ReferenceTypeBranch:
				fmt.Printf("  [Branch] %s -> %s\n", ref.Name, ref.Commit[:8])
			case types.ReferenceTypeTag:
				fmt.Printf("  [Tag]    %s -> %s\n", ref.Name, ref.Commit[:8])
			default:
				fmt.Printf("  [%s] %s -> %s\n", ref.Type, ref.Name, ref.Commit[:8])
			}
		}

		// Delete a reference
		fmt.Printf("\nDeleting reference 'refs/heads/v1.0'\n")
		err = gitInstance.DeleteRef("refs/heads/v1.0")
		if err != nil {
			log.Fatal(err)
		}

		// Verify deletion
		refs, err = gitInstance.ListRefs()
		if err != nil {
			log.Fatal(err)
		}

		deleted := true
		for _, ref := range refs {
			if ref.Name == "refs/heads/v1.0" {
				deleted = false
				break
			}
		}
		fmt.Printf("Reference successfully deleted: %v\n", deleted)
	}

	// Create a bare repository from scratch
	bareScratch := filepath.Join(tempDir, "scratch.git")
	fmt.Printf("\nCreating bare repository from scratch at: %s\n", bareScratch)
	
	err = gitInstance.Init(bareScratch, git.InitWithBare())
	if err != nil {
		log.Fatal(err)
	}

	gitInstance.SetWorkingDirectory(bareScratch)
	
	isBare, err = gitInstance.IsBareRepository()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Scratch repository is bare: %v\n", isBare)

	fmt.Printf("\nBare repository operations completed successfully!\n")
	fmt.Printf("\nKey concepts demonstrated:\n")
	fmt.Printf("- Detecting bare vs non-bare repositories\n")
	fmt.Printf("- Creating bare repositories with Init and Clone\n")
	fmt.Printf("- Managing references directly without working directory\n")
	fmt.Printf("- Creating, updating, and deleting refs\n")
	fmt.Printf("- Working with different reference types (branches, tags)\n")
}