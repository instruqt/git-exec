package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// setupConfigTestRepo creates a test repository for config tests
func setupConfigTestRepo(t *testing.T) string {
	tempDir := t.TempDir()
	
	gitInstance, err := NewGit()
	require.NoError(t, err)
	
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)
	
	gitInstance.SetWorkingDirectory(tempDir)
	
	err = gitInstance.SetConfig("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.email", "test@example.com")
	require.NoError(t, err)
	
	// Create initial commit
	initialFile := filepath.Join(tempDir, "README.md")
	err = os.WriteFile(initialFile, []byte("# Test Repo"), 0644)
	require.NoError(t, err)
	
	err = gitInstance.Add([]string{"README.md"})
	require.NoError(t, err)
	err = gitInstance.Commit("Initial commit")
	require.NoError(t, err)
	
	return tempDir
}

func TestConfigCommands(t *testing.T) {
	tempDir := setupConfigTestRepo(t)
	gitInstance, err := NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)

	// Test SetConfig and GetConfig
	t.Run("SetConfig and GetConfig", func(t *testing.T) {
		// Set a test configuration
		err := gitInstance.SetConfig("test.key", "test-value")
		require.NoError(t, err)

		// Get the configuration back
		value, err := gitInstance.GetConfig("test.key")
		require.NoError(t, err)
		assert.Equal(t, "test-value", value)

		// Set user config for later tests
		err = gitInstance.SetConfig("user.name", "Config Test User")
		require.NoError(t, err)
		err = gitInstance.SetConfig("user.email", "config@test.com")
		require.NoError(t, err)
	})

	// Test GetConfig with non-existent key
	t.Run("GetConfig non-existent key", func(t *testing.T) {
		_, err := gitInstance.GetConfig("non.existent.key")
		require.Error(t, err)
	})

	// Test UnsetConfig
	t.Run("UnsetConfig", func(t *testing.T) {
		// Set a config to unset
		err := gitInstance.SetConfig("test.unset", "value-to-unset")
		require.NoError(t, err)

		// Verify it exists
		value, err := gitInstance.GetConfig("test.unset")
		require.NoError(t, err)
		assert.Equal(t, "value-to-unset", value)

		// Unset it
		err = gitInstance.UnsetConfig("test.unset")
		require.NoError(t, err)

		// Verify it's gone
		_, err = gitInstance.GetConfig("test.unset")
		require.Error(t, err)
	})
}

func TestConfigListCommand(t *testing.T) {
	tempDir := setupConfigTestRepo(t)
	gitInstance, err := NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)

	// Set configs for this test
	err = gitInstance.SetConfig("test.key", "test-value")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.name", "Config Test User")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.email", "config@test.com")
	require.NoError(t, err)

	// Test ListConfig
	configs, err := gitInstance.ListConfig()
	require.NoError(t, err)
	require.NotEmpty(t, configs)

	// Should have at least our test configs
	var foundTestKey, foundLocalUserName, foundLocalUserEmail bool
	for _, config := range configs {
		switch config.Key {
		case "test.key":
			assert.Equal(t, "test-value", config.Value)
			assert.Equal(t, types.ConfigScopeLocal, config.Scope)
			foundTestKey = true
		case "user.name":
			// Check for the local scope config we set
			if config.Scope == types.ConfigScopeLocal && config.Value == "Config Test User" {
				foundLocalUserName = true
			}
		case "user.email":
			// Check for the local scope config we set
			if config.Scope == types.ConfigScopeLocal && config.Value == "config@test.com" {
				foundLocalUserEmail = true
			}
		}
	}

	assert.True(t, foundTestKey, "Should find test.key in config list")
	assert.True(t, foundLocalUserName, "Should find local user.name='Config Test User' in config list")
	assert.True(t, foundLocalUserEmail, "Should find local user.email='config@test.com' in config list")
}

func TestConfigOptions(t *testing.T) {
	tempDir := setupConfigTestRepo(t)
	gitInstance, err := NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)

	// Test scope options
	t.Run("ConfigWithLocalScope", func(t *testing.T) {
		err := gitInstance.SetConfig("test.local", "local-value", ConfigWithLocalScope())
		require.NoError(t, err)

		value, err := gitInstance.GetConfig("test.local", ConfigWithLocalScope())
		require.NoError(t, err)
		assert.Equal(t, "local-value", value)
	})

	t.Run("ConfigWithGlobalScope", func(t *testing.T) {
		// Set global config (might fail in test environment, that's ok)
		err := gitInstance.SetConfig("test.global", "global-value", ConfigWithGlobalScope())
		// Note: This might fail in CI environments without global git config access
		if err == nil {
			value, err := gitInstance.GetConfig("test.global", ConfigWithGlobalScope())
			if err == nil {
				assert.Equal(t, "global-value", value)
			}
		}
	})

	t.Run("ListConfig basic functionality", func(t *testing.T) {
		// Set a test config first
		err := gitInstance.SetConfig("test.options", "options-value")
		require.NoError(t, err)

		configs, err := gitInstance.ListConfig()
		require.NoError(t, err)
		require.NotEmpty(t, configs)

		// Should have configs from various scopes
		foundScopes := make(map[types.ConfigScope]bool)
		var foundTestConfig bool
		for _, config := range configs {
			foundScopes[config.Scope] = true
			if config.Key == "test.options" {
				assert.Equal(t, "options-value", config.Value)
				assert.NotEmpty(t, config.Source, "Should have source file information")
				foundTestConfig = true
			}
		}

		// Should at least have local scope
		assert.True(t, foundScopes[types.ConfigScopeLocal], "Should find local scope configs")
		assert.True(t, foundTestConfig, "Should find our test config")
	})
}

func TestConfigParsing(t *testing.T) {
	// Test parseConfigList function directly
	t.Run("parseConfigList", func(t *testing.T) {
		output := `local	file:/path/to/.git/config	user.name=John Doe
local	file:/path/to/.git/config	user.email=john@example.com
global	file:/home/user/.gitconfig	core.editor=vim
system	file:/etc/gitconfig	init.defaultBranch=main`

		entries, err := parseConfigList(output)
		require.NoError(t, err)
		require.Len(t, entries, 4)

		// Check first entry
		assert.Equal(t, "user.name", entries[0].Key)
		assert.Equal(t, "John Doe", entries[0].Value)
		assert.Equal(t, types.ConfigScopeLocal, entries[0].Scope)
		assert.Equal(t, "/path/to/.git/config", entries[0].Source)

		// Check global entry
		assert.Equal(t, "core.editor", entries[2].Key)
		assert.Equal(t, "vim", entries[2].Value)
		assert.Equal(t, types.ConfigScopeGlobal, entries[2].Scope)
		assert.Equal(t, "/home/user/.gitconfig", entries[2].Source)

		// Check system entry
		assert.Equal(t, "init.defaultBranch", entries[3].Key)
		assert.Equal(t, "main", entries[3].Value)
		assert.Equal(t, types.ConfigScopeSystem, entries[3].Scope)
		assert.Equal(t, "/etc/gitconfig", entries[3].Source)
	})

	t.Run("parseConfigList with malformed lines", func(t *testing.T) {
		output := `local	file:/path/to/.git/config	user.name=John Doe
malformed line without proper format
local	file:/path/to/.git/config	user.email=john@example.com`

		entries, err := parseConfigList(output)
		require.NoError(t, err)
		require.Len(t, entries, 2) // Should skip malformed line

		assert.Equal(t, "user.name", entries[0].Key)
		assert.Equal(t, "user.email", entries[1].Key)
	})

	t.Run("parseConfigList empty output", func(t *testing.T) {
		entries, err := parseConfigList("")
		require.NoError(t, err)
		require.Empty(t, entries)
	})
}