// Copyright (C) 2026 LeRedTeam
// SPDX-License-Identifier: AGPL-3.0-or-later

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

// UpdateLast applies a mutation to the last captured call under the lock.
func (c *Capturer) UpdateLast(fn func(*policy.ObservedCall)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.calls) > 0 {
		fn(&c.calls[len(c.calls)-1])
	}
}

// Calls returns all observed calls.
func (c *Capturer) Calls() []policy.ObservedCall {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := make([]policy.ObservedCall, len(c.calls))
	copy(result, c.calls)
	return result
}
