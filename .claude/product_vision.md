# Product Vision

## The Problem

Developers building on AWS face a constant friction point: IAM permissions.

**The over-permission trap:**
Most developers default to broad policies (`AdministratorAccess`, `PowerUserAccess`, or `*` wildcards) because:
- They don't know exactly what permissions their code needs
- They don't have time to debug AccessDenied errors
- IAM documentation is vast and confusing
- Trial-and-error is slow and frustrating

**The under-permission trap:**
When developers try to follow least-privilege:
- They spend hours debugging AccessDenied errors
- They add permissions one-by-one through iteration
- They often miss required permissions for edge cases
- The feedback loop is painfully slow

**The result:**
- Production environments with overly broad permissions (security risk)
- Developer frustration and lost productivity
- Security teams fighting a losing battle

---

## The Solution

**IAM Auto-Policy Generator** closes the feedback loop between code execution and policy authoring.

Instead of guessing permissions, developers:
1. Run their code through our CLI wrapper
2. We observe which AWS APIs are actually called
3. We generate a minimal policy granting exactly those permissions

**The value proposition:**
- From hours of debugging → seconds of generation
- From guessing → knowing
- From over-permission → least-privilege
- From manual → automated

---

## Why Now

1. **AWS complexity is increasing**: More services, more actions, more conditions
2. **Security requirements are tightening**: Zero-trust, least-privilege are now expected
3. **Developer tooling is maturing**: Developers expect automation for tedious tasks
4. **CI/CD is standard**: Policies can be generated and enforced in pipelines

---

## Target User

**Primary persona: Backend Developer**
- Building applications on AWS
- Uses 3-10 AWS services regularly
- Not a security expert
- Values time over perfection
- Runs code locally and in CI

**Secondary persona: DevOps Engineer**
- Manages CI/CD pipelines
- Responsible for deployment roles
- Wants reproducible, auditable policies
- Needs CI integration

**Not targeting:**
- Security auditors (they have specialized tools)
- Compliance officers (different problem space)
- Multi-cloud architects (AWS-only focus)
- Enterprise platform teams (too complex for our scope)

---

## Product Principles

1. **Observe, don't guess**: Base permissions on actual behavior
2. **Minimal by default**: Generate the smallest viable policy
3. **Local-first**: No cloud backend required
4. **Trust nothing**: Never store or transmit credentials
5. **Fit existing workflows**: CLI and CI integration
6. **Self-serve only**: No onboarding required

---

## Success Metrics

For a solo-founder micro-SaaS:

**Adoption:**
- GitHub stars (awareness)
- CLI downloads (usage)
- GitHub Action installations (stickiness)

**Revenue:**
- Paid license activations
- License renewals
- Revenue per user

**Sustainability:**
- Support tickets per user (lower is better)
- Time to resolve issues (lower is better)
- Infrastructure cost per user (should be ~zero)

---

## What This Product Is NOT

See `non_goals.md` for explicit exclusions.

In short: This is a narrow automation tool, not a platform, not a suite, not a service.
