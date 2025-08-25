package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/instruqt/git-exec/pkg/git"
	"github.com/instruqt/git-exec/pkg/git/commands"
	_ "github.com/instruqt/git-exec/pkg/git/commands" // Register session factory
)

func main() {
	tempDir := filepath.Join(os.TempDir(), "simple-session-demo")
	os.RemoveAll(tempDir) // Clean up
	
	fmt.Println("=== Simple Session API ===")
	
	// Example 1: New project (directory doesn't exist)
	fmt.Println("\n1. Creating new project session:")
	
	projectPath := filepath.Join(tempDir, "my-new-project")
	
	// NewSessionWithConfig automatically creates directory and initializes repository
	session, err := git.NewSessionWithConfig(projectPath,
		git.WithSessionUser("Alice Developer", "alice@company.com"),
		git.WithInstruqtMetadata("user-123", "session-456", time.Now()),
		git.WithSessionMetadata("project", "web-app"),
	)
	if err != nil {
		log.Fatalf("Failed to create new project session: %v", err)
	}
	
	fmt.Printf("✓ New project session created at: %s\n", projectPath)
	fmt.Printf("✓ Session ID: %s\n", session.GetSessionID())
	
	// Create and commit a file - user context applied automatically
	readmeFile := filepath.Join(projectPath, "README.md")
	err = os.WriteFile(readmeFile, []byte("# My New Project\n\nCreated with git-exec sessions!"), 0644)
	if err != nil {
		log.Fatalf("Failed to create README: %v", err)
	}
	
	err = session.Add([]string{"README.md"})
	if err != nil {
		log.Fatalf("Failed to add README: %v", err)
	}
	
	err = session.Commit("Initial commit")
	if err != nil {
		log.Fatalf("Failed to commit: %v", err)
	}
	
	fmt.Println("✓ Created and committed README.md")
	
	// Example 2: Working with existing repository
	fmt.Println("\n2. Opening existing project session:")
	
	// NewSessionWithConfig detects existing repo and loads it
	existingSession, err := git.NewSessionWithConfig(projectPath,
		git.WithSessionUser("Bob Collaborator", "bob@company.com"),
		git.WithInstruqtMetadata("user-789", "session-012", time.Now()),
		git.WithSessionMetadata("role", "reviewer"),
	)
	if err != nil {
		log.Fatalf("Failed to open existing project: %v", err)
	}
	
	fmt.Printf("✓ Opened existing project as new session\n")
	fmt.Printf("✓ New session ID: %s\n", existingSession.GetSessionID())
	
	// Add another file with the new user context
	codeFile := filepath.Join(projectPath, "main.go")
	codeContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello from", "Bob!")
}
`
	err = os.WriteFile(codeFile, []byte(codeContent), 0644)
	if err != nil {
		log.Fatalf("Failed to create main.go: %v", err)
	}
	
	err = existingSession.Add([]string{"main.go"})
	if err != nil {
		log.Fatalf("Failed to add main.go: %v", err)
	}
	
	// This commit will use Bob's user info automatically
	err = existingSession.Commit("Add main.go")
	if err != nil {
		log.Fatalf("Failed to commit: %v", err)
	}
	
	fmt.Println("✓ Bob added main.go with his user context")
	
	// Example 3: Clone into session directory
	fmt.Println("\n3. Clone into empty directory:")
	
	clonePath := filepath.Join(tempDir, "cloned-project")
	
	// Create session for empty directory
	cloneSession, err := git.NewSessionWithConfig(clonePath,
		git.WithSessionUser("Charlie Cloner", "charlie@company.com"),
		git.WithInstruqtMetadata("user-999", "session-999", time.Now()),
	)
	if err != nil {
		log.Fatalf("Failed to create clone session: %v", err)
	}
	
	// Since clonePath doesn't exist yet, NewSessionWithConfig creates it and inits repo
	// But we want to clone instead, so let's remove the .git directory
	gitDir := filepath.Join(clonePath, ".git")
	os.RemoveAll(gitDir)
	
	// Now clone our first project into this directory
	// Session automatically applies Charlie's user context
	err = cloneSession.Clone(projectPath, clonePath)
	if err != nil {
		log.Fatalf("Failed to clone: %v", err)
	}
	
	fmt.Printf("✓ Cloned project with Charlie's session context\n")
	
	// Example 4: Using WithConfig for one-off operations
	fmt.Println("\n4. One-off operations with WithConfig:")
	
	tempRepoPath := filepath.Join(tempDir, "temp-work")
	
	// Regular git operations with temporary config
	git, err := commands.NewGit()
	if err != nil {
		log.Fatalf("Failed to create git instance: %v", err)
	}
	
	err = git.Init(tempRepoPath)
	if err != nil {
		log.Fatalf("Failed to init temp repo: %v", err)
	}
	
	git.SetWorkingDirectory(tempRepoPath)
	
	tempFile := filepath.Join(tempRepoPath, "temp.txt")
	err = os.WriteFile(tempFile, []byte("Temporary work"), 0644)
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	
	err = git.Add([]string{"temp.txt"})
	if err != nil {
		log.Fatalf("Failed to add temp file: %v", err)
	}
	
	// One-off commit with temporary user config
	err = git.Commit("Temporary work commit",
		commands.WithConfig("user.name", "Temp Worker"),
		commands.WithConfig("user.email", "temp@example.com"),
	)
	if err != nil {
		log.Fatalf("Failed to commit temp work: %v", err)
	}
	
	fmt.Println("✓ One-off commit with temporary user config")
	
	// Show the results
	fmt.Println("\n=== Results ===")
	
	// Show git logs from the main project to see different authors
	fmt.Println("\nCommit history from main project:")
	logs, err := session.Log(commands.LogWithMaxCount("10"))
	if err == nil {
		for i, log := range logs {
			fmt.Printf("%d. %s - %s\n", i+1, log.Message, log.Author)
		}
	}
	
	fmt.Printf("\nAll projects created in: %s\n", tempDir)
	fmt.Println("\n✓ Session API is now clean and simple:")
	fmt.Println("  - NewSessionWithConfig() detects existing repos automatically")
	fmt.Println("  - session.Clone() applies session user context")
	fmt.Println("  - session.Commit() uses session user info")
	fmt.Println("  - WithConfig() for temporary one-off config")
	fmt.Println("  - No more WithoutRepository() or CloneIntoSession() needed!")
}