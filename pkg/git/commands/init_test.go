package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("init", "/path/to/repo")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"init", "/path/to/repo"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestInitWithBareOption(t *testing.T) {
	opt := InitWithBare()
	
	cmd := &Command{args: []string{"init"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--bare")
}

func TestInitWithQuietOption(t *testing.T) {
	opt := InitWithQuiet()
	
	cmd := &Command{args: []string{"init"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--quiet")
}

func TestInitWithBranchOption(t *testing.T) {
	opt := InitWithBranch("main")
	
	cmd := &Command{args: []string{"init"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--initial-branch")
	require.Contains(t, cmd.args, "main")
}

func TestInitWithTemplateOption(t *testing.T) {
	opt := InitWithTemplate("/path/to/template")
	
	cmd := &Command{args: []string{"init"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--template")
	require.Contains(t, cmd.args, "/path/to/template")
}

func TestInitWithSharedOption(t *testing.T) {
	opt := InitWithShared("group")
	
	cmd := &Command{args: []string{"init"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--shared=group")
}

func TestInitWithSharedNoPermissionsOption(t *testing.T) {
	opt := InitWithShared("")
	
	cmd := &Command{args: []string{"init"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--shared")
}