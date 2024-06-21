package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Git interface {
	SetWorkingDirectory(wd string)

	Init(options ...string) (string, error)
	AddRemote(url string, name string) error
	RemoveRemote(name string) error
	Clone(url string, options ...string) error
	Status() ([]File, error)
	Add(files ...string) error
	Reset(files ...string) error
	Commit(message string, author string, email string) error
	Diff() ([]Diff, error)
	Show() error
	Log() error
	Fetch() error
	Pull() error
	Push() error
	ListBranches() ([]Branch, error)
	CreateBranch(branch string) error
	Checkout() error
	Tag(name string) error
	Revert() error
	Merge() error
	Rebase() error
	Reflog() error
	Config() error
	Remove() error
}

type GitImpl struct {
	path string
	wd   string
}

func New() (Git, error) {
	path, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	return &GitImpl{
		path: path,
	}, nil
}

func (g *GitImpl) SetWorkingDirectory(wd string) {
	g.wd = wd
}

func (g *GitImpl) Init(options ...string) (string, error) {
	cmd := exec.Command(g.path, "init")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	cmd.Args = append(cmd.Args, options...)

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return "", fmt.Errorf(string(exitError.Stderr))
		}

		return "", err
	}

	return string(output), nil
}

func (g *GitImpl) AddRemote(name string, url string) error {
	cmd := exec.Command(g.path, "remote", "add", name, url)
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) RemoveRemote(name string) error {
	cmd := exec.Command(g.path, "remote", "rm", name)
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Clone(url string, options ...string) error {
	cmd := exec.Command(g.path, "clone", "-q", url)
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	cmd.Args = append(cmd.Args, options...)

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}

	return nil
}

func (g *GitImpl) Status() ([]File, error) {
	files := []File{}

	cmd := exec.Command(g.path, "status", "--porcelain")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}
			return files, fmt.Errorf(string(exitError.Stderr))
		}

		return files, err
	}

	lines := strings.Split(string(output), "\n")
	for _, file := range lines {
		if file == "" {
			continue
		}

		switch file[:2] {
		case "??":
			files = append(files, File{
				Status: FileUntracked,
				Name:   file[3:],
			})
		case "A ":
			files = append(files, File{
				Status: FileAdded,
				Name:   file[3:],
			})
		case "M ":
			files = append(files, File{
				Status: FileModified,
				Name:   file[3:],
			})
		case "D ":
			files = append(files, File{
				Status: FileDeleted,
				Name:   file[3:],
			})
		case "R ":
			parts := strings.Split(file[3:], " -> ")
			files = append(files, File{
				Status:      FileRenamed,
				Name:        parts[0],
				Destination: parts[1],
			})
		case "C ":
			parts := strings.Split(file[3:], " -> ")
			files = append(files, File{
				Status:      FileCopied,
				Name:        parts[0],
				Destination: parts[1],
			})
		case "U ":
			files = append(files, File{
				Status: FileUpdated,
				Name:   file[3:],
			})
		}

		// TODO add more cases
	}

	return files, nil
}

func (g *GitImpl) Commit(message string, author string, email string) error {
	cmd := exec.Command(g.path, "commit", "-q", "-m", message, fmt.Sprintf(`--author="%s <%s>"`, author, email))
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}

	return nil
}

func (g *GitImpl) Diff() ([]Diff, error) {
	diffs := []Diff{}
	cmd := exec.Command(g.path, "diff", "-U1000000", "--histogram")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}
			return diffs, fmt.Errorf(string(exitError.Stderr))
		}
		return diffs, err
	}

	// split the output of the individual diffs and capture the names of the files
	fileRegex := regexp.MustCompile(`diff --(.+) a/(.+) b/(.+)`)
	files := fileRegex.Split(string(output), -1)
	names := fileRegex.FindAllStringSubmatch(string(output), -1)

	nameIndex := 0
	for _, file := range files {
		if file == "" {
			continue
		}

		// split the diff into the header and the contents
		contentsRegex := regexp.MustCompile(`(?s)(.*)@@ -\d(?:,\d+)? \+\d(?:,\d+)? @@\n(?s)(.*)`)
		parts := contentsRegex.FindAllStringSubmatch(file, -1)

		if len(parts) == 0 {
			return diffs, fmt.Errorf("could not find diff")
		}

		if len(parts) > 1 {
			return diffs, fmt.Errorf("found more than one diff")
		}

		header := parts[0][1]
		contents := parts[0][2]

		dh, err := parseDiffHeader(header)
		if err != nil {
			return diffs, err
		}

		diffs = append(diffs, Diff{
			Format:   names[nameIndex][1],
			OldFile:  names[nameIndex][2],
			NewFile:  names[nameIndex][3],
			Header:   dh,
			Contents: contents,
		})

		nameIndex++
	}

	return diffs, nil
}

func (g *GitImpl) Show() error {
	cmd := exec.Command(g.path, "x")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Log() error {
	cmd := exec.Command(g.path, "x")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Fetch() error {
	cmd := exec.Command(g.path, "fetch")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Pull() error {
	cmd := exec.Command(g.path, "pull")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Push() error {
	cmd := exec.Command(g.path, "push")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) ListBranches() ([]Branch, error) {
	branches := []Branch{}

	cmd := exec.Command(g.path, "branch")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}
			return branches, fmt.Errorf(string(exitError.Stderr))
		}
		return branches, err
	}

	lines := strings.Split(string(output), "\n")
	for _, branch := range lines {
		if branch == "" {
			continue
		}

		branches = append(branches, Branch{
			Name:   branch[2:],
			Active: branch[0] == '*',
		})

	}

	return branches, nil
}

func (g *GitImpl) CreateBranch(branch string) error {
	cmd := exec.Command(g.path, "branch", branch)
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}
		return err
	}

	return nil
}

func (g *GitImpl) CheckoutFile(file string) error {
	cmd := exec.Command(g.path, "x")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Checkout() error {
	cmd := exec.Command(g.path, "x")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Tag(name string) error {
	cmd := exec.Command(g.path, "tag", name)
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Revert() error {
	cmd := exec.Command(g.path, "x")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Merge() error {
	cmd := exec.Command(g.path, "merge")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Rebase() error {
	cmd := exec.Command(g.path, "rebase")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Add(files ...string) error {
	cmd := exec.Command(g.path, "add")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	if len(files) > 0 {
		cmd.Args = append(cmd.Args, files...)
	} else {
		cmd.Args = append(cmd.Args, ".")
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}
		return err
	}

	return nil
}

func (g *GitImpl) Reset(files ...string) error {
	cmd := exec.Command(g.path, "reset")

	if len(files) > 0 {
		cmd.Args = append(cmd.Args, files...)
	}

	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}
		return err
	}

	return nil
}

func (g *GitImpl) Reflog() error {
	cmd := exec.Command(g.path, "x")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Config() error {
	cmd := exec.Command(g.path, "x")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *GitImpl) Remove() error {
	cmd := exec.Command(g.path, "x")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return fmt.Errorf(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func parseDiffHeader(header string) (DiffHeader, error) {
	dh := DiffHeader{}

	lines := strings.Split(header, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "old mode ") {
			value := strings.TrimPrefix(line, "old mode ")
			mode, err := strconv.Atoi(value)
			if err != nil {
				return dh, err
			}

			dh.OldMode = &mode
		}

		if strings.HasPrefix(line, "new mode ") {
			value := strings.TrimPrefix(line, "new mode ")
			mode, err := strconv.Atoi(value)
			if err != nil {
				return dh, err
			}

			dh.NewMode = &mode
		}

		if strings.HasPrefix(line, "deleted file mode ") {
			value := strings.TrimPrefix(line, "deleted file mode ")
			mode, err := strconv.Atoi(value)
			if err != nil {
				return dh, err
			}

			dh.DeletedFileMode = &mode
		}

		if strings.HasPrefix(line, "new file mode ") {
			value := strings.TrimPrefix(line, "new file mode ")
			mode, err := strconv.Atoi(value)
			if err != nil {
				return dh, err
			}

			dh.NewFileMode = &mode
		}

		if strings.HasPrefix(line, "copy from ") {
			value := strings.TrimPrefix(line, "copy from ")
			dh.CopyFrom = &value
		}

		if strings.HasPrefix(line, "copy to ") {
			value := strings.TrimPrefix(line, "copy to ")
			dh.CopyTo = &value
		}

		if strings.HasPrefix(line, "rename from ") {
			value := strings.TrimPrefix(line, "rename from ")
			dh.RenameFrom = &value
		}

		if strings.HasPrefix(line, "rename to ") {
			value := strings.TrimPrefix(line, "rename to ")
			dh.RenameTo = &value
		}

		if strings.HasPrefix(line, "similarity index ") {
			value := strings.TrimPrefix(line, "similarity index ")
			value = strings.TrimSuffix(value, "%")
			index, err := strconv.Atoi(value)
			if err != nil {
				return dh, err
			}

			dh.SimilarityIndex = &index
		}

		if strings.HasPrefix(line, "dissimilarity index ") {
			value := strings.TrimPrefix(line, "dissimilarity index ")
			value = strings.TrimSuffix(value, "%")
			index, err := strconv.Atoi(value)
			if err != nil {
				return dh, err
			}

			dh.DissimilarityIndex = &index
		}

		if strings.HasPrefix(line, "index ") {
			value := strings.TrimPrefix(line, "index ")
			dh.Index = &value
		}
	}

	return dh, nil
}
