package git

import (
	"strings"
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Status shows the working tree status
func (g *gitImpl) Status(opts ...Option) ([]types.File, error) {
	cmd := g.newCommand("status", "--porcelain")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}

	files := []types.File{}
	lines := strings.Split(string(output), "\n")
	for _, file := range lines {
		if file == "" {
			continue
		}

		switch file[:2] {
		case "??":
			files = append(files, types.File{
				Status: types.FileStatusUntracked,
				Name:   file[3:],
			})
		case "A ":
			files = append(files, types.File{
				Status: types.FileStatusAdded,
				Name:   file[3:],
			})
		case "M ":
			files = append(files, types.File{
				Status: types.FileStatusModified,
				Name:   file[3:],
			})
		case "D ":
			files = append(files, types.File{
				Status: types.FileStatusDeleted,
				Name:   file[3:],
			})
		case "R ":
			parts := strings.Split(file[3:], " -> ")
			files = append(files, types.File{
				Status:      types.FileStatusRenamed,
				Name:        parts[0],
				Destination: parts[1],
			})
		case "C ":
			parts := strings.Split(file[3:], " -> ")
			files = append(files, types.File{
				Status:      types.FileStatusCopied,
				Name:        parts[0],
				Destination: parts[1],
			})
		case "U ":
			files = append(files, types.File{
				Status: types.FileStatusUpdated,
				Name:   file[3:],
			})
		}
		// TODO: add more cases for other status codes
	}

	return files, nil
}