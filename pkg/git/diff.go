package git

import (
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Diff shows differences between commits, commit and working tree, etc
func (g *gitImpl) Diff(opts ...Option) ([]types.Diff, error) {
	cmd := g.newCommand("diff")
	cmd.ApplyOptions(opts...)
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}
	if len(output) == 0 {
		return []types.Diff{}, nil
	}
	// Simplified - would need proper diff parsing
	return []types.Diff{}, nil
}