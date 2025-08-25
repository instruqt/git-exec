package git

import (
	"strings"
	"github.com/instruqt/git-exec/pkg/git/errors"
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Merge merges branches
func (g *gitImpl) Merge(opts ...Option) (*types.MergeResult, error) {
	cmd := g.newCommand("merge")
	cmd.ApplyOptions(opts...)
	
	result := &types.MergeResult{}
	output, err := cmd.Execute()
	
	if err != nil {
		if gitErr, ok := err.(*errors.GitError); ok {
			if strings.Contains(gitErr.Stderr, "CONFLICT") || strings.Contains(gitErr.Stdout, "CONFLICT") {
				result.Success = false
				return result, nil
			}
			result.AbortReason = gitErr.Stderr
		}
		result.Success = false
		return result, err
	}
	
	result.Success = true
	outputStr := string(output)
	if strings.Contains(outputStr, "Fast-forward") {
		result.FastForward = true
	}
	return result, nil
}

// MergeAbort aborts a merge in progress
func (g *gitImpl) MergeAbort() error {
	cmd := g.newCommand("merge", "--abort")
	_, err := cmd.Execute()
	return err
}

// MergeContinue continues a merge after resolving conflicts
func (g *gitImpl) MergeContinue() error {
	cmd := g.newCommand("merge", "--continue")
	_, err := cmd.Execute()
	return err
}

// ResolveConflicts resolves merge conflicts using specified resolutions
func (g *gitImpl) ResolveConflicts(resolutions []types.ConflictResolution) error {
	for _, resolution := range resolutions {
		if resolution.UseOurs {
			cmd := g.newCommand("checkout", "--ours", resolution.FilePath)
			if _, err := cmd.Execute(); err != nil {
				return err
			}
		} else if resolution.UseTheirs {
			cmd := g.newCommand("checkout", "--theirs", resolution.FilePath)
			if _, err := cmd.Execute(); err != nil {
				return err
			}
		}
		addCmd := g.newCommand("add", resolution.FilePath)
		if _, err := addCmd.Execute(); err != nil {
			return err
		}
	}
	return nil
}