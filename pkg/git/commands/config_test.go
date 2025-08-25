package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	repoPath := filepath.Join(tempDir, "test-repo")
	
	t.Run("use WithConfig for temporary user config", func(t *testing.T) {
		g, err := NewGit()
		require.NoError(t, err)
		
		// Initialize repository
		err = g.Init(repoPath)
		require.NoError(t, err)
		
		// Set working directory
		g.SetWorkingDirectory(repoPath)
		
		// Create a test file
		testFile := filepath.Join(repoPath, "test.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0644)
		require.NoError(t, err)
		
		// Add the file
		err = g.Add([]string{"test.txt"})
		require.NoError(t, err)
		
		// Commit with temporary user config using WithConfig
		err = g.Commit("Test commit",
			WithConfig("user.name", "Temporary User"),
			WithConfig("user.email", "temp@example.com"),
		)
		require.NoError(t, err)
		
		// Verify commit was made with the temporary config
		logs, err := g.Log(LogWithMaxCount("1"))
		require.NoError(t, err)
		require.Len(t, logs, 1)
		
		// Check that author contains the temporary user info
		require.Contains(t, logs[0].Author, "Temporary User")
		require.Contains(t, logs[0].Author, "temp@example.com")
	})
	
	t.Run("use WithConfigs for multiple config values", func(t *testing.T) {
		g, err := NewGit()
		require.NoError(t, err)
		
		repoPath2 := filepath.Join(tempDir, "test-repo2")
		err = g.Init(repoPath2)
		require.NoError(t, err)
		
		g.SetWorkingDirectory(repoPath2)
		
		// Create a test file
		testFile := filepath.Join(repoPath2, "test.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0644)
		require.NoError(t, err)
		
		// Add the file
		err = g.Add([]string{"test.txt"})
		require.NoError(t, err)
		
		// Commit with multiple config values
		configs := map[string]string{
			"user.name":  "Multi Config User",
			"user.email": "multi@example.com",
		}
		
		err = g.Commit("Test commit with multiple configs",
			WithConfigs(configs),
		)
		require.NoError(t, err)
		
		// Verify commit was made with the config
		logs, err := g.Log(LogWithMaxCount("1"))
		require.NoError(t, err)
		require.Len(t, logs, 1)
		
		// Check that author contains the configured user info
		require.Contains(t, logs[0].Author, "Multi Config User")
		require.Contains(t, logs[0].Author, "multi@example.com")
	})
}