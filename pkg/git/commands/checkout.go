package commands

// Checkout switches branches or restores working tree files
func (g *git) Checkout(opts ...Option) error {
	cmd := g.newCommand("checkout")
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Checkout-specific options

// CheckoutWithBranch switches to a branch
func CheckoutWithBranch(branch string) Option {
	return WithArgs(branch)
}

// CheckoutWithNewBranch creates and switches to a new branch
func CheckoutWithNewBranch(branch string) Option {
	return WithArgs("-b", branch)
}

// CheckoutWithNewBranchForce creates and switches to a new branch, resetting if exists
func CheckoutWithNewBranchForce(branch string) Option {
	return WithArgs("-B", branch)
}

// CheckoutWithDetach detaches HEAD at the specified commit
func CheckoutWithDetach() Option {
	return WithArgs("--detach")
}

// CheckoutWithForce forces checkout (throw away local modifications)
func CheckoutWithForce() Option {
	return WithArgs("--force")
}

// CheckoutWithMerge performs a 3-way merge with the new branch
func CheckoutWithMerge() Option {
	return WithArgs("--merge")
}

// CheckoutWithConflict recreates conflicted merge in the specified paths
func CheckoutWithConflict(style string) Option {
	return WithArgs("--conflict=" + style)
}

// CheckoutWithPatch interactively select hunks in the diff
func CheckoutWithPatch() Option {
	return WithArgs("--patch")
}

// CheckoutWithIgnoreSkipWorktreeBits ignores sparse-checkout patterns
func CheckoutWithIgnoreSkipWorktreeBits() Option {
	return WithArgs("--ignore-skip-worktree-bits")
}

// CheckoutWithRecurseSubmodules updates submodules according to configuration
func CheckoutWithRecurseSubmodules() Option {
	return WithArgs("--recurse-submodules")
}

// CheckoutWithNoRecurseSubmodules don't update submodules
func CheckoutWithNoRecurseSubmodules() Option {
	return WithArgs("--no-recurse-submodules")
}

// CheckoutWithOverlayMode uses overlay mode (default)
func CheckoutWithOverlayMode() Option {
	return WithArgs("--overlay")
}

// CheckoutWithNoOverlayMode uses no-overlay mode
func CheckoutWithNoOverlayMode() Option {
	return WithArgs("--no-overlay")
}

// CheckoutWithQuiet suppresses feedback messages
func CheckoutWithQuiet() Option {
	return WithArgs("--quiet")
}

// CheckoutWithProgress shows progress status
func CheckoutWithProgress() Option {
	return WithArgs("--progress")
}

// CheckoutWithNoProgress hides progress status
func CheckoutWithNoProgress() Option {
	return WithArgs("--no-progress")
}

// CheckoutWithFiles checks out specific files from the index
func CheckoutWithFiles(files []string) Option {
	return func(c *Command) {
		c.args = append(c.args, files...)
	}
}