package test

import (
	"github.com/instruqt/git-exec/pkg/git/commands"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
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

	output, err := git.Log()
	require.NoError(t, err)

	require.Len(t, output, 1)
	require.Equal(t, "John Doe <john.doe@gmail.com>", output[0].Author)
	require.NotEmpty(t, output[0].AuthorDate)
	require.Equal(t, "John Doe <john.doe@gmail.com>", output[0].Committer)
	require.NotEmpty(t, output[0].CommitterDate)
	require.Equal(t, "Initial commit", output[0].Message)
}