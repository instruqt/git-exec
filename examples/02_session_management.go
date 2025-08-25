// Package main demonstrates session management with persistent user context
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/instruqt/git-exec/pkg/git"
)

func main() {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "git-exec-session-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create new session with user configuration
	sessionPath := filepath.Join(tempDir, "project-session")

	session, err := git.NewSession(sessionPath,
		git.WithSessionUser("Alice Developer", "alice@company.com"),
		git.WithInstruqtMetadata("user-123", "session-456", time.Now()),
		git.WithMetadata("project", "web-app"),
		git.WithMetadata("team", "frontend"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Session created: %s\n", session.GetSessionID())

	// Create and commit a file - user context is automatically applied
	readmeFile := filepath.Join(sessionPath, "README.md")
	err = os.WriteFile(readmeFile, []byte("# Web App Project\n\nCreated by Alice"), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = session.Add([]string{"README.md"})
	if err != nil {
		log.Fatal(err)
	}

	err = session.Commit("Initial project setup")
	if err != nil {
		log.Fatal(err)
	}

	// Load existing session
	loadedSession, err := git.LoadSession(sessionPath)
	if err != nil {
		log.Fatal(err)
	}

	config := loadedSession.GetConfig()
	fmt.Printf("Loaded session for: %s <%s>\n", config.UserName, config.UserEmail)
	fmt.Printf("Project: %s, Team: %s\n", config.Metadata["project"], config.Metadata["team"])

	// Update user information
	err = loadedSession.UpdateUser("Alice Smith", "alice.smith@company.com")
	if err != nil {
		log.Fatal(err)
	}

	// Add more content with updated user
	codeFile := filepath.Join(sessionPath, "app.js")
	err = os.WriteFile(codeFile, []byte("console.log('Hello from frontend team!');"), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = loadedSession.Add([]string{"app.js"})
	if err != nil {
		log.Fatal(err)
	}

	err = loadedSession.Commit("Add main application file")
	if err != nil {
		log.Fatal(err)
	}

	// Show commit history
	logs, err := loadedSession.Log(git.LogWithMaxCount("5"))
	if err == nil {
		fmt.Println("\nCommit history:")
		for i, logEntry := range logs {
			fmt.Printf("  %d. %s - %s\n", i+1, logEntry.Message, logEntry.Author)
		}
	}

	fmt.Printf("\nSession repository: %s\n", sessionPath)
}