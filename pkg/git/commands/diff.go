package commands

import (
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Diff shows changes between commits, commit and working tree, etc
func (g *git) Diff(opts ...Option) ([]types.Diff, error) {
	cmd := g.newCommand("diff", "-U1000000", "--histogram")
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}

	diffs, err := parseDiffs(string(output))
	if err != nil {
		return nil, err
	}

	return diffs, nil
}

// Diff-specific options

// DiffWithCached shows changes between index and HEAD
func DiffWithCached() Option {
	return WithArgs("--cached")
}

// DiffWithStaged shows changes between index and HEAD (alias for --cached)
func DiffWithStaged() Option {
	return WithArgs("--staged")
}

// DiffWithNameOnly shows only names of changed files
func DiffWithNameOnly() Option {
	return WithArgs("--name-only")
}

// DiffWithNameStatus shows names and status of changed files
func DiffWithNameStatus() Option {
	return WithArgs("--name-status")
}

// DiffWithStat shows diffstat
func DiffWithStat() Option {
	return WithArgs("--stat")
}

// DiffWithShortStat shows only the summary line of --stat
func DiffWithShortStat() Option {
	return WithArgs("--shortstat")
}

// DiffWithNumStat shows numeric diffstat
func DiffWithNumStat() Option {
	return WithArgs("--numstat")
}

// DiffWithPatch shows patch format (default)
func DiffWithPatch() Option {
	return WithArgs("--patch")
}

// DiffWithNoPatch suppresses diff output
func DiffWithNoPatch() Option {
	return WithArgs("--no-patch")
}

// DiffWithRaw shows raw format
func DiffWithRaw() Option {
	return WithArgs("--raw")
}

// DiffWithMinimal spends extra time to make sure smallest possible diff
func DiffWithMinimal() Option {
	return WithArgs("--minimal")
}

// DiffWithPatience uses patience diff algorithm
func DiffWithPatience() Option {
	return WithArgs("--patience")
}

// DiffWithHistogram uses histogram diff algorithm
func DiffWithHistogram() Option {
	return WithArgs("--histogram")
}

// DiffWithAlgorithm specifies diff algorithm
func DiffWithAlgorithm(algorithm string) Option {
	return WithArgs("--diff-algorithm=" + algorithm)
}

// DiffWithContext specifies number of context lines
func DiffWithContext(lines string) Option {
	return WithArgs("-U" + lines)
}

// DiffWithIgnoreSpaceAtEol ignores changes in whitespace at EOL
func DiffWithIgnoreSpaceAtEol() Option {
	return WithArgs("--ignore-space-at-eol")
}

// DiffWithIgnoreSpaceChange ignores changes in amount of whitespace
func DiffWithIgnoreSpaceChange() Option {
	return WithArgs("--ignore-space-change")
}

// DiffWithIgnoreAllSpace ignores whitespace when comparing lines
func DiffWithIgnoreAllSpace() Option {
	return WithArgs("--ignore-all-space")
}

// DiffWithIgnoreBlankLines ignores changes whose lines are all blank
func DiffWithIgnoreBlankLines() Option {
	return WithArgs("--ignore-blank-lines")
}

// DiffWithCommit compares with a specific commit
func DiffWithCommit(commit string) Option {
	return WithArgs(commit)
}

// DiffWithCommitRange compares a range of commits
func DiffWithCommitRange(from, to string) Option {
	return WithArgs(from + ".." + to)
}