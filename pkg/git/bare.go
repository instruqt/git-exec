package git

import (
	"strings"
)

// IsBareRepository checks if the current repository is a bare repository
func (g *gitImpl) IsBareRepository() (bool, error) {
	cmd := g.newCommand("rev-parse", "--is-bare-repository")
	output, err := cmd.Execute()
	if err != nil {
		return false, err
	}
	
	result := strings.TrimSpace(string(output))
	return result == "true", nil
}