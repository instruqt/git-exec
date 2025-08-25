package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	gitpkg "github.com/instruqt/git-exec/pkg/git"
	"github.com/instruqt/git-exec/pkg/git/errors"
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Merge joins two or more development histories together
func (g *git) Merge(opts ...gitpkg.Option) (*types.MergeResult, error) {
	cmd := g.newCommand("merge")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	// Extract branch name from command environment
	result := &types.MergeResult{}
	if gitCmd, ok := cmd.(*command); ok {
		if branchName, exists := gitCmd.env["MERGE_BRANCH"]; exists {
			result.MergedBranch = branchName
		}
	}
	
	output, err := cmd.Execute()
	
	// Parse merge output regardless of success/failure
	g.parseMergeOutput(string(output), result)
	
	if err != nil {
		// Check if it's a merge conflict (exit code 1)
		if gitErr, ok := err.(*errors.GitError); ok && gitErr.ExitCode == 1 {
			result.Success = false
			// Parse conflicts
			conflicts, conflictErr := g.parseConflicts()
			if conflictErr != nil {
				return result, fmt.Errorf("merge conflict occurred, but failed to parse conflicts: %w", conflictErr)
			}
			result.Conflicts = conflicts
			result.ConflictedFiles = make([]string, len(conflicts))
			for i, conflict := range conflicts {
				result.ConflictedFiles[i] = conflict.Path
			}
			return result, nil // Return result with conflicts, not an error
		}
		// Other types of errors
		result.Success = false
		result.AbortReason = err.Error()
		return result, err
	}
	
	result.Success = true
	return result, nil
}

// parseMergeOutput parses the git merge command output
func (g *git) parseMergeOutput(output string, result *types.MergeResult) {
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Fast-forward merge - extract branch name
		if strings.Contains(line, "Fast-forward") {
			result.FastForward = true
		}
		
		// Updating lines like "Updating abc123..def456" can give us branch info
		if strings.HasPrefix(line, "Updating ") {
			result.FastForward = true
		}
		
		// Look for merge commit message that includes branch name
		if strings.HasPrefix(line, "Merge branch '") && strings.Contains(line, "'") {
			// Extract branch name from merge commit message
			start := strings.Index(line, "'")
			if start != -1 {
				end := strings.Index(line[start+1:], "'")
				if end != -1 {
					result.MergedBranch = line[start+1 : start+1+end]
				}
			}
		}
		
		// Merge commit
		if strings.HasPrefix(line, "Merge made by") {
			// Extract strategy
			if strings.Contains(line, "recursive") {
				result.Strategy = "recursive"
			} else if strings.Contains(line, "ort") {
				result.Strategy = "ort"
			} else if strings.Contains(line, "octopus") {
				result.Strategy = "octopus"
			}
		}
		
		// Stats parsing
		if strings.Contains(line, "file") && (strings.Contains(line, "changed") || strings.Contains(line, "insertion") || strings.Contains(line, "deletion")) {
			g.parseStats(line, result)
		}
	}
}

// parseStats parses merge statistics from output
func (g *git) parseStats(line string, result *types.MergeResult) {
	// Parse lines like: "1 file changed, 2 insertions(+), 3 deletions(-)"
	re := regexp.MustCompile(`(\d+)\s+files?\s+changed(?:,\s+(\d+)\s+insertions?\(\+\))?(?:,\s+(\d+)\s+deletions?\(-\))?`)
	matches := re.FindStringSubmatch(line)
	
	if len(matches) > 1 {
		if files, err := strconv.Atoi(matches[1]); err == nil {
			result.Stats.FilesChanged = files
		}
		if len(matches) > 2 && matches[2] != "" {
			if insertions, err := strconv.Atoi(matches[2]); err == nil {
				result.Stats.Insertions = insertions
			}
		}
		if len(matches) > 3 && matches[3] != "" {
			if deletions, err := strconv.Atoi(matches[3]); err == nil {
				result.Stats.Deletions = deletions
			}
		}
	}
}

// parseConflicts parses merge conflicts from the repository
func (g *git) parseConflicts() ([]types.ConflictFile, error) {
	// Get list of conflicted files
	cmd := g.newCommand("diff", "--name-only", "--diff-filter=U")
	output, err := cmd.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get conflicted files: %w", err)
	}
	
	conflictedPaths := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(conflictedPaths) == 1 && conflictedPaths[0] == "" {
		return nil, nil // No conflicts
	}
	
	var conflicts []types.ConflictFile
	for _, path := range conflictedPaths {
		if path == "" {
			continue
		}
		
		conflict, err := g.parseConflictFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse conflict in %s: %w", path, err)
		}
		conflicts = append(conflicts, conflict)
	}
	
	return conflicts, nil
}

// parseConflictFile parses conflict markers in a specific file
func (g *git) parseConflictFile(filePath string) (types.ConflictFile, error) {
	fullPath := filepath.Join(g.wd, filePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return types.ConflictFile{}, fmt.Errorf("failed to read conflicted file: %w", err)
	}
	
	conflict := types.ConflictFile{
		Path:    filePath,
		Content: string(content),
		Status:  types.ConflictStatusBothModified, // Default, could be refined
	}
	
	// Parse conflict sections
	sections, err := g.parseConflictSections(string(content))
	if err != nil {
		return conflict, fmt.Errorf("failed to parse conflict sections: %w", err)
	}
	
	conflict.Sections = sections
	return conflict, nil
}

// parseConflictSections parses individual conflict sections within file content
func (g *git) parseConflictSections(content string) ([]types.ConflictSection, error) {
	lines := strings.Split(content, "\n")
	var sections []types.ConflictSection
	
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		
		// Look for conflict start marker: <<<<<<<
		if strings.HasPrefix(line, "<<<<<<<") {
			section := types.ConflictSection{
				StartLine: i + 1,
			}
			
			// Collect our content until separator ===
			var ourLines []string
			i++ // Move past start marker
			for i < len(lines) && !strings.HasPrefix(lines[i], "=======") {
				ourLines = append(ourLines, lines[i])
				i++
			}
			section.OurContent = strings.Join(ourLines, "\n")
			
			// Move past separator
			if i < len(lines) && strings.HasPrefix(lines[i], "=======") {
				i++ // Move past separator
			}
			
			// Collect their content until end marker >>>>>>>
			var theirLines []string
			for i < len(lines) && !strings.HasPrefix(lines[i], ">>>>>>>") {
				theirLines = append(theirLines, lines[i])
				i++
			}
			section.TheirContent = strings.Join(theirLines, "\n")
			
			if i < len(lines) {
				section.EndLine = i + 1
			}
			
			sections = append(sections, section)
		}
	}
	
	return sections, nil
}

// ResolveConflicts applies conflict resolutions to files
func (g *git) ResolveConflicts(resolutions []types.ConflictResolution) error {
	for _, resolution := range resolutions {
		if err := g.resolveConflictFile(resolution); err != nil {
			return fmt.Errorf("failed to resolve conflicts in %s: %w", resolution.FilePath, err)
		}
	}
	return nil
}

// resolveConflictFile resolves conflicts in a single file
func (g *git) resolveConflictFile(resolution types.ConflictResolution) error {
	filePath := filepath.Join(g.wd, resolution.FilePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	
	// Simple resolution strategies
	if resolution.UseOurs {
		resolved, err := g.resolveUseOurs(string(content))
		if err != nil {
			return err
		}
		return os.WriteFile(filePath, []byte(resolved), 0644)
	}
	
	if resolution.UseTheirs {
		resolved, err := g.resolveUseTheirs(string(content))
		if err != nil {
			return err
		}
		return os.WriteFile(filePath, []byte(resolved), 0644)
	}
	
	if resolution.Custom && len(resolution.Sections) > 0 {
		resolved, err := g.resolveCustom(string(content), resolution.Sections)
		if err != nil {
			return err
		}
		return os.WriteFile(filePath, []byte(resolved), 0644)
	}
	
	return fmt.Errorf("no resolution method specified")
}

// resolveUseOurs resolves conflicts by keeping "our" version
func (g *git) resolveUseOurs(content string) (string, error) {
	lines := strings.Split(content, "\n")
	var result []string
	
	inConflict := false
	inTheirSection := false
	
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "<<<<<<<"):
			inConflict = true
			inTheirSection = false
			// Skip conflict marker
		case strings.HasPrefix(line, "======="):
			inTheirSection = true
			// Skip separator
		case strings.HasPrefix(line, ">>>>>>>"):
			inConflict = false
			inTheirSection = false
			// Skip end marker
		default:
			if !inConflict {
				// Outside conflict - include line
				result = append(result, line)
			} else if !inTheirSection {
				// In our section - include line
				result = append(result, line)
			}
			// If inTheirSection is true, skip the line
		}
	}
	
	return strings.Join(result, "\n"), nil
}

// resolveUseTheirs resolves conflicts by keeping "their" version
func (g *git) resolveUseTheirs(content string) (string, error) {
	lines := strings.Split(content, "\n")
	var result []string
	
	inConflict := false
	inTheirSection := false
	
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "<<<<<<<"):
			inConflict = true
			inTheirSection = false
		case strings.HasPrefix(line, "======="):
			inTheirSection = true
		case strings.HasPrefix(line, ">>>>>>>"):
			inConflict = false
			inTheirSection = false
		default:
			if !inConflict {
				// Outside conflict
				result = append(result, line)
			} else if inTheirSection {
				// In "their" section
				result = append(result, line)
			}
			// Skip "our" section
		}
	}
	
	return strings.Join(result, "\n"), nil
}

// resolveCustom resolves conflicts using custom resolutions
func (g *git) resolveCustom(content string, resolutions []types.ResolvedSection) (string, error) {
	lines := strings.Split(content, "\n")
	var result []string
	
	sectionIndex := 0
	inConflict := false
	
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "<<<<<<<"):
			inConflict = true
			// Skip conflict start marker
		case strings.HasPrefix(line, ">>>>>>>"):
			// Apply custom resolution for this section
			if sectionIndex < len(resolutions) && resolutions[sectionIndex].SectionIndex == sectionIndex {
				result = append(result, resolutions[sectionIndex].Resolution)
			}
			sectionIndex++
			inConflict = false
		case inConflict:
			// Skip all conflict content, we'll use custom resolution
		default:
			// Outside conflict
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n"), nil
}

// MergeAbort aborts an ongoing merge
func (g *git) MergeAbort() error {
	cmd := g.newCommand("merge", "--abort")
	_, err := cmd.Execute()
	return err
}

// MergeContinue continues a merge after resolving conflicts
func (g *git) MergeContinue() error {
	cmd := g.newCommand("merge", "--continue")
	_, err := cmd.Execute()
	return err
}

// Merge-specific options

// MergeWithBranch merges the specified branch
func MergeWithBranch(branch string) gitpkg.Option {
	return func(c gitpkg.Command) {
		c.AddArgs(branch)
		// Store the branch name for the result
		if gitCmd, ok := c.(*command); ok {
			if gitCmd.env == nil {
				gitCmd.env = make(map[string]string)
			}
			gitCmd.env["MERGE_BRANCH"] = branch
		}
	}
}

// MergeWithCommit merges commits into current branch
func MergeWithCommit(commit string) gitpkg.Option {
	return WithArgs(commit)
}

// MergeWithNoFF creates a merge commit even when fast-forward is possible
func MergeWithNoFF() gitpkg.Option {
	return WithArgs("--no-ff")
}

// MergeWithFFOnly only updates if the merge can be resolved as fast-forward
func MergeWithFFOnly() gitpkg.Option {
	return WithArgs("--ff-only")
}

// MergeWithSquash creates a single commit instead of merging
func MergeWithSquash() gitpkg.Option {
	return WithArgs("--squash")
}

// MergeWithStrategy specifies merge strategy
func MergeWithStrategy(strategy string) gitpkg.Option {
	return WithArgs("--strategy", strategy)
}

// MergeWithStrategyOption passes option to merge strategy
func MergeWithStrategyOption(option string) gitpkg.Option {
	return WithArgs("--strategy-option", option)
}

// MergeWithEdit invokes editor to edit merge commit message
func MergeWithEdit() gitpkg.Option {
	return WithArgs("--edit")
}

// MergeWithNoEdit accepts auto-generated message
func MergeWithNoEdit() gitpkg.Option {
	return WithArgs("--no-edit")
}

// MergeWithMessage specifies merge commit message
func MergeWithMessage(message string) gitpkg.Option {
	return WithArgs("-m", message)
}

// MergeWithFile reads merge message from file
func MergeWithFile(file string) gitpkg.Option {
	return WithArgs("-F", file)
}

// MergeWithAbort aborts current merge and restores pre-merge state
func MergeWithAbort() gitpkg.Option {
	return WithArgs("--abort")
}

// MergeWithContinue continues merge after resolving conflicts
func MergeWithContinue() gitpkg.Option {
	return WithArgs("--continue")
}

// MergeWithQuiet suppresses output
func MergeWithQuiet() gitpkg.Option {
	return WithArgs("--quiet")
}

// MergeWithVerbose shows verbose output
func MergeWithVerbose() gitpkg.Option {
	return WithArgs("--verbose")
}

// MergeWithProgress shows progress
func MergeWithProgress() gitpkg.Option {
	return WithArgs("--progress")
}

// MergeWithNoProgress hides progress
func MergeWithNoProgress() gitpkg.Option {
	return WithArgs("--no-progress")
}

// MergeWithSign makes a GPG-signed merge commit
func MergeWithSign() gitpkg.Option {
	return WithArgs("--gpg-sign")
}

// MergeWithNoSign doesn't GPG-sign merge commit
func MergeWithNoSign() gitpkg.Option {
	return WithArgs("--no-gpg-sign")
}

// MergeWithLog includes one-line descriptions from merged commits
func MergeWithLog(n string) gitpkg.Option {
	if n == "" {
		return WithArgs("--log")
	}
	return WithArgs("--log=" + n)
}

// MergeWithNoLog doesn't include one-line descriptions
func MergeWithNoLog() gitpkg.Option {
	return WithArgs("--no-log")
}

// MergeWithStat shows diffstat at end of merge
func MergeWithStat() gitpkg.Option {
	return WithArgs("--stat")
}

// MergeWithNoStat doesn't show diffstat
func MergeWithNoStat() gitpkg.Option {
	return WithArgs("--no-stat")
}