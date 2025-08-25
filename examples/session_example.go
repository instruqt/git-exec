package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/instruqt/git-exec/pkg/git"
	_ "github.com/instruqt/git-exec/pkg/git/commands" // Register session factory
)

func main() {
	// Example 1: Create a new session with user configuration
	fmt.Println("=== Creating New Session ===")
	
	tempDir := filepath.Join(os.TempDir(), "git-session-example")
	sessionPath := filepath.Join(tempDir, "my-session")
	
	// Clean up from previous runs
	os.RemoveAll(tempDir)
	
	// Create a new session with persistent configuration
	session, err := git.NewSessionWithConfig(sessionPath,
		git.WithSessionUser("Jane Developer", "jane@instruqt.com"),
		git.WithInstruqtMetadata("user-123", "session-456", time.Now()),
		git.WithSessionMetadata("track", "git-basics"),
		git.WithSessionMetadata("environment", "production"),
	)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	
	fmt.Printf("Session created at: %s\n", sessionPath)
	fmt.Printf("Session ID: %s\n", session.GetSessionID())
	fmt.Printf("User ID: %s\n", session.GetUserID())
	
	// Create a test file
	testFile := filepath.Join(sessionPath, "README.md")
	content := `# Session Example

This file was created in a git-exec session.

Session Details:
- User: Jane Developer
- Email: jane@instruqt.com
- Session ID: session-456
- User ID: user-123
`
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Failed to create test file: %v", err)
	}
	
	// Add and commit the file
	err = session.Add([]string{"README.md"})
	if err != nil {
		log.Fatalf("Failed to add file: %v", err)
	}
	
	err = session.Commit("Initial commit from session")
	if err != nil {
		log.Fatalf("Failed to commit: %v", err)
	}
	
	fmt.Println("File committed successfully with session user context")
	
	// Example 2: Load an existing session
	fmt.Println("\n=== Loading Existing Session ===")
	
	loadedSession, err := git.LoadSession(sessionPath)
	if err != nil {
		log.Fatalf("Failed to load session: %v", err)
	}
	
	config := loadedSession.GetConfig()
	fmt.Printf("Loaded session for user: %s <%s>\n", config.UserName, config.UserEmail)
	fmt.Printf("Session ID: %s\n", config.SessionID)
	fmt.Printf("User ID: %s\n", config.UserID)
	fmt.Printf("Track: %s\n", config.Metadata["track"])
	fmt.Printf("Environment: %s\n", config.Metadata["environment"])
	
	// Example 3: Update user information
	fmt.Println("\n=== Updating User Information ===")
	
	err = loadedSession.UpdateUser("Jane Smith", "jane.smith@instruqt.com")
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	
	fmt.Println("User information updated")
	
	// Example 4: Validate session
	fmt.Println("\n=== Validating Session ===")
	
	err = git.ValidateSession(sessionPath)
	if err != nil {
		log.Printf("Session validation failed: %v", err)
	} else {
		fmt.Println("Session is valid")
	}
	
	// Example 5: Get session info without loading full session
	fmt.Println("\n=== Getting Session Info ===")
	
	info, err := git.GetSessionInfo(sessionPath)
	if err != nil {
		log.Fatalf("Failed to get session info: %v", err)
	}
	
	fmt.Printf("Session info:\n")
	fmt.Printf("  User: %s <%s>\n", info.UserName, info.UserEmail)
	fmt.Printf("  Session ID: %s\n", info.SessionID)
	fmt.Printf("  User ID: %s\n", info.UserID)
	fmt.Printf("  Created: %s\n", info.Created.Format(time.RFC3339))
	
	// Example 6: Clone into a new session
	fmt.Println("\n=== Cloning Into New Session ===")
	
	clonedSessionPath := filepath.Join(tempDir, "cloned-session")
	
	// For this example, we'll use the local session as the source
	// In practice, you'd use a GitHub URL
	clonedSession, err := git.CloneIntoSession(
		sessionPath,
		clonedSessionPath,
		git.WithSessionUser("John Developer", "john@instruqt.com"),
		git.WithInstruqtMetadata("user-789", "session-012", time.Now()),
	)
	if err != nil {
		log.Fatalf("Failed to clone into session: %v", err)
	}
	
	fmt.Printf("Cloned into new session at: %s\n", clonedSessionPath)
	fmt.Printf("New session ID: %s\n", clonedSession.GetSessionID())
	
	// Clean up (optional - comment out to inspect the created repositories)
	// fmt.Println("\n=== Cleaning Up ===")
	// os.RemoveAll(tempDir)
	// fmt.Println("Temporary files removed")
	
	fmt.Println("\n=== Example Complete ===")
	fmt.Printf("Session files are in: %s\n", tempDir)
	fmt.Println("You can inspect the git repositories and their .git/config files")
}