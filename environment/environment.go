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

	// parent holds any parent environment.  Our env. allows
	// nesting to implement scope.
	parent *Environment
}

// New is our constructor
func New() *Environment {
	return &Environment{vars: make(map[string]string)}
}

// NewEnclosedEnvironment create new environment by outer parameter
func NewEnclosedEnvironment(parent *Environment) *Environment {
	env := New()
	env.parent = parent
	return env
}

// Clear removes a variable
func (e *Environment) Clear(name string) {
	delete(e.vars, name)
}

// Get retrieves a variable.
func (e *Environment) Get(name string) (string, bool) {

	// Get from this scope
	x, ok := e.vars[name]

	// If it failed then look in the parent.
	if !ok && e.parent != nil {
		x, ok = e.parent.Get(name)
	}
	return x, ok
}

// Set stores a variable, or updates an existing one.
func (e *Environment) Set(name string, value string) {

	// If the variable is in this scope then set it
	_, ok := e.vars[name]
	if ok {
		e.vars[name] = value
		return
	}

	// If the variable is in the parent-scope then we'll
	// update it there
	if e.parent != nil {
		_, ok = e.parent.Get(name)

		if ok {
			e.parent.Set(name, value)
		}
	}

	// Wasn't in this scope.
	// Wasn't in the parent-scope.
	// Create it as local.
	e.vars[name] = value
}
