package git_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/instruqt/git-exec/pkg/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test tag operations - CRUD operations and edge cases
func TestTagOperations(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Create a commit to tag
	testFile := filepath.Join(tempDir, "tag-test.txt")
	err = os.WriteFile(testFile, []byte("tag test content"), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"tag-test.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Commit for tagging")
	require.NoError(t, err)
	
	// Test creating tags
	err = gitInstance.Tag("v1.0.0")
	require.NoError(t, err)
	
	err = gitInstance.Tag("v1.1.0")
	require.NoError(t, err)
	
	err = gitInstance.Tag("beta-1")
	require.NoError(t, err)
	
	// Test listing tags
	tags, err := gitInstance.ListTags()
	require.NoError(t, err)
	assert.Len(t, tags, 3)
	
	// Sort tags for consistent testing
	sort.Strings(tags)
	expected := []string{"beta-1", "v1.0.0", "v1.1.0"}
	sort.Strings(expected)
	assert.Equal(t, expected, tags)
	
	// Test deleting a tag
	err = gitInstance.DeleteTag("beta-1")
	require.NoError(t, err)
	
	// Verify tag was deleted
	tags, err = gitInstance.ListTags()
	require.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.NotContains(t, tags, "beta-1")
	
	// Test deleting non-existent tag
	err = gitInstance.DeleteTag("nonexistent")
	assert.Error(t, err, "Should fail to delete non-existent tag")
}

// Test tag edge cases and empty repository
func TestTagEdgeCases(t *testing.T) {
	tempDir := t.TempDir()
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	
	// Initialize empty repository
	err = gitInstance.Init(tempDir)
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Configure user
	err = gitInstance.Config("user.name", "Test User")
	require.NoError(t, err)
	err = gitInstance.Config("user.email", "test@example.com")
	require.NoError(t, err)
	
	// Test listing tags in empty repository
	tags, err := gitInstance.ListTags()
	require.NoError(t, err)
	assert.Empty(t, tags)
	
	// Test creating tag without commits (should fail)
	err = gitInstance.Tag("empty-tag")
	assert.Error(t, err, "Should fail to create tag without commits")
}

// Test tag naming validation through Git's behavior
func TestTagNaming(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Create a commit first
	testFile := filepath.Join(tempDir, "naming-test.txt")
	err = os.WriteFile(testFile, []byte("naming test"), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"naming-test.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Commit for naming test")
	require.NoError(t, err)
	
	// Test valid tag names
	validTags := []string{
		"v1.0",
		"release-2023",
		"feature/awesome",
		"1.0.0-beta.1",
	}
	
	for _, tagName := range validTags {
		err = gitInstance.Tag(tagName)
		assert.NoError(t, err, "Should create tag with name: %s", tagName)
	}
	
	// Verify all tags were created
	tags, err := gitInstance.ListTags()
	require.NoError(t, err)
	assert.Len(t, tags, len(validTags))
	
	// Test duplicate tag (should fail)
	err = gitInstance.Tag("v1.0")
	assert.Error(t, err, "Should fail to create duplicate tag")
}

// Test remote tag operations - interface testing
func TestRemoteTagOperations(t *testing.T) {
	tempDir := setupTestRepo(t)
	gitInstance, err := git.NewGit()
	require.NoError(t, err)
	gitInstance.SetWorkingDirectory(tempDir)
	
	// Add a remote for testing
	err = gitInstance.AddRemote("origin", "https://github.com/test/empty-repo.git")
	require.NoError(t, err)
	
	// Create a tag
	testFile := filepath.Join(tempDir, "remote-tag-test.txt")
	err = os.WriteFile(testFile, []byte("remote tag test"), 0644)
	require.NoError(t, err)
	err = gitInstance.Add([]string{"remote-tag-test.txt"})
	require.NoError(t, err)
	err = gitInstance.Commit("Commit for remote tag test")
	require.NoError(t, err)
	err = gitInstance.Tag("v2.0.0")
	require.NoError(t, err)
	
	// Test push tags - will likely fail due to network/auth, but tests interface
	_, err = gitInstance.PushTags("origin")
	// Don't require success - network operations are environment dependent
	
	// Test delete remote tag - tests interface
	err = gitInstance.DeleteRemoteTag("origin", "v2.0.0")
	// Don't require success - network operations are environment dependent
	// The value is in testing the interface works
}