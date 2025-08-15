package commands

// Rebase reapplies commits on top of another base tip
func (g *git) Rebase(opts ...Option) error {
	// TODO: implement - currently has basic placeholder implementation
	cmd := g.newCommand("rebase")
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Rebase-specific options

// RebaseWithUpstream rebases onto the specified upstream
func RebaseWithUpstream(upstream string) Option {
	return WithArgs(upstream)
}

// RebaseWithBranch rebases the specified branch
func RebaseWithBranch(branch string) Option {
	return WithArgs(branch)
}

// RebaseWithOnto transplants branch based on upstream onto another base
func RebaseWithOnto(newbase string) Option {
	return WithArgs("--onto", newbase)
}

// RebaseWithInteractive lets user edit the list of commits to rebase
func RebaseWithInteractive() Option {
	return WithArgs("--interactive")
}

// RebaseWithExec executes command on each commit
func RebaseWithExec(command string) Option {
	return WithArgs("--exec", command)
}

// RebaseWithRoot rebases all commits reachable from branch
func RebaseWithRoot() Option {
	return WithArgs("--root")
}

// RebaseWithAutosquash moves commits that begin with squash!/fixup! 
func RebaseWithAutosquash() Option {
	return WithArgs("--autosquash")
}

// RebaseWithNoAutosquash disables autosquash behavior
func RebaseWithNoAutosquash() Option {
	return WithArgs("--no-autosquash")
}

// RebaseWithAutostash stashes local changes before rebasing
func RebaseWithAutostash() Option {
	return WithArgs("--autostash")
}

// RebaseWithNoAutostash doesn't stash local changes
func RebaseWithNoAutostash() Option {
	return WithArgs("--no-autostash")
}

// RebaseWithKeepEmpty preserves empty commits
func RebaseWithKeepEmpty() Option {
	return WithArgs("--keep-empty")
}

// RebaseWithSkipEmpty skips empty commits
func RebaseWithSkipEmpty() Option {
	return WithArgs("--skip-empty")
}

// RebaseWithPreserveMerges preserves merge commits
func RebaseWithPreserveMerges() Option {
	return WithArgs("--preserve-merges")
}

// RebaseWithRebaseMerges rebases merge commits instead of ignoring them
func RebaseWithRebaseMerges() Option {
	return WithArgs("--rebase-merges")
}

// RebaseWithStrategy specifies merge strategy
func RebaseWithStrategy(strategy string) Option {
	return WithArgs("--strategy", strategy)
}

// RebaseWithStrategyOption passes option to merge strategy
func RebaseWithStrategyOption(option string) Option {
	return WithArgs("--strategy-option", option)
}

// RebaseWithQuiet suppresses output
func RebaseWithQuiet() Option {
	return WithArgs("--quiet")
}

// RebaseWithVerbose shows verbose output
func RebaseWithVerbose() Option {
	return WithArgs("--verbose")
}

// RebaseWithStat shows diffstat
func RebaseWithStat() Option {
	return WithArgs("--stat")
}

// RebaseWithNoStat doesn't show diffstat
func RebaseWithNoStat() Option {
	return WithArgs("--no-stat")
}

// RebaseWithVerify runs pre-rebase hook
func RebaseWithVerify() Option {
	return WithArgs("--verify")
}

// RebaseWithNoVerify skips pre-rebase hook
func RebaseWithNoVerify() Option {
	return WithArgs("--no-verify")
}

// RebaseWithContinue continues rebase after resolving conflicts
func RebaseWithContinue() Option {
	return WithArgs("--continue")
}

// RebaseWithSkip skips current patch and continues
func RebaseWithSkip() Option {
	return WithArgs("--skip")
}

// RebaseWithAbort aborts rebase and restores original branch
func RebaseWithAbort() Option {
	return WithArgs("--abort")
}

// RebaseWithQuit aborts rebase but leaves HEAD where it is
func RebaseWithQuit() Option {
	return WithArgs("--quit")
}

// RebaseWithEditTodo edits the todo list during interactive rebase
func RebaseWithEditTodo() Option {
	return WithArgs("--edit-todo")
}

// RebaseWithShowCurrentPatch shows current patch being applied
func RebaseWithShowCurrentPatch() Option {
	return WithArgs("--show-current-patch")
}