package commands

import (
	"testing"

	"github.com/instruqt/git-exec/pkg/git/types"
	"github.com/stretchr/testify/require"
)

func TestParseDiffHeader(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected types.DiffHeader
		wantErr  bool
	}{
		{
			name:   "old and new mode",
			header: "old mode 100644\nnew mode 100755",
			expected: types.DiffHeader{
				OldMode: intPtr(100644),
				NewMode: intPtr(100755),
			},
		},
		{
			name:   "deleted file mode",
			header: "deleted file mode 100644",
			expected: types.DiffHeader{
				DeletedFileMode: intPtr(100644),
			},
		},
		{
			name:   "new file mode",
			header: "new file mode 100644",
			expected: types.DiffHeader{
				NewFileMode: intPtr(100644),
			},
		},
		{
			name:   "copy from/to",
			header: "copy from file1.txt\ncopy to file2.txt\nsimilarity index 100%",
			expected: types.DiffHeader{
				CopyFrom:        stringPtr("file1.txt"),
				CopyTo:          stringPtr("file2.txt"),
				SimilarityIndex: intPtr(100),
			},
		},
		{
			name:   "rename from/to",
			header: "rename from old.txt\nrename to new.txt\nsimilarity index 95%",
			expected: types.DiffHeader{
				RenameFrom:      stringPtr("old.txt"),
				RenameTo:        stringPtr("new.txt"),
				SimilarityIndex: intPtr(95),
			},
		},
		{
			name:   "dissimilarity index",
			header: "dissimilarity index 80%",
			expected: types.DiffHeader{
				DissimilarityIndex: intPtr(80),
			},
		},
		{
			name:   "index",
			header: "index abc123..def456 100644",
			expected: types.DiffHeader{
				Index: stringPtr("abc123..def456 100644"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDiffHeader(tt.header)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseDiffStats(t *testing.T) {
	input := `file1.go    | 10 +++++-----
file2.go    | 5 +++++
dir/file3.go | 2 --
4 files changed, 15 insertions(+), 7 deletions(-)`

	expected := []types.DiffStat{
		{File: "file1.go", Changes: 10, Insertions: 5, Deletions: 5},
		{File: "file2.go", Changes: 5, Insertions: 5, Deletions: 0},
		{File: "dir/file3.go", Changes: 2, Insertions: 0, Deletions: 2},
	}

	result, err := parseDiffStats(input)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestParseDiffModes(t *testing.T) {
	input := `M	file1.go
A	file2.go
D	file3.go
R100	old.txt	new.txt
C90	original.txt	copy.txt`

	expected := []types.DiffMode{
		{Action: "M", Mode: 0, File: "file1.go"},
		{Action: "A", Mode: 0, File: "file2.go"},
		{Action: "D", Mode: 0, File: "file3.go"},
		{Action: "R", Mode: 100, File: "old.txt"},
		{Action: "C", Mode: 90, File: "original.txt"},
	}

	result, err := parseDiffModes(input)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestParseRefs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []types.Ref
		wantErr  bool
	}{
		{
			name: "simple format that matches regex",
			input: `   abc123..def456  main  -> origin/main
 * [new branch]    feature  -> origin/feature`,
			expected: []types.Ref{
				{
					Status:  types.RefStatusFastForward,
					Summary: "abc123..def456",
					From:    "main",
					To:      "origin/main",
				},
				{
					Status:  types.RefStatusNew,
					Summary: "[new branch]",
					From:    "feature",
					To:      "origin/feature",
				},
			},
		},
		{
			name: "forced update and rejected",
			input: ` + abc123...def456 main       -> origin/main (forced update)
 ! [rejected]      feature    -> origin/feature (non-fast-forward)`,
			expected: []types.Ref{
				{
					Status:  types.RefStatusForcedUpdate,
					Summary: "abc123...def456",
					From:    "main",
					To:      "origin/main",
					Reason:  stringPtr("forced update"),
				},
				{
					Status:  types.RefStatusRejected,
					Summary: "[rejected]",
					From:    "feature",
					To:      "origin/feature",
					Reason:  stringPtr("non-fast-forward"),
				},
			},
		},
		{
			name: "tag update and up to date",
			input: ` t [tag update]    v1.0.0     -> v1.0.0
 = [up to date]    main       -> origin/main`,
			expected: []types.Ref{
				{
					Status:  types.RefStatusTagUpdate,
					Summary: "[tag update]",
					From:    "v1.0.0",
					To:      "v1.0.0",
				},
				{
					Status:  types.RefStatusUpToDate,
					Summary: "[up to date]",
					From:    "main",
					To:      "origin/main",
				},
			},
		},
		{
			name: "pruned ref",
			input: ` - [deleted]       refs/remotes/origin/old-branch`,
			expected: []types.Ref{
				{
					Status:  types.RefStatusPruned,
					Summary: "[deleted] refs/remotes/origin/old-branch",
					From:    "",
					To:      "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRefs(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseRemoteOutput(t *testing.T) {
	input := `To github.com:user/repo.git
   abc123..def456  main -> main
 * [new branch]    feature -> feature`

	expected := []types.Remote{
		{
			Name: "github.com:user/repo.git",
			Refs: []types.Ref{
				{
					Status:  types.RefStatusFastForward,
					Summary: "abc123..def456",
					From:    "main",
					To:      "main",
				},
				{
					Status:  types.RefStatusNew,
					Summary: "[new branch]",
					From:    "feature",
					To:      "feature",
				},
			},
		},
	}

	result, err := parseRemoteOutput(input)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestParseRemoteOutputMultipleRemotes(t *testing.T) {
	input := `To github.com:user/repo1.git
   abc123..def456  main -> main
To github.com:user/repo2.git
 * [new branch]    feature -> feature`

	expected := []types.Remote{
		{
			Name: "github.com:user/repo1.git",
			Refs: []types.Ref{
				{
					Status:  types.RefStatusFastForward,
					Summary: "abc123..def456",
					From:    "main",
					To:      "main",
				},
			},
		},
		{
			Name: "github.com:user/repo2.git",
			Refs: []types.Ref{
				{
					Status:  types.RefStatusNew,
					Summary: "[new branch]",
					From:    "feature",
					To:      "feature",
				},
			},
		},
	}

	result, err := parseRemoteOutput(input)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

// Helper functions for test readability
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}