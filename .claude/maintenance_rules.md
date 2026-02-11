# Maintenance Rules

## Design for Maintainability

Every feature must pass: "Can a solo founder maintain this for 3 years?"

---

## Code Principles

1. **Simple over clever** — Readable code beats smart code
2. **Explicit over implicit** — No magic, no hidden behavior
3. **Standard library first** — Only add deps when necessary
4. **Few files over many** — Consolidate where sensible
5. **No dead code** — Delete unused code immediately

---

## Dependency Rules

### Before adding a dependency:
- [ ] Can this be done with stdlib?
- [ ] Is it actively maintained?
- [ ] Does it have minimal transitive deps?
- [ ] Is it security-audited?
- [ ] Is the license compatible?

### Dependency budget: <20 direct dependencies

---

## Technical Debt

### Allowed
- TODOs with issue numbers
- Temporary workarounds with expiration dates
- Known limitations documented

### Forbidden
- Undocumented hacks
- "Fix later" without tracking
- Commented-out code

---

## Release Cadence

- **Patches**: As needed for bugs/security
- **Minor**: Monthly at most
- **Major**: Annually at most

No release pressure. Stability over features.

---

## Documentation Requirements

- README: Installation and basic usage
- CLI help: All commands documented
- CHANGELOG: All changes tracked
- .claude/: Architecture decisions

No external docs site for MVP.

---

## Deprecation Policy

1. Announce deprecation in changelog
2. Warn in CLI output for 2 minor versions
3. Remove in next major version
4. Never break without major version bump
