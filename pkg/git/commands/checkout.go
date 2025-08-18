package commands

import gitpkg "github.com/instruqt/git-exec/pkg/git"

// Checkout switches branches or restores working tree files
func (g *git) Checkout(opts ...gitpkg.Option) error {
	cmd := g.newCommand("checkout")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Checkout-specific options

// CheckoutWithBranch switches to a branch
func CheckoutWithBranch(branch string) gitpkg.Option {
	return WithArgs(branch)
}

// CheckoutWithNewBranch creates and switches to a new branch
func CheckoutWithNewBranch(branch string) gitpkg.Option {
	return WithArgs("-b", branch)
}

// CheckoutWithNewBranchForce creates and switches to a new branch, resetting if exists
func CheckoutWithNewBranchForce(branch string) gitpkg.Option {
	return WithArgs("-B", branch)
}

// CheckoutWithDetach detaches HEAD at the specified commit
func CheckoutWithDetach() gitpkg.Option {
	return WithArgs("--detach")
}

// CheckoutWithForce forces checkout (throw away local modifications)
func CheckoutWithForce() gitpkg.Option {
	return WithArgs("--force")
}

// CheckoutWithMerge performs a 3-way merge with the new branch
func CheckoutWithMerge() gitpkg.Option {
	return WithArgs("--merge")
}

// CheckoutWithConflict recreates conflicted merge in the specified paths
func CheckoutWithConflict(style string) gitpkg.Option {
	return WithArgs("--conflict=" + style)
}

// CheckoutWithPatch interactively select hunks in the diff
func CheckoutWithPatch() gitpkg.Option {
	return WithArgs("--patch")
}

// CheckoutWithIgnoreSkipWorktreeBits ignores sparse-checkout patterns
func CheckoutWithIgnoreSkipWorktreeBits() gitpkg.Option {
	return WithArgs("--ignore-skip-worktree-bits")
}

// CheckoutWithRecurseSubmodules updates submodules according to configuration
func CheckoutWithRecurseSubmodules() gitpkg.Option {
	return WithArgs("--recurse-submodules")
}

// CheckoutWithNoRecurseSubmodules don't update submodules
func CheckoutWithNoRecurseSubmodules() gitpkg.Option {
	return WithArgs("--no-recurse-submodules")
}

// CheckoutWithOverlayMode uses overlay mode (default)
func CheckoutWithOverlayMode() gitpkg.Option {
	return WithArgs("--overlay")
}

// CheckoutWithNoOverlayMode uses no-overlay mode
func CheckoutWithNoOverlayMode() gitpkg.Option {
	return WithArgs("--no-overlay")
}

// CheckoutWithQuiet suppresses feedback messages
func CheckoutWithQuiet() gitpkg.Option {
	return WithArgs("--quiet")
}

// CheckoutWithProgress shows progress status
func CheckoutWithProgress() gitpkg.Option {
	return WithArgs("--progress")
}

// CheckoutWithNoProgress hides progress status
func CheckoutWithNoProgress() gitpkg.Option {
	return WithArgs("--no-progress")
}

// CheckoutWithFiles checks out specific files from the index
func CheckoutWithFiles(files []string) gitpkg.Option {
	return func(c gitpkg.Command) {
		c.AddArgs(files...)
	}
}