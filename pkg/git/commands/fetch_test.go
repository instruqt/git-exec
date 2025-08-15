package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	// Use current working directory (this repo) for testing
	path, err := os.Getwd()
	require.NoError(t, err)

	git, err := NewGit()
	require.NoError(t, err)

	git.SetWorkingDirectory(path)

	refs, err := git.Fetch()
	require.NoError(t, err)

	require.NotEmpty(t, refs)
}