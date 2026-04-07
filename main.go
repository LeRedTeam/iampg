// Copyright (C) 2026 LeRedTeam
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"os"

	"github.com/LeRedTeam/iampg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
