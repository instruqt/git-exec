// Package main demonstrates bare repository operations and reference management
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/instruqt/git-exec/pkg/git"
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

	// List branches and tags using existing methods
	branches, err := gitInstance.ListBranches()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nBranches in bare repository:\n")
	for _, branch := range branches {
		fmt.Printf("  %s\n", branch.Name)
	}

	// Create a new branch in bare repository
	err = gitInstance.CreateBranch("v1.0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nCreated new branch 'v1.0'\n")

	// Create a tag in bare repository  
	err = gitInstance.Tag("release-1.0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created tag 'release-1.0'\n")

	// List tags
	tags, err := gitInstance.ListTags()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nTags in bare repository:\n")
	for _, tag := range tags {
		fmt.Printf("  %s\n", tag)
	}

	// Delete the branch
	err = gitInstance.DeleteBranch("v1.0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nDeleted branch 'v1.0'\n")

	// Delete the tag
	err = gitInstance.DeleteTag("release-1.0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted tag 'release-1.0'\n")

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
	fmt.Printf("- Working with branches and tags in bare repositories\n")
	fmt.Printf("- All standard Git operations work in bare repos (no working directory needed)\n")
}