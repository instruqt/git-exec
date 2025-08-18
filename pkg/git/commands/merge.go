package commands
import gitpkg "github.com/instruqt/git-exec/pkg/git"

// Merge joins two or more development histories together
func (g *git) Merge(opts ...gitpkg.Option) error {
	// TODO: implement - currently has basic placeholder implementation
	cmd := g.newCommand("merge")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Merge-specific options

// MergeWithBranch merges the specified branch
func MergeWithBranch(branch string) gitpkg.Option {
	return WithArgs(branch)
}

// MergeWithCommit merges commits into current branch
func MergeWithCommit(commit string) gitpkg.Option {
	return WithArgs(commit)
}

// MergeWithNoFF creates a merge commit even when fast-forward is possible
func MergeWithNoFF() gitpkg.Option {
	return WithArgs("--no-ff")
}

// MergeWithFFOnly only updates if the merge can be resolved as fast-forward
func MergeWithFFOnly() gitpkg.Option {
	return WithArgs("--ff-only")
}

// MergeWithSquash creates a single commit instead of merging
func MergeWithSquash() gitpkg.Option {
	return WithArgs("--squash")
}

// MergeWithStrategy specifies merge strategy
func MergeWithStrategy(strategy string) gitpkg.Option {
	return WithArgs("--strategy", strategy)
}

// MergeWithStrategyOption passes option to merge strategy
func MergeWithStrategyOption(option string) gitpkg.Option {
	return WithArgs("--strategy-option", option)
}

// MergeWithEdit invokes editor to edit merge commit message
func MergeWithEdit() gitpkg.Option {
	return WithArgs("--edit")
}

// MergeWithNoEdit accepts auto-generated message
func MergeWithNoEdit() gitpkg.Option {
	return WithArgs("--no-edit")
}

// MergeWithMessage specifies merge commit message
func MergeWithMessage(message string) gitpkg.Option {
	return WithArgs("-m", message)
}

// MergeWithFile reads merge message from file
func MergeWithFile(file string) gitpkg.Option {
	return WithArgs("-F", file)
}

// MergeWithAbort aborts current merge and restores pre-merge state
func MergeWithAbort() gitpkg.Option {
	return WithArgs("--abort")
}

// MergeWithContinue continues merge after resolving conflicts
func MergeWithContinue() gitpkg.Option {
	return WithArgs("--continue")
}

// MergeWithQuiet suppresses output
func MergeWithQuiet() gitpkg.Option {
	return WithArgs("--quiet")
}

// MergeWithVerbose shows verbose output
func MergeWithVerbose() gitpkg.Option {
	return WithArgs("--verbose")
}

// MergeWithProgress shows progress
func MergeWithProgress() gitpkg.Option {
	return WithArgs("--progress")
}

// MergeWithNoProgress hides progress
func MergeWithNoProgress() gitpkg.Option {
	return WithArgs("--no-progress")
}

// MergeWithSign makes a GPG-signed merge commit
func MergeWithSign() gitpkg.Option {
	return WithArgs("--gpg-sign")
}

// MergeWithNoSign doesn't GPG-sign merge commit
func MergeWithNoSign() gitpkg.Option {
	return WithArgs("--no-gpg-sign")
}

// MergeWithLog includes one-line descriptions from merged commits
func MergeWithLog(n string) gitpkg.Option {
	if n == "" {
		return WithArgs("--log")
	}
	return WithArgs("--log=" + n)
}

// MergeWithNoLog doesn't include one-line descriptions
func MergeWithNoLog() gitpkg.Option {
	return WithArgs("--no-log")
}

// MergeWithStat shows diffstat at end of merge
func MergeWithStat() gitpkg.Option {
	return WithArgs("--stat")
}

// MergeWithNoStat doesn't show diffstat
func MergeWithNoStat() gitpkg.Option {
	return WithArgs("--no-stat")
}