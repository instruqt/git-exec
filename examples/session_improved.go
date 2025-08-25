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
	tempDir := filepath.Join(os.TempDir(), "git-session-improved")
	os.RemoveAll(tempDir) // Clean up from previous runs
	
	// Example 1: Create session then clone into it (preferred pattern)
	fmt.Println("=== Create Session Then Clone ===")
	
	sessionPath1 := filepath.Join(tempDir, "session1")
	
	// Create an empty session with configuration
	session1, err := git.CreateEmptySession(sessionPath1,
		git.WithSessionUser("Alice Developer", "alice@instruqt.com"),
		git.WithInstruqtMetadata("user-001", "session-001", time.Now()),
		git.WithSessionMetadata("track", "advanced-git"),
	)
	if err != nil {
		log.Fatalf("Failed to create empty session: %v", err)
	}
	
	fmt.Printf("Created empty session at: %s\n", sessionPath1)
	fmt.Printf("Session ID: %s\n", session1.GetSessionID())
	
	// Now clone a repository using the session
	// The session's user configuration will be automatically applied
	err = session1.Clone("https://github.com/golang/example", sessionPath1)
	if err != nil {
		// For this example, if we can't clone from GitHub, create a local repo
		fmt.Println("Could not clone from GitHub, initializing local repository instead")
		err = session1.InitRepository()
		if err != nil {
			log.Fatalf("Failed to initialize repository: %v", err)
		}
	}
	
	fmt.Println("Repository ready in session")
	
	// Example 2: Use WithConfig for one-off operations without session
	fmt.Println("\n=== Using WithConfig for Temporary User Context ===")
	
	repoPath := filepath.Join(tempDir, "temp-repo")
	
	// Create a regular git instance (no session)
	g, err := commands.NewGit()
	if err != nil {
		log.Fatalf("Failed to create git instance: %v", err)
	}
	
	// Initialize repository
	err = g.Init(repoPath)
	if err != nil {
		log.Fatalf("Failed to init repository: %v", err)
	}
	
	g.SetWorkingDirectory(repoPath)
	
	// Create a test file
	testFile := filepath.Join(repoPath, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		log.Fatalf("Failed to create test file: %v", err)
	}
	
	// Add and commit with temporary user config using WithConfig
	err = g.Add([]string{"test.txt"})
	if err != nil {
		log.Fatalf("Failed to add file: %v", err)
	}
	
	// Commit with temporary user configuration (no persistent session)
	err = g.Commit("Test commit",
		commands.WithConfig("user.name", "Bob Temporary"),
		commands.WithConfig("user.email", "bob@temp.com"),
	)
	if err != nil {
		log.Fatalf("Failed to commit: %v", err)
	}
	
	fmt.Println("Committed with temporary user config (no session)")
	
	// Example 3: Create session with repo initialization
	fmt.Println("\n=== Create Session With Repository ===")
	
	sessionPath2 := filepath.Join(tempDir, "session2")
	
	// This creates a session AND initializes a repository
	session2, err := git.NewSessionWithConfig(sessionPath2,
		git.WithSessionUser("Charlie Developer", "charlie@instruqt.com"),
		git.WithInstruqtMetadata("user-002", "session-002", time.Now()),
	)
	if err != nil {
		log.Fatalf("Failed to create session with config: %v", err)
	}
	
	fmt.Printf("Created session with repository at: %s\n", sessionPath2)
	
	// Create a file and commit using session
	testFile2 := filepath.Join(sessionPath2, "README.md")
	content := `# Session Example

This demonstrates the improved session API.`
	
	err = os.WriteFile(testFile2, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Failed to create README: %v", err)
	}
	
	err = session2.Add([]string{"README.md"})
	if err != nil {
		log.Fatalf("Failed to add README: %v", err)
	}
	
	// Commit automatically uses Charlie's user info from session
	err = session2.Commit("Add README")
	if err != nil {
		log.Fatalf("Failed to commit: %v", err)
	}
	
	fmt.Println("Committed with session user context")
	
	// Example 4: Use WithConfigs for multiple config values
	fmt.Println("\n=== Using WithConfigs for Multiple Settings ===")
	
	repoPath3 := filepath.Join(tempDir, "multi-config-repo")
	err = g.Init(repoPath3)
	if err != nil {
		log.Fatalf("Failed to init repository: %v", err)
	}
	
	g.SetWorkingDirectory(repoPath3)
	
	// Create a test file
	testFile3 := filepath.Join(repoPath3, "test.txt")
	err = os.WriteFile(testFile3, []byte("multi config test"), 0644)
	if err != nil {
		log.Fatalf("Failed to create test file: %v", err)
	}
	
	err = g.Add([]string{"test.txt"})
	if err != nil {
		log.Fatalf("Failed to add file: %v", err)
	}
	
	// Commit with multiple config values at once
	configs := map[string]string{
		"user.name":     "Diana Multi",
		"user.email":    "diana@multi.com",
		"core.editor":   "vim",
		"color.ui":      "auto",
	}
	
	err = g.Commit("Test with multiple configs",
		commands.WithConfigs(configs),
	)
	if err != nil {
		log.Fatalf("Failed to commit: %v", err)
	}
	
	fmt.Println("Committed with multiple config values")
	
	// Summary
	fmt.Println("\n=== Summary ===")
	fmt.Println("Session API improvements:")
	fmt.Println("1. CreateEmptySession + Clone: Create session first, then clone into it")
	fmt.Println("2. WithConfig: Apply temporary git config for any operation")
	fmt.Println("3. WithConfigs: Apply multiple config values at once")
	fmt.Println("4. Sessions automatically apply user context to all operations")
	
	fmt.Printf("\nAll examples created in: %s\n", tempDir)
}