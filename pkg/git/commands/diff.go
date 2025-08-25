package commands
import gitpkg "github.com/instruqt/git-exec/pkg/git"

import (
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Diff shows changes between commits, commit and working tree, etc
func (g *git) Diff(opts ...gitpkg.Option) ([]types.Diff, error) {
	cmd := g.newCommand("diff", "-U1000000", "--histogram")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
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
func DiffWithCached() gitpkg.Option {
	return WithArgs("--cached")
}

// DiffWithStaged shows changes between index and HEAD (alias for --cached)
func DiffWithStaged() gitpkg.Option {
	return WithArgs("--staged")
}

// DiffWithNameOnly shows only names of changed files
func DiffWithNameOnly() gitpkg.Option {
	return WithArgs("--name-only")
}

// DiffWithNameStatus shows names and status of changed files
func DiffWithNameStatus() gitpkg.Option {
	return WithArgs("--name-status")
}

// DiffWithStat shows diffstat
func DiffWithStat() gitpkg.Option {
	return WithArgs("--stat")
}

// DiffWithShortStat shows only the summary line of --stat
func DiffWithShortStat() gitpkg.Option {
	return WithArgs("--shortstat")
}

// DiffWithNumStat shows numeric diffstat
func DiffWithNumStat() gitpkg.Option {
	return WithArgs("--numstat")
}

// DiffWithPatch shows patch format (default)
func DiffWithPatch() gitpkg.Option {
	return WithArgs("--patch")
}

// DiffWithNoPatch suppresses diff output
func DiffWithNoPatch() gitpkg.Option {
	return WithArgs("--no-patch")
}

// DiffWithRaw shows raw format
func DiffWithRaw() gitpkg.Option {
	return WithArgs("--raw")
}

// DiffWithMinimal spends extra time to make sure smallest possible diff
func DiffWithMinimal() gitpkg.Option {
	return WithArgs("--minimal")
}

// DiffWithPatience uses patience diff algorithm
func DiffWithPatience() gitpkg.Option {
	return WithArgs("--patience")
}

// DiffWithHistogram uses histogram diff algorithm
func DiffWithHistogram() gitpkg.Option {
	return WithArgs("--histogram")
}

// DiffWithAlgorithm specifies diff algorithm
func DiffWithAlgorithm(algorithm string) gitpkg.Option {
	return WithArgs("--diff-algorithm=" + algorithm)
}

// DiffWithContext specifies number of context lines
func DiffWithContext(lines string) gitpkg.Option {
	return WithArgs("-U" + lines)
}

// DiffWithIgnoreSpaceAtEol ignores changes in whitespace at EOL
func DiffWithIgnoreSpaceAtEol() gitpkg.Option {
	return WithArgs("--ignore-space-at-eol")
}

// DiffWithIgnoreSpaceChange ignores changes in amount of whitespace
func DiffWithIgnoreSpaceChange() gitpkg.Option {
	return WithArgs("--ignore-space-change")
}

// DiffWithIgnoreAllSpace ignores whitespace when comparing lines
func DiffWithIgnoreAllSpace() gitpkg.Option {
	return WithArgs("--ignore-all-space")
}

// DiffWithIgnoreBlankLines ignores changes whose lines are all blank
func DiffWithIgnoreBlankLines() gitpkg.Option {
	return WithArgs("--ignore-blank-lines")
}

// DiffWithCommit compares with a specific commit
func DiffWithCommit(commit string) gitpkg.Option {
	return WithArgs(commit)
}

// DiffWithCommitRange compares a range of commits
func DiffWithCommitRange(from, to string) gitpkg.Option {
	return WithArgs(from + ".." + to)
}