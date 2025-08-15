package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("tag", "v1.0.0")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"tag", "v1.0.0"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestListTagsCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("tag", "-l")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"tag", "-l"}, cmd.args)
}

func TestDeleteTagCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("tag", "-d", "v1.0.0")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"tag", "-d", "v1.0.0"}, cmd.args)
}

func TestTagWithAnnotatedOption(t *testing.T) {
	opt := TagWithAnnotated()
	
	cmd := &Command{args: []string{"tag"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-a")
}

func TestTagWithMessageOption(t *testing.T) {
	opt := TagWithMessage("Release version 1.0.0")
	
	cmd := &Command{args: []string{"tag"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-m")
	require.Contains(t, cmd.args, "Release version 1.0.0")
}

func TestTagWithSignOption(t *testing.T) {
	opt := TagWithSign()
	
	cmd := &Command{args: []string{"tag"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-s")
}

func TestTagWithForceOption(t *testing.T) {
	opt := TagWithForce()
	
	cmd := &Command{args: []string{"tag"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-f")
}