package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}
	
	// Test add command construction with files
	files := []string{"file1.go", "file2.go"}
	cmd := git.newCommand("add")
	cmd.args = append(cmd.args, files...)
	
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"add", "file1.go", "file2.go"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestAddWithForceOption(t *testing.T) {
	opt := AddWithForce()
	
	cmd := &Command{args: []string{"add"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--force")
}

func TestAddWithAllOption(t *testing.T) {
	opt := AddWithAll()
	
	cmd := &Command{args: []string{"add"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--all")
}

func TestAddWithDryRunOption(t *testing.T) {
	opt := AddWithDryRun()
	
	cmd := &Command{args: []string{"add"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--dry-run")
}

func TestAddWithPatchOption(t *testing.T) {
	opt := AddWithPatch()
	
	cmd := &Command{args: []string{"add"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--patch")
}