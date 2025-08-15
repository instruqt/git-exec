package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPushCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("push")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"push"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestPushWithRemoteOption(t *testing.T) {
	opt := PushWithRemote("origin")
	
	cmd := &Command{args: []string{"push"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "origin")
}

func TestPushWithBranchOption(t *testing.T) {
	opt := PushWithBranch("main")
	
	cmd := &Command{args: []string{"push"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "main")
}

func TestPushWithRemoteAndBranchOption(t *testing.T) {
	opt := PushWithRemoteAndBranch("origin", "feature")
	
	cmd := &Command{args: []string{"push"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "origin")
	require.Contains(t, cmd.args, "feature")
}

func TestPushWithForceOption(t *testing.T) {
	opt := PushWithForce()
	
	cmd := &Command{args: []string{"push"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--force")
}

func TestPushWithForceWithLeaseOption(t *testing.T) {
	opt := PushWithForceWithLease()
	
	cmd := &Command{args: []string{"push"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--force-with-lease")
}

func TestPushWithAllOption(t *testing.T) {
	opt := PushWithAll()
	
	cmd := &Command{args: []string{"push"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--all")
}

func TestPushWithTagsOption(t *testing.T) {
	opt := PushWithTags()
	
	cmd := &Command{args: []string{"push"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--tags")
}

func TestPushWithSetUpstreamOption(t *testing.T) {
	opt := PushWithSetUpstream()
	
	cmd := &Command{args: []string{"push"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--set-upstream")
}

func TestPushWithDeleteOption(t *testing.T) {
	opt := PushWithDelete()
	
	cmd := &Command{args: []string{"push"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--delete")
}