package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShowCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("show", "--format=fuller", "HEAD")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"show", "--format=fuller", "HEAD"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestShowWithFormatOption(t *testing.T) {
	opt := ShowWithFormat("oneline")
	
	cmd := &Command{args: []string{"show"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--format=oneline")
}

func TestShowWithPrettyOption(t *testing.T) {
	opt := ShowWithPretty("short")
	
	cmd := &Command{args: []string{"show"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--pretty=short")
}

func TestShowWithOnelineOption(t *testing.T) {
	opt := ShowWithOneline()
	
	cmd := &Command{args: []string{"show"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--oneline")
}

func TestShowWithStatOption(t *testing.T) {
	opt := ShowWithStat()
	
	cmd := &Command{args: []string{"show"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--stat")
}

func TestShowWithNameOnlyOption(t *testing.T) {
	opt := ShowWithNameOnly()
	
	cmd := &Command{args: []string{"show"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--name-only")
}

func TestShowWithNameStatusOption(t *testing.T) {
	opt := ShowWithNameStatus()
	
	cmd := &Command{args: []string{"show"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--name-status")
}

func TestShowWithNoPatchOption(t *testing.T) {
	opt := ShowWithNoPatch()
	
	cmd := &Command{args: []string{"show"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--no-patch")
}

func TestShowWithPatchOption(t *testing.T) {
	opt := ShowWithPatch()
	
	cmd := &Command{args: []string{"show"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--patch")
}