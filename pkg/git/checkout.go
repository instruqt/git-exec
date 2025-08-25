package git

// Checkout checks out branches, commits, or files
func (g *gitImpl) Checkout(opts ...Option) error {
	cmd := g.newCommand("checkout")
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}