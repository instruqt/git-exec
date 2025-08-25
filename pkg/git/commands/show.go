package commands
import gitpkg "github.com/instruqt/git-exec/pkg/git"

import (
	"regexp"
	"time"

	"github.com/instruqt/git-exec/pkg/git/types"
)

// Show shows information about a Git object (commit, tag, etc)
func (g *git) Show(object string, opts ...gitpkg.Option) (*types.Log, error) {
	cmd := g.newCommand("show", "--format=fuller", object)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	output, err := cmd.Execute()
	if err != nil {
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

	return &types.Log{
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

// Show-specific options

// ShowWithFormat specifies output format
func ShowWithFormat(format string) gitpkg.Option {
	return WithArgs("--format=" + format)
}

// ShowWithPretty specifies pretty format
func ShowWithPretty(format string) gitpkg.Option {
	return WithArgs("--pretty=" + format)
}

// ShowWithOneline shows each commit on a single line
func ShowWithOneline() gitpkg.Option {
	return WithArgs("--oneline")
}

// ShowWithShort shows short format
func ShowWithShort() gitpkg.Option {
	return WithArgs("--short")
}

// ShowWithMedium shows medium format (default)
func ShowWithMedium() gitpkg.Option {
	return WithArgs("--medium")
}

// ShowWithFull shows full format
func ShowWithFull() gitpkg.Option {
	return WithArgs("--full")
}

// ShowWithFuller shows fuller format
func ShowWithFuller() gitpkg.Option {
	return WithArgs("--fuller")
}

// ShowWithRaw shows raw format
func ShowWithRaw() gitpkg.Option {
	return WithArgs("--raw")
}

// ShowWithStat shows diffstat
func ShowWithStat() gitpkg.Option {
	return WithArgs("--stat")
}

// ShowWithNameOnly shows only names of changed files
func ShowWithNameOnly() gitpkg.Option {
	return WithArgs("--name-only")
}

// ShowWithNameStatus shows names and status of changed files
func ShowWithNameStatus() gitpkg.Option {
	return WithArgs("--name-status")
}

// ShowWithNoPatch suppresses diff output
func ShowWithNoPatch() gitpkg.Option {
	return WithArgs("--no-patch")
}

// ShowWithPatch shows patch format
func ShowWithPatch() gitpkg.Option {
	return WithArgs("--patch")
}

// ShowWithAbbrevCommit shows abbreviated commit hash
func ShowWithAbbrevCommit() gitpkg.Option {
	return WithArgs("--abbrev-commit")
}

// ShowWithNoAbbrevCommit shows full commit hash
func ShowWithNoAbbrevCommit() gitpkg.Option {
	return WithArgs("--no-abbrev-commit")
}