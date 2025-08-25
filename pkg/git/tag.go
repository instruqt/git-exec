package git

import (
	"fmt"
	"strings"
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Tag creates a new tag
func (g *gitImpl) Tag(name string, opts ...Option) error {
	cmd := g.newCommand("tag", name)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}

// ListTags lists all tags in the repository
func (g *gitImpl) ListTags(opts ...Option) ([]string, error) {
	cmd := g.newCommand("tag", "-l")
	cmd.ApplyOptions(opts...)
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}
	if len(output) == 0 {
		return []string{}, nil
	}
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	result := []string{}
	for _, tag := range tags {
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result, nil
}

// DeleteTag deletes a tag
func (g *gitImpl) DeleteTag(name string, opts ...Option) error {
	cmd := g.newCommand("tag", "-d", name)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}

// PushTags pushes all tags to the remote
func (g *gitImpl) PushTags(remote string, opts ...Option) ([]types.Remote, error) {
	cmd := g.newCommand("push", remote, "--tags")
	cmd.ApplyOptions(opts...)
	_, err := cmd.ExecuteWithStderr()
	if err != nil {
		return nil, err
	}
	return []types.Remote{}, nil
}

// DeleteRemoteTag deletes a tag from the remote repository
func (g *gitImpl) DeleteRemoteTag(remote, tagName string, opts ...Option) error {
	refspec := fmt.Sprintf(":refs/tags/%s", tagName)
	cmd := g.newCommand("push", remote, refspec)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}