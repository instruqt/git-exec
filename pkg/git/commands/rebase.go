package commands
import gitpkg "github.com/instruqt/git-exec/pkg/git"

// Rebase reapplies commits on top of another base tip
func (g *git) Rebase(opts ...gitpkg.Option) error {
	// TODO: implement - currently has basic placeholder implementation
	cmd := g.newCommand("rebase")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Rebase-specific options

// RebaseWithUpstream rebases onto the specified upstream
func RebaseWithUpstream(upstream string) gitpkg.Option {
	return WithArgs(upstream)
}

// RebaseWithBranch rebases the specified branch
func RebaseWithBranch(branch string) gitpkg.Option {
	return WithArgs(branch)
}

// RebaseWithOnto transplants branch based on upstream onto another base
func RebaseWithOnto(newbase string) gitpkg.Option {
	return WithArgs("--onto", newbase)
}

// RebaseWithInteractive lets user edit the list of commits to rebase
func RebaseWithInteractive() gitpkg.Option {
	return WithArgs("--interactive")
}

// RebaseWithExec executes command on each commit
func RebaseWithExec(command string) gitpkg.Option {
	return WithArgs("--exec", command)
}

// RebaseWithRoot rebases all commits reachable from branch
func RebaseWithRoot() gitpkg.Option {
	return WithArgs("--root")
}

// RebaseWithAutosquash moves commits that begin with squash!/fixup! 
func RebaseWithAutosquash() gitpkg.Option {
	return WithArgs("--autosquash")
}

// RebaseWithNoAutosquash disables autosquash behavior
func RebaseWithNoAutosquash() gitpkg.Option {
	return WithArgs("--no-autosquash")
}

// RebaseWithAutostash stashes local changes before rebasing
func RebaseWithAutostash() gitpkg.Option {
	return WithArgs("--autostash")
}

// RebaseWithNoAutostash doesn't stash local changes
func RebaseWithNoAutostash() gitpkg.Option {
	return WithArgs("--no-autostash")
}

// RebaseWithKeepEmpty preserves empty commits
func RebaseWithKeepEmpty() gitpkg.Option {
	return WithArgs("--keep-empty")
}

// RebaseWithSkipEmpty skips empty commits
func RebaseWithSkipEmpty() gitpkg.Option {
	return WithArgs("--skip-empty")
}

// RebaseWithPreserveMerges preserves merge commits
func RebaseWithPreserveMerges() gitpkg.Option {
	return WithArgs("--preserve-merges")
}

// RebaseWithRebaseMerges rebases merge commits instead of ignoring them
func RebaseWithRebaseMerges() gitpkg.Option {
	return WithArgs("--rebase-merges")
}

// RebaseWithStrategy specifies merge strategy
func RebaseWithStrategy(strategy string) gitpkg.Option {
	return WithArgs("--strategy", strategy)
}

// RebaseWithStrategyOption passes option to merge strategy
func RebaseWithStrategyOption(option string) gitpkg.Option {
	return WithArgs("--strategy-option", option)
}

// RebaseWithQuiet suppresses output
func RebaseWithQuiet() gitpkg.Option {
	return WithArgs("--quiet")
}

// RebaseWithVerbose shows verbose output
func RebaseWithVerbose() gitpkg.Option {
	return WithArgs("--verbose")
}

// RebaseWithStat shows diffstat
func RebaseWithStat() gitpkg.Option {
	return WithArgs("--stat")
}

// RebaseWithNoStat doesn't show diffstat
func RebaseWithNoStat() gitpkg.Option {
	return WithArgs("--no-stat")
}

// RebaseWithVerify runs pre-rebase hook
func RebaseWithVerify() gitpkg.Option {
	return WithArgs("--verify")
}

// RebaseWithNoVerify skips pre-rebase hook
func RebaseWithNoVerify() gitpkg.Option {
	return WithArgs("--no-verify")
}

// RebaseWithContinue continues rebase after resolving conflicts
func RebaseWithContinue() gitpkg.Option {
	return WithArgs("--continue")
}

// RebaseWithSkip skips current patch and continues
func RebaseWithSkip() gitpkg.Option {
	return WithArgs("--skip")
}

// RebaseWithAbort aborts rebase and restores original branch
func RebaseWithAbort() gitpkg.Option {
	return WithArgs("--abort")
}

// RebaseWithQuit aborts rebase but leaves HEAD where it is
func RebaseWithQuit() gitpkg.Option {
	return WithArgs("--quit")
}

// RebaseWithEditTodo edits the todo list during interactive rebase
func RebaseWithEditTodo() gitpkg.Option {
	return WithArgs("--edit-todo")
}

// RebaseWithShowCurrentPatch shows current patch being applied
func RebaseWithShowCurrentPatch() gitpkg.Option {
	return WithArgs("--show-current-patch")
}