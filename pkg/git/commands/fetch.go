package commands
import gitpkg "github.com/instruqt/git-exec/pkg/git"

import (
	"fmt"
	"regexp"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// Fetch downloads objects and refs from another repository
func (g *git) Fetch(opts ...gitpkg.Option) ([]types.Remote, error) {
	cmd := g.newCommand("fetch", "-v")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	// Fetch writes its output to stderr
	output, err := cmd.ExecuteWithStderr()
	if err != nil {
		return nil, err
	}

	// Parse the fetch output
	remotesRegex := regexp.MustCompile(`From (?P<remote>.+)\n`)
	remoteList := remotesRegex.Split(string(output), -1)
	remoteMatches := remotesRegex.FindAllStringSubmatch(string(output), -1)

	remotes := []types.Remote{}

	remoteIndex := 0
	for _, remoteItem := range remoteList {
		if remoteItem == "" {
			continue
		}

		remote := types.Remote{
			Name: remoteMatches[remoteIndex][1],
			Refs: []types.Ref{},
		}

		refs, err := parseRefs(remoteItem)
		if err != nil {
			return nil, err
		}

		remote.Refs = refs

		remotes = append(remotes, remote)
		remoteIndex++
	}

	return remotes, nil
}

// Fetch-specific options

// FetchWithAll fetches all remotes
func FetchWithAll() gitpkg.Option {
	return WithArgs("--all")
}

// FetchWithPrune removes remote-tracking references that no longer exist on the remote
func FetchWithPrune() gitpkg.Option {
	return WithArgs("--prune")
}

// FetchWithPruneTags removes local tags that no longer exist on the remote
func FetchWithPruneTags() gitpkg.Option {
	return WithArgs("--prune-tags")
}

// FetchWithTags fetches all tags from the remote
func FetchWithTags() gitpkg.Option {
	return WithArgs("--tags")
}

// FetchWithNoTags doesn't fetch any tags
func FetchWithNoTags() gitpkg.Option {
	return WithArgs("--no-tags")
}

// FetchWithDepth limits fetching to the specified number of commits
func FetchWithDepth(depth int) gitpkg.Option {
	return WithArgs("--depth", fmt.Sprintf("%d", depth))
}

// FetchWithRemote specifies which remote to fetch from
func FetchWithRemote(remote string) gitpkg.Option {
	return WithArgs(remote)
}

// FetchWithForce forces updates of local branches
func FetchWithForce() gitpkg.Option {
	return WithArgs("--force")
}

// FetchWithDryRun shows what would be done without making changes
func FetchWithDryRun() gitpkg.Option {
	return WithArgs("--dry-run")
}

// FetchWithRefetch re-fetches the entire repository
func FetchWithRefetch() gitpkg.Option {
	return WithArgs("--refetch")
}