package git

// Reset resets files from the staging area
func (g *gitImpl) Reset(files []string, opts ...Option) error {
	cmd := g.newCommand("reset")
	cmd.ApplyOptions(opts...)
	if len(files) > 0 {
		cmd.AddArgs(files...)
	}
	_, err := cmd.Execute()
	return err
}