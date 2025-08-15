package commands

import (
	"strings"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// Status shows the working tree status
func (g *git) Status(opts ...Option) ([]types.File, error) {
	cmd := g.newCommand("status", "--porcelain")
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
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

// Status-specific options

// StatusWithShort gives output in short format
func StatusWithShort() Option {
	return WithArgs("--short")
}

// StatusWithBranch shows branch information
func StatusWithBranch() Option {
	return WithArgs("--branch")
}

// StatusWithPorcelain gives porcelain output (default for this implementation)
func StatusWithPorcelain() Option {
	return WithArgs("--porcelain")
}

// StatusWithLong gives output in long format (default Git behavior)
func StatusWithLong() Option {
	return WithArgs("--long")
}

// StatusWithShowStash shows stash information
func StatusWithShowStash() Option {
	return WithArgs("--show-stash")
}

// StatusWithAheadBehind shows ahead/behind counts
func StatusWithAheadBehind() Option {
	return WithArgs("--ahead-behind")
}

// StatusWithUntrackedFiles controls how untracked files are shown
func StatusWithUntrackedFiles(mode string) Option {
	return WithArgs("--untracked-files=" + mode)
}

// StatusWithIgnoredFiles shows ignored files
func StatusWithIgnoredFiles() Option {
	return WithArgs("--ignored")
}