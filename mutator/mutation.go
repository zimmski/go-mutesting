package mutator

// Mutation defines the behavior of one mutation
type Mutation struct {
	// Change is called before executing the exec command.
	Change func()
	// Reset is called after executing the exec command.
	Reset func()
}
