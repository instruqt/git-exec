package git

import (
	"strings"
	"github.com/instruqt/git-exec/pkg/git/errors"
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Pull pulls changes from the remote repository
func (g *gitImpl) Pull(opts ...Option) (*types.MergeResult, error) {
	cmd := g.newCommand("pull")
	cmd.ApplyOptions(opts...)
	output, err := cmd.Execute()
	
	result := &types.MergeResult{Success: err == nil}
	if err != nil {
		if gitErr, ok := err.(*errors.GitError); ok {
			if strings.Contains(gitErr.Stderr, "CONFLICT") {
				result.Success = false
				result.ConflictedFiles = strings.Fields(gitErr.Stderr)
			}
		}
		return result, err
	}
	
	outputStr := string(output)
	if strings.Contains(outputStr, "Fast-forward") {
		result.FastForward = true
	}
	return result, nil
}