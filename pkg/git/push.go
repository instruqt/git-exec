package git

import (
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Push pushes changes to the remote repository
func (g *gitImpl) Push(opts ...Option) ([]types.Remote, error) {
	cmd := g.newCommand("push")
	cmd.ApplyOptions(opts...)
	_, err := cmd.ExecuteWithStderr()
	if err != nil {
		return nil, err
	}
	return []types.Remote{}, nil
}