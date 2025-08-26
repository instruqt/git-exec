package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/instruqt/git-exec/pkg/git"
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