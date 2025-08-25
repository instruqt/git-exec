package git

import (
	"strings"
	"time"
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Log shows the commit logs
func (g *gitImpl) Log(opts ...Option) ([]types.Log, error) {
	// Use a custom format for easier parsing
	format := "--pretty=format:COMMIT:%H%nTREE:%T%nPARENT:%P%nAUTHOR:%an <%ae>%nAUTHOR_DATE:%ai%nCOMMITTER:%cn <%ce>%nCOMMITTER_DATE:%ci%nMESSAGE:%s%n---END---"
	cmd := g.newCommand("log", format)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}
	
	return parseLogOutput(string(output))
}

// parseLogOutput parses git log output into Log structs
func parseLogOutput(output string) ([]types.Log, error) {
	if output == "" {
		return []types.Log{}, nil
	}
	
	logs := []types.Log{}
	entries := strings.Split(output, "---END---")
	
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		
		log := types.Log{}
		lines := strings.Split(entry, "\n")
		
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			
			if strings.HasPrefix(line, "COMMIT:") {
				log.Commit = strings.TrimPrefix(line, "COMMIT:")
			} else if strings.HasPrefix(line, "TREE:") {
				log.Tree = strings.TrimPrefix(line, "TREE:")
			} else if strings.HasPrefix(line, "PARENT:") {
				log.Parent = strings.TrimPrefix(line, "PARENT:")
			} else if strings.HasPrefix(line, "AUTHOR:") {
				log.Author = strings.TrimPrefix(line, "AUTHOR:")
			} else if strings.HasPrefix(line, "AUTHOR_DATE:") {
				dateStr := strings.TrimPrefix(line, "AUTHOR_DATE:")
				if date, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr); err == nil {
					log.AuthorDate = date
				}
			} else if strings.HasPrefix(line, "COMMITTER:") {
				log.Committer = strings.TrimPrefix(line, "COMMITTER:")
			} else if strings.HasPrefix(line, "COMMITTER_DATE:") {
				dateStr := strings.TrimPrefix(line, "COMMITTER_DATE:")
				if date, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr); err == nil {
					log.CommitterDate = date
				}
			} else if strings.HasPrefix(line, "MESSAGE:") {
				log.Message = strings.TrimPrefix(line, "MESSAGE:")
			}
		}
		
		logs = append(logs, log)
	}
	
	return logs, nil
}