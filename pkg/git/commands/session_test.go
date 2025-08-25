package commands

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	gitpkg "github.com/instruqt/git-exec/pkg/git"
	"github.com/stretchr/testify/require"
)

func TestNewSessionWithConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-session-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	sessionPath := filepath.Join(tempDir, "test-session")
	
	t.Run("create new session with user config", func(t *testing.T) {
		session, err := NewSession(sessionPath,
			gitpkg.WithUser("Test User", "test@example.com"),
			gitpkg.WithInstruqtMetadata("user-123", "session-456", time.Now()),
		)
		require.NoError(t, err)
		require.NotNil(t, session)
		
		// Verify configuration
		config := session.GetConfig()
		require.Equal(t, "Test User", config.UserName)
		require.Equal(t, "test@example.com", config.UserEmail)
		require.Equal(t, "user-123", config.UserID)
		require.Equal(t, "session-456", config.SessionID)
		require.Equal(t, sessionPath, config.WorkingDirectory)
		
		// Verify session is valid
		require.True(t, session.IsValid())
	})
}

func TestLoadSession(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-session-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	sessionPath := filepath.Join(tempDir, "test-session")
	
	t.Run("load existing session", func(t *testing.T) {
		// First create a session
		created := time.Now().UTC().Truncate(time.Second)
		originalSession, err := NewSession(sessionPath,
			gitpkg.WithUser("Test User", "test@example.com"),
			gitpkg.WithInstruqtMetadata("user-123", "session-456", created),
			gitpkg.WithMetadata("environment", "test"),
			gitpkg.WithMetadata("track", "git-basics"),
		)
		require.NoError(t, err)
		
		// Load the session
		loadedSession, err := LoadSession(sessionPath)
		require.NoError(t, err)
		require.NotNil(t, loadedSession)
		
		// Verify loaded configuration matches original
		config := loadedSession.GetConfig()
		require.Equal(t, "Test User", config.UserName)
		require.Equal(t, "test@example.com", config.UserEmail)
		require.Equal(t, "user-123", config.UserID)
		require.Equal(t, "session-456", config.SessionID)
		require.Equal(t, "test", config.Metadata["environment"])
		require.Equal(t, "git-basics", config.Metadata["track"])
		
		// Verify session IDs match
		require.Equal(t, originalSession.GetSessionID(), loadedSession.GetSessionID())
		require.Equal(t, originalSession.GetUserID(), loadedSession.GetUserID())
	})
	
	t.Run("load non-existent session fails", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "non-existent")
		_, err := LoadSession(nonExistentPath)
		require.Error(t, err)
	})
}

func TestSessionUpdateUser(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-session-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	sessionPath := filepath.Join(tempDir, "test-session")
	
	t.Run("update user information", func(t *testing.T) {
		session, err := NewSession(sessionPath,
			gitpkg.WithUser("Original User", "original@example.com"),
		)
		require.NoError(t, err)
		
		// Update user
		err = session.UpdateUser("Updated User", "updated@example.com")
		require.NoError(t, err)
		
		// Verify update in memory
		config := session.GetConfig()
		require.Equal(t, "Updated User", config.UserName)
		require.Equal(t, "updated@example.com", config.UserEmail)
		
		// Load session again to verify persistence
		loadedSession, err := LoadSession(sessionPath)
		require.NoError(t, err)
		
		loadedConfig := loadedSession.GetConfig()
		require.Equal(t, "Updated User", loadedConfig.UserName)
		require.Equal(t, "updated@example.com", loadedConfig.UserEmail)
	})
}

func TestSessionDestroy(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-session-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	sessionPath := filepath.Join(tempDir, "test-session")
	
	t.Run("destroy removes session config", func(t *testing.T) {
		session, err := NewSession(sessionPath,
			gitpkg.WithUser("Test User", "test@example.com"),
			gitpkg.WithInstruqtMetadata("user-123", "session-456", time.Now()),
			gitpkg.WithMetadata("environment", "test"),
		)
		require.NoError(t, err)
		
		// Destroy session
		err = session.Destroy()
		require.NoError(t, err)
		
		// Load session and verify Instruqt config is gone
		loadedSession, err := LoadSession(sessionPath)
		require.NoError(t, err)
		
		config := loadedSession.GetConfig()
		require.Empty(t, config.UserID)
		require.Empty(t, config.SessionID)
		require.Empty(t, config.Metadata["environment"])
		
		// User config should still be there (it's standard git config)
		require.Equal(t, "Test User", config.UserName)
		require.Equal(t, "test@example.com", config.UserEmail)
	})
}

func TestCloneIntoSession(t *testing.T) {
	// This test requires a real git repository to clone from
	// We'll use the current repository for testing
	t.Skip("Skipping clone test - requires real repository")
	
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-session-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	sessionPath := filepath.Join(tempDir, "cloned-session")
	
	t.Run("clone repository into session", func(t *testing.T) {
		// Get current directory as source
		cwd, err := os.Getwd()
		require.NoError(t, err)
		
		// Create session and clone into it
		session, err := NewSession(sessionPath,
			gitpkg.WithUser("Clone User", "clone@example.com"),
			gitpkg.WithInstruqtMetadata("user-789", "session-012", time.Now()),
		)
		require.NoError(t, err)
		
		err = session.Clone(cwd, sessionPath)
		require.NoError(t, err)
		require.NotNil(t, session)
		
		// Verify session configuration
		config := session.GetConfig()
		require.Equal(t, "Clone User", config.UserName)
		require.Equal(t, "clone@example.com", config.UserEmail)
		require.Equal(t, "user-789", config.UserID)
		require.Equal(t, "session-012", config.SessionID)
		
		// Verify it's a valid git repository
		require.True(t, session.IsValid())
	})
}

func TestValidateSession(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-session-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	t.Run("validate valid session", func(t *testing.T) {
		sessionPath := filepath.Join(tempDir, "valid-session")
		
		// Create a valid session
		_, err := NewSession(sessionPath,
			gitpkg.WithInstruqtMetadata("user-123", "session-456", time.Now()),
		)
		require.NoError(t, err)
		
		// Validate it
		err = ValidateSession(sessionPath)
		require.NoError(t, err)
	})
	
	t.Run("validate session without ID fails", func(t *testing.T) {
		sessionPath := filepath.Join(tempDir, "invalid-session")
		
		// Create a session without session ID
		_, err := NewSession(sessionPath,
			gitpkg.WithUser("Test User", "test@example.com"),
		)
		require.NoError(t, err)
		
		// Validation should fail
		err = ValidateSession(sessionPath)
		require.Error(t, err)
		require.Contains(t, err.Error(), "session ID is missing")
	})
	
	t.Run("validate non-existent session fails", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "non-existent")
		err := ValidateSession(nonExistentPath)
		require.Error(t, err)
	})
}

func TestGetSessionInfo(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-session-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	sessionPath := filepath.Join(tempDir, "test-session")
	
	t.Run("get session info", func(t *testing.T) {
		// Create a session
		created := time.Now().UTC().Truncate(time.Second)
		_, err := NewSession(sessionPath,
			gitpkg.WithUser("Info User", "info@example.com"),
			gitpkg.WithInstruqtMetadata("user-999", "session-888", created),
			gitpkg.WithMetadata("key1", "value1"),
			gitpkg.WithMetadata("key2", "value2"),
		)
		require.NoError(t, err)
		
		// Get session info
		info, err := GetSessionInfo(sessionPath)
		require.NoError(t, err)
		require.NotNil(t, info)
		
		// Verify info
		require.Equal(t, "Info User", info.UserName)
		require.Equal(t, "info@example.com", info.UserEmail)
		require.Equal(t, "user-999", info.UserID)
		require.Equal(t, "session-888", info.SessionID)
		require.Equal(t, "value1", info.Metadata["key1"])
		require.Equal(t, "value2", info.Metadata["key2"])
	})
}

func TestSessionCommitWithUserContext(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-session-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	sessionPath := filepath.Join(tempDir, "test-session")
	
	t.Run("commit uses session user context", func(t *testing.T) {
		session, err := NewSession(sessionPath,
			gitpkg.WithUser("Commit User", "commit@example.com"),
			gitpkg.WithInstruqtMetadata("user-123", "session-456", time.Now()),
		)
		require.NoError(t, err)
		
		// Create a test file
		testFile := filepath.Join(sessionPath, "test.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0644)
		require.NoError(t, err)
		
		// Add the file
		err = session.Add([]string{"test.txt"})
		require.NoError(t, err)
		
		// Commit should automatically use session user context
		err = session.Commit("Test commit")
		require.NoError(t, err)
		
		// Verify commit author
		logs, err := session.Log(LogWithMaxCount("1"))
		require.NoError(t, err)
		require.Len(t, logs, 1)
		
		// Check author information (Author is stored as "Name <email>" format)
		require.Contains(t, logs[0].Author, "Commit User")
		require.Contains(t, logs[0].Author, "commit@example.com")
	})
}

func TestSessionMetadata(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-session-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	sessionPath := filepath.Join(tempDir, "test-session")
	
	t.Run("multiple metadata entries", func(t *testing.T) {
		session, err := NewSession(sessionPath,
			gitpkg.WithMetadata("track", "git-advanced"),
			gitpkg.WithMetadata("environment", "production"),
			gitpkg.WithMetadata("version", "1.2.3"),
			gitpkg.WithMetadata("feature", "session-management"),
		)
		require.NoError(t, err)
		
		// Verify all metadata is stored
		config := session.GetConfig()
		require.Equal(t, "git-advanced", config.Metadata["track"])
		require.Equal(t, "production", config.Metadata["environment"])
		require.Equal(t, "1.2.3", config.Metadata["version"])
		require.Equal(t, "session-management", config.Metadata["feature"])
		
		// Load session and verify persistence
		loadedSession, err := LoadSession(sessionPath)
		require.NoError(t, err)
		
		loadedConfig := loadedSession.GetConfig()
		require.Equal(t, "git-advanced", loadedConfig.Metadata["track"])
		require.Equal(t, "production", loadedConfig.Metadata["environment"])
		require.Equal(t, "1.2.3", loadedConfig.Metadata["version"])
		require.Equal(t, "session-management", loadedConfig.Metadata["feature"])
	})
}