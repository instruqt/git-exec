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
	tempDir := filepath.Join(os.TempDir(), "clean-session-api")
	os.RemoveAll(tempDir) // Clean up from previous runs
	
	fmt.Println("=== Clean Session API Examples ===")
	
	// Example 1: Create session and clone into it (recommended pattern)
	fmt.Println("\n1. Create session, then clone:")
	
	sessionPath := filepath.Join(tempDir, "my-project")
	
	// Create session without initializing repository (we'll clone instead)
	session, err := git.NewSessionWithConfig(sessionPath,
		git.WithSessionUser("Alice Developer", "alice@instruqt.com"),
		git.WithInstruqtMetadata("user-001", "session-001", time.Now()),
		git.WithSessionMetadata("track", "git-mastery"),
		git.WithoutRepository(), // Don't init repo, we'll clone
	)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	
	// Clone into the session directory
	// Session automatically applies user context to the clone operation
	err = session.Clone("https://github.com/golang/example", sessionPath)
	if err != nil {
		// If remote clone fails, create local repo for demo
		fmt.Println("Remote clone failed, creating local repo...")
		err = session.InitRepository()
		if err != nil {
			log.Fatalf("Failed to init repository: %v", err)
		}
		
		// Create a test file
		testFile := filepath.Join(sessionPath, "main.go")
		content := `package main

import "fmt"

func main() {
	fmt.Println("Hello from session!")
}
`
		err = os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			log.Fatalf("Failed to create test file: %v", err)
		}
		
		// Add and commit with session user context (automatic)
		err = session.Add([]string{"main.go"})
		if err != nil {
			log.Fatalf("Failed to add file: %v", err)
		}
		
		err = session.Commit("Initial commit from session")
		if err != nil {
			log.Fatalf("Failed to commit: %v", err)
		}
	}
	
	fmt.Printf("✓ Session created at: %s\n", sessionPath)
	fmt.Printf("✓ Session ID: %s, User: %s\n", session.GetSessionID(), session.GetConfig().UserName)
	
	// Example 2: Using WithConfig for one-off operations  
	fmt.Println("\n2. One-off operations with WithConfig:")
	
	quickRepoPath := filepath.Join(tempDir, "quick-repo")
	
	// Regular git instance for quick operations
	git, err := commands.NewGit()
	if err != nil {
		log.Fatalf("Failed to create git instance: %v", err)
	}
	
	err = git.Init(quickRepoPath)
	if err != nil {
		log.Fatalf("Failed to init repo: %v", err)
	}
	
	git.SetWorkingDirectory(quickRepoPath)
	
	// Create and commit file with temporary config
	quickFile := filepath.Join(quickRepoPath, "quick.txt")
	err = os.WriteFile(quickFile, []byte("Quick operation"), 0644)
	if err != nil {
		log.Fatalf("Failed to create quick file: %v", err)
	}
	
	err = git.Add([]string{"quick.txt"})
	if err != nil {
		log.Fatalf("Failed to add file: %v", err)
	}
	
	// Use WithConfig for temporary user context (no persistent session)
	err = git.Commit("Quick commit",
		commands.WithConfig("user.name", "Quick User"),
		commands.WithConfig("user.email", "quick@temp.com"),
	)
	if err != nil {
		log.Fatalf("Failed to commit: %v", err)
	}
	
	fmt.Printf("✓ Quick repo created with temporary config\n")
	
	// Example 3: Session with repository initialization
	fmt.Println("\n3. Session with repository initialization:")
	
	sessionPath2 := filepath.Join(tempDir, "new-project")
	
	// Create session AND initialize repository (default behavior)
	session2, err := git.NewSessionWithConfig(sessionPath2,
		git.WithSessionUser("Bob Builder", "bob@instruqt.com"),
		git.WithInstruqtMetadata("user-002", "session-002", time.Now()),
		git.WithSessionMetadata("environment", "production"),
	)
	if err != nil {
		log.Fatalf("Failed to create session with repo: %v", err)
	}
	
	// Create and commit file (session user context applied automatically)
	projectFile := filepath.Join(sessionPath2, "README.md")
	content := fmt.Sprintf(`# %s

This project was created in a git-exec session.

Session Details:
- User: %s
- Session ID: %s
- Created: %s
`, 
		"New Project",
		session2.GetConfig().UserName,
		session2.GetSessionID(),
		session2.GetConfig().Created.Format(time.RFC3339),
	)
	
	err = os.WriteFile(projectFile, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Failed to create project file: %v", err)
	}
	
	err = session2.Add([]string{"README.md"})
	if err != nil {
		log.Fatalf("Failed to add file: %v", err)
	}
	
	err = session2.Commit("Add project README")
	if err != nil {
		log.Fatalf("Failed to commit: %v", err)
	}
	
	fmt.Printf("✓ Project session created and committed\n")
	
	// Example 4: Multiple config values with WithConfigs
	fmt.Println("\n4. Multiple config values:")
	
	multiRepoPath := filepath.Join(tempDir, "multi-config")
	err = git.Init(multiRepoPath)
	if err != nil {
		log.Fatalf("Failed to init multi repo: %v", err)
	}
	
	git.SetWorkingDirectory(multiRepoPath)
	
	multiFile := filepath.Join(multiRepoPath, "config-test.txt")
	err = os.WriteFile(multiFile, []byte("Multiple config test"), 0644)
	if err != nil {
		log.Fatalf("Failed to create multi file: %v", err)
	}
	
	err = git.Add([]string{"config-test.txt"})
	if err != nil {
		log.Fatalf("Failed to add file: %v", err)
	}
	
	// Apply multiple config values at once
	configs := map[string]string{
		"user.name":        "Charlie Config",
		"user.email":       "charlie@config.com",
		"commit.gpgsign":   "false",
		"core.autocrlf":    "input",
	}
	
	err = git.Commit("Multi-config commit", commands.WithConfigs(configs))
	if err != nil {
		log.Fatalf("Failed to commit with configs: %v", err)
	}
	
	fmt.Printf("✓ Committed with multiple config values\n")
	
	// Summary
	fmt.Println("\n=== Clean API Summary ===")
	fmt.Println("✓ NewSessionWithConfig() - creates sessions (use WithoutRepository() to skip repo init)")
	fmt.Println("✓ session.Clone() - clones with session user context")
	fmt.Println("✓ session.Commit() - commits with session user context")
	fmt.Println("✓ WithConfig() - applies temporary git config to any operation")
	fmt.Println("✓ WithConfigs() - applies multiple config values at once")
	fmt.Println("✓ Sessions automatically persist user context across operations")
	fmt.Println("✓ No more CloneIntoSession() - just create session then clone")
	
	fmt.Printf("\nAll examples in: %s\n", tempDir)
	
	// Show the session configurations
	fmt.Println("\n=== Session Configurations ===")
	
	// Check first session config
	configFile1 := filepath.Join(sessionPath, ".git", "config")
	if data, err := os.ReadFile(configFile1); err == nil {
		fmt.Printf("\nSession 1 config (%s):\n%s\n", sessionPath, string(data))
	}
	
	// Check second session config  
	configFile2 := filepath.Join(sessionPath2, ".git", "config")
	if data, err := os.ReadFile(configFile2); err == nil {
		fmt.Printf("\nSession 2 config (%s):\n%s\n", sessionPath2, string(data))
	}
}