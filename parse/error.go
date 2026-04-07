// Copyright (C) 2026 LeRedTeam
// SPDX-License-Identifier: AGPL-3.0-or-later

package parse

import (
	"regexp"
	"strings"

	"github.com/LeRedTeam/iampg/policy"
)

var (
	// Common AccessDenied patterns
	// Pattern: "User: arn:... is not authorized to perform: service:Action on resource: arn:..."
	accessDeniedPattern1 = regexp.MustCompile(
		`not authorized to perform:\s*([a-z0-9-]+):([A-Za-z0-9]+)\s+on resource:\s*(arn:[^\s]+)`,
	)
	// Pattern: "AccessDenied: service:Action on arn:..."
	accessDeniedPattern2 = regexp.MustCompile(
		`AccessDenied[:\s]+([a-z0-9-]+):([A-Za-z0-9]+)\s+on\s+(arn:[^\s]+)`,
	)
	// Pattern: "when calling the ActionName operation"
	accessDeniedPattern3 = regexp.MustCompile(
		`when calling the ([A-Za-z0-9]+) operation`,
	)
	// Pattern: "Action: service:Action"
	actionPattern = regexp.MustCompile(
		`[Aa]ction:\s*([a-z0-9-]+):([A-Za-z0-9]+)`,
	)
	// Pattern: "Resource: arn:..."
	resourcePattern = regexp.MustCompile(
		`[Rr]esource:\s*(arn:[^\s]+)`,
	)
	// Service from error context
	servicePattern = regexp.MustCompile(
		`([a-z0-9]+)\.amazonaws\.com`,
	)
)

// cleanARN strips trailing punctuation that may have been captured from sentence context.
func cleanARN(arn string) string {
	return strings.TrimRight(arn, ".,;)\"'")
}

// ParseAccessDenied parses an AccessDenied error message and returns observed calls.
func ParseAccessDenied(message string) []policy.ObservedCall {
	var calls []policy.ObservedCall

	// Try pattern 1: full IAM-style error
	if matches := accessDeniedPattern1.FindAllStringSubmatch(message, -1); matches != nil {
		for _, m := range matches {
			calls = append(calls, policy.ObservedCall{
				Service:  m[1],
				Action:   m[2],
				Resource: cleanARN(m[3]),
			})
		}
		return calls
	}

	// Try pattern 2: short form
	if matches := accessDeniedPattern2.FindAllStringSubmatch(message, -1); matches != nil {
		for _, m := range matches {
			calls = append(calls, policy.ObservedCall{
				Service:  m[1],
				Action:   m[2],
				Resource: cleanARN(m[3]),
			})
		}
		return calls
	}

	// Try to extract components separately
	var service, action, resource string

	// Find action
	if m := actionPattern.FindStringSubmatch(message); m != nil {
		service = m[1]
		action = m[2]
	} else if m := accessDeniedPattern3.FindStringSubmatch(message); m != nil {
		action = m[1]
	}

	// Find resource
	if m := resourcePattern.FindStringSubmatch(message); m != nil {
		resource = cleanARN(m[1])
	}

	// Find service from context if not found
	if service == "" {
		if m := servicePattern.FindStringSubmatch(message); m != nil {
			service = m[1]
		}
	}

	// Extract service from resource ARN if still missing
	if service == "" && resource != "" {
		service = extractServiceFromARN(resource)
	}

	if action != "" && service != "" {
		if resource == "" {
			resource = "*"
		}
		calls = append(calls, policy.ObservedCall{
			Service:  service,
			Action:   action,
			Resource: resource,
		})
	}

	return calls
}

// ParseMultipleErrors parses multiple error messages (one per line).
func ParseMultipleErrors(input string) []policy.ObservedCall {
	var allCalls []policy.ObservedCall
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		calls := ParseAccessDenied(line)
		allCalls = append(allCalls, calls...)
	}

	return allCalls
}

func extractServiceFromARN(arn string) string {
	// ARN format: arn:partition:service:region:account:resource
	parts := strings.Split(arn, ":")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}
