package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddRemoteCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("remote", "add", "origin", "https://github.com/user/repo.git")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"remote", "add", "origin", "https://github.com/user/repo.git"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestRemoveRemoteCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("remote", "rm", "origin")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"remote", "rm", "origin"}, cmd.args)
}

func TestSetRemoteURLCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("remote", "set-url", "origin", "https://github.com/user/new-repo.git")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"remote", "set-url", "origin", "https://github.com/user/new-repo.git"}, cmd.args)
}

func TestListRemotesCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("remote", "-v")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"remote", "-v"}, cmd.args)
}

func TestRemoteAddWithTrackOption(t *testing.T) {
	opt := RemoteAddWithTrack("main")
	
	cmd := &Command{args: []string{"remote", "add"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-t")
	require.Contains(t, cmd.args, "main")
}

func TestRemoteAddWithMasterOption(t *testing.T) {
	opt := RemoteAddWithMaster("develop")
	
	cmd := &Command{args: []string{"remote", "add"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-m")
	require.Contains(t, cmd.args, "develop")
}

func TestSetRemoteURLWithPushOption(t *testing.T) {
	opt := SetRemoteURLWithPush()
	
	cmd := &Command{args: []string{"remote", "set-url"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--push")
}

func TestSetRemoteURLWithAddOption(t *testing.T) {
	opt := SetRemoteURLWithAdd()
	
	cmd := &Command{args: []string{"remote", "set-url"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--add")
}