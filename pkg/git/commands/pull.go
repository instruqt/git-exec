package commands

import (
	"regexp"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// Pull incorporates changes from a remote repository into the current branch
func (g *git) Pull(opts ...Option) (*types.MergeResult, error) {
	cmd := g.newCommand("pull")
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}

	// Parse the pull output
	pullRegex := regexp.MustCompile(`Updating\s(?P<start_commit>\w+)..(?P<end_commit>\w+)\n(?P<method>\S+)\n(?P<files>(?:\s.+\|\s+\d+\s[\+\-]+\n)*)\s(?P<changes>\d+) file(?:s)? changed(?:, (?P<insertions>\d+) insertion(?:s)?\(\+\))?(?:, (?P<deletions>\d+) deletion(?:s)?\(\-\))?\n(?P<modes>(?:\s(?:create|delete) mode \d+ \S+\n?)*)?`)
	pullMatches := pullRegex.FindAllStringSubmatch(string(output), -1)

	result := &types.MergeResult{
		DiffStats: []types.DiffStat{},
		DiffModes: []types.DiffMode{},
	}

	for _, pullMatch := range pullMatches {
		pull := make(map[string]string)
		for i, name := range pullRegex.SubexpNames() {
			if i != 0 && name != "" {
				pull[name] = pullMatch[i]
			}
		}

		result.StartCommit = pull["start_commit"]
		result.EndCommit = pull["end_commit"]
		result.Method = pull["method"]

		diffStats, err := parseDiffStats(pull["files"])
		if err != nil {
			return nil, err
		}
		result.DiffStats = append(result.DiffStats, diffStats...)

		diffModes, err := parseDiffModes(pull["modes"])
		if err != nil {
			return nil, err
		}
		result.DiffModes = append(result.DiffModes, diffModes...)
	}

	return result, nil
}

// Pull-specific options

// PullWithRemote specifies which remote to pull from
func PullWithRemote(remote string) Option {
	return WithArgs(remote)
}

// PullWithBranch specifies which branch to pull
func PullWithBranch(branch string) Option {
	return WithArgs(branch)
}

// PullWithRemoteAndBranch specifies both remote and branch
func PullWithRemoteAndBranch(remote, branch string) Option {
	return WithArgs(remote, branch)
}

// PullWithRebase rebases the current branch on top of the upstream
func PullWithRebase() Option {
	return WithArgs("--rebase")
}

// PullWithNoRebase merges the upstream into the current branch (default)
func PullWithNoRebase() Option {
	return WithArgs("--no-rebase")
}

// PullWithFFOnly only updates if the merge can be resolved as a fast-forward
func PullWithFFOnly() Option {
	return WithArgs("--ff-only")
}

// PullWithNoFF creates a merge commit even when fast-forward is possible
func PullWithNoFF() Option {
	return WithArgs("--no-ff")
}

// PullWithSquash creates a single commit instead of merging
func PullWithSquash() Option {
	return WithArgs("--squash")
}

// PullWithStrategy specifies merge strategy
func PullWithStrategy(strategy string) Option {
	return WithArgs("--strategy", strategy)
}

// PullWithAll pulls from all remotes
func PullWithAll() Option {
	return WithArgs("--all")
}

// PullWithTags fetches tags as well
func PullWithTags() Option {
	return WithArgs("--tags")
}

// PullWithPrune removes remote-tracking references that no longer exist
func PullWithPrune() Option {
	return WithArgs("--prune")
}