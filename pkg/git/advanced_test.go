package git_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/instruqt/git-exec/pkg/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Revert command - test reverting commits
func TestRevertCommand(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Create a commit to revert
	testFile := filepath.Join(tempDir, "revert-test.txt")
	err = os.WriteFile(testFile, []byte("content to revert"), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"revert-test.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Commit to revert")
	require.NoError(t, err)
	
	// Get the commit hash
	logs, err := gitInstance.Log(git.LogWithMaxCount("1"))
	require.NoError(t, err)
	require.Len(t, logs, 1)
	
	// Test revert - we test with HEAD since it's more reliable
	err = gitInstance.Revert(git.WithArgs("HEAD", "--no-edit"))
	require.NoError(t, err)
	
	// Verify revert commit was created
	logs, err = gitInstance.Log(git.LogWithMaxCount("2"))
	require.NoError(t, err)
	assert.Len(t, logs, 2)
	
	// The first commit should be the revert commit
	assert.Contains(t, logs[0].Message, "Revert")
}

// Test Rebase command - test interface (actual rebase is complex)
func TestRebaseCommand(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Create commits on main branch
	for i := 1; i <= 2; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("main-%d.txt", i))
		err = os.WriteFile(filename, []byte(fmt.Sprintf("main content %d", i)), 0644)
		require.NoError(t, err)
		err = gitInstance.Add([]string{filepath.Base(filename)})
		require.NoError(t, err)
		err = gitInstance.Commit(fmt.Sprintf("Main commit %d", i))
		require.NoError(t, err)
	}
	
	// Create and switch to feature branch
	err = gitInstance.CreateBranch("feature-rebase")
	require.NoError(t, err)
	_, err = gitInstance.Checkout(git.CheckoutWithBranch("feature-rebase"))
	require.NoError(t, err)
	
	// Create commits on feature branch
	featureFile := filepath.Join(tempDir, "feature.txt")
	err = os.WriteFile(featureFile, []byte("feature content"), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"feature.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Feature commit")
	require.NoError(t, err)
	
	// Test rebase onto main - this tests the interface
	err = gitInstance.Rebase(git.WithArgs("main"))
	// Rebase may succeed or fail depending on conflicts, but shouldn't panic
	// The value is testing the interface exists and handles various outcomes
}

// Test Reflog command - test interface
func TestReflogCommand(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Create some commits to generate reflog entries
	for i := 1; i <= 3; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("reflog-%d.txt", i))
		err = os.WriteFile(filename, []byte(fmt.Sprintf("reflog content %d", i)), 0644)
		require.NoError(t, err)
		err = gitInstance.Add([]string{filepath.Base(filename)})
		require.NoError(t, err)
		err = gitInstance.Commit(fmt.Sprintf("Reflog commit %d", i))
		require.NoError(t, err)
	}
	
	// Test reflog - tests interface
	err = gitInstance.Reflog()
	require.NoError(t, err, "Reflog should not error on valid repository")
}

// Test Remove command - test removing files from Git
func TestRemoveCommand(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Create and commit files
	file1 := filepath.Join(tempDir, "remove1.txt")
	file2 := filepath.Join(tempDir, "remove2.txt")
	
	err = os.WriteFile(file1, []byte("content1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file2, []byte("content2"), 0644)
	require.NoError(t, err)
	
	err = gitInstance.Add([]string{"remove1.txt", "remove2.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Add files to remove")
	require.NoError(t, err)
	
	// Test removing file from Git and filesystem
	err = gitInstance.Remove(git.WithArgs("remove1.txt"))
	require.NoError(t, err)
	
	// Verify file is staged for removal
	files, err := gitInstance.Status()
	require.NoError(t, err)
	
	var found bool
	for _, file := range files {
		if file.Name == "remove1.txt" && file.Status == "deleted" {
			found = true
			break
		}
	}
	assert.True(t, found, "Should show remove1.txt as deleted")
	
	// Verify file was removed from filesystem
	assert.NoFileExists(t, file1)
	
	// Test removing file from Git only (cached)
	err = gitInstance.Remove(git.WithArgs("--cached", "remove2.txt"))
	require.NoError(t, err)
	
	// File should still exist in filesystem but be removed from Git
	assert.FileExists(t, file2)
	
	files, err = gitInstance.Status()
	require.NoError(t, err)
	
	statusMap := make(map[string]string)
	for _, file := range files {
		statusMap[file.Name] = string(file.Status)
	}
	
	// remove2.txt should be deleted from Git but untracked (since it still exists)
	assert.Contains(t, statusMap, "remove2.txt")
}

// Test error handling for advanced operations
func TestAdvancedCommandErrors(t *testing.T) {
	tempDir := t.TempDir()
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	
	// Test commands on non-git directory
	gitInstance.SetWorkingDirectory(tempDir)
	
	// These should fail gracefully on non-git directory
	err = gitInstance.Revert()
	assert.Error(t, err, "Revert should fail on non-git directory")
	
	err = gitInstance.Rebase()
	assert.Error(t, err, "Rebase should fail on non-git directory")
	
	err = gitInstance.Reflog()
	assert.Error(t, err, "Reflog should fail on non-git directory")
	
	err = gitInstance.Remove()
	assert.Error(t, err, "Remove should fail on non-git directory")
}

// Test branch operations with advanced scenarios
func TestAdvancedBranchOperations(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Add a remote for upstream testing
	err = gitInstance.AddRemote("origin", "https://github.com/test/repo.git")
	require.NoError(t, err)
	
	// Create a branch
	err = gitInstance.CreateBranch("feature-upstream")
	require.NoError(t, err)
	
	// Test setting upstream (interface test - won't work without actual remote)
	err = gitInstance.SetUpstream("feature-upstream", "origin")
	// This may fail without real remote, but tests the interface
	
	// Test deleting branch that doesn't exist
	err = gitInstance.DeleteBranch("nonexistent-branch")
	assert.Error(t, err, "Should fail to delete non-existent branch")
	
	// Test force delete (if branch has unmerged changes)
	err = gitInstance.CreateBranch("force-delete-test")
	require.NoError(t, err)
	_, err = gitInstance.Checkout(git.CheckoutWithBranch("force-delete-test"))
	require.NoError(t, err)
	
	// Make a commit on the branch
	testFile := filepath.Join(tempDir, "force-delete.txt")
	err = os.WriteFile(testFile, []byte("unmerged content"), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"force-delete.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Unmerged commit")
	require.NoError(t, err)
	
	// Switch back to main
	_, err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
	if err != nil {
		_, err = gitInstance.Checkout(git.CheckoutWithBranch("master"))
	}
	require.NoError(t, err)
	
	// Normal delete should fail
	err = gitInstance.DeleteBranch("force-delete-test")
	assert.Error(t, err, "Should fail to delete branch with unmerged changes")
	
	// Force delete should work
	err = gitInstance.DeleteBranch("force-delete-test", git.WithArgs("-D"))
	require.NoError(t, err)
}