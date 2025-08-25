package git

import (
	"github.com/instruqt/git-exec/pkg/git/types"
)

// Show shows information about a git object
func (g *gitImpl) Show(object string, opts ...Option) (*types.Log, error) {
	format := "--pretty=format:COMMIT:%H%nTREE:%T%nPARENT:%P%nAUTHOR:%an <%ae>%nAUTHOR_DATE:%ai%nCOMMITTER:%cn <%ce>%nCOMMITTER_DATE:%ci%nMESSAGE:%s%n---END---"
	cmd := g.newCommand("show", format, object)
	cmd.ApplyOptions(opts...)
	output, err := cmd.Execute()
	if err != nil {
		return nil, err
	}
	
	logs, err := parseLogOutput(string(output))
	if err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return nil, nil
	}
	return &logs[0], nil
}