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
