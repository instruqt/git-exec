package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPush(t *testing.T) {
	// create a repository with a commit
	first, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := NewGit()
	require.NoError(t, err)

	git.SetWorkingDirectory(first)

	err = git.Init(first, InitWithBare())
	require.NoError(t, err)

	// clone the repository to a second directory
	second, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	err = git.Clone(first, second)
	require.NoError(t, err)

	// add a commit to the second repository
	git.SetWorkingDirectory(second)

	err = os.WriteFile(filepath.Join(second, "new.txt"), []byte("New\n"), 0644)
	require.NoError(t, err)

	err = git.Add([]string{"new.txt"})
	require.NoError(t, err)

	err = git.Commit("Add new.txt", WithUser("John Doe", "john.doe@gmail.com"))
	require.NoError(t, err)

	// push the changes from the second repository to the first repository
	remotes, err := git.Push()
	require.NoError(t, err)

	require.NotEmpty(t, remotes)
}