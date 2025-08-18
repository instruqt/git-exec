package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommitCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("commit", "-m", "test message")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"commit", "-m", "test message"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestCommitWithAuthorOption(t *testing.T) {
	opt := CommitWithAuthor("John Doe", "john@example.com")
	
	cmd := &Command{env: make(map[string]string)}
	opt(cmd)
	
	require.Equal(t, "John Doe", cmd.env["GIT_AUTHOR_NAME"])
	require.Equal(t, "john@example.com", cmd.env["GIT_AUTHOR_EMAIL"])
	require.Equal(t, "John Doe", cmd.env["GIT_COMMITTER_NAME"])
	require.Equal(t, "john@example.com", cmd.env["GIT_COMMITTER_EMAIL"])
}

func TestCommitWithAllOption(t *testing.T) {
	opt := CommitWithAll()
	
	cmd := &Command{args: []string{"commit"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--all")
}

func TestCommitWithAmendOption(t *testing.T) {
	opt := CommitWithAmend()
	
	cmd := &Command{args: []string{"commit"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--amend")
}

func TestCommitWithSignoffOption(t *testing.T) {
	opt := CommitWithSignoff()
	
	cmd := &Command{args: []string{"commit"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--signoff")
}

func TestCommitWithGPGSignOption(t *testing.T) {
	opt := CommitWithGPGSign("keyid123")
	
	cmd := &Command{args: []string{"commit"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--gpg-sign=keyid123")
}

func TestCommitWithNoVerifyOption(t *testing.T) {
	opt := CommitWithNoVerify()
	
	cmd := &Command{args: []string{"commit"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--no-verify")
}