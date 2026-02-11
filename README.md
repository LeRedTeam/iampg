# iampg - IAM Auto-Policy Generator

Generate minimal IAM policies by observing AWS API calls or parsing logs.

Stop guessing IAM permissions. Run your code, get your policy.

## GitHub Action

```yaml
- name: Generate IAM Policy
  uses: LeRedTeam/iampg@v1
  with:
    mode: parse
    cloudtrail: ./cloudtrail-logs.json
    output: policy.json

- name: Upload policy
  uses: actions/upload-artifact@v4
  with:
    name: iam-policy
    path: policy.json
```

## Installation

```bash
# Download latest release
curl -sSL https://github.com/LeRedTeam/iampg/releases/latest/download/iampg-$(uname -s)-$(uname -m) -o iampg
chmod +x iampg
sudo mv iampg /usr/local/bin/

# Or build from source
go install github.com/LeRedTeam/iampg@latest
```

## Usage

### Capture from AWS CLI commands

```bash
# Run an AWS command and generate the required policy
iampg run -- aws s3 ls s3://my-bucket/
iampg run -- aws dynamodb put-item --table-name Users --item '{"id":{"S":"1"}}'

# Save to file
iampg run --output policy.json -- aws s3 cp file.txt s3://bucket/

# Verbose mode (show captured calls)
iampg run -v -- aws lambda invoke --function-name MyFunc out.json
```

### Parse AccessDenied errors

```bash
# Parse a single error message
iampg parse --error "User: arn:aws:iam::123:user/dev is not authorized to perform: s3:GetObject on resource: arn:aws:s3:::bucket/key"

# Parse multiple errors from a file
cat errors.log | iampg parse --stdin

# Parse CloudTrail logs
iampg parse --cloudtrail trail.json
```

### Example output

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": "arn:aws:s3:::my-bucket/*"
    }
  ]
}
```

## Commands

| Command | Description |
|---------|-------------|
| `run -- <cmd>` | Execute command and capture AWS API calls |
| `parse` | Parse CloudTrail logs or AccessDenied errors |

### Run flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Write policy to file (default: stdout) |
| `-f, --format` | Output format: json (default) |
| `-v, --verbose` | Show captured AWS calls |

### Parse flags

| Flag | Description |
|------|-------------|
| `--cloudtrail` | CloudTrail JSON log file |
| `--error` | AccessDenied error message |
| `--stdin` | Read from stdin |
| `-o, --output` | Write policy to file |

## How it works

**Run mode:** Parses AWS CLI arguments to determine which IAM actions are being invoked and extracts resource ARNs from the command.

**Parse mode:** Uses regex patterns to extract service, action, and resource from CloudTrail events or AccessDenied error messages.

## Supported services

- S3
- DynamoDB
- Lambda
- SQS
- SNS
- STS
- IAM
- And more (generic action parsing)

## Security

- **No credential storage** - Uses your existing AWS credentials
- **No network calls** - All processing is local
- **No telemetry** - Nothing is sent anywhere
