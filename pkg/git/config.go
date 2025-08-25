package git

import (
	"regexp"
	"strings"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// SetConfig sets a git configuration value
func (g *gitImpl) SetConfig(key string, value string, opts ...Option) error {
	cmd := g.newCommand("config")
	cmd.ApplyOptions(opts...)
	cmd.AddArgs(key, value)
	_, err := cmd.Execute()
	return err
}

// GetConfig retrieves a git configuration value
func (g *gitImpl) GetConfig(key string, opts ...Option) (string, error) {
	cmd := g.newCommand("config")
	cmd.ApplyOptions(opts...)
	cmd.AddArgs("--get", key)
	output, err := cmd.Execute()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// ListConfig lists all git configuration entries
func (g *gitImpl) ListConfig(opts ...Option) ([]types.ConfigEntry, error) {
	cmd := g.newCommand("config")
	cmd.ApplyOptions(opts...)
	cmd.AddArgs("--list", "--show-origin", "--show-scope")
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}

	return parseConfigList(string(output))
}

// UnsetConfig removes a git configuration value
func (g *gitImpl) UnsetConfig(key string, opts ...Option) error {
	cmd := g.newCommand("config")
	cmd.ApplyOptions(opts...)
	cmd.AddArgs("--unset", key)
	_, err := cmd.Execute()
	return err
}

// parseConfigList parses the output of `git config --list --show-origin --show-scope`
func parseConfigList(output string) ([]types.ConfigEntry, error) {
	var entries []types.ConfigEntry
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Pattern to match: scope	file:source	key=value
	// Example: local	file:.git/config	user.name=John Doe
	pattern := regexp.MustCompile(`^(local|global|system)\t([^\t]+)\t([^=]+)=(.*)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := pattern.FindStringSubmatch(line)
		if len(matches) != 5 {
			// Skip malformed lines
			continue
		}

		scope := types.ConfigScope(matches[1])
		source := strings.TrimPrefix(matches[2], "file:")
		key := matches[3]
		value := matches[4]

		entries = append(entries, types.ConfigEntry{
			Key:    key,
			Value:  value,
			Scope:  scope,
			Source: source,
		})
	}

	return entries, nil
}