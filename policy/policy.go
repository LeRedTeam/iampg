// Copyright (C) 2026 LeRedTeam
// SPDX-License-Identifier: AGPL-3.0-or-later

package policy

import (
	"encoding/json"
	"sort"
)

// Statement represents a single IAM policy statement.
type Statement struct {
	Effect   string   `json:"Effect" yaml:"Effect"`
	Action   []string `json:"Action" yaml:"Action"`
	Resource string   `json:"Resource,omitempty" yaml:"Resource,omitempty"`
}

// UnmarshalJSON handles IAM's flexible format where Action/Resource can be string or array.
func (s *Statement) UnmarshalJSON(data []byte) error {
	type Alias Statement
	aux := &struct {
		Action   interface{} `json:"Action"`
		Resource interface{} `json:"Resource"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle Action as string or array
	switch v := aux.Action.(type) {
	case string:
		s.Action = []string{v}
	case []interface{}:
		s.Action = make([]string, 0, len(v))
		for _, a := range v {
			if str, ok := a.(string); ok {
				s.Action = append(s.Action, str)
			}
		}
	}

	// Handle Resource as string or array (take first if array)
	switch v := aux.Resource.(type) {
	case string:
		s.Resource = v
	case []interface{}:
		if len(v) > 0 {
			if str, ok := v[0].(string); ok {
				s.Resource = str
			}
		}
	}

	return nil
}

// Document represents an IAM policy document.
type Document struct {
	Version   string      `json:"Version" yaml:"Version"`
	Statement []Statement `json:"Statement" yaml:"Statement"`
}

// ObservedCall represents a captured AWS API call.
type ObservedCall struct {
	Service  string
	Action   string
	Resource string
	Region   string
}

// Generate creates an IAM policy document from observed calls.
func Generate(calls []ObservedCall) *Document {
	if len(calls) == 0 {
		return &Document{
			Version:   "2012-10-17",
			Statement: []Statement{},
		}
	}

	// Group by service and resource
	type key struct {
		service  string
		resource string
	}
	grouped := make(map[key][]string)

	for _, c := range calls {
		k := key{service: c.Service, resource: c.Resource}
		action := c.Service + ":" + c.Action
		grouped[k] = appendUnique(grouped[k], action)
	}

	var statements []Statement
	for k, actions := range grouped {
		sort.Strings(actions)
		stmt := Statement{
			Effect:   "Allow",
			Action:   actions,
			Resource: k.resource,
		}
		if stmt.Resource == "" {
			stmt.Resource = "*"
		}
		statements = append(statements, stmt)
	}

	// Sort statements for determinism
	sort.Slice(statements, func(i, j int) bool {
		if statements[i].Resource != statements[j].Resource {
			return statements[i].Resource < statements[j].Resource
		}
		if len(statements[i].Action) > 0 && len(statements[j].Action) > 0 {
			return statements[i].Action[0] < statements[j].Action[0]
		}
		return false
	})

	return &Document{
		Version:   "2012-10-17",
		Statement: statements,
	}
}

// ToJSON converts the policy document to JSON.
func (d *Document) ToJSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}
