package commands

// Merge joins two or more development histories together
func (g *git) Merge(opts ...Option) error {
	// TODO: implement - currently has basic placeholder implementation
	cmd := g.newCommand("merge")
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Merge-specific options

// MergeWithBranch merges the specified branch
func MergeWithBranch(branch string) Option {
	return WithArgs(branch)
}

// MergeWithCommit merges commits into current branch
func MergeWithCommit(commit string) Option {
	return WithArgs(commit)
}

// MergeWithNoFF creates a merge commit even when fast-forward is possible
func MergeWithNoFF() Option {
	return WithArgs("--no-ff")
}

// MergeWithFFOnly only updates if the merge can be resolved as fast-forward
func MergeWithFFOnly() Option {
	return WithArgs("--ff-only")
}

// MergeWithSquash creates a single commit instead of merging
func MergeWithSquash() Option {
	return WithArgs("--squash")
}

// MergeWithStrategy specifies merge strategy
func MergeWithStrategy(strategy string) Option {
	return WithArgs("--strategy", strategy)
}

// MergeWithStrategyOption passes option to merge strategy
func MergeWithStrategyOption(option string) Option {
	return WithArgs("--strategy-option", option)
}

// MergeWithEdit invokes editor to edit merge commit message
func MergeWithEdit() Option {
	return WithArgs("--edit")
}

// MergeWithNoEdit accepts auto-generated message
func MergeWithNoEdit() Option {
	return WithArgs("--no-edit")
}

// MergeWithMessage specifies merge commit message
func MergeWithMessage(message string) Option {
	return WithArgs("-m", message)
}

// MergeWithFile reads merge message from file
func MergeWithFile(file string) Option {
	return WithArgs("-F", file)
}

// MergeWithAbort aborts current merge and restores pre-merge state
func MergeWithAbort() Option {
	return WithArgs("--abort")
}

// MergeWithContinue continues merge after resolving conflicts
func MergeWithContinue() Option {
	return WithArgs("--continue")
}

// MergeWithQuiet suppresses output
func MergeWithQuiet() Option {
	return WithArgs("--quiet")
}

// MergeWithVerbose shows verbose output
func MergeWithVerbose() Option {
	return WithArgs("--verbose")
}

// MergeWithProgress shows progress
func MergeWithProgress() Option {
	return WithArgs("--progress")
}

// MergeWithNoProgress hides progress
func MergeWithNoProgress() Option {
	return WithArgs("--no-progress")
}

// MergeWithSign makes a GPG-signed merge commit
func MergeWithSign() Option {
	return WithArgs("--gpg-sign")
}

// MergeWithNoSign doesn't GPG-sign merge commit
func MergeWithNoSign() Option {
	return WithArgs("--no-gpg-sign")
}

// MergeWithLog includes one-line descriptions from merged commits
func MergeWithLog(n string) Option {
	if n == "" {
		return WithArgs("--log")
	}
	return WithArgs("--log=" + n)
}

// MergeWithNoLog doesn't include one-line descriptions
func MergeWithNoLog() Option {
	return WithArgs("--no-log")
}

// MergeWithStat shows diffstat at end of merge
func MergeWithStat() Option {
	return WithArgs("--stat")
}

// MergeWithNoStat doesn't show diffstat
func MergeWithNoStat() Option {
	return WithArgs("--no-stat")
}