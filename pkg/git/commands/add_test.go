package commands

import (
	"testing"

	gitpkg "github.com/instruqt/git-exec/pkg/git"
	"github.com/instruqt/git-exec/pkg/git/mocks"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	// Test that the Add function exists and has correct signature
	// We're testing in a real git repo, so some operations may succeed
	
	t.Run("add function signature", func(t *testing.T) {
		g := &git{path: "git", wd: ""}
		
		// Test that Add function accepts the expected parameters
		// and returns an error (which is the expected return type)
		err := g.Add([]string{"nonexistent-file.txt"})
		
		// We expect an error for nonexistent files
		require.Error(t, err)
	})
	
	t.Run("add with options", func(t *testing.T) {
		g := &git{path: "git", wd: ""}
		
		// Test that Add function accepts options
		err := g.Add([]string{"nonexistent-file.txt"}, AddWithDryRun())
		
		// We expect an error for nonexistent files, even with dry run
		require.Error(t, err)
	})
}

// gitWithMockCommand wraps git with mock command creation
type gitWithMockCommand struct {
	*git
	mockCommand gitpkg.Command
}

// Override newCommand to return mock
func (g *gitWithMockCommand) newCommand(operation string, args ...string) gitpkg.Command {
	println("DEBUG: Using mock command for operation:", operation)
	return g.mockCommand
}

func TestAddCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}
	
	// Test add command construction with files
	files := []string{"file1.go", "file2.go"}
	cmd := git.newCommand("add")
	cmd.AddArgs(files...)
	
	require.Equal(t, "git", cmd.(*command).gitPath)
	require.Contains(t, cmd.GetArgs(), "add")
	require.Contains(t, cmd.GetArgs(), "file1.go")
	require.Contains(t, cmd.GetArgs(), "file2.go")
	require.Equal(t, "/test", cmd.(*command).workingDir)
}

func TestAddWithForceOption(t *testing.T) {
	opt := AddWithForce()
	
	mockCmd := mocks.NewCommand(t)
	mockCmd.On("AddArgs", "--force").Once()
	
	opt(mockCmd)
}

func TestAddWithAllOption(t *testing.T) {
	opt := AddWithAll()
	
	mockCmd := mocks.NewCommand(t)
	mockCmd.On("AddArgs", "--all").Once()
	
	opt(mockCmd)
}

func TestAddWithDryRunOption(t *testing.T) {
	opt := AddWithDryRun()
	
	mockCmd := mocks.NewCommand(t)
	mockCmd.On("AddArgs", "--dry-run").Once()
	
	opt(mockCmd)
}

func TestAddWithPatchOption(t *testing.T) {
	opt := AddWithPatch()
	
	mockCmd := mocks.NewCommand(t)
	mockCmd.On("AddArgs", "--patch").Once()
	
	opt(mockCmd)
}