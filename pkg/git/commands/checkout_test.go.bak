package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckoutWithBranchOption(t *testing.T) {
	opt := CheckoutWithBranch("feature-branch")
	
	cmd := &Command{args: []string{"checkout"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "feature-branch")
}

func TestCheckoutWithFilesOption(t *testing.T) {
	files := []string{"file1.go", "file2.go", "dir/file3.txt"}
	opt := CheckoutWithFiles(files)
	
	cmd := &Command{args: []string{"checkout"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "file1.go")
	require.Contains(t, cmd.args, "file2.go")
	require.Contains(t, cmd.args, "dir/file3.txt")
}

func TestCheckoutWithNewBranchOption(t *testing.T) {
	opt := CheckoutWithNewBranch("new-feature")
	
	cmd := &Command{args: []string{"checkout"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "-b")
	require.Contains(t, cmd.args, "new-feature")
}

func TestCheckoutWithForceOption(t *testing.T) {
	opt := CheckoutWithForce()
	
	cmd := &Command{args: []string{"checkout"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--force")
}