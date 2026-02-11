package capture

import (
	"sync"

	"github.com/LeRedTeam/iampg/policy"
)

// Capturer observes AWS API calls during command execution.
type Capturer struct {
	calls []policy.ObservedCall
	mu    sync.Mutex
}

// New creates a new Capturer.
func New() *Capturer {
	return &Capturer{
		calls: make([]policy.ObservedCall, 0),
	}
}

// AddCall records an observed AWS API call.
func (c *Capturer) AddCall(call policy.ObservedCall) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.calls = append(c.calls, call)
}

// Calls returns all observed calls.
func (c *Capturer) Calls() []policy.ObservedCall {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := make([]policy.ObservedCall, len(c.calls))
	copy(result, c.calls)
	return result
}
