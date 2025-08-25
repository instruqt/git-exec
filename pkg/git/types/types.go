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
	Success          bool
	FastForward      bool
	MergeCommit      string
	MergedBranch     string
	BaseBranch       string
	Strategy         string
	ConflictedFiles  []string
	Conflicts        []ConflictFile
	Stats            MergeStats
	AbortReason      string
}

type MergeStats struct {
	FilesChanged int
	Insertions   int
	Deletions    int
}

type ConflictFile struct {
	Path      string
	Status    ConflictStatus
	Sections  []ConflictSection
	Content   string // Raw file content with conflict markers
}

type ConflictStatus string

const (
	ConflictStatusBothModified ConflictStatus = "both_modified"
	ConflictStatusAddedByUs    ConflictStatus = "added_by_us"
	ConflictStatusAddedByThem  ConflictStatus = "added_by_them"
	ConflictStatusDeletedByUs  ConflictStatus = "deleted_by_us"
	ConflictStatusDeletedByThem ConflictStatus = "deleted_by_them"
	ConflictStatusBothAdded    ConflictStatus = "both_added"
	ConflictStatusBothDeleted  ConflictStatus = "both_deleted"
)

type ConflictSection struct {
	StartLine    int
	EndLine      int
	OurContent   string
	TheirContent string
	BaseContent  string // Available with diff3 conflict style
	Resolved     bool
	Resolution   string // User's resolution
}

type ConflictResolution struct {
	FilePath   string
	Sections   []ResolvedSection
	UseOurs    bool   // Use our version entirely
	UseTheirs  bool   // Use their version entirely
	Custom     bool   // Use custom resolution
}

type ResolvedSection struct {
	SectionIndex int
	Resolution   string // The resolved content for this section
}

// ConfigEntry represents a git configuration entry
type ConfigEntry struct {
	Key    string      // Configuration key (e.g., "user.name")
	Value  string      // Configuration value
	Scope  ConfigScope // Scope where this config is defined
	Source string      // File path where config is defined
}

// ConfigScope represents the scope of a git configuration
type ConfigScope string

const (
	ConfigScopeLocal  ConfigScope = "local"  // Repository-specific config
	ConfigScopeGlobal ConfigScope = "global" // User-specific config
	ConfigScopeSystem ConfigScope = "system" // System-wide config
)

// CheckoutResult represents the result of a checkout operation
type CheckoutResult struct {
	Success          bool     // Whether checkout was successful
	PreviousHEAD     string   // Previous HEAD reference
	NewHEAD          string   // New HEAD reference
	Branch           string   // Branch name (if checking out a branch)
	Commit           string   // Commit hash (if checking out a commit)
	DetachedHEAD     bool     // Whether in detached HEAD state
	NewBranch        bool     // Whether a new branch was created
	ModifiedFiles    []string // Files that were modified during checkout
	RestoredFiles    []string // Files that were restored during checkout
	UntrackedFiles   []string // Untracked files that would be overwritten
	Warning          string   // Any warning messages
	UpstreamBranch   string   // Upstream branch set for new branches
}

// ReferenceType represents the type of git reference
type ReferenceType string

const (
	ReferenceTypeBranch ReferenceType = "branch"
	ReferenceTypeTag    ReferenceType = "tag"
	ReferenceTypeRemote ReferenceType = "remote"
	ReferenceTypeNote   ReferenceType = "note"
	ReferenceTypeOther  ReferenceType = "other"
)

// Reference represents a git reference
type Reference struct {
	Name   string        // Reference name (e.g., refs/heads/main)
	Commit string        // Commit hash the reference points to
	Type   ReferenceType // Reference type (tag, branch, etc.)
}
