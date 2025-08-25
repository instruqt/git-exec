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
	tempDir := filepath.Join(os.TempDir(), "final-clean-api")
	os.RemoveAll(tempDir) // Clean up
	
	fmt.Println("=== Final Clean and Consistent API ===")
	
	// Example 1: Session with consistent options pattern
	fmt.Println("\n1. Sessions with consistent options:")
	
	projectPath := filepath.Join(tempDir, "my-project")
	
	// NewSession with options - consistent with rest of API
	session, err := git.NewSession(projectPath,
		git.WithUser("Alice Developer", "alice@company.com"),
		git.WithInstruqtMetadata("user-123", "session-456", time.Now()),
		git.WithMetadata("project", "web-app"),
		git.WithMetadata("environment", "development"),
	)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	
	fmt.Printf("✓ Session created with consistent API\n")
	fmt.Printf("✓ Session ID: %s, User: %s\n", 
		session.GetSessionID(), 
		session.GetConfig().UserName)
	
	// Create and commit file - session user context applied automatically
	readmeFile := filepath.Join(projectPath, "README.md")
	content := `# My Project

This project demonstrates the clean, consistent git-exec API.

## Features

- NewSession() with options (not NewSessionWithConfig)
- WithUser() instead of WithSessionUser()
- WithMetadata() instead of WithSessionMetadata()  
- session.Clone() applies session context automatically
- session.Commit() uses session user info automatically
`
	
	err = os.WriteFile(readmeFile, []byte(content), 0644)
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
	
	fmt.Println("✓ Committed with session context")
	
	// Example 2: One-off operations with WithConfig
	fmt.Println("\n2. One-off operations with WithConfig:")
	
	quickRepoPath := filepath.Join(tempDir, "quick-repo")
	
	g, err := commands.NewGit()
	if err != nil {
		log.Fatalf("Failed to create git: %v", err)
	}
	
	err = g.Init(quickRepoPath)
	if err != nil {
		log.Fatalf("Failed to init: %v", err)
	}
	
	g.SetWorkingDirectory(quickRepoPath)
	
	quickFile := filepath.Join(quickRepoPath, "quick.txt")
	err = os.WriteFile(quickFile, []byte("Quick work"), 0644)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	
	err = g.Add([]string{"quick.txt"})
	if err != nil {
		log.Fatalf("Failed to add: %v", err)
	}
	
	// One-off commit with temporary config - consistent with options pattern
	err = g.Commit("Quick commit",
		commands.WithConfig("user.name", "Quick Worker"),
		commands.WithConfig("user.email", "quick@temp.com"),
	)
	if err != nil {
		log.Fatalf("Failed to commit: %v", err)
	}
	
	fmt.Println("✓ One-off commit with temporary config")
	
	// Example 3: Clone into session (the clean way)
	fmt.Println("\n3. Clone into session directory:")
	
	clonePath := filepath.Join(tempDir, "cloned-project")
	
	// Create session for directory that will be cloned into
	cloneSession, err := git.NewSession(clonePath,
		git.WithUser("Bob Cloner", "bob@company.com"),
		git.WithInstruqtMetadata("user-789", "session-789", time.Now()),
		git.WithMetadata("type", "clone"),
	)
	if err != nil {
		log.Fatalf("Failed to create clone session: %v", err)
	}
	
	// Remove the initialized repo since we want to clone
	os.RemoveAll(filepath.Join(clonePath, ".git"))
	
	// Clone with session context applied
	err = cloneSession.Clone(projectPath, clonePath)
	if err != nil {
		log.Fatalf("Failed to clone: %v", err)
	}
	
	fmt.Printf("✓ Cloned with session context\n")
	
	// Example 4: Working with existing repository
	fmt.Println("\n4. Working with existing repository:")
	
	// NewSession automatically detects existing repository
	existingSession, err := git.NewSession(projectPath,
		git.WithUser("Charlie Contributor", "charlie@company.com"),
		git.WithInstruqtMetadata("user-999", "session-999", time.Now()),
		git.WithMetadata("role", "contributor"),
	)
	if err != nil {
		log.Fatalf("Failed to open existing repo: %v", err)
	}
	
	fmt.Printf("✓ Opened existing repo with new session context\n")
	fmt.Printf("✓ New session ID: %s\n", existingSession.GetSessionID())
	
	// Add a file with new user context
	contributionFile := filepath.Join(projectPath, "CONTRIBUTING.md")
	contribution := `# Contributing

Guidelines for contributing to this project.

## Getting Started

1. Create a session: git.NewSession(path, git.WithUser(...))
2. Make changes: session.Add(...), session.Commit(...)
3. Push: session.Push(...)
`
	
	err = os.WriteFile(contributionFile, []byte(contribution), 0644)
	if err != nil {
		log.Fatalf("Failed to create contributing file: %v", err)
	}
	
	err = existingSession.Add([]string{"CONTRIBUTING.md"})
	if err != nil {
		log.Fatalf("Failed to add contributing: %v", err)
	}
	
	err = existingSession.Commit("Add contributing guidelines")
	if err != nil {
		log.Fatalf("Failed to commit contributing: %v", err)
	}
	
	fmt.Println("✓ Added contribution with Charlie's context")
	
	// Example 5: Multiple configs at once
	fmt.Println("\n5. Multiple configs with WithConfigs:")
	
	multiRepoPath := filepath.Join(tempDir, "multi-config")
	err = g.Init(multiRepoPath)
	if err != nil {
		log.Fatalf("Failed to init multi repo: %v", err)
	}
	
	g.SetWorkingDirectory(multiRepoPath)
	
	multiFile := filepath.Join(multiRepoPath, "config-demo.txt")
	err = os.WriteFile(multiFile, []byte("Multiple config demo"), 0644)
	if err != nil {
		log.Fatalf("Failed to create multi file: %v", err)
	}
	
	err = g.Add([]string{"config-demo.txt"})
	if err != nil {
		log.Fatalf("Failed to add multi file: %v", err)
	}
	
	// Multiple configs applied at once
	configs := map[string]string{
		"user.name":       "Diana Multi",
		"user.email":      "diana@example.com",
		"commit.gpgsign":  "false",
		"core.autocrlf":   "input",
	}
	
	err = g.Commit("Multi-config commit", commands.WithConfigs(configs))
	if err != nil {
		log.Fatalf("Failed to multi-config commit: %v", err)
	}
	
	fmt.Println("✓ Committed with multiple config values")
	
	// Show commit history
	fmt.Println("\n=== Commit History ===")
	logs, err := session.Log(commands.LogWithMaxCount("5"))
	if err == nil {
		fmt.Println("\nProject commits:")
		for i, log := range logs {
			fmt.Printf("%d. \"%s\" - %s\n", i+1, log.Message, log.Author)
		}
	}
	
	// Final summary
	fmt.Println("\n=== API Consistency Summary ===")
	fmt.Println("✓ NewSession(path, options...) - consistent with other functions")
	fmt.Println("✓ WithUser() - shorter, cleaner than WithSessionUser()")
	fmt.Println("✓ WithMetadata() - shorter than WithSessionMetadata()")  
	fmt.Println("✓ WithConfig() - temporary config for any git operation")
	fmt.Println("✓ WithConfigs() - multiple config values at once")
	fmt.Println("✓ session.Clone() - applies session context automatically")
	fmt.Println("✓ session.Commit() - uses session user info automatically")
	fmt.Println("✓ Auto-detection - NewSession() detects existing repos")
	fmt.Println("✓ No special cases - all functions follow same option pattern")
	
	fmt.Printf("\nAll examples in: %s\n", tempDir)
}