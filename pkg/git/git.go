package git

import (
	"github.com/instruqt/git-exec/pkg/git/commands"
	"github.com/instruqt/git-exec/pkg/git/types"
)

type Git interface {
	SetWorkingDirectory(wd string)

	Init(path string, options ...commands.Option) error
	AddRemote(name, url string, options ...commands.Option) error
	RemoveRemote(name string, options ...commands.Option) error
	SetRemoteURL(name, url string, options ...commands.Option) error
	ListRemotes(options ...commands.Option) ([]types.Remote, error)
	Clone(url, destination string, options ...commands.Option) error
	Status(options ...commands.Option) ([]types.File, error)
	Add(files []string, options ...commands.Option) error
	Reset(files []string, options ...commands.Option) error
	Commit(message string, options ...commands.Option) error
	Diff(options ...commands.Option) ([]types.Diff, error)
	Show(object string, options ...commands.Option) (*types.Log, error)
	Log(options ...commands.Option) ([]types.Log, error)
	Fetch(options ...commands.Option) ([]types.Remote, error)
	Pull(options ...commands.Option) (*types.MergeResult, error)
	Push(options ...commands.Option) ([]types.Remote, error)
	ListBranches(options ...commands.Option) ([]types.Branch, error)
	CreateBranch(branch string, options ...commands.Option) error
	DeleteBranch(branch string, options ...commands.Option) error
	SetUpstream(branch string, remote string, options ...commands.Option) error
	Checkout(options ...commands.Option) error
	Tag(name string, options ...commands.Option) error
	ListTags(options ...commands.Option) ([]string, error)
	DeleteTag(name string, options ...commands.Option) error
	PushTags(remote string, options ...commands.Option) ([]types.Remote, error)
	DeleteRemoteTag(remote, tagName string, options ...commands.Option) error
	Revert(options ...commands.Option) error
	Merge(options ...commands.Option) error
	Rebase(options ...commands.Option) error
	Reflog(options ...commands.Option) error
	Config(key string, value string, options ...commands.Option) error
	Remove(options ...commands.Option) error
}

// NewSession creates a new Git session
func NewSession() (Git, error) {
	return commands.NewGit()
}

