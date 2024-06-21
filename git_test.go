package git

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitRepositoryInEmptyDirectory(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	g, err := New()
	require.NoError(t, err)

	output, err := g.Init(path)
	require.NoError(t, err)
	require.Equal(t, "Initialized empty Git repository in "+path+"/.git/\n", output)
	require.DirExists(t, filepath.Join(path, ".git"))
}

func TestInitRepositoryInExistingRepository(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	g, err := New()
	require.NoError(t, err)

	output, err := g.Init(path)
	require.NoError(t, err)
	require.Equal(t, "Initialized empty Git repository in "+path+"/.git/\n", output)
	require.DirExists(t, filepath.Join(path, ".git"))

	output, err = g.Init(path)
	require.NoError(t, err)
	require.Equal(t, "Reinitialized existing Git repository in "+path+"/.git/\n", output)
	require.DirExists(t, filepath.Join(path, ".git"))
}

// TODO: what happens if origin or url have invalid chars? -> add test case
func TestAddRemote(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	g, err := New()
	require.NoError(t, err)

	_, err = g.Init(path)
	require.NoError(t, err)

	g.SetWorkingDirectory(path)

	err = g.AddRemote("origin", "git@github.com:instruqt/git-exec.git")
	require.NoError(t, err)
}

func TestRemoveRemote(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	g, err := New()
	require.NoError(t, err)

	_, err = g.Init(path)
	require.NoError(t, err)

	g.SetWorkingDirectory(path)

	err = g.AddRemote("origin", "git@github.com:instruqt/git-exec.git")
	require.NoError(t, err)

	err = g.RemoveRemote("origin")
	require.NoError(t, err)
}

func TestCloneIntoEmptyDirectory(t *testing.T) {
	g, err := New()
	require.NoError(t, err)

	sourcePath := t.TempDir()
	_, err = g.Init(sourcePath, "--bare")
	require.NoError(t, err)

	destinationPath := t.TempDir()
	err = g.Clone(sourcePath, destinationPath)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(destinationPath, ".git"))
}

func TestCloneIntoExistingDirectory(t *testing.T) {
	g, err := New()
	require.NoError(t, err)

	sourcePath := t.TempDir()
	_, err = g.Init(sourcePath, "--bare")
	require.NoError(t, err)

	err = g.Clone(sourcePath, sourcePath)
	require.Error(t, err)
}
