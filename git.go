package git

import (
	"errors"
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
	ListRemotes() ([]Remote, error)
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
	Config(key string, value string) error
	Remove() error
}

type git struct {
	path string
	wd   string
}

func New() (Git, error) {
	path, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	return &git{
		path: path,
	}, nil
}

func (g *git) SetWorkingDirectory(wd string) {
	g.wd = wd
}

func (g *git) Init(options ...string) (string, error) {
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

			return "", errors.New(string(exitError.Stderr))
		}

		return "", err
	}

	return string(output), nil
}

func (g *git) AddRemote(name string, url string) error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *git) RemoveRemote(name string) error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *git) ListRemotes() ([]Remote, error) {
	remotes := []Remote{}

	cmd := exec.Command(g.path, "remote", "-v")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return remotes, errors.New(string(exitError.Stderr))
		}

		return remotes, err
	}

	regex := regexp.MustCompile(`(.+)\t(.+)\s\(fetch\)\n(?:.+)\t(?:.+)\s\(push\)`)
	results := regex.FindAllStringSubmatch(string(output), -1)

	for index := range results {
		name := results[index][1]
		url := results[index][2]

		remotes = append(remotes, Remote{
			Name: name,
			URL:  url,
		})
	}

	return remotes, nil
}

func (g *git) Clone(url string, options ...string) error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}

	return nil
}

func (g *git) Status() ([]File, error) {
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
			return files, errors.New(string(exitError.Stderr))
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
				Status: FileStatusUntracked,
				Name:   file[3:],
			})
		case "A ":
			files = append(files, File{
				Status: FileStatusAdded,
				Name:   file[3:],
			})
		case "M ":
			files = append(files, File{
				Status: FileStatusModified,
				Name:   file[3:],
			})
		case "D ":
			files = append(files, File{
				Status: FileStatusDeleted,
				Name:   file[3:],
			})
		case "R ":
			parts := strings.Split(file[3:], " -> ")
			files = append(files, File{
				Status:      FileStatusRenamed,
				Name:        parts[0],
				Destination: parts[1],
			})
		case "C ":
			parts := strings.Split(file[3:], " -> ")
			files = append(files, File{
				Status:      FileStatusCopied,
				Name:        parts[0],
				Destination: parts[1],
			})
		case "U ":
			files = append(files, File{
				Status: FileStatusUpdated,
				Name:   file[3:],
			})
		}

		// TODO: add more cases
	}

	return files, nil
}

func (g *git) Commit(message string, author string, email string) error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}

	return nil
}

func (g *git) Diff() ([]Diff, error) {
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
			return diffs, errors.New(string(exitError.Stderr))
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
			return diffs, errors.New("could not find diff")
		}

		if len(parts) > 1 {
			return diffs, errors.New("found more than one diff")
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

// TODO: implement
func (g *git) Show() error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Log() error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Fetch() error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Pull() error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Push() error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *git) ListBranches() ([]Branch, error) {
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
			return branches, errors.New(string(exitError.Stderr))
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

// TODO: implement
func (g *git) CreateBranch(branch string) error {
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

			return errors.New(string(exitError.Stderr))
		}
		return err
	}

	return nil
}

// TODO: implement
func (g *git) CheckoutFile(file string) error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Checkout() error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Tag(name string) error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Revert() error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Merge() error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Rebase() error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

func (g *git) Add(files ...string) error {
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

			return errors.New(string(exitError.Stderr))
		}
		return err
	}

	return nil
}

func (g *git) Reset(files ...string) error {
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

			return errors.New(string(exitError.Stderr))
		}
		return err
	}

	return nil
}

// TODO: implement
func (g *git) Reflog() error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Config(key string, value string) error {
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

			return errors.New(string(exitError.Stderr))
		}

		return err
	}
	return nil
}

// TODO: implement
func (g *git) Remove() error {
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

			return errors.New(string(exitError.Stderr))
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
