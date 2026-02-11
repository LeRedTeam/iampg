# Decision Framework

## Feature Evaluation Checklist

Before proposing or implementing any feature:

```
[ ] Fits MVP definition (or explicitly approved expansion)
[ ] Maintains local-first architecture
[ ] Avoids credential storage
[ ] Minimizes infrastructure cost
[ ] Minimizes maintenance burden
[ ] Minimizes support requirements
[ ] Can be built in weeks, not months
[ ] Passes anti-overengineering test
```

---

## Anti-Overengineering Test

1. Does this reduce friction? → Must be yes
2. Does this reduce risk? → Should be yes/neutral
3. Does this reduce maintenance? → Should be yes/neutral
4. Does this increase complexity? → Must be minimal
5. Does this introduce recurring cost? → Must be no

**If complexity > value, reject.**

---

## Enterprise Request Evaluation

For any enterprise-oriented request:

1. Does this increase infrastructure cost? → Reject
2. Does this increase maintenance cost? → Reject
3. Does this increase support burden? → Reject
4. Does this violate product simplicity? → Reject

Enterprise must adapt to product, not vice versa.

---

## Scope Change Process

1. Document the proposed change
2. Evaluate against all checklists
3. Estimate build time and maintenance cost
4. Require explicit founder approval
5. Update .claude/ documentation
6. Implement only after approval

Default answer to scope expansion: **NO**

---

## Priority Framework

| Priority | Criteria |
|----------|----------|
| P0 | Security issues, data loss bugs |
| P1 | Core functionality broken |
| P2 | Significant UX issues |
| P3 | Nice-to-have improvements |
| P4 | Cosmetic issues |

Focus on P0-P1. P3-P4 only when P0-P2 are clear.

---

## Build vs Buy vs Skip

- **Build** if: Core to value prop, simple to maintain
- **Buy/Use** if: Commodity problem, good library exists
- **Skip** if: Out of scope, high maintenance, low value

Default: **Skip**

---

## Technical Decision Record

For significant decisions, document:

1. Context: What problem are we solving?
2. Options: What alternatives exist?
3. Decision: What did we choose?
4. Rationale: Why this option?
5. Consequences: What are the tradeoffs?

Store in `.claude/decisions/` if many accumulate.
