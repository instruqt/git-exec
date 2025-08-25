package git_test

import (
	"testing"

	"github.com/instruqt/git-exec/pkg/git"
	"github.com/instruqt/git-exec/pkg/git/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Example test demonstrating how easy it is to use the generated mocks
func TestGitMockUsage(t *testing.T) {
	// Create a mock Git instance
	mockGit := mocks.NewMockGit(t)

	// Set up expectations using testify/mock syntax
	mockGit.On("Add", []string{"test.txt"}).Return(nil).Once()
	mockGit.On("Commit", "test commit", mock.AnythingOfType("git.Option")).Return(nil).Once()

	// Use the mock in your code
	err := mockGit.Add([]string{"test.txt"})
	assert.NoError(t, err)

	err = mockGit.Commit("test commit", git.WithQuiet())
	assert.NoError(t, err)

	// Verify all expectations were met
	mockGit.AssertExpectations(t)
}

// Example test showing Session mock usage
func TestSessionMockUsage(t *testing.T) {
	// Create a mock Session instance  
	mockSession := mocks.NewMockSession(t)

	// Set up expectations
	mockConfig := &git.SessionConfig{
		UserName: "Test User",
		Metadata: map[string]string{"session-id": "test-session-123"},
	}
	mockSession.On("GetSessionConfig").Return(mockConfig).Once()
	mockSession.On("Add", []string{"README.md"}).Return(nil).Once()

	// Use the mock
	config := mockSession.GetSessionConfig()
	assert.Equal(t, "Test User", config.UserName)
	assert.Equal(t, "test-session-123", config.Metadata["session-id"])

	err := mockSession.Add([]string{"README.md"})
	assert.NoError(t, err)

	// Verify all expectations were met
	mockSession.AssertExpectations(t)
}