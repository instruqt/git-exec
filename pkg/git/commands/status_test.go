package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatusCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("status", "--porcelain")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"status", "--porcelain"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestStatusWithShortOption(t *testing.T) {
	opt := StatusWithShort()
	
	cmd := &Command{args: []string{"status"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--short")
}

func TestStatusWithBranchOption(t *testing.T) {
	opt := StatusWithBranch()
	
	cmd := &Command{args: []string{"status"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--branch")
}

func TestStatusWithPorcelainOption(t *testing.T) {
	opt := StatusWithPorcelain()
	
	cmd := &Command{args: []string{"status"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--porcelain")
}

func TestStatusWithLongOption(t *testing.T) {
	opt := StatusWithLong()
	
	cmd := &Command{args: []string{"status"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--long")
}

func TestStatusWithUntrackedFilesOption(t *testing.T) {
	opt := StatusWithUntrackedFiles("all")
	
	cmd := &Command{args: []string{"status"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--untracked-files=all")
}

func TestStatusWithIgnoredFilesOption(t *testing.T) {
	opt := StatusWithIgnoredFiles()
	
	cmd := &Command{args: []string{"status"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--ignored")
}