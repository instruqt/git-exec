package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloneCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	// Test basic clone command construction
	cmd := git.newCommand("clone", "https://github.com/user/repo", "/path/to/dest")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"clone", "https://github.com/user/repo", "/path/to/dest"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestCloneWithBareOption(t *testing.T) {
	opt := CloneWithBare()
	
	cmd := &Command{args: []string{"clone"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--bare")
}

func TestCloneWithBranchOption(t *testing.T) {
	opt := CloneWithBranch("main")
	
	cmd := &Command{args: []string{"clone"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--branch")
	require.Contains(t, cmd.args, "main")
}

func TestCloneWithDepthOption(t *testing.T) {
	opt := CloneWithDepth(5)
	
	cmd := &Command{args: []string{"clone"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--depth")
	require.Contains(t, cmd.args, "5")
}