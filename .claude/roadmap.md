# Roadmap

## Current Status: Ready for v1.0 Release

All planned features are implemented and tested.

## Completed Features

### Free Tier
- [x] `run` command - Capture AWS CLI calls
- [x] `parse` command - Parse CloudTrail and AccessDenied errors
- [x] JSON output format

### Pro Tier
- [x] License validation (Ed25519 signed keys)
- [x] YAML output format
- [x] Terraform output format
- [x] SARIF output format
- [x] `refine` command - Policy analysis
- [x] Wildcard detection
- [x] Scoping suggestions
- [x] Dangerous permission detection
- [x] Policy diff/comparison
- [x] CI enforcement mode (--enforce)
- [x] `aggregate` command - Multi-policy merging

### Infrastructure
- [x] GitHub Actions CI (tests on every push)
- [x] GitHub Actions Release (GoReleaser)
- [x] GitHub Action for Marketplace
- [x] Cross-platform builds (linux/darwin/windows, amd64/arm64)
- [x] Unit tests (29 tests passing)

---

## Release Checklist

### Before v1.0.0
- [ ] Create production keypair ✅
- [ ] Set IAMPG_PUBLIC_KEY secret ✅
- [ ] Make repo public
- [ ] Create v1.0.0 tag
- [ ] Publish to GitHub Marketplace
- [ ] Add Marketplace categories (Security, CI/CD)

### Post-Launch
- Bug fixes based on user feedback
- Documentation improvements
- Additional AWS service support

---

## Future Considerations (Not Committed)

These may be considered based on user demand:

| Feature | Consideration |
|---------|---------------|
| SDK capture | HTTP proxy with MITM for non-CLI usage |
| More services | Extended action mappings for niche services |
| Policy templates | Pre-built policies for common patterns |
| VS Code extension | IDE integration |

**Constraint:** Any future feature must pass the anti-overengineering test in `decision_framework.md`.

---

## Hard Boundaries (Never on Roadmap)

- Web UI/dashboard
- SaaS backend
- User accounts/database
- Multi-cloud support
- Real-time monitoring
- Compliance frameworks
- Enterprise contracts

See `non_goals.md` for complete list.
