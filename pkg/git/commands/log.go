package commands
import gitpkg "github.com/instruqt/git-exec/pkg/git"

import (
	"regexp"
	"time"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// Log shows commit logs
func (g *git) Log(opts ...gitpkg.Option) ([]types.Log, error) {
	cmd := g.newCommand("log", "--format=fuller")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
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
func LogWithOneline() gitpkg.Option {
	return WithArgs("--oneline")
}

// LogWithGraph shows text-based graphical representation
func LogWithGraph() gitpkg.Option {
	return WithArgs("--graph")
}

// LogWithDecorate shows ref names
func LogWithDecorate() gitpkg.Option {
	return WithArgs("--decorate")
}

// LogWithAll shows all refs
func LogWithAll() gitpkg.Option {
	return WithArgs("--all")
}

// LogWithStat shows diffstat for each commit
func LogWithStat() gitpkg.Option {
	return WithArgs("--stat")
}

// LogWithShortStat shows only summary line of diffstat
func LogWithShortStat() gitpkg.Option {
	return WithArgs("--shortstat")
}

// LogWithNameOnly shows only names of changed files
func LogWithNameOnly() gitpkg.Option {
	return WithArgs("--name-only")
}

// LogWithNameStatus shows names and status of changed files
func LogWithNameStatus() gitpkg.Option {
	return WithArgs("--name-status")
}

// LogWithAbbrevCommit shows abbreviated commit hashes
func LogWithAbbrevCommit() gitpkg.Option {
	return WithArgs("--abbrev-commit")
}

// LogWithMaxCount limits number of commits
func LogWithMaxCount(count string) gitpkg.Option {
	return WithArgs("--max-count=" + count)
}

// LogWithSkip skips first N commits
func LogWithSkip(count string) gitpkg.Option {
	return WithArgs("--skip=" + count)
}

// LogWithSince shows commits after date
func LogWithSince(date string) gitpkg.Option {
	return WithArgs("--since=" + date)
}

// LogWithUntil shows commits before date
func LogWithUntil(date string) gitpkg.Option {
	return WithArgs("--until=" + date)
}

// LogWithAuthor shows commits by author
func LogWithAuthor(author string) gitpkg.Option {
	return WithArgs("--author=" + author)
}

// LogWithCommitter shows commits by committer
func LogWithCommitter(committer string) gitpkg.Option {
	return WithArgs("--committer=" + committer)
}

// LogWithGrep searches commit messages
func LogWithGrep(pattern string) gitpkg.Option {
	return WithArgs("--grep=" + pattern)
}

// LogWithFormat specifies output format
func LogWithFormat(format string) gitpkg.Option {
	return WithArgs("--format=" + format)
}

// LogWithPretty specifies pretty format
func LogWithPretty(format string) gitpkg.Option {
	return WithArgs("--pretty=" + format)
}

// LogWithNoMerges excludes merge commits
func LogWithNoMerges() gitpkg.Option {
	return WithArgs("--no-merges")
}

// LogWithMerges shows only merge commits
func LogWithMerges() gitpkg.Option {
	return WithArgs("--merges")
}

// LogWithFirstParent follows only first parent of merge commits
func LogWithFirstParent() gitpkg.Option {
	return WithArgs("--first-parent")
}

// LogWithReverse shows commits in reverse order
func LogWithReverse() gitpkg.Option {
	return WithArgs("--reverse")
}

// LogWithFollow continues listing history of file beyond renames
func LogWithFollow() gitpkg.Option {
	return WithArgs("--follow")
}

// LogWithPath limits to specific path
func LogWithPath(path string) gitpkg.Option {
	return WithArgs("--", path)
}