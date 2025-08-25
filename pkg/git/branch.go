package git

import (
	"fmt"
	"strings"
	"github.com/instruqt/git-exec/pkg/git/types"
)

// ListBranches lists all branches in the repository
func (g *gitImpl) ListBranches(opts ...Option) ([]types.Branch, error) {
	cmd := g.newCommand("branch")
	cmd.ApplyOptions(opts...)
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}
	
	branches := []types.Branch{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		active := false
		name := line
		if strings.HasPrefix(line, "* ") {
			active = true
			name = line[2:]
		} else if strings.HasPrefix(line, "  ") {
			name = line[2:]
		}
		branches = append(branches, types.Branch{Name: strings.TrimSpace(name), Active: active})
	}
	return branches, nil
}

// CreateBranch creates a new branch
func (g *gitImpl) CreateBranch(branch string, opts ...Option) error {
	cmd := g.newCommand("branch", branch)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}

// DeleteBranch deletes a branch
func (g *gitImpl) DeleteBranch(branch string, opts ...Option) error {
	cmd := g.newCommand("branch", "-d", branch)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}

// SetUpstream sets the upstream branch for tracking
func (g *gitImpl) SetUpstream(branch string, remote string, opts ...Option) error {
	upstreamRef := fmt.Sprintf("%s/%s", remote, branch)
	cmd := g.newCommand("branch", "--set-upstream-to", upstreamRef, branch)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}