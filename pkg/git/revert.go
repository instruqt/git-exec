package git

// Revert reverts commits
func (g *gitImpl) Revert(opts ...Option) error {
	cmd := g.newCommand("revert")
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}