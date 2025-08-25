package commands
import gitpkg "github.com/instruqt/git-exec/pkg/git"

import "fmt"

// Init initializes a new Git repository
func (g *git) Init(path string, opts ...gitpkg.Option) error {
	cmd := g.newCommand("init", path)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Init-specific options

// InitWithBare creates a bare repository
func InitWithBare() gitpkg.Option {
	return WithArgs("--bare")
}

// InitWithQuiet suppresses output
func InitWithQuiet() gitpkg.Option {
	return WithArgs("--quiet")
}

// InitWithBranch sets the name of the initial branch
func InitWithBranch(branch string) gitpkg.Option {
	return WithArgs("--initial-branch", branch)
}

// InitWithTemplate specifies the template directory
func InitWithTemplate(template string) gitpkg.Option {
	return WithArgs("--template", template)
}

// InitWithSeparateGitDir creates the repository with a separate git directory
func InitWithSeparateGitDir(gitdir string) gitpkg.Option {
	return WithArgs("--separate-git-dir", gitdir)
}

// InitWithShared specifies that the repository is to be shared amongst group members
func InitWithShared(permissions string) gitpkg.Option {
	if permissions == "" {
		return WithArgs("--shared")
	}
	return WithArgs(fmt.Sprintf("--shared=%s", permissions))
}