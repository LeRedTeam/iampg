# Roadmap

## Phase 1: MVP (v1.0)

**Goal:** Working CLI that generates policies from command execution and log parsing.

**Scope:**
- `run` command with proxy capture
- `parse` command for CloudTrail and errors
- JSON output
- Cross-platform binaries

**Timeline:** 4 weeks

---

## Phase 2: Stabilization (v1.x)

**Goal:** Fix bugs, respond to community feedback, harden core functionality.

**Scope:**
- Bug fixes from user reports
- Error message improvements
- Documentation improvements
- Edge case handling

**Timeline:** Ongoing, 2-3 months post-launch

**Constraint:** No new features, stability only.

---

## Phase 3: Paid Features (v2.0)

**Goal:** Introduce paid tier with refinement features.

**Scope:**
- `refine` command
- Wildcard detection
- Scoping suggestions
- Policy diff/comparison
- Terraform output format
- YAML output format
- License key validation

**Timeline:** 3-6 months post-launch

**Prerequisite:** MVP adoption validates demand.

---

## Phase 4: CI Integration (v2.x)

**Goal:** First-class CI/CD support.

**Scope:**
- SARIF output for security scanners
- CI enforcement mode (fail on broad policies)
- Multi-run aggregation
- GitHub Action wrapper

**Timeline:** 6-9 months post-launch

---

## Phase 5: Polish (v3.0)

**Goal:** Quality-of-life improvements based on user feedback.

**Scope:** TBD based on usage patterns

**Timeline:** 12+ months

---

## Hard Boundaries (Never on Roadmap)

- Web UI/dashboard
- SaaS backend
- User accounts
- Database
- Multi-cloud support
- Real-time monitoring
- Compliance frameworks
- Enterprise contracts

See `non_goals.md` for complete list.

---

## Roadmap Evolution Rules

1. Features move from future phases to current only with explicit approval
2. Features never move from "Never" to roadmap without documented justification
3. Timeline is indicative, not commitment
4. Scope reduction is allowed; scope expansion requires approval
