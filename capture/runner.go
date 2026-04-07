// Copyright (C) 2026 LeRedTeam
// SPDX-License-Identifier: AGPL-3.0-or-later

package capture

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/LeRedTeam/iampg/policy"
)

// Runner executes a command and captures AWS API calls.
type Runner struct {
	capturer *Capturer
	verbose  bool
}

// NewRunner creates a new Runner.
func NewRunner(verbose bool) *Runner {
	return &Runner{
		capturer: New(),
		verbose:  verbose,
	}
}

// Run executes the command and returns observed calls and the exit code.
func (r *Runner) Run(args []string) ([]policy.ObservedCall, int, error) {
	if len(args) == 0 {
		return nil, 1, fmt.Errorf("no command provided")
	}

	cmd := exec.Command(args[0], args[1:]...)

	// Set up environment with AWS debug logging
	cmd.Env = append(os.Environ(),
		"AWS_DEBUG=true",
	)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, 1, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Pass through stdout
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return nil, 1, fmt.Errorf("failed to start command: %w", err)
	}

	// Read stderr and parse for AWS calls
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stderrPipe)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				// Pass through any remaining partial line
				if line != "" {
					fmt.Fprint(os.Stderr, line)
				}
				break
			}

			// Parse AWS debug output
			if call := r.parseDebugLine(line); call != nil {
				r.capturer.AddCall(*call)
				if r.verbose {
					fmt.Fprintf(os.Stderr, "[capture] %s:%s on %s\n", call.Service, call.Action, call.Resource)
				}
			} else {
				// Pass through non-capture stderr
				fmt.Fprint(os.Stderr, line)
			}
		}
	}()

	err = cmd.Wait()
	wg.Wait() // Ensure goroutine finishes before reading calls
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, 1, fmt.Errorf("command failed: %w", err)
		}
	}

	return r.capturer.Calls(), exitCode, nil
}

// AWS CLI debug output patterns
var (
	// Pattern: "2024-01-01 12:00:00,000 - MainThread - botocore.endpoint - DEBUG - Making request for OperationName"
	botocoreOpPattern = regexp.MustCompile(`Making request for (\w+)`)
	// Pattern: "2024-01-01 12:00:00,000 - MainThread - botocore.endpoint - DEBUG - https://service.region.amazonaws.com/"
	botocoreURLPattern = regexp.MustCompile(`https://([a-z0-9-]+)\.([a-z0-9-]+)\.amazonaws\.com`)
)

func (r *Runner) parseDebugLine(line string) *policy.ObservedCall {
	// Try to extract operation name
	if matches := botocoreOpPattern.FindStringSubmatch(line); matches != nil {
		return &policy.ObservedCall{
			Action: matches[1],
		}
	}

	// Try to extract service and region from URL
	if matches := botocoreURLPattern.FindStringSubmatch(line); matches != nil {
		service := matches[1]
		region := matches[2]

		// Update the last incomplete call with service info
		r.capturer.UpdateLast(func(call *policy.ObservedCall) {
			if call.Service == "" {
				call.Service = service
				call.Region = region
			}
		})
	}

	return nil
}

// RunWithProxy executes with an HTTP proxy (for non-TLS or with MITM).
func (r *Runner) RunWithProxy(args []string) ([]policy.ObservedCall, int, error) {
	proxy := NewProxy(r.capturer, r.verbose)
	addr, err := proxy.Start()
	if err != nil {
		return nil, 1, fmt.Errorf("failed to start proxy: %w", err)
	}
	defer proxy.Stop()

	if len(args) == 0 {
		return nil, 1, fmt.Errorf("no command provided")
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(),
		"HTTP_PROXY=http://"+addr,
		"HTTPS_PROXY=http://"+addr,
		"http_proxy=http://"+addr,
		"https_proxy=http://"+addr,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	return r.capturer.Calls(), exitCode, nil
}

// RunWithCloudTrailSim simulates by parsing AWS CLI commands directly.
// This is the most reliable method for AWS CLI.
func (r *Runner) RunWithCloudTrailSim(args []string) ([]policy.ObservedCall, int, error) {
	// Check if this is an AWS CLI command
	if len(args) > 0 && (args[0] == "aws" || strings.HasSuffix(args[0], "/aws")) {
		call := parseAWSCLIArgs(args)
		if call != nil {
			r.capturer.AddCall(*call)
			// s3 mv needs additional permissions beyond cp
			if call.Service == "s3" {
				isMv := false
				for _, a := range args {
					if a == "mv" {
						isMv = true
						break
					}
				}
				if isMv {
					// mv always needs DeleteObject on the source
					r.capturer.AddCall(policy.ObservedCall{
						Service:  "s3",
						Action:   "DeleteObject",
						Resource: call.Resource,
					})
				}
			}
			if r.verbose {
				fmt.Fprintf(os.Stderr, "[capture] %s:%s on %s\n", call.Service, call.Action, call.Resource)
			}
		}
	} else if len(args) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: '%s' is not an AWS CLI command. Only 'aws' CLI commands are captured in this mode.\n", args[0])
	}

	// Still run the command
	if len(args) == 0 {
		return nil, 1, fmt.Errorf("no command provided")
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			// Command not found, permission denied, etc.
			return r.capturer.Calls(), 1, fmt.Errorf("failed to execute command: %w", err)
		}
	}

	return r.capturer.Calls(), exitCode, nil
}

// parseAWSCLIArgs parses AWS CLI arguments to determine the API call.
func parseAWSCLIArgs(args []string) *policy.ObservedCall {
	if len(args) < 3 {
		return nil
	}

	// Skip 'aws' and find service and command
	var service, command string
	var positionalArgs []string

	// AWS CLI boolean flags that don't take a value
	booleanFlags := map[string]bool{
		"--recursive": true, "--no-sign-request": true, "--debug": true,
		"--no-verify-ssl": true, "--no-paginate": true, "--dryrun": true,
		"--dry-run": true, "--no-include-email": true, "--generate-cli-skeleton": true,
		"--cli-auto-prompt": true, "--no-cli-auto-prompt": true,
		"--only-show-errors": true, "--no-progress": true, "--quiet": true,
		"--force": true, "--exact-timestamps": true, "--delete": true,
		"--follow-symlinks": true, "--no-follow-symlinks": true,
		"--summarize": true, "--human-readable": true, "--paginate": true,
	}

	i := 1
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			// Handle --flag=value syntax
			if strings.Contains(arg, "=") {
				i++ // skip the flag=value, value is part of the flag
				continue
			}
			if booleanFlags[arg] {
				// Boolean flag: skip only the flag itself
				i++
			} else {
				// Flag with value: skip flag and its value
				i++
				if i < len(args) && !strings.HasPrefix(args[i], "--") {
					i++
				}
			}
			continue
		}
		if service == "" {
			service = arg
		} else if command == "" {
			command = arg
		} else {
			positionalArgs = append(positionalArgs, arg)
		}
		i++
	}

	if service == "" || command == "" {
		return nil
	}

	// Map CLI service names to IAM service names
	serviceMapping := map[string]string{
		"s3api": "s3",
	}
	if mapped, ok := serviceMapping[service]; ok {
		service = mapped
	}

	// Extract resource from arguments
	resource := extractResourceFromArgs(service, command, args)

	// Map CLI command to IAM action
	action := cliCommandToAction(service, command)

	// Special cases for S3 commands that depend on arguments
	if service == "s3" {
		switch command {
		case "ls":
			hasS3Path := false
			for _, arg := range args {
				if strings.HasPrefix(arg, "s3://") {
					hasS3Path = true
					break
				}
			}
			if hasS3Path {
				action = "ListBucket"
			} else {
				action = "ListAllMyBuckets"
				resource = "*"
			}
		case "cp", "mv", "sync":
			// Determine direction: if first s3:// path is the source, it's a download
			firstS3Idx := -1
			for idx, arg := range positionalArgs {
				if strings.HasPrefix(arg, "s3://") {
					firstS3Idx = idx
					break
				}
			}
			if firstS3Idx == 0 {
				// Source is S3 → download
				action = "GetObject"
			}
			// else destination is S3 → PutObject (already the default)
		}
	}

	return &policy.ObservedCall{
		Service:  service,
		Action:   action,
		Resource: resource,
	}
}

// cliCommandToAction maps AWS CLI commands to IAM actions.
// For S3 ls, we need args to determine if it's ListAllMyBuckets or ListBucket
func cliCommandToAction(service, command string) string {
	// Special mappings for services where CLI commands don't match IAM actions
	// Note: s3 ls is handled specially in parseAWSCLIArgs
	s3Actions := map[string]string{
		"cp":      "PutObject", // Could be GetObject too, determined by direction
		"mv":      "PutObject",
		"rm":      "DeleteObject",
		"mb":      "CreateBucket",
		"rb":      "DeleteBucket",
		"sync":    "PutObject",
		"presign": "GetObject",
		"website": "PutBucketWebsite",
	}

	if service == "s3" {
		if action, ok := s3Actions[command]; ok {
			return action
		}
	}

	// Convert kebab-case to PascalCase for standard mappings
	parts := strings.Split(command, "-")
	var result strings.Builder
	for _, p := range parts {
		if len(p) > 0 {
			result.WriteString(strings.ToUpper(p[:1]))
			if len(p) > 1 {
				result.WriteString(p[1:])
			}
		}
	}
	return result.String()
}

// extractResourceFromArgs extracts resource ARN from CLI arguments.
// getFlagValue returns the value of a CLI flag, handling both "--flag value" and "--flag=value" syntax.
func getFlagValue(args []string, flag string) string {
	for i, arg := range args {
		if arg == flag && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(arg, flag+"=") {
			return strings.TrimPrefix(arg, flag+"=")
		}
	}
	return ""
}

func extractResourceFromArgs(service, command string, args []string) string {
	switch service {
	case "s3":
		return extractS3Resource(args)
	case "dynamodb":
		return extractDynamoDBResource(args)
	case "lambda":
		return extractLambdaResource(args)
	case "sqs":
		return extractSQSResource(args)
	case "sns":
		return extractSNSResource(args)
	case "sts":
		return extractSTSResource(args)
	case "iam":
		return extractIAMResource(args)
	case "secretsmanager":
		return extractSecretsManagerResource(args)
	case "ssm":
		return extractSSMResource(args)
	case "logs":
		return extractCloudWatchLogsResource(args)
	case "kms":
		return extractKMSResource(args)
	default:
		return "*"
	}
}

func extractS3Resource(args []string) string {
	// Check for s3:// paths first (s3 CLI style)
	var s3Paths []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "s3://") {
			s3Paths = append(s3Paths, arg)
		}
	}

	if len(s3Paths) > 0 {
		// Parse first S3 path (for ls, this is the target)
		path := strings.TrimPrefix(s3Paths[0], "s3://")
		path = strings.TrimSuffix(path, "/")
		parts := strings.SplitN(path, "/", 2)
		bucket := parts[0]

		if bucket == "" {
			return "arn:aws:s3:::*"
		}
		if len(parts) > 1 && parts[1] != "" {
			return "arn:aws:s3:::" + bucket + "/" + parts[1]
		}
		return "arn:aws:s3:::" + bucket + "/*"
	}

	// Check for --bucket and --key flags (s3api CLI style)
	bucket := getFlagValue(args, "--bucket")
	key := getFlagValue(args, "--key")

	if bucket != "" {
		if key != "" {
			return "arn:aws:s3:::" + bucket + "/" + key
		}
		return "arn:aws:s3:::" + bucket
	}

	return "arn:aws:s3:::*"
}

func extractDynamoDBResource(args []string) string {
	table := getFlagValue(args, "--table-name")
	if table != "" {
		return "arn:aws:dynamodb:*:*:table/" + table
	}
	return "*"
}

func extractLambdaResource(args []string) string {
	fn := getFlagValue(args, "--function-name")
	if fn != "" {
		return "arn:aws:lambda:*:*:function:" + fn
	}
	return "*"
}

func extractSQSResource(args []string) string {
	queueURL := getFlagValue(args, "--queue-url")
	if queueURL != "" {
		parts := strings.Split(queueURL, "/")
		if len(parts) >= 5 {
			return "arn:aws:sqs:*:" + parts[3] + ":" + parts[4]
		}
	}
	return "*"
}

func extractSNSResource(args []string) string {
	if arn := getFlagValue(args, "--topic-arn"); arn != "" {
		return arn
	}
	if arn := getFlagValue(args, "--target-arn"); arn != "" {
		return arn
	}
	return "*"
}

func extractSTSResource(args []string) string {
	if arn := getFlagValue(args, "--role-arn"); arn != "" {
		return arn
	}
	return "*"
}

func extractIAMResource(args []string) string {
	if arn := getFlagValue(args, "--policy-arn"); arn != "" {
		return arn
	}
	if name := getFlagValue(args, "--role-name"); name != "" {
		return "arn:aws:iam::*:role/" + name
	}
	if name := getFlagValue(args, "--user-name"); name != "" {
		return "arn:aws:iam::*:user/" + name
	}
	if name := getFlagValue(args, "--group-name"); name != "" {
		return "arn:aws:iam::*:group/" + name
	}
	return "*"
}

func extractSecretsManagerResource(args []string) string {
	secretID := getFlagValue(args, "--secret-id")
	if secretID != "" {
		if strings.HasPrefix(secretID, "arn:") {
			return secretID
		}
		return "arn:aws:secretsmanager:*:*:secret:" + secretID
	}
	return "*"
}

func extractSSMResource(args []string) string {
	name := getFlagValue(args, "--name")
	if name == "" {
		namesList := getFlagValue(args, "--names")
		if namesList != "" {
			name = strings.Split(namesList, ",")[0]
		}
	}
	if name != "" {
		if strings.HasPrefix(name, "arn:") {
			return name
		}
		if !strings.HasPrefix(name, "/") {
			name = "/" + name
		}
		return "arn:aws:ssm:*:*:parameter" + name
	}
	return "*"
}

func extractCloudWatchLogsResource(args []string) string {
	logGroup := getFlagValue(args, "--log-group-name")
	if logGroup != "" {
		return "arn:aws:logs:*:*:log-group:" + logGroup
	}
	return "*"
}

func extractKMSResource(args []string) string {
	keyID := getFlagValue(args, "--key-id")
	if keyID != "" {
		if strings.HasPrefix(keyID, "arn:") {
			return keyID
		}
		return "arn:aws:kms:*:*:key/" + keyID
	}
	if alias := getFlagValue(args, "--alias-name"); alias != "" {
		return "arn:aws:kms:*:*:" + alias
	}
	return "*"
}
