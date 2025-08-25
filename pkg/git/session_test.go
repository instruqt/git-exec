package git_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/instruqt/git-exec/pkg/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test session creation and persistence - the core value of sessions
func TestSessionPersistence(t *testing.T) {
	tempDir := t.TempDir()
	sessionPath := filepath.Join(tempDir, "test-session")
	
	// Create session with configuration
	session, err := git.NewSession(sessionPath,
		git.WithSessionUser("Test User", "test@example.com"),
		git.WithInstruqtMetadata("user-123", "session-456", time.Now()),
		git.WithMetadata("project", "test-project"),
		git.WithMetadata("team", "test-team"),
	)
	require.NoError(t, err)
	
	// Verify session was created and configured
	config := session.GetSessionConfig()
	assert.Equal(t, "Test User", config.UserName)
	assert.Equal(t, "test@example.com", config.UserEmail)
	assert.Equal(t, "user-123", config.UserID)
	assert.Equal(t, "session-456", config.SessionID)
	assert.Equal(t, "test-project", config.Metadata["project"])
	assert.Equal(t, "test-team", config.Metadata["team"])
	
	// Create a commit to test automatic user context
	testFile := filepath.Join(sessionPath, "test.txt")
	err = os.WriteFile(testFile, []byte("session test"), 0644)
	require.NoError(t, err)
	
	err = session.Add([]string{"test.txt"})
	require.NoError(t, err)
	
	err = session.Commit("Test session commit")
	require.NoError(t, err)
	
	// Verify commit has correct author
	logs, err := session.Log(git.LogWithMaxCount("1"))
	require.NoError(t, err)
	require.Len(t, logs, 1)
	assert.Equal(t, "Test session commit", logs[0].Message)
	assert.Contains(t, logs[0].Author, "Test User")
	assert.Contains(t, logs[0].Author, "test@example.com")
}

// Test session loading and configuration persistence
func TestSessionReload(t *testing.T) {
	tempDir := t.TempDir()
	sessionPath := filepath.Join(tempDir, "persistent-session")
	
	originalTime := time.Now().Truncate(time.Second)
	
	// Create session with configuration
	session1, err := git.NewSession(sessionPath,
		git.WithSessionUser("Original User", "original@example.com"),
		git.WithInstruqtMetadata("user-original", "session-original", originalTime),
		git.WithMetadata("version", "1.0"),
	)
	require.NoError(t, err)
	
	// Make a commit to ensure repository exists
	testFile := filepath.Join(sessionPath, "persistent.txt")
	err = os.WriteFile(testFile, []byte("persistent content"), 0644)
	require.NoError(t, err)
	
	err = session1.Add([]string{"persistent.txt"})
	require.NoError(t, err)
	err = session1.Commit("Persistent commit")
	require.NoError(t, err)
	
	// Load existing session
	session2, err := git.LoadSession(sessionPath)
	require.NoError(t, err)
	
	// Verify configuration was persisted
	config := session2.GetSessionConfig()
	assert.Equal(t, "Original User", config.UserName)
	assert.Equal(t, "original@example.com", config.UserEmail)
	assert.Equal(t, "user-original", config.UserID)
	assert.Equal(t, "session-original", config.SessionID)
	assert.Equal(t, "1.0", config.Metadata["version"])
	assert.Equal(t, originalTime, config.Created)
}

// Test session user updates and configuration changes
func TestSessionUserUpdate(t *testing.T) {
	tempDir := t.TempDir()
	sessionPath := filepath.Join(tempDir, "update-session")
	
	// Create session
	session, err := git.NewSession(sessionPath,
		git.WithSessionUser("Old User", "old@example.com"),
	)
	require.NoError(t, err)
	
	// Update user information
	err = session.UpdateUser("New User", "new@example.com")
	require.NoError(t, err)
	
	// Verify update
	config := session.GetSessionConfig()
	assert.Equal(t, "New User", config.UserName)
	assert.Equal(t, "new@example.com", config.UserEmail)
	
	// Test that new commits use updated user
	testFile := filepath.Join(sessionPath, "update-test.txt")
	err = os.WriteFile(testFile, []byte("updated content"), 0644)
	require.NoError(t, err)
	
	err = session.Add([]string{"update-test.txt"})
	require.NoError(t, err)
	err = session.Commit("Updated user commit")
	require.NoError(t, err)
	
	// Verify commit has new author
	logs, err := session.Log(git.LogWithMaxCount("1"))
	require.NoError(t, err)
	require.Len(t, logs, 1)
	assert.Contains(t, logs[0].Author, "New User")
	assert.Contains(t, logs[0].Author, "new@example.com")
}

// Test session validation and error handling
func TestSessionValidation(t *testing.T) {
	tempDir := t.TempDir()
	
	// Test loading non-existent session
	nonExistentPath := filepath.Join(tempDir, "does-not-exist")
	_, err := git.LoadSession(nonExistentPath)
	assert.Error(t, err)
	
	// Test validating non-existent session
	err = git.ValidateSession(nonExistentPath)
	assert.Error(t, err)
	
	// Test creating session in existing directory (should work)
	sessionPath := filepath.Join(tempDir, "existing-dir")
	err = os.MkdirAll(sessionPath, 0755)
	require.NoError(t, err)
	
	session, err := git.NewSession(sessionPath,
		git.WithSessionUser("Test User", "test@example.com"),
		git.WithInstruqtMetadata("user-test", "session-test", time.Now()),
	)
	require.NoError(t, err)
	
	// Validate session
	err = git.ValidateSession(sessionPath)
	assert.NoError(t, err)
	
	// Test session validity
	assert.True(t, session.IsValid())
	
	// Test session info retrieval
	info, err := git.GetSessionInfo(sessionPath)
	require.NoError(t, err)
	assert.Equal(t, "Test User", info.UserName)
	assert.Equal(t, "test@example.com", info.UserEmail)
	assert.Equal(t, "user-test", info.UserID)
	assert.Equal(t, "session-test", info.SessionID)
}

// Test session destroy functionality
func TestSessionDestroy(t *testing.T) {
	tempDir := t.TempDir()
	sessionPath := filepath.Join(tempDir, "destroy-session")
	
	// Create session with metadata
	session, err := git.NewSession(sessionPath,
		git.WithSessionUser("Test User", "test@example.com"),
		git.WithInstruqtMetadata("user-destroy", "session-destroy", time.Now()),
		git.WithMetadata("cleanup", "test"),
	)
	require.NoError(t, err)
	
	// Make a commit so session is fully established
	testFile := filepath.Join(sessionPath, "destroy-test.txt")
	err = os.WriteFile(testFile, []byte("will be cleaned"), 0644)
	require.NoError(t, err)
	
	err = session.Add([]string{"destroy-test.txt"})
	require.NoError(t, err)
	err = session.Commit("Commit before destroy")
	require.NoError(t, err)
	
	// Verify session is valid before destroy
	assert.True(t, session.IsValid())
	
	// Destroy session-specific configuration
	err = session.Destroy()
	assert.NoError(t, err)
	
	// Repository should still exist but session config should be cleaned
	assert.DirExists(t, filepath.Join(sessionPath, ".git"))
	
	// Session should still be valid (repository exists) but metadata is gone
	assert.True(t, session.IsValid())
}