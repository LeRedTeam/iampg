# Architecture Principles

This document defines the architectural constraints and patterns for IAM Auto-Policy Generator.

---

## Core Architectural Constraints

### 1. Local-First Execution

All processing happens on the user's machine or CI runner.

**Implications:**
- No server-side computation
- No API calls to our infrastructure
- No telemetry unless explicitly opted-in
- Binary must be self-contained

**Allowed:**
- License validation (offline-capable)
- Opt-in anonymous usage stats (future consideration)

### 2. Stateless Operation

The CLI maintains no state between invocations.

**Implications:**
- No database
- No local cache files (except license)
- No configuration files required (optional only)
- Each run is independent

**Why:**
- Eliminates state corruption bugs
- Simplifies debugging
- Enables easy CI integration
- Reduces maintenance burden

### 3. Zero Credential Storage

Credentials are never stored, logged, or transmitted.

**Implications:**
- Use ambient AWS credentials (environment, profiles, instance roles)
- Never write credentials to disk
- Never include credentials in error messages
- Never transmit credentials over network

### 4. Single Binary Distribution

The tool is distributed as a single, self-contained binary.

**Implications:**
- No runtime dependencies
- No package manager conflicts
- Simple installation (`curl | sh` or download)
- Cross-platform builds (Linux, macOS, Windows)

**Build targets:**
- `linux/amd64`
- `linux/arm64`
- `darwin/amd64`
- `darwin/arm64`
- `windows/amd64`

### 5. Deterministic Output

Given the same input, the output must be identical.

**Implications:**
- Sort policy statements consistently
- Sort actions alphabetically
- No timestamps in output (unless requested)
- No random elements

**Why:**
- Enables diffing between runs
- Enables caching in CI
- Reduces user confusion
- Simplifies testing

---

## Technical Architecture

### CLI Structure

```
iampg (binary)
├── cmd/           # Command implementations
│   ├── run.go     # Wrapper mode
│   ├── parse.go   # Log parsing mode
│   ├── refine.go  # Policy refinement (paid)
│   └── root.go    # Root command, help
├── capture/       # AWS call capture logic
│   ├── proxy.go   # HTTP proxy for SDK calls
│   └── trace.go   # Call tracing
├── parse/         # Log parsing logic
│   ├── cloudtrail.go
│   └── error.go
├── policy/        # Policy generation
│   ├── generate.go
│   ├── merge.go
│   └── format.go
├── refine/        # Policy analysis (paid)
│   ├── wildcards.go
│   └── scope.go
└── license/       # License validation
    └── validate.go
```

### Capture Mechanism

**Approach: AWS CLI Command Parsing (MVP)**

For MVP, we use direct CLI argument parsing:

1. Detect if command is AWS CLI (`aws ...`)
2. Parse service and subcommand from arguments
3. Map CLI commands to IAM actions (e.g., `s3 ls` → `s3:ListBucket`)
4. Extract resource ARNs from arguments (e.g., `s3://bucket/key`)
5. Execute command normally
6. Generate policy from parsed calls

**Why CLI parsing for MVP:**
- Zero runtime overhead
- No proxy/MITM complexity
- No certificate handling
- Works immediately with AWS CLI

**Limitations:**
- Only captures AWS CLI commands (not SDK calls)
- Cannot capture programmatic AWS SDK usage

**Future enhancements:**
- HTTP proxy with MITM for SDK capture
- eBPF tracing for more robust capture (Linux only)
- AWS SDK debug log parsing

### Policy Generation

**Input:** List of observed AWS API calls
```go
type ObservedCall struct {
    Service  string // e.g., "s3"
    Action   string // e.g., "GetObject"
    Resource string // e.g., "arn:aws:s3:::bucket/key"
    Region   string // e.g., "us-east-1"
}
```

**Output:** IAM policy document
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject"],
      "Resource": "arn:aws:s3:::bucket/key"
    }
  ]
}
```

**Generation rules:**
1. Group by service
2. Deduplicate actions
3. Preserve specific resource ARNs
4. Sort alphabetically for determinism
5. Never add wildcards (user's choice to broaden)

### Log Parsing

**CloudTrail format:**
```json
{
  "Records": [
    {
      "eventSource": "s3.amazonaws.com",
      "eventName": "GetObject",
      "resources": [
        {"ARN": "arn:aws:s3:::bucket/key"}
      ]
    }
  ]
}
```

**AccessDenied format:**
```
User: arn:aws:iam::123:user/dev is not authorized to perform: s3:GetObject on resource: arn:aws:s3:::bucket/key
```

**Parsing strategy:**
- Regex extraction for error messages
- JSON unmarshaling for CloudTrail
- Graceful handling of unknown formats

---

## Output Formats

### JSON (default)
Standard IAM policy document.

### YAML (optional)
For users who prefer YAML.

### Terraform (paid)
```hcl
resource "aws_iam_policy" "generated" {
  name   = "generated-policy"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [...]
  })
}
```

### SARIF (paid)
For CI integration with security scanners.

---

## Error Handling

### Principle: Explicit and Actionable

Every error must:
1. State what went wrong
2. Suggest how to fix it
3. Exit with appropriate code

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Wrapped command failed (policy still generated) |
| 4 | License validation failed |
| 5 | Capture failed (permissions, proxy issues) |

### Error Message Format

```
Error: <what happened>
Cause: <why it happened>
Fix: <how to resolve>
```

Example:
```
Error: Failed to start capture proxy
Cause: Port 8080 is already in use
Fix: Set IAMPG_PROXY_PORT to use a different port, or free port 8080
```

---

## Configuration

### No Required Configuration

The tool must work with zero configuration.

### Optional Configuration

Environment variables only (no config files):

| Variable | Default | Description |
|----------|---------|-------------|
| `IAMPG_PROXY_PORT` | random | Proxy port |
| `IAMPG_OUTPUT_FORMAT` | json | Default output format |
| `IAMPG_LICENSE_KEY` | none | License for paid features |
| `IAMPG_VERBOSE` | false | Verbose output |

### Why No Config File

- One less thing to manage
- One less thing to document
- One less thing to debug
- Environment variables work in CI
- Environment variables are explicit

---

## Testing Strategy

### Unit Tests (Implemented)
- `policy/policy_test.go` - Policy generation logic
- `parse/error_test.go` - AccessDenied parsing
- `parse/cloudtrail_test.go` - CloudTrail parsing
- `capture/runner_test.go` - AWS CLI argument parsing

### Running Tests
```bash
go test ./... -v
```

### CI Pipeline
Tests run automatically on:
- Every push to main
- Every pull request

See `.github/workflows/ci.yaml`

### Test Coverage
- Policy generation: determinism, deduplication, grouping
- Error parsing: multiple formats, edge cases
- CloudTrail parsing: single/multiple records, different formats
- CLI parsing: S3, DynamoDB, Lambda commands with flags

### Test Principle
Prefer real AWS calls in integration tests over mocks. Use test accounts or LocalStack where appropriate.

---

## Dependency Policy

### Allowed Dependencies
- Standard library
- Well-maintained, minimal-dependency libraries
- AWS SDK (for parsing ARNs, understanding service names)

### Forbidden Dependencies
- Web frameworks
- Database drivers
- Message queue clients
- Anything requiring runtime services

### Dependency Evaluation
Before adding a dependency, answer:
1. Can this be done with standard library?
2. Is this dependency actively maintained?
3. Does this dependency have minimal transitive dependencies?
4. Does this dependency have security vulnerabilities?
5. Does this dependency have a compatible license?

---

## Build and Release

### Build System
- Go modules for dependency management
- GoReleaser for cross-platform builds
- GitHub Actions for CI/CD

### Release Artifacts
- Binary for each platform
- SHA256 checksums
- Source tarball

### Distribution Channels
- GitHub Releases
- Homebrew tap
- Direct download from docs site (if any)
