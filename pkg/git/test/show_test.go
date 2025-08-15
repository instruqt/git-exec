package test

import (
	"github.com/instruqt/git-exec/pkg/git/commands"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShow(t *testing.T) {
	path, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	git, err := commands.NewGit()
	require.NoError(t, err)

	err = git.Init(path)
	require.NoError(t, err)

	git.SetWorkingDirectory(path)

	err = os.WriteFile(filepath.Join(path, "file.txt"), []byte("Hello, World!"), 0644)
	require.NoError(t, err)

	err = git.Add([]string{"file.txt"})
	require.NoError(t, err)

	err = git.Commit("Initial commit", commands.WithUser("John Doe", "john.doe@gmail.com"))
	require.NoError(t, err)

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