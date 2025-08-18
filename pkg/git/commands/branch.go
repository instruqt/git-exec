package commands

import (
	"strings"

	gitpkg "github.com/instruqt/git-exec/pkg/git"
	"github.com/instruqt/git-exec/pkg/git/types"
)

// ListBranches lists all local branches
func (g *git) ListBranches(opts ...gitpkg.Option) ([]types.Branch, error) {
	cmd := g.newCommand("branch")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}

	branches := []types.Branch{}
	lines := strings.Split(string(output), "\n")
	for _, branch := range lines {
		if branch == "" {
			continue
		}

		branches = append(branches, types.Branch{
			Name:   branch[2:],
			Active: branch[0] == '*',
		})
	}

	return branches, nil
}

// CreateBranch creates a new branch
func (g *git) CreateBranch(branch string, opts ...gitpkg.Option) error {
	cmd := g.newCommand("branch", branch)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// SetUpstream sets the upstream for a branch
func (g *git) SetUpstream(branch string, remote string, opts ...gitpkg.Option) error {
	cmd := g.newCommand("branch", branch, "-u", remote+"/"+branch)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// DeleteBranch deletes a branch
func (g *git) DeleteBranch(branch string, opts ...gitpkg.Option) error {
	cmd := g.newCommand("branch", "-d", branch)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Branch-specific options

// BranchWithRemote lists remote-tracking branches
func BranchWithRemote() gitpkg.Option {
	return WithArgs("-r")
}

// BranchWithAll lists both local and remote-tracking branches
func BranchWithAll() gitpkg.Option {
	return WithArgs("-a")
}

// BranchWithVerbose shows hash and commit subject line for each head
func BranchWithVerbose() gitpkg.Option {
	return WithArgs("-v")
}

// BranchWithVeryVerbose shows hash, commit subject line and upstream branch
func BranchWithVeryVerbose() gitpkg.Option {
	return WithArgs("-vv")
}

// BranchWithColor uses colors in output
func BranchWithColor(when string) gitpkg.Option {
	if when == "" {
		return WithArgs("--color")
	}
	return WithArgs("--color=" + when)
}

// BranchWithNoColor disables colors in output
func BranchWithNoColor() gitpkg.Option {
	return WithArgs("--no-color")
}

// BranchWithMerged shows only branches merged into the named commit
func BranchWithMerged(commit string) gitpkg.Option {
	if commit == "" {
		return WithArgs("--merged")
	}
	return WithArgs("--merged", commit)
}

// BranchWithNoMerged shows only branches not merged into the named commit
func BranchWithNoMerged(commit string) gitpkg.Option {
	if commit == "" {
		return WithArgs("--no-merged")
	}
	return WithArgs("--no-merged", commit)
}

// BranchWithContains shows only branches that contain the commit
func BranchWithContains(commit string) gitpkg.Option {
	return WithArgs("--contains", commit)
}

// BranchWithNoContains shows only branches that don't contain the commit
func BranchWithNoContains(commit string) gitpkg.Option {
	return WithArgs("--no-contains", commit)
}

// CreateBranchWithStartPoint creates branch starting from the specified commit
func CreateBranchWithStartPoint(startPoint string) gitpkg.Option {
	return WithArgs(startPoint)
}

// CreateBranchWithTrack sets up tracking when creating a branch
func CreateBranchWithTrack() gitpkg.Option {
	return WithArgs("--track")
}

// CreateBranchWithNoTrack doesn't set up tracking when creating a branch
func CreateBranchWithNoTrack() gitpkg.Option {
	return WithArgs("--no-track")
}

// CreateBranchWithForce forces creation (resets existing branch to start point)
func CreateBranchWithForce() gitpkg.Option {
	return WithArgs("--force")
}

// Delete branch-specific options

// DeleteBranchWithForce forces deletion of branch (even if not merged)
func DeleteBranchWithForce() gitpkg.Option {
	return WithArgs("-D")
}

// DeleteBranchWithRemote deletes remote-tracking branch
func DeleteBranchWithRemote() gitpkg.Option {
	return WithArgs("-r")
}