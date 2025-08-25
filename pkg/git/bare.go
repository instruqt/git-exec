package git

import (
	"strings"

	"github.com/instruqt/git-exec/pkg/git/types"
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

// UpdateRef updates a reference to point to a specific commit
func (g *gitImpl) UpdateRef(ref string, commit string, opts ...Option) error {
	cmd := g.newCommand("update-ref", ref, commit)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}

// DeleteRef deletes a reference
func (g *gitImpl) DeleteRef(ref string, opts ...Option) error {
	cmd := g.newCommand("update-ref", "-d", ref)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}

// ListRefs lists all references in the repository
func (g *gitImpl) ListRefs(opts ...Option) ([]types.Reference, error) {
	cmd := g.newCommand("for-each-ref", "--format=%(refname)%00%(objectname)")
	cmd.ApplyOptions(opts...)
	
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var refs []types.Reference
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.Split(line, "\x00")
		if len(parts) >= 2 {
			ref := types.Reference{
				Name:   parts[0],
				Commit: parts[1],
				Type:   determineReferenceType(parts[0]),
			}
			refs = append(refs, ref)
		}
	}
	
	return refs, nil
}

// determineReferenceType determines the type of reference based on its name
func determineReferenceType(refName string) types.ReferenceType {
	switch {
	case strings.HasPrefix(refName, "refs/heads/"):
		return types.ReferenceTypeBranch
	case strings.HasPrefix(refName, "refs/tags/"):
		return types.ReferenceTypeTag
	case strings.HasPrefix(refName, "refs/remotes/"):
		return types.ReferenceTypeRemote
	case strings.HasPrefix(refName, "refs/notes/"):
		return types.ReferenceTypeNote
	default:
		return types.ReferenceTypeOther
	}
}