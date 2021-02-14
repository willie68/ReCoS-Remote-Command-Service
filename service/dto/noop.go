package dto

// NoopCommand is a command to do nothing.
type NoopCommand struct {
	Parameters map[string]interface{}
}

// Execute a delay in the actual context
func (d *NoopCommand) Execute(a *Action) (bool, error) {
	return true, nil
}
