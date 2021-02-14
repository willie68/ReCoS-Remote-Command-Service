package dto

// NoopCommand is a command to do nothing.
type NoopCommand struct {
	Parameters map[string]interface{}
}

// Init nothing
func (d *NoopCommand) Init(a *Action) (bool, error) {
	return true, nil
}

// Stop nothing
func (d *NoopCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (d *NoopCommand) Execute(a *Action) (bool, error) {
	return true, nil
}
