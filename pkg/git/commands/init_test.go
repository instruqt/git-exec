package commands

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitRepositoryInEmptyDirectory(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := NewGit()
	require.NoError(t, err)

	err = git.Init(path)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(path, ".git"))
}

func TestInitRepositoryInExistingRepository(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := NewGit()
	require.NoError(t, err)

	err = git.Init(path)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(path, ".git"))

	err = git.Init(path)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(path, ".git"))
}