package types

import "time"

type FileStatus string

type RefStatus string

const (
	FileStatusUnspecified FileStatus = "unspecified"
	FileStatusUntracked   FileStatus = "untracked"
	FileStatusUnmodified  FileStatus = "unmodified"
	FileStatusModified    FileStatus = "modified"
	FileStatusAdded       FileStatus = "added"
	FileStatusDeleted     FileStatus = "deleted"
	FileStatusRenamed     FileStatus = "renamed"
	FileStatusCopied      FileStatus = "copied"
	FileStatusUpdated     FileStatus = "updated"

	RefStatusUnspecified  RefStatus = "unspecified"
	RefStatusFastForward  RefStatus = "fast_forward"  // " "
	RefStatusForcedUpdate RefStatus = "forced_update" // "+"
	RefStatusPruned       RefStatus = "pruned"        // "-"
	RefStatusTagUpdate    RefStatus = "tag_update"    // "t"
	RefStatusNew          RefStatus = "new"           // "*"
	RefStatusRejected     RefStatus = "rejected"      // "!"
	RefStatusUpToDate     RefStatus = "up_to_date"    // "="
)

type File struct {
	Status      FileStatus
	Name        string
	Destination string
}

type Branch struct {
	Name   string
	Active bool
}

type Diff struct {
	Format   string
	OldFile  string
	NewFile  string
	Header   DiffHeader
	Contents string
}

type DiffHeader struct {
	OldMode            *int
	NewMode            *int
	DeletedFileMode    *int
	NewFileMode        *int
	CopyFrom           *string
	CopyTo             *string
	RenameFrom         *string
	RenameTo           *string
	SimilarityIndex    *int
	DissimilarityIndex *int
	Index              *string
}

type DiffStat struct {
	File       string
	Changes    int
	Insertions int
	Deletions  int
}

type DiffMode struct {
	Action string
	Mode   int
	File   string
}

type Remote struct {
	Name string
	URL  string
	Refs []Ref
}

type Log struct {
	Commit        string
	Tree          string
	Parent        string
	Author        string
	AuthorDate    time.Time
	Message       string
	Committer     string
	CommitterDate time.Time
	Diffs         []Diff
}

type Ref struct {
	Status  RefStatus
	Summary string
	From    string
	To      string
	Reason  *string
}

type MergeResult struct {
	StartCommit string
	EndCommit   string
	Method      string
	DiffStats   []DiffStat
	DiffModes   []DiffMode
}
