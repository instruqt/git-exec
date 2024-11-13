package git

type FileStatus string

const (
	FileStatusUntracked  FileStatus = "untracked"
	FileStatusUnmodified FileStatus = "unmodified"
	FileStatusModified   FileStatus = "modified"
	FileStatusAdded      FileStatus = "added"
	FileStatusDeleted    FileStatus = "deleted"
	FileStatusRenamed    FileStatus = "renamed"
	FileStatusCopied     FileStatus = "copied"
	FileStatusUpdated    FileStatus = "updated"
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

type Remote struct {
	Name string
	URL  string
}
