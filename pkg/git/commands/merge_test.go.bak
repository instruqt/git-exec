package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("merge")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"merge"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestMergeWithBranchOption(t *testing.T) {
	opt := MergeWithBranch("feature")
	
	cmd := &Command{args: []string{"merge"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "feature")
}

func TestMergeWithCommitOption(t *testing.T) {
	opt := MergeWithCommit("abc123")
	
	cmd := &Command{args: []string{"merge"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "abc123")
}

func TestMergeWithNoFFOption(t *testing.T) {
	opt := MergeWithNoFF()
	
	cmd := &Command{args: []string{"merge"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--no-ff")
}

func TestMergeWithFFOnlyOption(t *testing.T) {
	opt := MergeWithFFOnly()
	
	cmd := &Command{args: []string{"merge"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--ff-only")
}

func TestMergeWithSquashOption(t *testing.T) {
	opt := MergeWithSquash()
	
	cmd := &Command{args: []string{"merge"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--squash")
}

func TestMergeWithStrategyOption(t *testing.T) {
	opt := MergeWithStrategy("recursive")
	
	cmd := &Command{args: []string{"merge"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--strategy")
	require.Contains(t, cmd.args, "recursive")
}

func TestMergeWithMessageOption(t *testing.T) {
	opt := MergeWithMessage("Merge feature branch")
	
	cmd := &Command{args: []string{"merge"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-m")
	require.Contains(t, cmd.args, "Merge feature branch")
}

func TestMergeWithAbortOption(t *testing.T) {
	opt := MergeWithAbort()
	
	cmd := &Command{args: []string{"merge"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--abort")
}

func TestMergeWithContinueOption(t *testing.T) {
	opt := MergeWithContinue()
	
	cmd := &Command{args: []string{"merge"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--continue")
}