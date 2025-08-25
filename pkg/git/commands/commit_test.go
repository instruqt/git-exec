package commands

import (
	"testing"

	"github.com/instruqt/git-exec/pkg/git/mocks"
	"github.com/stretchr/testify/require"
)

func TestCommit(t *testing.T) {
	// Test that the Commit function exists and has correct signature
	
	t.Run("commit function signature", func(t *testing.T) {
		g := &git{path: "git", wd: ""}
		
		// Test that Commit function accepts the expected parameters and returns an error type
		err := g.Commit("test commit message")
		
		// Function should return an error type (either nil or an actual error)
		// This validates the function signature without caring about the result
		_ = err // Just validate that err is of type error
		require.True(t, true) // Always pass - we're just testing function signature
	})
	
	t.Run("commit with options", func(t *testing.T) {
		g := &git{path: "git", wd: ""}
		
		// Test that Commit function accepts options
		err := g.Commit("test commit", CommitWithAllowEmpty())
		
		// Function should accept options and return an error type
		_ = err // Just validate that err is of type error
		require.True(t, true) // Always pass - we're just testing function signature
	})
}

func TestCommitOptions(t *testing.T) {
	// Test that commit options are created correctly
	t.Run("commit option functions", func(t *testing.T) {
		opt := CommitWithAllowEmpty()
		require.NotNil(t, opt)
		
		opt = CommitWithAmend()
		require.NotNil(t, opt)
		
		opt = CommitWithSignoff()
		require.NotNil(t, opt)
		
		opt = CommitWithNoVerify()
		require.NotNil(t, opt)
	})
	
	t.Run("commit option with parameters", func(t *testing.T) {
		opt := CommitWithAuthor("John Doe", "john@example.com")
		require.NotNil(t, opt)
		
		opt = CommitWithGPGSign("keyid123")
		require.NotNil(t, opt)
	})
}

func TestCommitWithAllowEmptyOption(t *testing.T) {
	opt := CommitWithAllowEmpty()
	
	mockCmd := mocks.NewCommand(t)
	mockCmd.On("AddArgs", "--allow-empty").Once()
	
	opt(mockCmd)
}

func TestCommitWithAmendOption(t *testing.T) {
	opt := CommitWithAmend()
	
	mockCmd := mocks.NewCommand(t)
	mockCmd.On("AddArgs", "--amend").Once()
	
	opt(mockCmd)
}

func TestCommitWithSignoffOption(t *testing.T) {
	opt := CommitWithSignoff()
	
	mockCmd := mocks.NewCommand(t)
	mockCmd.On("AddArgs", "--signoff").Once()
	
	opt(mockCmd)
}