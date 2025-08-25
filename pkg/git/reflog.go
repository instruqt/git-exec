package git

// Reflog manages the reference logs
func (g *gitImpl) Reflog(opts ...Option) error {
	cmd := g.newCommand("reflog")
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}