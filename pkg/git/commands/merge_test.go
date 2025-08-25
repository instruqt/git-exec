package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/instruqt/git-exec/pkg/git/types"
	"github.com/stretchr/testify/require"
)

func TestMergeSuccessful(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-merge-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	repoPath := filepath.Join(tempDir, "test-repo")
	
	t.Run("successful merge without conflicts", func(t *testing.T) {
		// Initialize repository
		g, err := NewGit()
		require.NoError(t, err)
		
		err = g.Init(repoPath)
		require.NoError(t, err)
		
		g.SetWorkingDirectory(repoPath)
		
		// Create initial commit on main
		mainFile := filepath.Join(repoPath, "main.txt")
		err = os.WriteFile(mainFile, []byte("main content\n"), 0644)
		require.NoError(t, err)
		
		err = g.Add([]string{"main.txt"})
		require.NoError(t, err)
		
		err = g.Commit("Initial commit",
			WithConfig("user.name", "Test User"),
			WithConfig("user.email", "test@example.com"))
		require.NoError(t, err)
		
		// Create and checkout feature branch
		err = g.CreateBranch("feature-branch")
		require.NoError(t, err)
		
		err = g.Checkout(CheckoutWithBranch("feature-branch"))
		require.NoError(t, err)
		
		// Create feature file
		featureFile := filepath.Join(repoPath, "feature.txt")
		err = os.WriteFile(featureFile, []byte("feature content\n"), 0644)
		require.NoError(t, err)
		
		err = g.Add([]string{"feature.txt"})
		require.NoError(t, err)
		
		err = g.Commit("Add feature",
			WithConfig("user.name", "Test User"),
			WithConfig("user.email", "test@example.com"))
		require.NoError(t, err)
		
		// Switch back to main
		err = g.Checkout(CheckoutWithBranch("main"))
		require.NoError(t, err)
		
		// Merge feature branch
		result, err := g.Merge(MergeWithBranch("feature-branch"))
		require.NoError(t, err)
		require.NotNil(t, result)
		require.True(t, result.Success)
		require.Equal(t, "feature-branch", result.MergedBranch)
		require.Empty(t, result.ConflictedFiles)
	})
}

func TestMergeWithConflicts(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-merge-conflict-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	repoPath := filepath.Join(tempDir, "conflict-repo")
	
	t.Run("merge with conflicts", func(t *testing.T) {
		// Initialize repository
		g, err := NewGit()
		require.NoError(t, err)
		
		err = g.Init(repoPath)
		require.NoError(t, err)
		
		g.SetWorkingDirectory(repoPath)
		
		// Create initial file
		conflictFile := filepath.Join(repoPath, "conflict.txt")
		initialContent := `line 1
line 2
line 3
`
		err = os.WriteFile(conflictFile, []byte(initialContent), 0644)
		require.NoError(t, err)
		
		err = g.Add([]string{"conflict.txt"})
		require.NoError(t, err)
		
		err = g.Commit("Initial commit",
			WithConfig("user.name", "Test User"),
			WithConfig("user.email", "test@example.com"))
		require.NoError(t, err)
		
		// Create branch and modify file
		err = g.CreateBranch("feature")
		require.NoError(t, err)
		
		err = g.Checkout(CheckoutWithBranch("feature"))
		require.NoError(t, err)
		
		// Modify file in feature branch
		featureContent := `line 1
feature line 2
line 3
`
		err = os.WriteFile(conflictFile, []byte(featureContent), 0644)
		require.NoError(t, err)
		
		err = g.Add([]string{"conflict.txt"})
		require.NoError(t, err)
		
		err = g.Commit("Feature changes",
			WithConfig("user.name", "Test User"),
			WithConfig("user.email", "test@example.com"))
		require.NoError(t, err)
		
		// Switch to main and modify the same line
		err = g.Checkout(CheckoutWithBranch("main"))
		require.NoError(t, err)
		
		mainContent := `line 1
main line 2
line 3
`
		err = os.WriteFile(conflictFile, []byte(mainContent), 0644)
		require.NoError(t, err)
		
		err = g.Add([]string{"conflict.txt"})
		require.NoError(t, err)
		
		err = g.Commit("Main changes",
			WithConfig("user.name", "Test User"),
			WithConfig("user.email", "test@example.com"))
		require.NoError(t, err)
		
		// Attempt to merge - should create conflicts
		result, err := g.Merge(MergeWithBranch("feature"))
		require.NoError(t, err) // No error, but conflicts detected
		require.NotNil(t, result)
		require.False(t, result.Success)
		require.NotEmpty(t, result.ConflictedFiles)
		require.Contains(t, result.ConflictedFiles, "conflict.txt")
		require.NotEmpty(t, result.Conflicts)
		
		// Verify conflict details
		conflict := result.Conflicts[0]
		require.Equal(t, "conflict.txt", conflict.Path)
		require.NotEmpty(t, conflict.Sections)
		require.Contains(t, conflict.Sections[0].OurContent, "main line 2")
		require.Contains(t, conflict.Sections[0].TheirContent, "feature line 2")
	})
}

func TestConflictResolution(t *testing.T) {
	// This test creates a conflict and then resolves it
	tempDir, err := os.MkdirTemp("", "git-resolve-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	repoPath := filepath.Join(tempDir, "resolve-repo")
	
	t.Run("resolve conflicts use ours", func(t *testing.T) {
		// Set up repository with conflicts (similar to previous test)
		g := setupConflictedRepo(t, repoPath)
		
		// Create conflict content manually for testing resolution
		conflictContent := `line 1
<<<<<<< HEAD
main line 2
=======
feature line 2
>>>>>>> feature
line 3
`
		conflictFile := filepath.Join(repoPath, "conflict.txt")
		err := os.WriteFile(conflictFile, []byte(conflictContent), 0644)
		require.NoError(t, err)
		
		// Resolve conflicts using "ours" strategy
		resolutions := []types.ConflictResolution{
			{
				FilePath: "conflict.txt",
				UseOurs:  true,
			},
		}
		
		err = g.ResolveConflicts(resolutions)
		require.NoError(t, err)
		
		// Verify resolution
		resolved, err := os.ReadFile(conflictFile)
		require.NoError(t, err)
		resolvedContent := string(resolved)
		require.Contains(t, resolvedContent, "main line 2")
		require.NotContains(t, resolvedContent, "feature line 2")
		require.NotContains(t, resolvedContent, "<<<<<<")
		require.NotContains(t, resolvedContent, ">>>>>>")
		require.NotContains(t, resolvedContent, "======")
	})
	
	t.Run("resolve conflicts use theirs", func(t *testing.T) {
		// Set up repository with conflicts
		g := setupConflictedRepo(t, repoPath+"-theirs")
		
		conflictContent := `line 1
<<<<<<< HEAD
main line 2
=======
feature line 2
>>>>>>> feature
line 3
`
		conflictFile := filepath.Join(repoPath+"-theirs", "conflict.txt")
		err := os.WriteFile(conflictFile, []byte(conflictContent), 0644)
		require.NoError(t, err)
		
		// Resolve conflicts using "theirs" strategy
		resolutions := []types.ConflictResolution{
			{
				FilePath:  "conflict.txt",
				UseTheirs: true,
			},
		}
		
		err = g.ResolveConflicts(resolutions)
		require.NoError(t, err)
		
		// Verify resolution
		resolved, err := os.ReadFile(conflictFile)
		require.NoError(t, err)
		resolvedContent := string(resolved)
		require.Contains(t, resolvedContent, "feature line 2")
		require.NotContains(t, resolvedContent, "main line 2")
		require.NotContains(t, resolvedContent, "<<<<<<")
		require.NotContains(t, resolvedContent, ">>>>>>")
		require.NotContains(t, resolvedContent, "======")
	})
}

func TestMergeAbortAndContinue(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-merge-abort-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	repoPath := filepath.Join(tempDir, "abort-repo")
	
	t.Run("merge abort", func(t *testing.T) {
		g := setupConflictedRepo(t, repoPath)
		
		// Start a merge that will conflict
		_, err := g.Merge(MergeWithBranch("feature"))
		require.NoError(t, err) // Conflicts are not errors, just unsuccessful merges
		
		// Abort the merge
		err = g.MergeAbort()
		require.NoError(t, err)
		
		// Verify we're back to clean state
		files, err := g.Status()
		require.NoError(t, err)
		// Should have no conflicted files
		for _, file := range files {
			require.NotEqual(t, "UU", file.Status) // UU indicates merge conflict
		}
	})
}

func TestConflictParsing(t *testing.T) {
	g := &git{}
	
	t.Run("parse conflict sections", func(t *testing.T) {
		conflictContent := `line before conflict
<<<<<<< HEAD
our content line 1
our content line 2
=======
their content line 1
their content line 2
>>>>>>> branch-name
line after conflict
<<<<<<< HEAD
another our section
=======
another their section
>>>>>>> branch-name
final line`
		
		sections, err := g.parseConflictSections(conflictContent)
		require.NoError(t, err)
		require.Len(t, sections, 2)
		
		// First section
		require.Equal(t, 2, sections[0].StartLine)
		require.Equal(t, 8, sections[0].EndLine)
		require.Contains(t, sections[0].OurContent, "our content line 1")
		require.Contains(t, sections[0].OurContent, "our content line 2")
		require.Contains(t, sections[0].TheirContent, "their content line 1")
		require.Contains(t, sections[0].TheirContent, "their content line 2")
		
		// Second section
		require.Equal(t, 10, sections[1].StartLine)
		require.Equal(t, 14, sections[1].EndLine)
		require.Contains(t, sections[1].OurContent, "another our section")
		require.Contains(t, sections[1].TheirContent, "another their section")
	})
}

// setupConflictedRepo creates a repository with a conflicted merge scenario
func setupConflictedRepo(t *testing.T, repoPath string) *git {
	g, err := NewGit()
	require.NoError(t, err)
	
	err = g.Init(repoPath)
	require.NoError(t, err)
	g.SetWorkingDirectory(repoPath)
	
	// Create initial commit
	conflictFile := filepath.Join(repoPath, "conflict.txt")
	err = os.WriteFile(conflictFile, []byte("line 1\nline 2\nline 3\n"), 0644)
	require.NoError(t, err)
	
	err = g.Add([]string{"conflict.txt"})
	require.NoError(t, err)
	
	err = g.Commit("Initial commit",
		WithConfig("user.name", "Test User"),
		WithConfig("user.email", "test@example.com"))
	require.NoError(t, err)
	
	// Create feature branch
	err = g.CreateBranch("feature")
	require.NoError(t, err)
	
	err = g.Checkout(CheckoutWithBranch("feature"))
	require.NoError(t, err)
	
	// Modify in feature
	err = os.WriteFile(conflictFile, []byte("line 1\nfeature line 2\nline 3\n"), 0644)
	require.NoError(t, err)
	
	err = g.Add([]string{"conflict.txt"})
	require.NoError(t, err)
	
	err = g.Commit("Feature changes",
		WithConfig("user.name", "Test User"),
		WithConfig("user.email", "test@example.com"))
	require.NoError(t, err)
	
	// Switch to main and modify
	err = g.Checkout(CheckoutWithBranch("main"))
	require.NoError(t, err)
	
	err = os.WriteFile(conflictFile, []byte("line 1\nmain line 2\nline 3\n"), 0644)
	require.NoError(t, err)
	
	err = g.Add([]string{"conflict.txt"})
	require.NoError(t, err)
	
	err = g.Commit("Main changes",
		WithConfig("user.name", "Test User"),
		WithConfig("user.email", "test@example.com"))
	require.NoError(t, err)
	
	return g
}