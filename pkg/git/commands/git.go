package commands

import (
	"os/exec"

	gitpkg "github.com/instruqt/git-exec/pkg/git"
)

// git implements the Git interface
type git struct {
	path string
	wd   string
}

// NewGit creates a new git implementation
func NewGit() (*git, error) {
	path, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	return &git{
		path: path,
	}, nil
}

// NewGitInstance creates a new Git instance (basic, no session)
func NewGitInstance() (gitpkg.Git, error) {
	return NewGit()
}

// SetWorkingDirectory sets the working directory for git operations
func (g *git) SetWorkingDirectory(wd string) {
	g.wd = wd
}

