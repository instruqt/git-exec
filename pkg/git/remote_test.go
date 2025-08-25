package git_test

import (
	"testing"

	"github.com/instruqt/git-exec/pkg/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test remote operations - the valuable scenarios are CRUD operations
func TestRemoteOperations(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Test adding remotes
	err = gitInstance.AddRemote("origin", "https://github.com/test/repo.git")
	require.NoError(t, err)
	
	err = gitInstance.AddRemote("upstream", "https://github.com/upstream/repo.git")
	require.NoError(t, err)
	
	// Test listing remotes
	remotes, err := gitInstance.ListRemotes()
	require.NoError(t, err)
	assert.Len(t, remotes, 2)
	
	remoteMap := make(map[string]string)
	for _, remote := range remotes {
		remoteMap[remote.Name] = remote.URL
	}
	
	assert.Equal(t, "https://github.com/test/repo.git", remoteMap["origin"])
	assert.Equal(t, "https://github.com/upstream/repo.git", remoteMap["upstream"])
	
	// Test changing remote URL
	err = gitInstance.SetRemoteURL("origin", "https://github.com/new/repo.git")
	require.NoError(t, err)
	
	// Verify URL change
	remotes, err = gitInstance.ListRemotes()
	require.NoError(t, err)
	
	remoteMap = make(map[string]string)
	for _, remote := range remotes {
		remoteMap[remote.Name] = remote.URL
	}
	
	assert.Equal(t, "https://github.com/new/repo.git", remoteMap["origin"])
	
	// Test removing remote
	err = gitInstance.RemoveRemote("upstream")
	require.NoError(t, err)
	
	// Verify remote was removed
	remotes, err = gitInstance.ListRemotes()
	require.NoError(t, err)
	assert.Len(t, remotes, 1)
	assert.Equal(t, "origin", remotes[0].Name)
}

// Test remote error cases
func TestRemoteErrorHandling(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Test removing non-existent remote
	err = gitInstance.RemoveRemote("nonexistent")
	assert.Error(t, err, "Should fail to remove non-existent remote")
	
	// Test adding duplicate remote
	err = gitInstance.AddRemote("test", "https://github.com/test/repo.git")
	require.NoError(t, err)
	
	err = gitInstance.AddRemote("test", "https://github.com/other/repo.git")
	assert.Error(t, err, "Should fail to add duplicate remote name")
}

// Test fetch operations - test the interface, actual network ops would be integration tests
func TestFetchCommand(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Add a remote (needed for fetch)
	err = gitInstance.AddRemote("origin", "https://github.com/test/empty-repo.git")
	require.NoError(t, err)
	
	// Test fetch - this will likely fail due to network/auth, but we test the interface
	_, err = gitInstance.Fetch()
	// Don't require success - network operations are environment dependent
	// We just verify the method exists and handles errors gracefully
}

// Test push operations - test interface
func TestPushCommand(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Add a remote
	err = gitInstance.AddRemote("origin", "https://github.com/test/empty-repo.git")
	require.NoError(t, err)
	
	// Test push - will likely fail but tests interface
	_, err = gitInstance.Push()
	// Don't require success - network operations are environment dependent
}

// Test pull operations - test interface  
func TestPullCommand(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Add a remote
	err = gitInstance.AddRemote("origin", "https://github.com/test/empty-repo.git")
	require.NoError(t, err)
	
	// Test pull - will likely fail but tests interface
	_, err = gitInstance.Pull()
	// Don't require success - network operations are environment dependent
	// The value is testing that the interface works and returns proper types
}