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

func parseDiffStats(input string) (types.DiffStat, error) {
	ds := types.DiffStat{}

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

func parseDiffModes(input string) (types.DiffMode, error) {
	dm := types.DiffMode{}

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

func parseRefs(input string) ([]types.Ref, error) {
	refs := []types.Ref{}
	refRegex := regexp.MustCompile(`\s+(?P<status>.{1})\s+(?P<summary>\[up to date\]|\[new branch\]|\[new tag]|\S+)\s+\s(?P<from>\S+)\s+->\s+(?P<to>\S+)\n`)
	refMatches := refRegex.FindAllStringSubmatch(input, -1)

	for _, refMatch := range refMatches {
		ref := make(map[string]string)
		for i, name := range refRegex.SubexpNames() {
			if i != 0 && name != "" {
				ref[name] = refMatch[i]
			}
		}

		status := types.RefStatusUnspecified
		switch ref["status"] {
		case " ":
			status = types.RefStatusFastForward
		case "+":
			status = types.RefStatusForcedUpdate
		case "*":
			status = types.RefStatusNew
		case "-":
			status = types.RefStatusPruned
		case "!":
			status = types.RefStatusRejected
		case "=":
			status = types.RefStatusUpToDate
		case "t":
			status = types.RefStatusTagUpdate
		}

		refs = append(refs, types.Ref{
			Status:  status,
			Summary: ref["summary"],
			From:    ref["from"],
			To:      ref["to"],
			Reason:  nil,
		})
	}

	return refs, nil
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