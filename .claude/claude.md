# Claude Guidelines for IAM Auto-Policy Generator

This document defines how Claude must operate when assisting with this project.

---

## Product Philosophy

IAM Auto-Policy Generator is a **narrow automation tool** that solves one problem:
developers don't know what IAM permissions their code actually needs.

The product observes real AWS API calls and generates minimal IAM policies.

It is:
- Developer-first
- CLI-first
- Local-first
- Privacy-first

It is NOT:
- A governance platform
- A security suite
- A SaaS product
- A consulting tool

---

## Scope Guardrails

Claude must NEVER propose features that:

1. Require an always-on backend
2. Store AWS credentials
3. Persist sensitive logs without explicit opt-in
4. Add UI/dashboard components
5. Support clouds other than AWS
6. Implement multi-account orchestration
7. Add compliance/audit frameworks
8. Create role lifecycle management

If a user or requirement suggests any of these, Claude must:
1. Acknowledge the request
2. Explain why it violates product constraints
3. Suggest an alternative that fits within scope (if possible)
4. Decline to implement if no alternative exists

---

## Enterprise Guardrail

Enterprise customers may use this product.
Enterprise-driven complexity is NOT allowed.

Before proposing any enterprise-oriented feature, Claude must evaluate:

- Does this increase infrastructure cost? → If yes, reject.
- Does this increase maintenance burden? → If yes, reject.
- Does this increase support requirements? → If yes, reject.
- Does this violate product simplicity? → If yes, reject.

Enterprise customers must accept product constraints.
No custom development for enterprise use cases unless explicitly approved.

---

## Cost Discipline Rules

Every architectural decision must consider:

1. **Infrastructure cost**: Prefer zero-cost solutions (local execution, no backend)
2. **Maintenance cost**: Prefer simple, stateless designs
3. **Support cost**: Prefer self-documenting, deterministic behavior
4. **Development cost**: Prefer weeks-to-build over months-to-build

Claude must never propose:
- Database dependencies
- Message queue dependencies
- Container orchestration
- Serverless backends (unless explicitly approved)
- Third-party service integrations

---

## Maintenance Discipline

Code must be:
- Simple enough to maintain solo
- Stateless where possible
- Deterministic in behavior
- Explicit in failure modes

Claude must:
- Prefer standard library over dependencies
- Prefer fewer files over more files
- Prefer straightforward code over clever code
- Avoid abstraction layers unless they reduce complexity

---

## Security Doctrine

1. **Never store credentials**: Not in files, not in memory longer than execution, not anywhere.
2. **Never transmit credentials**: No network calls with credentials.
3. **Never log sensitive data by default**: Raw API responses, request bodies, etc.
4. **Explicit opt-in only**: Any feature that touches sensitive data requires explicit user action.
5. **Local-first**: All processing happens on user's machine or CI runner.

---

## Anti-Overengineering Rule

Before proposing ANY feature, Claude must answer:

1. Does this reduce friction? → Must be yes.
2. Does this reduce risk? → Should be yes or neutral.
3. Does this reduce maintenance? → Should be yes or neutral.
4. Does this increase complexity? → Must be minimal.
5. Does this introduce recurring cost? → Must be no (or explicitly justified).

If complexity > value, reject the feature.

**Examples of rejected patterns:**
- Plugin architectures
- Configuration file formats
- Database schemas
- API versioning strategies
- Feature flags
- A/B testing infrastructure

---

## Documentation Update Rule

When Claude makes architectural changes, it must:

1. Update relevant documentation files in `.claude/`
2. Ensure `mvp_definition.md` reflects current scope
3. Update `roadmap.md` if timeline changes
4. Note any new constraints in `architecture_principles.md`

Documentation must stay synchronized with implementation.

---

## Decision Evaluation Checklist

For any significant decision, Claude must complete this checklist:

```
[ ] Does this fit the MVP definition?
[ ] Does this avoid scope expansion?
[ ] Does this maintain local-first architecture?
[ ] Does this avoid credential storage?
[ ] Does this minimize infrastructure cost?
[ ] Does this minimize maintenance burden?
[ ] Does this minimize support requirements?
[ ] Is this simple enough to maintain solo?
[ ] Can this be built in weeks, not months?
[ ] Does this pass the anti-overengineering test?
```

If any answer is NO, the decision must be reconsidered or explicitly justified.

---

## Explicit Claude Directives

Claude must:

1. **Never expand scope without explicit approval** — Even if a feature seems useful, if it's not in MVP, don't propose it without flagging scope expansion.

2. **Always evaluate infrastructure cost impact** — Before suggesting any architecture, estimate ongoing costs.

3. **Always evaluate maintenance impact** — Consider: "Can a solo founder maintain this for 3 years?"

4. **Always evaluate support burden** — Will this generate support tickets? How many? How complex?

5. **Prioritize simplicity over completeness** — A working 80% solution beats a complex 100% solution.

6. **Update documentation when architecture evolves** — Every significant change must be reflected in `.claude/` docs.

7. **Reject features that increase complexity without proportional value** — The burden of proof is on the feature, not the rejection.

---

## Communication Style

When discussing this project, Claude must:

- Be direct and concise
- Avoid enterprise jargon
- Focus on implementation over theory
- Provide concrete examples
- Flag scope violations immediately
- Suggest simpler alternatives when possible

---

## File Reference

Related documentation:
- `product_vision.md` — Why this product exists
- `non_goals.md` — What this product will never be
- `architecture_principles.md` — How to build it
- `monetization_guardrails.md` — How to charge for it
- `security_doctrine.md` — How to keep it safe
- `maintenance_rules.md` — How to keep it maintainable
- `support_policy.md` — How to handle support
- `decision_framework.md` — How to make decisions
- `mvp_definition.md` — What version 1 includes
- `roadmap.md` — What comes after version 1
