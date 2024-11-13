package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	Show(object string) (*Log, error)
	Log() ([]Log, error)
	Fetch() ([]Ref, error)
	Pull() (*MergeResult, error)
	ListBranches() ([]Branch, error)

	Push() error
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
	cmd := exec.Command(g.path, "commit", "-q", "-m", message)

	cmd.Env = append(cmd.Env, fmt.Sprintf("GIT_AUTHOR_NAME=%s", author))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GIT_AUTHOR_EMAIL=%s", email))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GIT_COMMITTER_NAME=%s", author))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GIT_COMMITTER_EMAIL=%s", email))

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
			return nil, errors.New(string(exitError.Stderr))
		}
		return nil, err
	}

	diffs, err := parseDiffs(string(output))
	if err != nil {
		return nil, err
	}

	return diffs, nil
}

func (g *git) Show(object string) (*Log, error) {
	cmd := exec.Command(g.path, "show", "--format=fuller", object)
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return nil, errors.New(string(exitError.Stderr))
		}

		return nil, err
	}

	detailsRegex := regexp.MustCompile(`(?s)commit (?P<commit>.+)\nAuthor:\s+(?P<author>.+)\nAuthorDate:\s+(?P<author_date>.+)\nCommit:\s+(?P<committer>.+)\nCommitDate:\s+(?P<committer_date>.+)\n\n\s+(?P<message>.*)\n\n(?P<diff>.*)`)
	matches := detailsRegex.FindAllStringSubmatch(string(output), -1)

	details := make(map[string]string)
	for i, name := range detailsRegex.SubexpNames() {
		if i != 0 && name != "" {
			details[name] = matches[0][i]
		}
	}

	authorDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", details["author_date"])
	if err != nil {
		return nil, err
	}

	committerDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", details["committer_date"])
	if err != nil {
		return nil, err
	}

	diffs, err := parseDiffs(details["diff"])
	if err != nil {
		return nil, err
	}

	return &Log{
		Commit:        details["commit"],
		Tree:          "", // These show up on the cli..but not when running with golang..
		Parent:        "", // These show up on the cli..but not when running with golang..
		Author:        details["author"],
		AuthorDate:    authorDate,
		Message:       details["message"],
		Committer:     details["committer"],
		CommitterDate: committerDate,
		Diffs:         diffs,
	}, nil
}

func (g *git) Log() ([]Log, error) {
	cmd := exec.Command(g.path, "log", "--format=fuller")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = output
			}

			return nil, errors.New(string(exitError.Stderr))
		}

		return nil, err
	}

	detailsRegex := regexp.MustCompile(`commit (?P<commit>.+)\nAuthor:\s+(?P<author>.+)\nAuthorDate:\s+(?P<author_date>.+)\nCommit:\s+(?P<committer>.+)\nCommitDate:\s+(?P<committer_date>.+)\n\n\s+(?P<message>.*)(?:\n\n)?`)
	matches := detailsRegex.FindAllStringSubmatch(string(output), -1)

	logs := []Log{}
	for _, match := range matches {
		details := make(map[string]string)
		for i, name := range detailsRegex.SubexpNames() {
			if i != 0 && name != "" {
				details[name] = match[i]
			}
		}

		authorDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", details["author_date"])
		if err != nil {
			return nil, err
		}

		committerDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", details["committer_date"])
		if err != nil {
			return nil, err
		}

		logs = append(logs, Log{
			Commit:        details["commit"],
			Tree:          "", // These show up on the cli..but not when running with golang..
			Parent:        "", // These show up on the cli..but not when running with golang..
			Author:        details["author"],
			AuthorDate:    authorDate,
			Message:       details["message"],
			Committer:     details["committer"],
			CommitterDate: committerDate,
		})
	}

	return logs, nil
}

// TODO: implement
func (g *git) Fetch() ([]Ref, error) {
	cmd := exec.Command(g.path, "fetch", "-v") //, "--refetch")
	if g.wd != "" {
		cmd.Dir = g.wd
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Stderr == nil {
				exitError.Stderr = stderr.Bytes()
			}

			return nil, errors.New(string(exitError.Stderr))
		}

		return nil, err
	}

	// fetch writes its output on stderr
	output := stderr.Bytes()

	remotesRegex := regexp.MustCompile(`From (?P<remote>.+)\n`)
	remotes := remotesRegex.Split(string(output), -1)
	remoteMatches := remotesRegex.FindAllStringSubmatch(string(output), -1)

	refs := []Ref{}

	remoteIndex := 0
	for _, remote := range remotes {
		if remote == "" {
			continue
		}

		remoteName := remoteMatches[remoteIndex][1]

		// TODO: handle reason in case of a rejected ref
		refRegex := regexp.MustCompile(`\s+(?P<status>.{1})\s+(?P<summary>\[up to date\]|\[new branch\]|\[new tag]|\S+)\s+\s(?P<from>\S+)\s+->\s+(?P<to>\S+)\n`)
		refMatches := refRegex.FindAllStringSubmatch(remote, -1)

		for _, refMatch := range refMatches {
			ref := make(map[string]string)
			for i, name := range refRegex.SubexpNames() {
				if i != 0 && name != "" {
					ref[name] = refMatch[i]
				}
			}

			status := RefStatusUnspecified
			switch ref["status"] {
			case " ":
				status = RefStatusFastForward
			case "+":
				status = RefStatusForcedUpdate
			case "*":
				status = RefStatusNew
			case "-":
				status = RefStatusPruned
			case "!":
				status = RefStatusRejected
			case "=":
				status = RefStatusUpToDate
			case "t":
				status = RefStatusTagUpdate
			}

			refs = append(refs, Ref{
				Remote:  remoteName,
				Status:  status,
				Summary: ref["summary"],
				From:    ref["from"],
				To:      ref["to"],
				Reason:  nil,
			})
		}

		remoteIndex++
	}

	return refs, nil
}

func (g *git) Pull() (*MergeResult, error) {
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

			return nil, errors.New(string(exitError.Stderr))
		}

		return nil, err
	}

	pullRegex := regexp.MustCompile(`Updating\s(?P<start_commit>\w+)..(?P<end_commit>\w+)\n(?P<method>\S+)\n(?P<files>(?:\s.+\|\s+\d+\s[\+\-]+\n)*)\s(?P<changes>\d+) file(?:s)? changed(?:, (?P<insertions>\d+) insertion(?:s)?\(\+\))?(?:, (?P<deletions>\d+) deletion(?:s)?\(\-\))?\n(?P<modes>(?:\s(?:create|delete) mode \d+ \S+\n?)*)?`)
	pullMatches := pullRegex.FindAllStringSubmatch(string(output), -1)

	result := &MergeResult{
		DiffStats: []DiffStat{},
		DiffModes: []DiffMode{},
	}

	for _, pullMatch := range pullMatches {
		pull := make(map[string]string)
		for i, name := range pullRegex.SubexpNames() {
			if i != 0 && name != "" {
				pull[name] = pullMatch[i]
			}
		}

		result.StartCommit = pull["start_commit"]
		result.EndCommit = pull["end_commit"]
		result.Method = pull["method"]

		diffStat, err := parseDiffStats(pull["files"])
		if err != nil {
			return nil, err
		}
		result.DiffStats = append(result.DiffStats, diffStat)

		diffMode, err := parseDiffModes(pull["modes"])
		if err != nil {
			return nil, err
		}
		result.DiffModes = append(result.DiffModes, diffMode)
	}

	return result, nil
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

func parseDiffs(input string) ([]Diff, error) {
	diffs := []Diff{}
	// split the output of the individual diffs and capture the names of the files
	fileRegex := regexp.MustCompile(`diff --(.+) a/(.+) b/(.+)`)
	files := fileRegex.Split(input, -1)
	names := fileRegex.FindAllStringSubmatch(input, -1)

	nameIndex := 0
	for _, file := range files {
		if file == "" {
			continue
		}

		// split the diff into the header and the contents
		contentsRegex := regexp.MustCompile(`(?s)(.*)@@ -\d+(?:,\d+)? \+\d+(?:,\d+)? @@\n?(?s)(?P<contents>.*)`)
		parts := contentsRegex.FindAllStringSubmatch(file, -1)

		if len(parts) == 0 {
			return nil, errors.New("could not find diff")
		}

		if len(parts) > 1 {
			return nil, errors.New("found more than one diff")
		}

		header := parts[0][1]
		contents := parts[0][2]

		dh, err := parseDiffHeader(header)
		if err != nil {
			return nil, err
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

func parseDiffStats(input string) (DiffStat, error) {
	ds := DiffStat{}

	filesRegex := regexp.MustCompile(`\s(?P<file>.+)\|\s+(?P<changes>\d+)\s(?P<insertions>\+*)?(?P<deletions>-*)?\n`)
	filesMatches := filesRegex.FindAllStringSubmatch(input, -1)

	for _, filesMatch := range filesMatches {
		files := make(map[string]string)
		for i, name := range filesRegex.SubexpNames() {
			if i != 0 && name != "" {
				files[name] = filesMatch[i]
			}
		}

		ds.File = strings.Trim(files["file"], " ")
		changes, err := strconv.Atoi(files["changes"])
		if err != nil {
			return ds, err
		}

		ds.Changes = changes
		ds.Insertions = len(files["insertions"])
		ds.Deletions = len(files["deletions"])
	}

	return ds, nil
}

func parseDiffModes(input string) (DiffMode, error) {
	dm := DiffMode{}

	modesRegex := regexp.MustCompile(`\s(?P<action>create|delete)\s+mode\s(?P<mode>\d+)\s(?P<file>\S+)\n`)
	modesMatches := modesRegex.FindAllStringSubmatch(input, -1)

	for _, modesMatch := range modesMatches {
		modes := make(map[string]string)
		for i, name := range modesRegex.SubexpNames() {
			if i != 0 && name != "" {
				modes[name] = modesMatch[i]
			}
		}

		mode, err := strconv.Atoi(modes["mode"])
		if err != nil {
			return dm, err
		}

		dm.Action = modes["action"]
		dm.Mode = mode
		dm.File = modes["file"]
	}

	return dm, nil
}
