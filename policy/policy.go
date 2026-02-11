package policy

import (
	"encoding/json"
	"sort"
)

// Statement represents a single IAM policy statement.
type Statement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource string   `json:"Resource,omitempty"`
}

// Document represents an IAM policy document.
type Document struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
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
