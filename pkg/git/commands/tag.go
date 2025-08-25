package commands
import gitpkg "github.com/instruqt/git-exec/pkg/git"

import (
	"strings"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// Tag creates, lists, deletes or verifies a tag object
func (g *git) Tag(name string, opts ...gitpkg.Option) error {
	cmd := g.newCommand("tag", name)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// ListTags lists all tags in the repository
func (g *git) ListTags(opts ...gitpkg.Option) ([]string, error) {
	cmd := g.newCommand("tag", "-l")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}

	tags := []string{}
	lines := strings.Split(string(output), "\n")
	for _, tag := range lines {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

// DeleteTag deletes a local tag
func (g *git) DeleteTag(name string, opts ...gitpkg.Option) error {
	cmd := g.newCommand("tag", "-d", name)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// PushTags pushes tags to remote
func (g *git) PushTags(remote string, opts ...gitpkg.Option) ([]types.Remote, error) {
	cmd := g.newCommand("push", remote, "--tags")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	// Push writes its output to stderr
	output, err := cmd.ExecuteWithStderr()
	if err != nil {
		return nil, err
	}

	// Parse push output (reuse logic from push.go)
	// TODO: Extract parseRemoteRefs to shared parser function
	return parseRemoteOutput(string(output))
}

// DeleteRemoteTag deletes a tag from remote repository
func (g *git) DeleteRemoteTag(remote, tagName string, opts ...gitpkg.Option) error {
	cmd := g.newCommand("push", remote, "--delete", "refs/tags/"+tagName)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.ExecuteWithStderr()
	return err
}

// Tag-specific options

// TagWithAnnotated creates an annotated tag
func TagWithAnnotated() gitpkg.Option {
	return WithArgs("-a")
}

// TagWithMessage specifies tag message
func TagWithMessage(message string) gitpkg.Option {
	return WithArgs("-m", message)
}

// TagWithFile reads message from file
func TagWithFile(file string) gitpkg.Option {
	return WithArgs("-F", file)
}

// TagWithSign makes a GPG-signed tag
func TagWithSign() gitpkg.Option {
	return WithArgs("-s")
}

// TagWithLocalUser uses specific GPG key
func TagWithLocalUser(keyid string) gitpkg.Option {
	return WithArgs("-u", keyid)
}

// TagWithForce replaces existing tag
func TagWithForce() gitpkg.Option {
	return WithArgs("-f")
}

// TagWithDelete deletes existing tags
func TagWithDelete() gitpkg.Option {
	return WithArgs("-d")
}

// TagWithVerify verifies GPG signature of given tags
func TagWithVerify() gitpkg.Option {
	return WithArgs("-v")
}

// TagWithList lists tags
func TagWithList() gitpkg.Option {
	return WithArgs("-l")
}

// TagWithSort sorts tags
func TagWithSort(key string) gitpkg.Option {
	return WithArgs("--sort=" + key)
}

// TagWithMerged shows only tags merged into the named commit
func TagWithMerged(commit string) gitpkg.Option {
	if commit == "" {
		return WithArgs("--merged")
	}
	return WithArgs("--merged", commit)
}

// TagWithNoMerged shows only tags not merged into the named commit
func TagWithNoMerged(commit string) gitpkg.Option {
	if commit == "" {
		return WithArgs("--no-merged")
	}
	return WithArgs("--no-merged", commit)
}

// TagWithContains shows only tags that contain the commit
func TagWithContains(commit string) gitpkg.Option {
	return WithArgs("--contains", commit)
}

// TagWithNoContains shows only tags that don't contain the commit
func TagWithNoContains(commit string) gitpkg.Option {
	return WithArgs("--no-contains", commit)
}

// TagWithPoints shows only tags that point at the object
func TagWithPoints(object string) gitpkg.Option {
	return WithArgs("--points-at", object)
}

// TagWithFormat specifies output format
func TagWithFormat(format string) gitpkg.Option {
	return WithArgs("--format=" + format)
}

// TagWithColor uses colors in output
func TagWithColor(when string) gitpkg.Option {
	if when == "" {
		return WithArgs("--color")
	}
	return WithArgs("--color=" + when)
}

// TagWithNoColor disables colors in output
func TagWithNoColor() gitpkg.Option {
	return WithArgs("--no-color")
}

// TagWithColumn displays tags in columns
func TagWithColumn(options string) gitpkg.Option {
	if options == "" {
		return WithArgs("--column")
	}
	return WithArgs("--column=" + options)
}

// TagWithNoColumn disables column output
func TagWithNoColumn() gitpkg.Option {
	return WithArgs("--no-column")
}

// TagWithObject creates tag for specific object
func TagWithObject(object string) gitpkg.Option {
	return WithArgs(object)
}