package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateBranchCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("branch", "feature-branch")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"branch", "feature-branch"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestDeleteBranchCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("branch", "-d", "old-branch")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"branch", "-d", "old-branch"}, cmd.args)
}

func TestBranchWithRemoteOption(t *testing.T) {
	opt := BranchWithRemote()
	
	cmd := &Command{args: []string{"branch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-r")
}

func TestBranchWithAllOption(t *testing.T) {
	opt := BranchWithAll()
	
	cmd := &Command{args: []string{"branch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-a")
}

func TestBranchWithVerboseOption(t *testing.T) {
	opt := BranchWithVerbose()
	
	cmd := &Command{args: []string{"branch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-v")
}

func TestDeleteBranchWithForceOption(t *testing.T) {
	opt := DeleteBranchWithForce()
	
	cmd := &Command{args: []string{"branch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-D")
}

func TestCreateBranchWithTrackOption(t *testing.T) {
	opt := CreateBranchWithTrack()
	
	cmd := &Command{args: []string{"branch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--track")
}

func TestCreateBranchWithForceOption(t *testing.T) {
	opt := CreateBranchWithForce()
	
	cmd := &Command{args: []string{"branch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--force")
}