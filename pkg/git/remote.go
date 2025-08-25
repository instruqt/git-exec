package git

import (
	"strings"
	"github.com/instruqt/git-exec/pkg/git/types"
)

// AddRemote adds a new remote repository
func (g *gitImpl) AddRemote(name, url string, opts ...Option) error {
	cmd := g.newCommand("remote", "add", name, url)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}

// RemoveRemote removes a remote repository
func (g *gitImpl) RemoveRemote(name string, opts ...Option) error {
	cmd := g.newCommand("remote", "remove", name)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}

// ListRemotes lists all remote repositories
func (g *gitImpl) ListRemotes(opts ...Option) ([]types.Remote, error) {
	cmd := g.newCommand("remote", "-v")
	cmd.ApplyOptions(opts...)
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}
	
	remotes := make(map[string]string)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			url := parts[1]
			remotes[name] = url
		}
	}
	
	result := []types.Remote{}
	for name, url := range remotes {
		result = append(result, types.Remote{Name: name, URL: url})
	}
	return result, nil
}

// SetRemoteURL sets the URL for a remote repository
func (g *gitImpl) SetRemoteURL(name, url string, opts ...Option) error {
	cmd := g.newCommand("remote", "set-url", name, url)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}