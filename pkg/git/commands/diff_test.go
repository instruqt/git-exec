package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiffWithCachedOption(t *testing.T) {
	opt := DiffWithCached()
	
	cmd := &Command{args: []string{"diff"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--cached")
}

func TestDiffWithStagedOption(t *testing.T) {
	opt := DiffWithStaged()
	
	cmd := &Command{args: []string{"diff"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--staged")
}

func TestDiffWithNameOnlyOption(t *testing.T) {
	opt := DiffWithNameOnly()
	
	cmd := &Command{args: []string{"diff"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--name-only")
}

func TestDiffWithNameStatusOption(t *testing.T) {
	opt := DiffWithNameStatus()
	
	cmd := &Command{args: []string{"diff"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--name-status")
}

func TestDiffWithStatOption(t *testing.T) {
	opt := DiffWithStat()
	
	cmd := &Command{args: []string{"diff"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--stat")
}

func TestDiffWithContextOption(t *testing.T) {
	opt := DiffWithContext("5")
	
	cmd := &Command{args: []string{"diff"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-U5")
}

func TestDiffWithCommitOption(t *testing.T) {
	opt := DiffWithCommit("abc123")
	
	cmd := &Command{args: []string{"diff"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "abc123")
}

func TestDiffWithCommitRangeOption(t *testing.T) {
	opt := DiffWithCommitRange("abc123", "def456")
	
	cmd := &Command{args: []string{"diff"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "abc123..def456")
}