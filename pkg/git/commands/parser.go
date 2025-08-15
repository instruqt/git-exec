package commands

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/instruqt/git-exec/pkg/git/types"
)

func parseDiffHeader(header string) (types.DiffHeader, error) {
	dh := types.DiffHeader{}

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

func parseDiffs(input string) ([]types.Diff, error) {
	diffs := []types.Diff{}
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

		diffs = append(diffs, types.Diff{
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

func parseDiffStats(input string) ([]types.DiffStat, error) {
	var stats []types.DiffStat

	// Match lines like: "file1.go    | 10 +++++-----"
	filesRegex := regexp.MustCompile(`(?m)^(?P<file>\S+.*?)\s+\|\s+(?P<changes>\d+)\s+(?P<insertions>\+*)(?P<deletions>-*)\s*$`)
	filesMatches := filesRegex.FindAllStringSubmatch(input, -1)

	for _, filesMatch := range filesMatches {
		files := make(map[string]string)
		for i, name := range filesRegex.SubexpNames() {
			if i != 0 && name != "" && i < len(filesMatch) {
				files[name] = filesMatch[i]
			}
		}

		changes, err := strconv.Atoi(files["changes"])
		if err != nil {
			return nil, err
		}

		ds := types.DiffStat{
			File:       strings.TrimSpace(files["file"]),
			Changes:    changes,
			Insertions: len(files["insertions"]),
			Deletions:  len(files["deletions"]),
		}
		
		stats = append(stats, ds)
	}

	return stats, nil
}

func parseDiffModes(input string) ([]types.DiffMode, error) {
	var modes []types.DiffMode

	// Handle Git's name-status output format
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Format: "M\tfile.go" or "R100\told.txt\tnew.txt"
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		status := parts[0]
		file := parts[1]

		// Extract numeric value from status (e.g., "R100" -> action="R", mode=100)
		var action string
		var mode int
		if len(status) > 1 && status[0] >= 'A' && status[0] <= 'Z' {
			action = string(status[0])
			if modeVal, err := strconv.Atoi(status[1:]); err == nil {
				mode = modeVal
			}
		} else {
			action = status
		}

		dm := types.DiffMode{
			Action: action,
			Mode:   mode,
			File:   file,
		}

		modes = append(modes, dm)
	}

	return modes, nil
}

func parseRefs(input string) ([]types.Ref, error) {
	refs := []types.Ref{}
	
	// Handle two different Git output formats:
	// 1. With arrows: " * [new branch]    feature -> origin/feature"
	// 2. Without arrows: " - [deleted]     refs/remotes/origin/old-branch"
	
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		// First try format with arrows (push/fetch refs)
		arrowRegex := regexp.MustCompile(`^\s*(?P<status>[ +*\-!=t])?\s*(?P<summary>\[[^\]]+\]|\S+(?:\.\.\.\S+)?)\s+(?P<from>\S+)\s+->\s+(?P<to>\S+)(?:\s+\((?P<reason>[^)]+)\))?`)
		if matches := arrowRegex.FindStringSubmatch(line); matches != nil {
			ref := make(map[string]string)
			for i, name := range arrowRegex.SubexpNames() {
				if i != 0 && name != "" && i < len(matches) {
					ref[name] = matches[i]
				}
			}
			
			status := parseRefStatus(ref["status"])
			refs = append(refs, types.Ref{
				Status:  status,
				Summary: ref["summary"],
				From:    ref["from"],
				To:      ref["to"],
				Reason:  stringToPtr(ref["reason"]),
			})
			continue
		}
		
		// Try format without arrows (deleted/pruned refs)
		noArrowRegex := regexp.MustCompile(`^\s*(?P<status>.)\s*(?P<summary>\[[^\]]+\])\s+(?P<ref>\S+)`)
		if matches := noArrowRegex.FindStringSubmatch(line); matches != nil {
			ref := make(map[string]string)
			for i, name := range noArrowRegex.SubexpNames() {
				if i != 0 && name != "" && i < len(matches) {
					ref[name] = matches[i]
				}
			}
			
			status := parseRefStatus(ref["status"])
			refs = append(refs, types.Ref{
				Status:  status,
				Summary: ref["summary"] + " " + ref["ref"],
				From:    "",
				To:      "",
				Reason:  nil,
			})
		}
	}

	return refs, nil
}

func parseRefStatus(statusChar string) types.RefStatus {
	switch statusChar {
	case " ", "":
		return types.RefStatusFastForward
	case "+":
		return types.RefStatusForcedUpdate
	case "*":
		return types.RefStatusNew
	case "-":
		return types.RefStatusPruned
	case "!":
		return types.RefStatusRejected
	case "=":
		return types.RefStatusUpToDate
	case "t":
		return types.RefStatusTagUpdate
	default:
		return types.RefStatusUnspecified
	}
}

func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// parseRemoteOutput parses push/fetch output to extract remote information
func parseRemoteOutput(output string) ([]types.Remote, error) {
	// Parse the push/fetch output that contains remote refs
	// This is used by both push and tag operations that interact with remotes
	remotesRegex := regexp.MustCompile(`To (?P<remote>.+)\n`)
	remoteList := remotesRegex.Split(output, -1)
	remoteMatches := remotesRegex.FindAllStringSubmatch(output, -1)

	remotes := []types.Remote{}

	remoteIndex := 0
	for _, remoteItem := range remoteList {
		if remoteItem == "" {
			continue
		}

		remote := types.Remote{
			Name: remoteMatches[remoteIndex][1],
			Refs: []types.Ref{},
		}

		refs, err := parseRefs(remoteItem)
		if err != nil {
			return nil, err
		}

		remote.Refs = refs
		remotes = append(remotes, remote)
		remoteIndex++
	}

	return remotes, nil
}