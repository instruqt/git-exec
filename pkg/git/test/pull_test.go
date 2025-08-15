package test

import (
	"github.com/instruqt/git-exec/pkg/git/commands"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPull(t *testing.T) {
	// create a repository with a commit
	first, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := commands.NewGit()
	require.NoError(t, err)

	git.SetWorkingDirectory(first)

	err = git.Init(first)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(first, "file.txt"), []byte("Hello, World!"), 0644)
	require.NoError(t, err)

	err = git.Add([]string{"file.txt"})
	require.NoError(t, err)

	err = git.Commit("Initial commit", commands.WithUser("John Doe", "john.doe@gmail.com"))
	require.NoError(t, err)

	// clone the repository to a second directory
	second, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	err = git.Clone(first, second)
	require.NoError(t, err)

	// add a commit to the first repository
	git.SetWorkingDirectory(first)

	err = os.WriteFile(filepath.Join(first, "new.txt"), []byte("New\n"), 0644)
	require.NoError(t, err)

	err = git.Add([]string{"new.txt"})
	require.NoError(t, err)

	err = git.Commit("Add new.txt", commands.WithUser("John Doe", "john.doe@gmail.com"))
	require.NoError(t, err)

	// pull the changes from the first repository to the second repository
	git.SetWorkingDirectory(second)

	_, err = git.Pull()
	require.NoError(t, err)
}