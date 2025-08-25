package git

import (
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
	Checkout(options ...Option) error
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
	Config(key string, value string, options ...Option) error
	Remove(options ...Option) error
}

