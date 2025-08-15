package commands

import (
	"regexp"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// AddRemote adds a remote named <name> for the repository at <url>
func (g *git) AddRemote(name, url string, opts ...Option) error {
	cmd := g.newCommand("remote", "add", name, url)
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// RemoveRemote removes the remote named <name>
func (g *git) RemoveRemote(name string, opts ...Option) error {
	cmd := g.newCommand("remote", "rm", name)
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// ListRemotes shows the remote repositories along with their URLs
func (g *git) ListRemotes(opts ...Option) ([]types.Remote, error) {
	cmd := g.newCommand("remote", "-v")
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}

	remotes := []types.Remote{}
	
	// Parse remote output in format: "name\turl (fetch/push)"
	regex := regexp.MustCompile(`(.+)\t(.+)\s\(fetch\)\n(?:.+)\t(?:.+)\s\(push\)`)
	results := regex.FindAllStringSubmatch(string(output), -1)

	for index := range results {
		name := results[index][1]
		url := results[index][2]

		remotes = append(remotes, types.Remote{
			Name: name,
			URL:  url,
		})
	}

	return remotes, nil
}

// SetRemoteURL changes the URL of a remote repository
func (g *git) SetRemoteURL(name, url string, opts ...Option) error {
	cmd := g.newCommand("remote", "set-url", name, url)
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Remote-specific options

// RemoteWithVerbose shows URLs after name
func RemoteWithVerbose() Option {
	return WithArgs("-v")
}

// RemoteAddWithTrack sets up branch tracking
func RemoteAddWithTrack(branch string) Option {
	return WithArgs("-t", branch)
}

// RemoteAddWithMaster sets the default branch
func RemoteAddWithMaster(branch string) Option {
	return WithArgs("-m", branch)
}

// RemoteAddWithFetch runs git fetch immediately after adding
func RemoteAddWithFetch() Option {
	return WithArgs("-f")
}

// SetRemoteURL-specific options

// SetRemoteURLWithPush sets push URL instead of fetch URL
func SetRemoteURLWithPush() Option {
	return WithArgs("--push")
}

// SetRemoteURLWithAdd adds URL instead of changing it
func SetRemoteURLWithAdd() Option {
	return WithArgs("--add")
}

// SetRemoteURLWithDelete deletes URL instead of changing it
func SetRemoteURLWithDelete() Option {
	return WithArgs("--delete")
}

// RemoteAddWithTags imports tags from the remote
func RemoteAddWithTags() Option {
	return WithArgs("--tags")
}

// RemoteAddWithNoTags doesn't import tags from the remote
func RemoteAddWithNoTags() Option {
	return WithArgs("--no-tags")
}

// RemoteAddWithMirror sets up mirroring mode
func RemoteAddWithMirror(mode string) Option {
	if mode == "" {
		return WithArgs("--mirror")
	}
	return WithArgs("--mirror=" + mode)
}