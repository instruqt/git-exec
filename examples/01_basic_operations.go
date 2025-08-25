// Package main demonstrates basic Git operations using git-exec
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/instruqt/git-exec/pkg/git/commands"
)

func main() {
	// Create temporary directory for examples
	tempDir, err := os.MkdirTemp("", "git-exec-basic-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	repoPath := filepath.Join(tempDir, "my-repo")

	// Create new Git instance
	git, err := commands.NewGit()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repository
	err = git.Init(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	git.SetWorkingDirectory(repoPath)

	// Configure user
	err = git.Config("user.name", "Example User")
	if err != nil {
		log.Fatal(err)
	}
	err = git.Config("user.email", "user@example.com")
	if err != nil {
		log.Fatal(err)
	}

	// Create and add files
	readmeFile := filepath.Join(repoPath, "README.md")
	err = os.WriteFile(readmeFile, []byte("# My Project\n\nThis is an example project."), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = git.Add([]string{"README.md"})
	if err != nil {
		log.Fatal(err)
	}

	// Check status
	files, err := git.Status()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Repository status (%d files):\n", len(files))
	for _, file := range files {
		fmt.Printf("  %s %s\n", file.Status, file.Name)
	}

	// Create initial commit
	err = git.Commit("Initial commit with README")
	if err != nil {
		log.Fatal(err)
	}

	// Add more files
	codeFile := filepath.Join(repoPath, "main.go")
	codeContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, Git!")
}
`
	err = os.WriteFile(codeFile, []byte(codeContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = git.Add([]string{"main.go"})
	if err != nil {
		log.Fatal(err)
	}

	err = git.Commit("Add main.go")
	if err != nil {
		log.Fatal(err)
	}

	// Show commit history
	logs, err := git.Log(commands.LogWithMaxCount("3"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nRecent commits:")
	for i, logEntry := range logs {
		fmt.Printf("  %d. %s - %s\n", i+1, logEntry.Message, logEntry.Author)
	}

	fmt.Printf("\nRepository created at: %s\n", repoPath)
}