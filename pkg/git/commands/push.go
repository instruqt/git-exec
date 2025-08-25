package commands
import gitpkg "github.com/instruqt/git-exec/pkg/git"

import (
	"regexp"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// Push updates remote refs along with associated objects
func (g *git) Push(opts ...gitpkg.Option) ([]types.Remote, error) {
	cmd := g.newCommand("push")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	// Push writes its output to stderr
	output, err := cmd.ExecuteWithStderr()
	if err != nil {
		return nil, err
	}

	// Parse the push output
	remotesRegex := regexp.MustCompile(`To (?P<remote>.+)\n`)
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

// Push-specific options

// PushWithRemote specifies which remote to push to
func PushWithRemote(remote string) gitpkg.Option {
	return WithArgs(remote)
}

// PushWithBranch specifies which branch to push
func PushWithBranch(branch string) gitpkg.Option {
	return WithArgs(branch)
}

// PushWithRemoteAndBranch specifies both remote and branch
func PushWithRemoteAndBranch(remote, branch string) gitpkg.Option {
	return WithArgs(remote, branch)
}

// PushWithForce forces the push even if it would overwrite commits
func PushWithForce() gitpkg.Option {
	return WithArgs("--force")
}

// PushWithForceWithLease forces push but ensures no work is lost
func PushWithForceWithLease() gitpkg.Option {
	return WithArgs("--force-with-lease")
}

// PushWithAll pushes all branches
func PushWithAll() gitpkg.Option {
	return WithArgs("--all")
}

// PushWithTags pushes all tags
func PushWithTags() gitpkg.Option {
	return WithArgs("--tags")
}

// PushWithFollowTags pushes tags that are reachable from the pushed commits
func PushWithFollowTags() gitpkg.Option {
	return WithArgs("--follow-tags")
}

// PushWithSetUpstream sets up tracking relationship
func PushWithSetUpstream() gitpkg.Option {
	return WithArgs("--set-upstream")
}

// PushWithDryRun shows what would be pushed without actually pushing
func PushWithDryRun() gitpkg.Option {
	return WithArgs("--dry-run")
}

// PushWithDelete deletes the specified branch/tag from the remote
func PushWithDelete() gitpkg.Option {
	return WithArgs("--delete")
}

// PushWithAtomic ensures all refs are updated or none are
func PushWithAtomic() gitpkg.Option {
	return WithArgs("--atomic")
}

// PushWithPorcelain produces machine-readable output
func PushWithPorcelain() gitpkg.Option {
	return WithArgs("--porcelain")
}