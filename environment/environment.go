// Package environment holds variable names/values.
//
// It isn't used so much at the moment, but in the future we'll want
// to use scopes/sub-frames for when we allow users to bind procedures,
// or write functions of their own.
package environment

// Environment is our state-holding object
type Environment struct {

	// vars holds the actual variables
	vars map[string]string
}

// New is our constructor
func New() *Environment {
	return &Environment{vars: make(map[string]string)}
}

// Clear removes a variable
func (e *Environment) Clear(name string) {
	delete(e.vars, name)
}

// Get retrieves a variable.
func (e *Environment) Get(name string) (string, bool) {
	x, ok := e.vars[name]
	return x, ok
}

// Set stores a variable, or updates an existing one.
func (e *Environment) Set(name string, value string) {
	e.vars[name] = value
}
