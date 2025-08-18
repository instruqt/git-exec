package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogCommand(t *testing.T) {
	git := &git{path: "git", wd: "/test"}

	cmd := git.newCommand("log", "--format=fuller")
	require.Equal(t, "git", cmd.gitPath)
	require.Equal(t, []string{"log", "--format=fuller"}, cmd.args)
	require.Equal(t, "/test", cmd.workingDir)
}

func TestLogWithOnelineOption(t *testing.T) {
	opt := LogWithOneline()
	
	cmd := &Command{args: []string{"log"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--oneline")
}

func TestLogWithGraphOption(t *testing.T) {
	opt := LogWithGraph()
	
	cmd := &Command{args: []string{"log"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--graph")
}

func TestLogWithDecorateOption(t *testing.T) {
	opt := LogWithDecorate()
	
	cmd := &Command{args: []string{"log"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--decorate")
}

func TestLogWithAllOption(t *testing.T) {
	opt := LogWithAll()
	
	cmd := &Command{args: []string{"log"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--all")
}

func TestLogWithMaxCountOption(t *testing.T) {
	opt := LogWithMaxCount("10")
	
	cmd := &Command{args: []string{"log"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--max-count=10")
}

func TestLogWithSkipOption(t *testing.T) {
	opt := LogWithSkip("5")
	
	cmd := &Command{args: []string{"log"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--skip=5")
}

func TestLogWithSinceOption(t *testing.T) {
	opt := LogWithSince("2023-01-01")
	
	cmd := &Command{args: []string{"log"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--since=2023-01-01")
}

func TestLogWithAuthorOption(t *testing.T) {
	opt := LogWithAuthor("john.doe")
	
	cmd := &Command{args: []string{"log"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--author=john.doe")
}

func TestLogWithFormatOption(t *testing.T) {
	opt := LogWithFormat("oneline")
	
	cmd := &Command{args: []string{"log"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--format=oneline")
}

func TestLogWithNoMergesOption(t *testing.T) {
	opt := LogWithNoMerges()
	
	cmd := &Command{args: []string{"log"}}
	opt(cmd)
	
	require.Contains(t, cmd.args, "--no-merges")
}