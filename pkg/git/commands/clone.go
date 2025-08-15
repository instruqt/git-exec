package commands

import "fmt"

// Clone clones a repository to a destination path
func (g *git) Clone(url, destination string, opts ...Option) error {
	cmd := g.newCommand("clone", url, destination)
	
	// Apply quiet by default
	cmd.args = append(cmd.args, "-q")
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Clone-specific options

// CloneWithBare adds the --bare flag for bare repository clones
func CloneWithBare() Option {
	return WithArgs("--bare")
}

// CloneWithBranch specifies which branch to clone
func CloneWithBranch(branch string) Option {
	return WithArgs("--branch", branch)
}

// CloneWithDepth creates a shallow clone with history truncated to the specified number of commits
func CloneWithDepth(depth int) Option {
	return WithArgs("--depth", fmt.Sprintf("%d", depth))
}

// CloneWithSingleBranch clones only the history leading to the tip of a single branch
func CloneWithSingleBranch() Option {
	return WithArgs("--single-branch")
}

// CloneWithNoCheckout performs a clone without checking out a working tree
func CloneWithNoCheckout() Option {
	return WithArgs("--no-checkout")
}

// CloneWithRecurseSubmodules initializes and clones submodules recursively
func CloneWithRecurseSubmodules() Option {
	return WithArgs("--recurse-submodules")
}

// CloneWithShallow creates a shallow clone (depth 1)
func CloneWithShallow() Option {
	return CloneWithDepth(1)
}

// CloneWithMirror sets up the clone as a mirror (implies --bare)
func CloneWithMirror() Option {
	return WithArgs("--mirror")
}

// CloneWithReference uses a local repository as a reference to reduce network transfer
func CloneWithReference(repo string) Option {
	return WithArgs("--reference", repo)
}

// CloneWithConfig sets a configuration variable in the newly created clone
func CloneWithConfig(key, value string) Option {
	return WithArgs("--config", fmt.Sprintf("%s=%s", key, value))
}