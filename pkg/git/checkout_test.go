package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/instruqt/git-exec/pkg/git"
	"github.com/stretchr/testify/require"
)

func TestCheckout_BranchSwitching(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-checkout-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	gitInstance, err := git.NewGit()
	require.NoError(t, err)

	// Initialize repository
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(tempDir)

	// Configure user
	err = gitInstance.SetConfig("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.email", "test@example.com")
	require.NoError(t, err)

	// Create initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)

	err = gitInstance.Commit("Initial commit")
	require.NoError(t, err)

	// Create a new branch
	err = gitInstance.CreateBranch("feature")
	require.NoError(t, err)

	// Test checkout to existing branch
	result, err := gitInstance.Checkout(git.CheckoutWithBranch("feature"))
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Equal(t, "feature", result.Branch)
	require.False(t, result.DetachedHEAD)
	require.False(t, result.NewBranch)

	// Test checkout back to main
	result, err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Equal(t, "main", result.Branch)
	require.False(t, result.DetachedHEAD)
}

func TestCheckout_CreateNewBranch(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-checkout-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	gitInstance, err := git.NewGit()
	require.NoError(t, err)

	// Initialize repository
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(tempDir)

	// Configure user
	err = gitInstance.SetConfig("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.email", "test@example.com")
	require.NoError(t, err)

	// Create initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)

	err = gitInstance.Commit("Initial commit")
	require.NoError(t, err)

	// Test checkout with create new branch
	result, err := gitInstance.Checkout(git.CheckoutWithCreateBranch("new-feature"))
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Equal(t, "new-feature", result.Branch)
	require.True(t, result.NewBranch)
	require.False(t, result.DetachedHEAD)
}

func TestCheckout_DetachedHEAD(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-checkout-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	gitInstance, err := git.NewGit()
	require.NoError(t, err)

	// Initialize repository
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(tempDir)

	// Configure user
	err = gitInstance.SetConfig("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.email", "test@example.com")
	require.NoError(t, err)

	// Create initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)

	err = gitInstance.Commit("Initial commit")
	require.NoError(t, err)

	// Get the commit hash to checkout
	logs, err := gitInstance.Log(git.LogWithMaxCount("1"))
	require.NoError(t, err)
	require.Len(t, logs, 1)

	commitHash := logs[0].Commit

	// Test checkout to specific commit (detached HEAD)
	result, err := gitInstance.Checkout(git.CheckoutWithCommit(commitHash))
	require.NoError(t, err)
	require.True(t, result.Success)
	require.True(t, result.DetachedHEAD)
	require.Equal(t, commitHash, result.Commit)
	require.Equal(t, "", result.Branch)
}

func TestCheckout_WithModifiedFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-checkout-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	gitInstance, err := git.NewGit()
	require.NoError(t, err)

	// Initialize repository
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(tempDir)

	// Configure user
	err = gitInstance.SetConfig("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.email", "test@example.com")
	require.NoError(t, err)

	// Create initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("original content"), 0644)
	require.NoError(t, err)

	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)

	err = gitInstance.Commit("Initial commit")
	require.NoError(t, err)

	// Create and switch to feature branch with different content
	err = gitInstance.CreateBranch("feature")
	require.NoError(t, err)

	result, err := gitInstance.Checkout(git.CheckoutWithBranch("feature"))
	require.NoError(t, err)
	require.True(t, result.Success)

	// Modify file on feature branch
	err = os.WriteFile(testFile, []byte("feature content"), 0644)
	require.NoError(t, err)

	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)

	err = gitInstance.Commit("Feature changes")
	require.NoError(t, err)

	// Checkout back to main - this should show file modifications
	result, err = gitInstance.Checkout(git.CheckoutWithBranch("main"))
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Equal(t, "main", result.Branch)

	// Verify file content changed back
	content, err := os.ReadFile(testFile)
	require.NoError(t, err)
	require.Equal(t, "original content", string(content))
}

func TestCheckout_ErrorCases(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-checkout-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	gitInstance, err := git.NewGit()
	require.NoError(t, err)

	// Initialize repository
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(tempDir)

	// Configure user
	err = gitInstance.SetConfig("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.email", "test@example.com")
	require.NoError(t, err)

	// Create initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)

	err = gitInstance.Commit("Initial commit")
	require.NoError(t, err)

	// Test checkout to non-existent branch
	result, err := gitInstance.Checkout(git.CheckoutWithBranch("nonexistent"))
	require.Error(t, err)
	require.False(t, result.Success)
}

func TestCheckout_WithOrphanBranch(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-checkout-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	gitInstance, err := git.NewGit()
	require.NoError(t, err)

	// Initialize repository
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(tempDir)

	// Configure user
	err = gitInstance.SetConfig("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.email", "test@example.com")
	require.NoError(t, err)

	// Create initial commit
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)

	err = gitInstance.Commit("Initial commit")
	require.NoError(t, err)

	// Test checkout with orphan branch
	result, err := gitInstance.Checkout(git.CheckoutWithOrphan("orphan-branch"))
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Equal(t, "orphan-branch", result.Branch)
	require.True(t, result.NewBranch)
}