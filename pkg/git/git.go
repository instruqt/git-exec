package git

import (
	"os/exec"
	"time"

	"github.com/instruqt/git-exec/pkg/git/types"
)

type Git interface {
	SetWorkingDirectory(wd string)

	Init(path string, options ...Option) error
	AddRemote(name, url string, options ...Option) error
	RemoveRemote(name string, options ...Option) error
	SetRemoteURL(name, url string, options ...Option) error
	ListRemotes(options ...Option) ([]types.Remote, error)
	Clone(url, destination string, options ...Option) error
	Status(options ...Option) ([]types.File, error)
	Add(files []string, options ...Option) error
	Reset(files []string, options ...Option) error
	Commit(message string, options ...Option) error
	Diff(options ...Option) ([]types.Diff, error)
	Show(object string, options ...Option) (*types.Log, error)
	Log(options ...Option) ([]types.Log, error)
	Fetch(options ...Option) ([]types.Remote, error)
	Pull(options ...Option) (*types.MergeResult, error)
	Push(options ...Option) ([]types.Remote, error)
	ListBranches(options ...Option) ([]types.Branch, error)
	CreateBranch(branch string, options ...Option) error
	DeleteBranch(branch string, options ...Option) error
	SetUpstream(branch string, remote string, options ...Option) error
	Checkout(options ...Option) (*types.CheckoutResult, error)
	Tag(name string, options ...Option) error
	ListTags(options ...Option) ([]string, error)
	DeleteTag(name string, options ...Option) error
	PushTags(remote string, options ...Option) ([]types.Remote, error)
	DeleteRemoteTag(remote, tagName string, options ...Option) error
	Revert(options ...Option) error
	Merge(options ...Option) (*types.MergeResult, error)
	MergeAbort() error
	MergeContinue() error
	ResolveConflicts(resolutions []types.ConflictResolution) error
	Rebase(options ...Option) error
	Reflog(options ...Option) error
	SetConfig(key string, value string, options ...Option) error
	GetConfig(key string, options ...Option) (string, error)
	ListConfig(options ...Option) ([]types.ConfigEntry, error)
	UnsetConfig(key string, options ...Option) error
	Remove(options ...Option) error
	
	// Bare repository operations
	IsBareRepository() (bool, error)
}

// Command interface defines the contract for git command execution
type Command interface {
	Execute() ([]byte, error)
	ExecuteCombined() ([]byte, error)
	ExecuteWithStderr() ([]byte, error)
	ApplyOptions(opts ...Option)
	// Internal methods for option configuration
	SetTimeout(timeout time.Duration)
	SetEnv(key, value string)
	SetWorkingDir(dir string)
	SetStdin(input string)
	AddArgs(args ...string)
	// Internal access methods
	GetArgs() []string
	SetArgs(args []string)
}

// Option is a functional option for configuring git commands
type Option func(Command)

// gitImpl implements the Git interface
type gitImpl struct {
	path string
	wd   string
}

// NewGit creates a new git implementation
func NewGit() (*gitImpl, error) {
	path, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	return &gitImpl{
		path: path,
	}, nil
}

// NewGitInstance creates a new Git instance (basic, no session)
func NewGitInstance() (Git, error) {
	return NewGit()
}

// SetWorkingDirectory sets the working directory for git operations
func (g *gitImpl) SetWorkingDirectory(wd string) {
	g.wd = wd
}