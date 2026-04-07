// Copyright (C) 2026 LeRedTeam
// SPDX-License-Identifier: AGPL-3.0-or-later

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "iampg",
	Short: "IAM Auto-Policy Generator",
	Long:  `Generate minimal IAM policies by observing AWS API calls or parsing logs.`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SilenceUsage = true
}

