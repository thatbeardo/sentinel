package mocks

import "time"

// Context is the mocked version of go context
type Context struct {
	value interface{}
	err   error
	t     time.Time
	ok    bool
}

// Value function returnsany values set and associated with this context
func (ctx Context) Value(key interface{}) interface{} {
	return ctx.value
}

// Err represents any error associated with this context
func (ctx Context) Err() error {
	return ctx.err
}

// Done method to represent that the context needs to be cancelled
func (ctx Context) Done() <-chan struct{} {
	return make(<-chan struct{})
}

// Deadline function to represent hard deadline for the context to stop
func (ctx Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.t, ctx.ok
}
