package commands

import "fmt"

// Init initializes a new Git repository
func (g *git) Init(path string, opts ...Option) error {
	cmd := g.newCommand("init", path)
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Init-specific options

// InitWithBare creates a bare repository
func InitWithBare() Option {
	return WithArgs("--bare")
}

// InitWithQuiet suppresses output
func InitWithQuiet() Option {
	return WithArgs("--quiet")
}

// InitWithBranch sets the name of the initial branch
func InitWithBranch(branch string) Option {
	return WithArgs("--initial-branch", branch)
}

// InitWithTemplate specifies the template directory
func InitWithTemplate(template string) Option {
	return WithArgs("--template", template)
}

// InitWithSeparateGitDir creates the repository with a separate git directory
func InitWithSeparateGitDir(gitdir string) Option {
	return WithArgs("--separate-git-dir", gitdir)
}

// InitWithShared specifies that the repository is to be shared amongst group members
func InitWithShared(permissions string) Option {
	if permissions == "" {
		return WithArgs("--shared")
	}
	return WithArgs(fmt.Sprintf("--shared=%s", permissions))
}