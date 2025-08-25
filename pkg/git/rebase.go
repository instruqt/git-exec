package git

// Rebase rebases branches
func (g *gitImpl) Rebase(opts ...Option) error {
	cmd := g.newCommand("rebase")
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}