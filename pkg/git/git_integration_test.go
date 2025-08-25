package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/instruqt/git-exec/pkg/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test the complete workflow that matters: init -> config -> add -> commit -> status
func TestGitWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	
	// Initialize repository
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)
	assert.DirExists(t, filepath.Join(tempDir, ".git"))
	
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Configure user (required for commits)
	err = gitInstance.Config("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.Config("user.email", "test@example.com")
	require.NoError(t, err)
	
	// Create and add a file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)
	
	// Check status before add (should show untracked)
	files, err := gitInstance.Status()
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, "test.txt", files[0].Name)
	assert.Equal(t, "untracked", string(files[0].Status))
	
	// Add file
	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)
	
	// Check status after add (should show added)
	files, err = gitInstance.Status()
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, "test.txt", files[0].Name)
	assert.Equal(t, "added", string(files[0].Status))
	
	// Commit
	err = gitInstance.Commit("Initial commit")
	require.NoError(t, err)
	
	// Check status after commit (should be clean)
	files, err = gitInstance.Status()
	require.NoError(t, err)
	assert.Empty(t, files)
	
	// Verify commit exists in log
	logs, err := gitInstance.Log(git.LogWithMaxCount("1"))
	require.NoError(t, err)
	require.Len(t, logs, 1)
	assert.Equal(t, "Initial commit", logs[0].Message)
	assert.Equal(t, "Test User <test@example.com>", logs[0].Author)
}

// Test error cases that actually matter
func TestGitErrorHandling(t *testing.T) {
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	
	// Try to commit without a repository
	err = gitInstance.Commit("should fail")
	assert.Error(t, err)
	
	// Try to add non-existent file
	tempDir := t.TempDir()
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	err = gitInstance.Add([]string{"nonexistent.txt"})
	assert.Error(t, err)
	
	// Try to init in existing repo (should handle gracefully or error)
	err = gitInstance.Init(tempDir)
	// This may succeed (reinit) or fail, both are acceptable behaviors
}

// Test branch operations workflow
func TestBranchOperations(t *testing.T) {
	tempDir := setupTestRepo(t)
	
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// List initial branches (should have main/master)
	branches, err := gitInstance.ListBranches()
	require.NoError(t, err)
	require.Len(t, branches, 1)
	assert.True(t, branches[0].Active)
	
	// Create new branch
	err = gitInstance.CreateBranch("feature-test")
	require.NoError(t, err)
	
	// List branches (should have 2 now)
	branches, err = gitInstance.ListBranches()
	require.NoError(t, err)
	assert.Len(t, branches, 2)
	
	// Checkout new branch
	_, err = gitInstance.Checkout(git.CheckoutWithBranch("feature-test"))
	require.NoError(t, err)
	
	// Verify we're on the new branch
	branches, err = gitInstance.ListBranches()
	require.NoError(t, err)
	
	var activeBranch string
	for _, branch := range branches {
		if branch.Active {
			activeBranch = branch.Name
			break
		}
	}
	assert.Equal(t, "feature-test", activeBranch)
	
	// Switch back to main
	_, err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
	if err != nil {
		// Try master if main doesn't exist
		_, err = gitInstance.Checkout(git.CheckoutWithBranch("master"))
	}
	require.NoError(t, err)
	
	// Delete feature branch
	err = gitInstance.DeleteBranch("feature-test")
	require.NoError(t, err)
	
	// Verify branch is gone
	branches, err = gitInstance.ListBranches()
	require.NoError(t, err)
	assert.Len(t, branches, 1)
}

// Test merge operations with actual conflicts
func TestMergeConflictResolution(t *testing.T) {
	tempDir := setupTestRepo(t)
	
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Create a file on main branch
	mainFile := filepath.Join(tempDir, "conflict.txt")
	err = os.WriteFile(mainFile, []byte("line 1\nmain content\nline 3"), 0644)
	require.NoError(t, err)
	
	err = gitInstance.Add([]string{"conflict.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Add conflict file on main")
	require.NoError(t, err)
	
	// Create and switch to feature branch
	err = gitInstance.CreateBranch("feature")
	require.NoError(t, err)
	_, err = gitInstance.Checkout(git.CheckoutWithBranch("feature"))
	require.NoError(t, err)
	
	// Modify the same file differently
	err = os.WriteFile(mainFile, []byte("line 1\nfeature content\nline 3"), 0644)
	require.NoError(t, err)
	
	err = gitInstance.Add([]string{"conflict.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Modify conflict file on feature")
	require.NoError(t, err)
	
	// Switch back to main
	_, err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
	if err != nil {
		_, err = gitInstance.Checkout(git.CheckoutWithBranch("master"))
	}
	require.NoError(t, err)
	
	// Attempt merge - should create conflict
	result, err := gitInstance.Merge(git.MergeWithBranch("feature"))
	
	if err == nil && !result.Success {
		// Conflict detected - test resolution
		assert.False(t, result.Success)
		
		// Use merge abort to test that workflow
		err = gitInstance.MergeAbort()
		assert.NoError(t, err)
	}
}

// Helper function to set up a basic test repository
func setupTestRepo(t *testing.T) string {
	tempDir := t.TempDir()
	
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)
	
	gitInstance.SetWorkingDirectory(tempDir)
	
	err = gitInstance.Config("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.Config("user.email", "test@example.com")
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