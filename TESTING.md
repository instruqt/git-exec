# Testing Strategy

## Comprehensive Command Coverage with Value-Driven Tests

We've created **focused tests for every Git command** that test real functionality, edge cases, and error scenarios.

### Test Files

#### `commands_test.go` - Individual Command Tests
- **`TestInitCommand`**: Normal vs bare repository initialization
- **`TestCloneCommand`**: Clone scenarios and error cases
- **`TestAddCommand`**: File staging with different patterns
- **`TestCommitCommand`**: Commit with options and custom authors
- **`TestStatusCommand`**: Status parsing for different file states
- **`TestResetCommand`**: Unstaging files and reset modes
- **`TestConfigCommand`**: Configuration setting and usage
- **`TestLogCommand`**: Log parsing and options
- **`TestShowCommand`**: Show specific commits
- **`TestCheckoutCommand`**: Branch switching and creation
- **`TestDiffCommand`**: Diff interface testing

#### `remote_test.go` - Remote Operations
- **`TestRemoteOperations`**: CRUD operations (add, list, change, remove)
- **`TestRemoteErrorHandling`**: Error cases for remote operations
- **`TestFetchCommand`**: Fetch interface testing
- **`TestPushCommand`**: Push interface testing  
- **`TestPullCommand`**: Pull interface testing

#### `tag_test.go` - Tag Operations
- **`TestTagOperations`**: Tag CRUD lifecycle
- **`TestTagEdgeCases`**: Empty repo and error scenarios
- **`TestTagNaming`**: Valid tag name patterns
- **`TestRemoteTagOperations`**: Remote tag push/delete interfaces

#### `advanced_test.go` - Advanced Operations
- **`TestRevertCommand`**: Commit reverting
- **`TestRebaseCommand`**: Rebase interface testing
- **`TestReflogCommand`**: Reflog interface
- **`TestRemoveCommand`**: File removal from Git
- **`TestAdvancedCommandErrors`**: Error handling
- **`TestAdvancedBranchOperations`**: Complex branch scenarios

#### `git_integration_test.go` - Workflow Integration
- **`TestGitWorkflow`**: Complete init → config → add → commit → status workflow
- **`TestGitErrorHandling`**: Error cases that actually matter
- **`TestBranchOperations`**: Full branch lifecycle
- **`TestMergeConflictResolution`**: Complex merge scenarios

#### `session_test.go` - Session Management
- **`TestSessionPersistence`**: Session creation with full configuration
- **`TestSessionReload`**: Configuration persistence across sessions
- **`TestSessionUserUpdate`**: Dynamic user configuration changes
- **`TestSessionValidation`**: Error handling and validation workflows
- **`TestSessionDestroy`**: Session cleanup functionality

#### `merge_test.go` - Merge Operations
- **`TestMergeConflictWorkflow`**: Real conflict creation and resolution
- **`TestFastForwardMerge`**: Clean merge scenarios
- **`TestNoFastForwardMerge`**: Explicit merge commit creation
- **`TestMergeAbortAndContinue`**: Merge state management

#### `git_test.go` - Mock Demonstrations
- **`TestGitMockUsage`**: Generated Git interface mocks
- **`TestSessionMockUsage`**: Generated Session interface mocks

## What We Test (and Why)

### ✅ **Integration Workflows**
- Complete Git workflows that users actually perform
- Error states that matter in real usage
- Configuration persistence and session management
- Complex merge scenarios with conflicts

### ✅ **Edge Cases That Matter**
- Missing files/repositories
- Configuration changes
- Session persistence across restarts
- Merge conflicts and resolutions

### ❌ **What We Don't Test**
- Basic method existence (the compiler catches this)
- Trivial happy-path scenarios without complexity
- Implementation details that don't affect behavior
- Redundant coverage of the same workflows

## Test Principles

1. **Value-Driven**: Every test validates important functionality
2. **Integration-Focused**: Tests real workflows, not isolated units
3. **Error-Aware**: Tests error conditions that users will encounter
4. **Maintainable**: Tests are clear, focused, and don't duplicate coverage

## Running Tests

```bash
# Run all tests
go test ./pkg/git -v

# Run specific test categories
go test ./pkg/git -v -run "TestGitWorkflow"
go test ./pkg/git -v -run "TestSession"
go test ./pkg/git -v -run "TestMerge"
```

## Test Coverage

- **Core Git Operations**: ✅ Covered with integration workflows
- **Session Management**: ✅ Covered with persistence and configuration
- **Merge Operations**: ✅ Covered with conflict scenarios
- **Error Handling**: ✅ Covered for meaningful error cases
- **Mock Generation**: ✅ Covered with usage examples