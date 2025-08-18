package commands

import (
	"fmt"
	
	gitpkg "github.com/instruqt/git-exec/pkg/git"
)

// Clone clones a repository to a destination path
func (g *git) Clone(url, destination string, opts ...gitpkg.Option) error {
	cmd := g.newCommand("clone", url, destination)
	
	// Apply quiet by default
	cmd.AddArgs("-q")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Clone-specific options

// CloneWithBare adds the --bare flag for bare repository clones
func CloneWithBare() gitpkg.Option {
	return WithArgs("--bare")
}

// CloneWithBranch specifies which branch to clone
func CloneWithBranch(branch string) gitpkg.Option {
	return WithArgs("--branch", branch)
}

// CloneWithDepth creates a shallow clone with history truncated to the specified number of commits
func CloneWithDepth(depth int) gitpkg.Option {
	return WithArgs("--depth", fmt.Sprintf("%d", depth))
}

// CloneWithSingleBranch clones only the history leading to the tip of a single branch
func CloneWithSingleBranch() gitpkg.Option {
	return WithArgs("--single-branch")
}

// CloneWithNoCheckout performs a clone without checking out a working tree
func CloneWithNoCheckout() gitpkg.Option {
	return WithArgs("--no-checkout")
}

// CloneWithRecurseSubmodules initializes and clones submodules recursively
func CloneWithRecurseSubmodules() gitpkg.Option {
	return WithArgs("--recurse-submodules")
}

// CloneWithShallow creates a shallow clone (depth 1)
func CloneWithShallow() gitpkg.Option {
	return CloneWithDepth(1)
}

// CloneWithMirror sets up the clone as a mirror (implies --bare)
func CloneWithMirror() gitpkg.Option {
	return WithArgs("--mirror")
}

// CloneWithReference uses a local repository as a reference to reduce network transfer
func CloneWithReference(repo string) gitpkg.Option {
	return WithArgs("--reference", repo)
}

// CloneWithConfig sets a configuration variable in the newly created clone
func CloneWithConfig(key, value string) gitpkg.Option {
	return WithArgs("--config", fmt.Sprintf("%s=%s", key, value))
}