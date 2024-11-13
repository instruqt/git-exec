package git

import (
	"fmt"
	"os"
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
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := New()
	require.NoError(t, err)

	_, err = git.Init(path, "--bare")
	require.NoError(t, err)

	destinationPath := t.TempDir()
	err = git.Clone(path, destinationPath)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(destinationPath, ".git"))
}

func TestCloneIntoExistingDirectory(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := New()
	require.NoError(t, err)

	_, err = git.Init(path, "--bare")
	require.NoError(t, err)

	err = git.Clone(path, path)
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("fatal: destination path '%s' already exists and is not an empty directory.\n", path))
}

func TestShow(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := New()
	require.NoError(t, err)

	_, err = git.Init(path)
	require.NoError(t, err)

	git.SetWorkingDirectory(path)

	err = os.WriteFile(filepath.Join(path, "file.txt"), []byte("Hello, World!"), 0644)
	require.NoError(t, err)

	err = git.Add("file.txt")
	require.NoError(t, err)

	err = git.Commit("Initial commit", "John Doe", "john.doe@gmail.com")
	require.NoError(t, err)

	// TODO: Add test case for show in a repo without commits -> "fatal: ambiguous argument 'HEAD': unknown revision or path not in the working tree."

	output, err := git.Show("HEAD")
	require.NoError(t, err)

	require.NotEmpty(t, output.Commit)
	require.Equal(t, "John Doe <john.doe@gmail.com>", output.Author)
	require.NotEmpty(t, output.AuthorDate)
	require.Equal(t, "John Doe <john.doe@gmail.com>", output.Committer)
	require.NotEmpty(t, output.CommitterDate)
	require.Equal(t, "Initial commit", output.Message)

	require.Len(t, output.Diffs, 1)
	require.Equal(t, "+Hello, World!\n\\ No newline at end of file\n", output.Diffs[0].Contents)
}

func TestLog(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := New()
	require.NoError(t, err)

	_, err = git.Init(path)
	require.NoError(t, err)

	git.SetWorkingDirectory(path)

	err = os.WriteFile(filepath.Join(path, "file.txt"), []byte("Hello, World!"), 0644)
	require.NoError(t, err)

	err = git.Add("file.txt")
	require.NoError(t, err)

	err = git.Commit("Initial commit", "John Doe", "john.doe@gmail.com")
	require.NoError(t, err)

	output, err := git.Log()
	require.NoError(t, err)

	require.Len(t, output, 1)
	require.Equal(t, "John Doe <john.doe@gmail.com>", output[0].Author)
	require.NotEmpty(t, output[0].AuthorDate)
	require.Equal(t, "John Doe <john.doe@gmail.com>", output[0].Committer)
	require.NotEmpty(t, output[0].CommitterDate)
	require.Equal(t, "Initial commit", output[0].Message)
}

// TODO: implement real test cases
func TestFetch(t *testing.T) {
	// path, err := filepath.EvalSymlinks(t.TempDir())
	path, err := os.Getwd()
	require.NoError(t, err)

	git, err := New()
	require.NoError(t, err)

	git.SetWorkingDirectory(path)

	refs, err := git.Fetch()
	require.NoError(t, err)

	require.NotEmpty(t, refs)
}

func TestPull(t *testing.T) {
	// create a repository with a commit
	first, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := New()
	require.NoError(t, err)

	git.SetWorkingDirectory(first)

	_, err = git.Init(first)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(first, "file.txt"), []byte("Hello, World!"), 0644)
	require.NoError(t, err)

	err = git.Add("file.txt")
	require.NoError(t, err)

	err = git.Commit("Initial commit", "John Doe", "john.doe@gmail.com")
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

	err = git.Add("new.txt")
	require.NoError(t, err)

	err = git.Commit("Add new.txt", "John Doe", "john.doe@gmail.com")
	require.NoError(t, err)

	// pull the changes from the first repository to the second repository
	git.SetWorkingDirectory(second)

	_, err = git.Pull()
	require.NoError(t, err)
}
