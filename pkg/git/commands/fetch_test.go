package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("fetch", "-v")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"fetch", "-v"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestFetchWithAllOption(t *testing.T) {
	opt := FetchWithAll()
	
	cmd := &Command{args: []string{"fetch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--all")
}

func TestFetchWithPruneOption(t *testing.T) {
	opt := FetchWithPrune()
	
	cmd := &Command{args: []string{"fetch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--prune")
}

func TestFetchWithPruneTagsOption(t *testing.T) {
	opt := FetchWithPruneTags()
	
	cmd := &Command{args: []string{"fetch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--prune-tags")
}

func TestFetchWithTagsOption(t *testing.T) {
	opt := FetchWithTags()
	
	cmd := &Command{args: []string{"fetch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--tags")
}

func TestFetchWithNoTagsOption(t *testing.T) {
	opt := FetchWithNoTags()
	
	cmd := &Command{args: []string{"fetch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--no-tags")
}

func TestFetchWithDepthOption(t *testing.T) {
	opt := FetchWithDepth(10)
	
	cmd := &Command{args: []string{"fetch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--depth")
	require.Contains(t, cmd.args, "10")
}

func TestFetchWithRemoteOption(t *testing.T) {
	opt := FetchWithRemote("upstream")
	
	cmd := &Command{args: []string{"fetch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "upstream")
}

func TestFetchWithForceOption(t *testing.T) {
	opt := FetchWithForce()
	
	cmd := &Command{args: []string{"fetch"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--force")
}