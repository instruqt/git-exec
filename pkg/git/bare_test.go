package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/instruqt/git-exec/pkg/git"
	"github.com/instruqt/git-exec/pkg/git/types"
	"github.com/stretchr/testify/require"
)

func TestBareRepository_Init(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-bare-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	bareRepo := filepath.Join(tempDir, "bare.git")

	gitInstance, err := git.NewGit()
	require.NoError(t, err)

	// Initialize bare repository
	err = gitInstance.Init(bareRepo, git.InitWithBare())
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(bareRepo)

	// Check if it's a bare repository
	isBare, err := gitInstance.IsBareRepository()
	require.NoError(t, err)
	require.True(t, isBare)
}

func TestBareRepository_Clone(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-bare-clone-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source repository
	sourceRepo := filepath.Join(tempDir, "source")
	gitInstance, err := git.NewGit()
	require.NoError(t, err)

	err = gitInstance.Init(sourceRepo)
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(sourceRepo)

	// Configure user
	err = gitInstance.SetConfig("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.email", "test@example.com")
	require.NoError(t, err)

	// Create initial commit
	testFile := filepath.Join(sourceRepo, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)

	err = gitInstance.Commit("Initial commit")
	require.NoError(t, err)

	// Clone as bare repository
	bareRepo := filepath.Join(tempDir, "bare.git")
	err = gitInstance.Clone(sourceRepo, bareRepo, git.CloneWithBare())
	require.NoError(t, err)

	// Check if cloned repo is bare
	gitInstance.SetWorkingDirectory(bareRepo)
	isBare, err := gitInstance.IsBareRepository()
	require.NoError(t, err)
	require.True(t, isBare)
}

func TestBareRepository_ReferenceManagement(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-bare-refs-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source repository with commits
	sourceRepo := filepath.Join(tempDir, "source")
	gitInstance, err := git.NewGit()
	require.NoError(t, err)

	err = gitInstance.Init(sourceRepo)
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(sourceRepo)

	// Configure user
	err = gitInstance.SetConfig("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.SetConfig("user.email", "test@example.com")
	require.NoError(t, err)

	// Create commits
	testFile := filepath.Join(sourceRepo, "test.txt")
	err = os.WriteFile(testFile, []byte("content 1"), 0644)
	require.NoError(t, err)

	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)

	err = gitInstance.Commit("First commit")
	require.NoError(t, err)

	// Get first commit hash
	logs1, err := gitInstance.Log(git.LogWithMaxCount("1"))
	require.NoError(t, err)
	require.Len(t, logs1, 1)
	firstCommit := logs1[0].Commit

	// Create second commit
	err = os.WriteFile(testFile, []byte("content 2"), 0644)
	require.NoError(t, err)

	err = gitInstance.Add([]string{"test.txt"})
	require.NoError(t, err)

	err = gitInstance.Commit("Second commit")
	require.NoError(t, err)

	// Clone as bare repository
	bareRepo := filepath.Join(tempDir, "bare.git")
	err = gitInstance.Clone(sourceRepo, bareRepo, git.CloneWithBare())
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(bareRepo)

	// List references
	refs, err := gitInstance.ListRefs()
	require.NoError(t, err)
	require.NotEmpty(t, refs)

	// Should have at least HEAD and main branch
	var hasMain bool
	for _, ref := range refs {
		if ref.Name == "refs/heads/main" {
			hasMain = true
			require.Equal(t, types.ReferenceTypeBranch, ref.Type)
		}
	}
	require.True(t, hasMain)

	// Update a reference
	err = gitInstance.UpdateRef("refs/heads/test-branch", firstCommit)
	require.NoError(t, err)

	// List references again
	refs, err = gitInstance.ListRefs()
	require.NoError(t, err)

	// Check new branch exists
	var hasTestBranch bool
	for _, ref := range refs {
		if ref.Name == "refs/heads/test-branch" {
			hasTestBranch = true
			require.Equal(t, firstCommit, ref.Commit)
			require.Equal(t, types.ReferenceTypeBranch, ref.Type)
		}
	}
	require.True(t, hasTestBranch)

	// Delete the reference
	err = gitInstance.DeleteRef("refs/heads/test-branch")
	require.NoError(t, err)

	// List references again
	refs, err = gitInstance.ListRefs()
	require.NoError(t, err)

	// Check branch is deleted
	for _, ref := range refs {
		require.NotEqual(t, "refs/heads/test-branch", ref.Name)
	}
}

func TestBareRepository_NonBareCheck(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-nonbare-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	normalRepo := filepath.Join(tempDir, "normal")

	gitInstance, err := git.NewGit()
	require.NoError(t, err)

	// Initialize normal (non-bare) repository
	err = gitInstance.Init(normalRepo)
	require.NoError(t, err)

	gitInstance.SetWorkingDirectory(normalRepo)

	// Check if it's a bare repository (should be false)
	isBare, err := gitInstance.IsBareRepository()
	require.NoError(t, err)
	require.False(t, isBare)
}