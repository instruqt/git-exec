package git

import (
	"regexp"
	"strings"
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Checkout checks out branches, commits, or files
func (g *gitImpl) Checkout(opts ...Option) (*types.CheckoutResult, error) {
	// Get current HEAD before checkout
	headCmd := g.newCommand("rev-parse", "HEAD")
	headOutput, _ := headCmd.Execute()
	previousHEAD := strings.TrimSpace(string(headOutput))

	// Get current branch before checkout (unused but might be useful for future enhancements)
	_ = "" // placeholder to avoid unused variable warning

	// Execute checkout command
	cmd := g.newCommand("checkout")
	cmd.ApplyOptions(opts...)
	output, err := cmd.ExecuteCombined()

	result := &types.CheckoutResult{
		PreviousHEAD: previousHEAD,
	}

	outputStr := string(output)

	if err != nil {
		result.Success = false
		return result, err
	}

	result.Success = true

	// Debug: print actual git output (disabled)
	// fmt.Printf("Git checkout output: %q\n", outputStr)

	// Parse checkout output
	result = parseCheckoutOutput(result, outputStr)

	// Get new HEAD after checkout
	newHeadCmd := g.newCommand("rev-parse", "HEAD")
	newHeadOutput, _ := newHeadCmd.Execute()
	result.NewHEAD = strings.TrimSpace(string(newHeadOutput))

	// Check if in detached HEAD state
	newBranchCmd := g.newCommand("rev-parse", "--abbrev-ref", "HEAD")
	newBranchOutput, branchErr := newBranchCmd.Execute()
	currentRef := strings.TrimSpace(string(newBranchOutput))

	if branchErr != nil || currentRef == "" {
		// This might be an orphan branch - try to get the branch name differently
		statusCmd := g.newCommand("status", "--porcelain=v1", "--branch")
		statusOutput, statusErr := statusCmd.Execute()
		if statusErr == nil {
			statusLines := strings.Split(string(statusOutput), "\n")
			for _, line := range statusLines {
				if strings.HasPrefix(line, "## ") {
					// Extract branch name from "## branch-name" or "## No commits yet on branch-name"
					branchInfo := strings.TrimPrefix(line, "## ")
					
					// Handle "No commits yet on branch-name" case
					if strings.HasPrefix(branchInfo, "No commits yet on ") {
						result.Branch = strings.TrimPrefix(branchInfo, "No commits yet on ")
						break
					}
					
					// Handle regular "branch-name..." case
					parts := strings.Split(branchInfo, "...")
					if len(parts) > 0 && parts[0] != "" {
						result.Branch = parts[0]
						break
					}
				}
			}
		}
	} else if currentRef == "HEAD" {
		result.DetachedHEAD = true
		result.Commit = result.NewHEAD
	} else {
		result.Branch = currentRef
	}

	return result, nil
}

// parseCheckoutOutput parses git checkout output and populates CheckoutResult
func parseCheckoutOutput(result *types.CheckoutResult, output string) *types.CheckoutResult {
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Branch switching patterns
		if strings.HasPrefix(line, "Switched to branch '") {
			result.Branch = extractQuotedString(line, "Switched to branch '", "'")
		} else if strings.HasPrefix(line, "Switched to a new branch '") {
			result.Branch = extractQuotedString(line, "Switched to a new branch '", "'")
			result.NewBranch = true
		} else if strings.Contains(line, "HEAD is now at") {
			result.DetachedHEAD = true
			parts := strings.Fields(line)
			if len(parts) >= 5 {
				result.Commit = parts[4] // Extract commit hash
			}
		}

		// Branch tracking
		if strings.Contains(line, "set up to track remote branch") {
			result.UpstreamBranch = extractUpstreamBranch(line)
		}

		// File modifications
		if strings.HasPrefix(line, "M\t") {
			result.ModifiedFiles = append(result.ModifiedFiles, strings.TrimPrefix(line, "M\t"))
		} else if strings.HasPrefix(line, "A\t") {
			result.RestoredFiles = append(result.RestoredFiles, strings.TrimPrefix(line, "A\t"))
		} else if strings.HasPrefix(line, "D\t") {
			result.ModifiedFiles = append(result.ModifiedFiles, strings.TrimPrefix(line, "D\t"))
		}

		// Warning messages
		if strings.Contains(line, "warning:") || strings.Contains(line, "Warning:") {
			if result.Warning == "" {
				result.Warning = line
			} else {
				result.Warning += "; " + line
			}
		}

		// Untracked files that would be overwritten
		if strings.Contains(line, "would be overwritten by checkout") {
			result.UntrackedFiles = append(result.UntrackedFiles, extractUntrackedFile(line))
		}
	}

	return result
}

// extractQuotedString extracts a string between specified quotes
func extractQuotedString(text, prefix, suffix string) string {
	start := strings.Index(text, prefix)
	if start == -1 {
		return ""
	}
	start += len(prefix)
	
	end := strings.Index(text[start:], suffix)
	if end == -1 {
		return ""
	}
	
	return text[start : start+end]
}

// extractUpstreamBranch extracts upstream branch from tracking message
func extractUpstreamBranch(line string) string {
	re := regexp.MustCompile(`'([^']+/[^']+)'`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractUntrackedFile extracts filename from overwrite warning
func extractUntrackedFile(line string) string {
	re := regexp.MustCompile(`error: The following untracked working tree files would be overwritten by checkout:\s+(.+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}