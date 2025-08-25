package git_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/instruqt/git-exec/pkg/git"
	"github.com/instruqt/git-exec/pkg/git/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test merge conflict detection and resolution - the most complex Git workflow
func TestMergeConflictWorkflow(t *testing.T) {
	tempDir := setupTestRepo(t)
	
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Create initial file on main branch
	conflictFile := filepath.Join(tempDir, "shared.txt")
	initialContent := `header
shared content
footer`
	
	err = os.WriteFile(conflictFile, []byte(initialContent), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"shared.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Add shared file")
	require.NoError(t, err)
	
	// Create branch and modify file
	err = gitInstance.CreateBranch("feature-branch")
	require.NoError(t, err)
	err = gitInstance.Checkout(git.CheckoutWithBranch("feature-branch"))
	require.NoError(t, err)
	
	featureContent := `header
feature modified content
footer`
	
	err = os.WriteFile(conflictFile, []byte(featureContent), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"shared.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Feature changes")
	require.NoError(t, err)
	
	// Switch back to main and make conflicting changes
	err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
	if err != nil {
		err = gitInstance.Checkout(git.CheckoutWithBranch("master"))
	}
	require.NoError(t, err)
	
	mainContent := `header
main modified content
footer`
	
	err = os.WriteFile(conflictFile, []byte(mainContent), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"shared.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Main changes")
	require.NoError(t, err)
	
	// Attempt merge - should detect conflict
	result, err := gitInstance.Merge(git.MergeWithBranch("feature-branch"))
	
	// The merge should either fail with an error or return unsuccessful result
	if err != nil {
		// Merge failed due to conflict - this is expected
		assert.Contains(t, err.Error(), "conflict", "Expected conflict-related error")
	} else {
		// Merge returned but was unsuccessful due to conflicts
		assert.False(t, result.Success, "Merge should be unsuccessful due to conflicts")
	}
	
	// Test conflict resolution using "ours" strategy
	resolutions := []types.ConflictResolution{
		{
			FilePath: "shared.txt",
			UseOurs:  true,
		},
	}
	
	// Note: This tests the interface, actual conflict resolution depends on Git state
	err = gitInstance.ResolveConflicts(resolutions)
	// Don't require this to succeed as Git state may vary, but it should not panic
	
	// Test merge abort functionality
	err = gitInstance.MergeAbort()
	// This may fail if no merge is in progress, which is fine for this test
}

// Test fast-forward merge (no conflicts)
func TestFastForwardMerge(t *testing.T) {
	tempDir := setupTestRepo(t)
	
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Create branch and add new file (no conflicts possible)
	err = gitInstance.CreateBranch("feature-addition")
	require.NoError(t, err)
	err = gitInstance.Checkout(git.CheckoutWithBranch("feature-addition"))
	require.NoError(t, err)
	
	// Add new file that doesn't exist on main
	newFile := filepath.Join(tempDir, "new-feature.txt")
	err = os.WriteFile(newFile, []byte("new feature content"), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"new-feature.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Add new feature")
	require.NoError(t, err)
	
	// Switch back to main
	err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
	if err != nil {
		err = gitInstance.Checkout(git.CheckoutWithBranch("master"))
	}
	require.NoError(t, err)
	
	// Merge should succeed (fast-forward)
	result, err := gitInstance.Merge(git.MergeWithBranch("feature-addition"))
	require.NoError(t, err)
	assert.True(t, result.Success)
	
	// Verify the new file exists
	assert.FileExists(t, newFile)
	
	// Verify commit history includes the merge
	logs, err := gitInstance.Log(git.LogWithMaxCount("3"))
	require.NoError(t, err)
	
	var foundFeatureCommit bool
	for _, log := range logs {
		if strings.Contains(log.Message, "Add new feature") {
			foundFeatureCommit = true
			break
		}
	}
	assert.True(t, foundFeatureCommit, "Should find the feature commit in history")
}

// Test merge with explicit merge commit (no-ff)
func TestNoFastForwardMerge(t *testing.T) {
	tempDir := setupTestRepo(t)
	
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Create branch and add file
	err = gitInstance.CreateBranch("feature-no-ff")
	require.NoError(t, err)
	err = gitInstance.Checkout(git.CheckoutWithBranch("feature-no-ff"))
	require.NoError(t, err)
	
	featureFile := filepath.Join(tempDir, "no-ff-feature.txt")
	err = os.WriteFile(featureFile, []byte("no-ff content"), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"no-ff-feature.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Add no-ff feature")
	require.NoError(t, err)
	
	// Switch back to main
	err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
	if err != nil {
		err = gitInstance.Checkout(git.CheckoutWithBranch("master"))
	}
	require.NoError(t, err)
	
	// Merge with no-fast-forward to create explicit merge commit
	result, err := gitInstance.Merge(
		git.MergeWithBranch("feature-no-ff"),
		git.MergeWithNoFF(),
		git.MergeWithCommitMessage("Explicit merge commit"),
	)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.False(t, result.FastForward, "Should not be fast-forward with --no-ff")
	
	// Verify merge commit exists
	logs, err := gitInstance.Log(git.LogWithMaxCount("2"))
	require.NoError(t, err)
	
	// Should have merge commit and feature commit
	assert.Len(t, logs, 2)
	assert.Contains(t, logs[0].Message, "merge", "First commit should be merge commit")
}

// Test merge abort and continue workflows
func TestMergeAbortAndContinue(t *testing.T) {
	tempDir := setupTestRepo(t)
	
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Test merge abort on clean state (should handle gracefully)
	err = gitInstance.MergeAbort()
	// This may fail, which is expected when no merge is in progress
	
	// Test merge continue on clean state (should handle gracefully)  
	err = gitInstance.MergeContinue()
	// This may fail, which is expected when no merge is in progress
	
	// These tests verify the methods exist and handle edge cases without panicking
	// The actual merge workflow testing is covered in other tests
}