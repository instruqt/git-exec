# Mock Generation with Mockery

With the flattened structure, mockery can now easily generate mocks for all interfaces without any circular dependency issues.

## File-Based Mock Naming

Mocks are named after their corresponding interface files:
- `git.go` → Contains `MockGit` for the Git interface
- `command.go` → Contains `MockCommand` for the Command interface  
- `session.go` → Contains `MockSession` for the Session interface

## Generate Individual Mocks

```bash
# Generate Git interface mock → pkg/git/mocks/git.go
mockery --dir=pkg/git --name=Git --output=pkg/git/mocks --filename=git.go

# Generate Session interface mock → pkg/git/mocks/session.go
mockery --dir=pkg/git --name=Session --output=pkg/git/mocks --filename=session.go

# Generate Command interface mock → pkg/git/mocks/command.go
mockery --dir=pkg/git --name=Command --output=pkg/git/mocks --filename=command.go
```

## Generate All Mocks with Config

Use the provided `.mockery.yaml` config:

```bash
mockery
```

This generates:
- `pkg/git/mocks/git.go`
- `pkg/git/mocks/command.go`
- `pkg/git/mocks/session.go`

## Usage in Tests

```go
package mypackage_test

import (
    "testing"
    "github.com/instruqt/git-exec/pkg/git"
    "github.com/instruqt/git-exec/pkg/git/mocks"
    "github.com/stretchr/testify/assert"
)

func TestWithMocks(t *testing.T) {
    // Create mock (note the Mock prefix)
    mockGit := mocks.NewMockGit(t)
    
    // Set up expectations
    mockGit.On("Add", []string{"file.txt"}).Return(nil).Once()
    
    // Use mock
    err := mockGit.Add([]string{"file.txt"})
    assert.NoError(t, err)
    
    // Verify expectations (automatic cleanup)
    mockGit.AssertExpectations(t)
}
```

## Benefits of Flattened Structure + File-Based Naming

1. **No Circular Dependencies**: Clean imports allow mockery to work flawlessly
2. **Organized Mock Files**: Each interface gets its own mock file matching the source
3. **Simple Commands**: Just `mockery --dir=pkg/git --name=InterfaceName`
4. **All Interfaces Mockable**: Git, Session, Command - all can be mocked
5. **Clean Integration**: Works perfectly with testify/mock patterns
6. **File Consistency**: Mock files mirror the structure of interface files
7. **No Factory Patterns Needed**: Direct interface usage for clean testing