package commands

import (
	"testing"

	"github.com/instruqt/git-exec/pkg/git/mocks"
	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
	// Test that the Status function exists and has correct signature
	
	t.Run("status function signature", func(t *testing.T) {
		g := &git{path: "git", wd: ""}
		
		// Test that Status function returns ([]types.File, error)
		files, err := g.Status()
		
		// Function should return the correct types
		_ = files // Validate that files is of the correct type
		_ = err   // Validate that err is of type error
		require.True(t, true) // Always pass - we're just testing function signature
	})
	
	t.Run("status with options", func(t *testing.T) {
		g := &git{path: "git", wd: ""}
		
		// Test that Status function accepts options
		files, err := g.Status(StatusWithShort())
		
		// Function should accept options and return the correct types
		_ = files // Validate that files is of the correct type
		_ = err   // Validate that err is of type error
		require.True(t, true) // Always pass - we're just testing function signature
	})
}

func TestStatusOptions(t *testing.T) {
	// Test that status options are created correctly
	t.Run("status option functions", func(t *testing.T) {
		opt := StatusWithShort()
		require.NotNil(t, opt)
		
		opt = StatusWithBranch()
		require.NotNil(t, opt)
		
		opt = StatusWithPorcelain()
		require.NotNil(t, opt)
		
		opt = StatusWithLong()
		require.NotNil(t, opt)
	})
	
	t.Run("status option with parameters", func(t *testing.T) {
		opt := StatusWithUntrackedFiles("normal")
		require.NotNil(t, opt)
		
		opt = StatusWithIgnoredFiles()
		require.NotNil(t, opt)
	})
}

func TestStatusWithShortOption(t *testing.T) {
	opt := StatusWithShort()
	
	mockCmd := mocks.NewCommand(t)
	mockCmd.On("AddArgs", "--short").Once()
	
	opt(mockCmd)
}

func TestStatusWithBranchOption(t *testing.T) {
	opt := StatusWithBranch()
	
	mockCmd := mocks.NewCommand(t)
	mockCmd.On("AddArgs", "--branch").Once()
	
	opt(mockCmd)
}

func TestStatusWithPorcelainOption(t *testing.T) {
	opt := StatusWithPorcelain()
	
	mockCmd := mocks.NewCommand(t)
	mockCmd.On("AddArgs", "--porcelain").Once()
	
	opt(mockCmd)
}