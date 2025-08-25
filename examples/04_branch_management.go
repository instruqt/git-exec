// Package main demonstrates branch operations and management
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
	tempDir, err := os.MkdirTemp("", "git-exec-branch-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	repoPath := filepath.Join(tempDir, "branch-repo")

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
	err = gitInstance.Config("user.name", "Branch Example")
	if err != nil {
		log.Fatal(err)
	}
	err = gitInstance.Config("user.email", "branch@example.com")
	if err != nil {
		log.Fatal(err)
	}

	// Create initial commit
	mainFile := filepath.Join(repoPath, "main.go")
	mainContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	err = os.WriteFile(mainFile, []byte(mainContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Add([]string{"main.go"})
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Commit("Initial commit")
	if err != nil {
		log.Fatal(err)
	}

	// Create new branches
	branchNames := []string{"feature/auth", "feature/database", "bugfix/login-error"}

	for _, branchName := range branchNames {
		err = gitInstance.CreateBranch(branchName)
		if err != nil {
			log.Fatal(err)
		}
	}

	// List branches
	branches, err := gitInstance.ListBranches()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created %d branches:\n", len(branches))
	for _, branch := range branches {
		marker := " "
		if branch.Active {
			marker = "*"
		}
		fmt.Printf("  %s %s\n", marker, branch.Name)
	}

	// Work on feature branch
	err = gitInstance.Checkout(git.CheckoutWithBranch("feature/auth"))
	if err != nil {
		log.Fatal(err)
	}

	authFile := filepath.Join(repoPath, "auth.go")
	authContent := `package main

func authenticate(user, pass string) bool {
	return user == "admin" && pass == "secret"
}
`
	err = os.WriteFile(authFile, []byte(authContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Add([]string{"auth.go"})
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Commit("Add authentication module")
	if err != nil {
		log.Fatal(err)
	}

	// Work on database branch
	err = gitInstance.Checkout(git.CheckoutWithBranch("feature/database"))
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(repoPath, "database.go")
	dbContent := `package main

import "database/sql"

func connect(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite3", dsn)
}
`
	err = os.WriteFile(dbFile, []byte(dbContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Add([]string{"database.go"})
	if err != nil {
		log.Fatal(err)
	}

	err = gitInstance.Commit("Add database connection")
	if err != nil {
		log.Fatal(err)
	}

	// Merge branches back to main
	err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
	if err != nil {
		log.Fatal(err)
	}

	// Merge feature branches
	for _, branchName := range []string{"feature/auth", "feature/database"} {
		result, err := gitInstance.Merge(git.MergeWithBranch(branchName))
		if err != nil {
			log.Fatal(err)
		}
		if result.Success {
			fmt.Printf("Merged %s into main\n", result.MergedBranch)
		}
	}

	// Clean up merged branches
	for _, branchName := range []string{"feature/auth", "feature/database"} {
		err = gitInstance.DeleteBranch(branchName)
		if err != nil {
			log.Printf("Could not delete branch %s: %v\n", branchName, err)
		}
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