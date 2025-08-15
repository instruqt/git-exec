package commands

import (
	"regexp"
	"time"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// Log shows commit logs
func (g *git) Log(opts ...Option) ([]types.Log, error) {
	cmd := g.newCommand("log", "--format=fuller")
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}

	detailsRegex := regexp.MustCompile(`commit (?P<commit>.+)\nAuthor:\s+(?P<author>.+)\nAuthorDate:\s+(?P<author_date>.+)\nCommit:\s+(?P<committer>.+)\nCommitDate:\s+(?P<committer_date>.+)\n\n\s+(?P<message>.*)(?:\n\n)?`)
	matches := detailsRegex.FindAllStringSubmatch(string(output), -1)

	logs := []types.Log{}
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

		logs = append(logs, types.Log{
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

// Log-specific options

// LogWithOneline shows each commit on a single line
func LogWithOneline() Option {
	return WithArgs("--oneline")
}

// LogWithGraph shows text-based graphical representation
func LogWithGraph() Option {
	return WithArgs("--graph")
}

// LogWithDecorate shows ref names
func LogWithDecorate() Option {
	return WithArgs("--decorate")
}

// LogWithAll shows all refs
func LogWithAll() Option {
	return WithArgs("--all")
}

// LogWithStat shows diffstat for each commit
func LogWithStat() Option {
	return WithArgs("--stat")
}

// LogWithShortStat shows only summary line of diffstat
func LogWithShortStat() Option {
	return WithArgs("--shortstat")
}

// LogWithNameOnly shows only names of changed files
func LogWithNameOnly() Option {
	return WithArgs("--name-only")
}

// LogWithNameStatus shows names and status of changed files
func LogWithNameStatus() Option {
	return WithArgs("--name-status")
}

// LogWithAbbrevCommit shows abbreviated commit hashes
func LogWithAbbrevCommit() Option {
	return WithArgs("--abbrev-commit")
}

// LogWithMaxCount limits number of commits
func LogWithMaxCount(count string) Option {
	return WithArgs("--max-count=" + count)
}

// LogWithSkip skips first N commits
func LogWithSkip(count string) Option {
	return WithArgs("--skip=" + count)
}

// LogWithSince shows commits after date
func LogWithSince(date string) Option {
	return WithArgs("--since=" + date)
}

// LogWithUntil shows commits before date
func LogWithUntil(date string) Option {
	return WithArgs("--until=" + date)
}

// LogWithAuthor shows commits by author
func LogWithAuthor(author string) Option {
	return WithArgs("--author=" + author)
}

// LogWithCommitter shows commits by committer
func LogWithCommitter(committer string) Option {
	return WithArgs("--committer=" + committer)
}

// LogWithGrep searches commit messages
func LogWithGrep(pattern string) Option {
	return WithArgs("--grep=" + pattern)
}

// LogWithFormat specifies output format
func LogWithFormat(format string) Option {
	return WithArgs("--format=" + format)
}

// LogWithPretty specifies pretty format
func LogWithPretty(format string) Option {
	return WithArgs("--pretty=" + format)
}

// LogWithNoMerges excludes merge commits
func LogWithNoMerges() Option {
	return WithArgs("--no-merges")
}

// LogWithMerges shows only merge commits
func LogWithMerges() Option {
	return WithArgs("--merges")
}

// LogWithFirstParent follows only first parent of merge commits
func LogWithFirstParent() Option {
	return WithArgs("--first-parent")
}

// LogWithReverse shows commits in reverse order
func LogWithReverse() Option {
	return WithArgs("--reverse")
}

// LogWithFollow continues listing history of file beyond renames
func LogWithFollow() Option {
	return WithArgs("--follow")
}

// LogWithPath limits to specific path
func LogWithPath(path string) Option {
	return WithArgs("--", path)
}