package git

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitRepositoryInEmptyDirectory(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := New()
	require.NoError(t, err)

	output, err := git.Init(path)
	require.NoError(t, err)
	require.Equal(t, "Initialized empty Git repository in "+path+"/.git/\n", output)
	require.DirExists(t, filepath.Join(path, ".git"))
}

func TestInitRepositoryInExistingRepository(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := New()
	require.NoError(t, err)

	output, err := git.Init(path)
	require.NoError(t, err)
	require.Equal(t, "Initialized empty Git repository in "+path+"/.git/\n", output)
	require.DirExists(t, filepath.Join(path, ".git"))

	output, err = git.Init(path)
	require.NoError(t, err)
	require.Equal(t, "Reinitialized existing Git repository in "+path+"/.git/\n", output)
	require.DirExists(t, filepath.Join(path, ".git"))
}

// TODO: what happens if origin or url have invalid chars? -> add test case
func TestAddRemote(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := New()
	require.NoError(t, err)

	_, err = git.Init(path)
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

	git, err := New()
	require.NoError(t, err)

	_, err = git.Init(path)
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

	git, err := New()
	require.NoError(t, err)

	_, err = git.Init(path)
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

func TestCloneIntoEmptyDirectory(t *testing.T) {
	git, err := New()
	require.NoError(t, err)

	sourcePath := t.TempDir()
	_, err = git.Init(sourcePath, "--bare")
	require.NoError(t, err)

	destinationPath := t.TempDir()
	err = git.Clone(sourcePath, destinationPath)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(destinationPath, ".git"))
}

func TestCloneIntoExistingDirectory(t *testing.T) {
	git, err := New()
	require.NoError(t, err)

	sourcePath := t.TempDir()
	_, err = git.Init(sourcePath, "--bare")
	require.NoError(t, err)

	err = git.Clone(sourcePath, sourcePath)
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("fatal: destination path '%s' already exists and is not an empty directory.\n", sourcePath))
}
