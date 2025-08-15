package commands

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddRemote(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := NewGit()
	require.NoError(t, err)

	err = git.Init(path)
	require.NoError(t, err)

	git.SetWorkingDirectory(path)

	err = git.AddRemote("origin", "git@github.com:instruqt/git-exec.git")
	require.NoError(t, err)

	remotes, err := git.ListRemotes()
	require.NoError(t, err)
	require.Len(t, remotes, 1)
	require.Equal(t, "origin", remotes[0].Name)
	require.Equal(t, "git@github.com:instruqt/git-exec.git", remotes[0].URL)
}

func TestRemoveRemote(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := NewGit()
	require.NoError(t, err)

	err = git.Init(path)
	require.NoError(t, err)

	git.SetWorkingDirectory(path)

	err = git.AddRemote("origin", "git@github.com:instruqt/git-exec.git")
	require.NoError(t, err)

	err = git.RemoveRemote("origin")
	require.NoError(t, err)
}

func TestListRemotes(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := NewGit()
	require.NoError(t, err)

	err = git.Init(path)
	require.NoError(t, err)

	git.SetWorkingDirectory(path)

	err = git.AddRemote("first", "first-url")
	require.NoError(t, err)

	err = git.AddRemote("second", "second-url")
	require.NoError(t, err)

	remotes, err := git.ListRemotes()
	require.NoError(t, err)
	require.Len(t, remotes, 2)
	require.Equal(t, "first", remotes[0].Name)
	require.Equal(t, "first-url", remotes[0].URL)
	require.Equal(t, "second", remotes[1].Name)
	require.Equal(t, "second-url", remotes[1].URL)
}