package test

import (
	"github.com/instruqt/git-exec/pkg/git/commands"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloneIntoEmptyDirectory(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := commands.NewGit()
	require.NoError(t, err)

	err = git.Init(path, commands.InitWithBare())
	require.NoError(t, err)

	destinationPath := t.TempDir()
	err = git.Clone(path, destinationPath)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(destinationPath, ".git"))
}

func TestCloneIntoExistingDirectory(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := commands.NewGit()
	require.NoError(t, err)

	err = git.Init(path, commands.InitWithBare())
	require.NoError(t, err)

	err = git.Clone(path, path)
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("destination path '%s' already exists and is not an empty directory", path))
}