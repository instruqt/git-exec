package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBranchOptions(t *testing.T) {
	// Test that branch options are created correctly
	opt := BranchWithRemote()
	require.NotNil(t, opt)
	
	opt = BranchWithAll()
	require.NotNil(t, opt)
	
	opt = BranchWithVerbose()
	require.NotNil(t, opt)
}