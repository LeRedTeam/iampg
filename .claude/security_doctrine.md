# Security Doctrine

## Core Principles

1. **Never store credentials** — Not in memory beyond execution, not on disk, not anywhere
2. **Never transmit credentials** — No network calls containing credentials
3. **Never log sensitive data** — No API responses, request bodies, or ARN paths by default
4. **Explicit opt-in** — Any sensitive feature requires user action
5. **Local-first** — All processing on user's machine

---

## Credential Handling

### Allowed
- Read credentials from environment (AWS_ACCESS_KEY_ID, etc.)
- Read credentials from AWS credential files (~/.aws/credentials)
- Use instance roles, container roles, SSO
- Pass credentials through to wrapped commands

### Forbidden
- Write credentials to any file
- Log credentials in any form
- Include credentials in error messages
- Transmit credentials to any endpoint
- Cache credentials between runs

---

## Data Handling

### Generated Policies
- Output to stdout by default
- User controls where output goes
- No telemetry about policy contents

### CloudTrail Logs
- Read-only access
- Never modify input files
- Never transmit log contents

### Error Messages
- Sanitize ARNs (redact account IDs optionally)
- Never include credentials
- Never include full request/response bodies

---

## Network Security

### Default: No outbound connections
The CLI makes no network calls except:
1. Proxied AWS API calls (user's credentials, user's traffic)
2. License validation (if online check is ever added — currently offline-only)

### Proxy Mode
- HTTPS interception requires user acknowledgment
- CA certificate handling must be explicit
- No MITM by default — warn and require flag

---

## Supply Chain Security

### Dependencies
- Minimal dependencies
- Audit dependencies for vulnerabilities
- Pin dependency versions
- Use Go modules for reproducibility

### Build Process
- Reproducible builds where possible
- Sign release binaries
- Publish checksums
- Build in CI (auditable)

---

## Threat Model

| Threat | Mitigation |
|--------|------------|
| Credential theft | Never store, never log, never transmit |
| Policy exfiltration | No network calls, local-only processing |
| Malicious input | Validate all inputs, no shell injection |
| Supply chain attack | Minimal deps, audited deps, signed releases |
| Local privilege escalation | Run with minimal permissions, no sudo required |

---

## Security Incident Response

If a security issue is discovered:
1. Acknowledge within 24 hours
2. Assess severity
3. Patch and release
4. Notify users via GitHub and changelog
5. CVE if applicable
