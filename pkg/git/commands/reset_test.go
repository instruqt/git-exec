package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResetCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("reset")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"reset"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestResetWithSoftOption(t *testing.T) {
	opt := ResetWithSoft()
	
	cmd := &Command{args: []string{"reset"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--soft")
}

func TestResetWithMixedOption(t *testing.T) {
	opt := ResetWithMixed()
	
	cmd := &Command{args: []string{"reset"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--mixed")
}

func TestResetWithHardOption(t *testing.T) {
	opt := ResetWithHard()
	
	cmd := &Command{args: []string{"reset"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--hard")
}

func TestResetWithMergeOption(t *testing.T) {
	opt := ResetWithMerge()
	
	cmd := &Command{args: []string{"reset"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--merge")
}

func TestResetWithKeepOption(t *testing.T) {
	opt := ResetWithKeep()
	
	cmd := &Command{args: []string{"reset"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--keep")
}

func TestResetWithRecurseSubmodulesOption(t *testing.T) {
	opt := ResetWithRecurseSubmodules()
	
	cmd := &Command{args: []string{"reset"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--recurse-submodules")
}