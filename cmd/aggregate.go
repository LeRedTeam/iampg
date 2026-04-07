// Copyright (C) 2026 LeRedTeam
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/LeRedTeam/iampg/license"
	"github.com/LeRedTeam/iampg/policy"
	"github.com/spf13/cobra"
)

var aggregateFiles []string
var aggregateOutput string
var aggregateFormat string
var aggregateResourceName string

var aggregateCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Combine multiple policies into one (pro)",
	Long: `Aggregate multiple policy files into a single merged policy.

Useful for combining policies from multiple test runs or environments.

Examples:
  iampg aggregate --files policy1.json,policy2.json
  iampg aggregate --files policy1.json --files policy2.json --output combined.json
  iampg aggregate --files "*.json" --format terraform`,
	RunE: runAggregate,
}

func init() {
	rootCmd.AddCommand(aggregateCmd)
	aggregateCmd.Flags().StringSliceVar(&aggregateFiles, "files", []string{}, "Policy files to aggregate (required)")
	aggregateCmd.Flags().StringVarP(&aggregateOutput, "output", "o", "", "Output file (default: stdout)")
	aggregateCmd.Flags().StringVar(&aggregateFormat, "format", "json", "Output format: json, yaml, terraform")
	aggregateCmd.Flags().StringVar(&aggregateResourceName, "resource-name", "aggregated_policy", "Terraform resource name")
	aggregateCmd.MarkFlagRequired("files")
}

func runAggregate(cmd *cobra.Command, args []string) error {
	// Check license
	if err := license.RequireFeature("aggregate"); err != nil {
		return err
	}

	if len(aggregateFiles) == 0 {
		return fmt.Errorf("at least one policy file is required")
	}

	// Load and merge all policies
	// Group by (Effect, Resource) to preserve Deny statements
	type stmtKey struct {
		Effect   string
		Resource string
	}
	grouped := make(map[stmtKey]map[string]bool)

	for _, file := range aggregateFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		var doc policy.Document
		if err := json.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf("failed to parse %s: %w", file, err)
		}

		for _, stmt := range doc.Statement {
			key := stmtKey{Effect: stmt.Effect, Resource: stmt.Resource}
			if key.Resource == "" {
				key.Resource = "*"
			}
			if grouped[key] == nil {
				grouped[key] = make(map[string]bool)
			}
			for _, action := range stmt.Action {
				grouped[key][action] = true
			}
		}

		fmt.Fprintf(os.Stderr, "Loaded %d statements from %s\n", len(doc.Statement), file)
	}

	// Build merged document
	var statements []policy.Statement
	for key, actions := range grouped {
		sorted := make([]string, 0, len(actions))
		for a := range actions {
			sorted = append(sorted, a)
		}
		sort.Strings(sorted)
		statements = append(statements, policy.Statement{
			Effect:   key.Effect,
			Action:   sorted,
			Resource: key.Resource,
		})
	}
	sort.Slice(statements, func(i, j int) bool {
		if statements[i].Effect != statements[j].Effect {
			return statements[i].Effect < statements[j].Effect
		}
		return statements[i].Resource < statements[j].Resource
	})
	merged := &policy.Document{Version: "2012-10-17", Statement: statements}

	// Output
	if err := outputPolicy(merged, aggregateFormat, aggregateOutput, aggregateResourceName); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Aggregated %d files into %d statements.\n", len(aggregateFiles), len(merged.Statement))

	return nil
}

