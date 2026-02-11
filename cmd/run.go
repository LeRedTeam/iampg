package cmd

import (
	"fmt"
	"os"

	"github.com/LeRedTeam/iampg/capture"
	"github.com/LeRedTeam/iampg/policy"
	"github.com/spf13/cobra"
)

var runOutput string
var runFormat string
var runVerbose bool

var runCmd = &cobra.Command{
	Use:   "run -- <command>",
	Short: "Capture AWS calls from a command and generate IAM policy",
	Long: `Execute a command while capturing AWS API calls made during execution.
Generates a minimal IAM policy granting only the observed permissions.

Example:
  iampg run -- aws s3 ls
  iampg run -- aws s3 cp file.txt s3://bucket/
  iampg run --output policy.json -- terraform apply
  iampg run -v -- python deploy.py`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	RunE:                  runRun,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&runOutput, "output", "o", "", "Write policy to file (default: stdout)")
	runCmd.Flags().StringVarP(&runFormat, "format", "f", "json", "Output format: json")
	runCmd.Flags().BoolVarP(&runVerbose, "verbose", "v", false, "Show captured AWS calls")
}

func runRun(cmd *cobra.Command, args []string) error {
	runner := capture.NewRunner(runVerbose)

	// Run the command and capture calls
	calls, exitCode, err := runner.RunWithCloudTrailSim(args)
	if err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}

	// Generate policy from observed calls
	doc := policy.Generate(calls)

	// Output the policy
	output, err := doc.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to generate policy JSON: %w", err)
	}

	if runOutput != "" {
		if err := os.WriteFile(runOutput, output, 0644); err != nil {
			return fmt.Errorf("failed to write policy to %s: %w", runOutput, err)
		}
		fmt.Fprintf(os.Stderr, "Policy written to %s\n", runOutput)
	} else {
		fmt.Println(string(output))
	}

	// Report on captured calls
	if len(calls) == 0 {
		fmt.Fprintln(os.Stderr, "No AWS API calls detected.")
	} else {
		fmt.Fprintf(os.Stderr, "Captured %d AWS API call(s).\n", len(calls))
	}

	// Exit with the wrapped command's exit code if it failed
	if exitCode != 0 {
		fmt.Fprintf(os.Stderr, "Command exited with code %d. Policy generated from observed calls.\n", exitCode)
		os.Exit(exitCode)
	}

	return nil
}
