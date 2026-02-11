# Monetization Guardrails

This document defines the monetization strategy and constraints for IAM Auto-Policy Generator.

---

## Pricing Philosophy

1. **Indie-friendly**: Pricing that individual developers can afford
2. **Sustainable**: Revenue that covers development time and minimal infra
3. **Simple**: No complex tiers, no enterprise negotiations
4. **Self-serve**: No sales calls, no custom contracts

---

## Tier Structure

### Free Tier

**Price:** $0

**Includes:**
- `run` command — Capture and generate policy from command execution
- `parse` command — Parse CloudTrail logs and AccessDenied errors
- JSON output format
- Unlimited local usage
- Community support (GitHub issues)

**Why free:**
- Lower barrier to adoption
- Builds trust and awareness
- Free tier users become paid tier advocates
- Core functionality has near-zero marginal cost

### Pro Tier

**Price:** $19/month or $149/year (2 months free)

**Includes everything in Free, plus:**
- `refine` command — Policy analysis and improvement suggestions
- Wildcard detection and scoping suggestions
- Policy diff/drift comparison
- Terraform output format
- YAML output format
- SARIF output for CI integration
- CI enforcement mode (fail on overly-broad policies)
- Multi-run aggregation
- Priority support (48h response time)

**Why this price:**
- Affordable for individual developers ($0.63/day)
- Significant value over free tier
- Sustainable for solo founder
- No negotiation needed

### Team Tier (Future Consideration)

**Price:** $49/month per seat or $399/year

**Includes everything in Pro, plus:**
- Volume licensing
- Invoice billing
- Company-wide license key
- Shared policy baselines

**Constraints:**
- No custom features
- No SLA beyond response time
- No dedicated support
- No on-premise deployment

**When to introduce:**
- Only after Pro tier proves sustainable
- Only if demand is clear
- Only if it doesn't increase support burden significantly

---

## License Model

### Key-Based Licensing (Implemented)

- License key unlocks paid features
- Key validated locally (no phone-home required)
- Key tied to email address
- Key has expiration date
- Set via `IAMPG_LICENSE_KEY` environment variable

### License Key Format

License keys are Ed25519-signed JSON payloads containing:
- Email address
- Tier (pro, team)
- Issued date
- Expiration date

The public key is embedded in the binary at build time.

### License Enforcement

| Context | Enforcement |
|---------|-------------|
| Local development | Honor system |
| CI/CD pipeline | License check on each run |
| Output metadata | License info in generated policy comments (optional) |

### Grace Period

- 7-day grace period when license expires
- Clear warning messages during grace period
- Features disabled after grace period

### No Phone-Home

The CLI does not require internet access for license validation.

**Implementation:**
- Ed25519 cryptographic signatures
- Public key embedded in binary
- Private key stored securely (GitHub secret)
- Local validation only

### Admin Commands (Hidden)

```bash
# Generate keypair (one-time setup)
iampg license generate-keypair

# Generate license key for customer
iampg license generate --email user@example.com --tier pro --days 365 --private-key $KEY

# Check license status
iampg license status
```

### GitHub Secrets Required

- `IAMPG_PUBLIC_KEY` - Public key for release builds
- `IAMPG_PRIVATE_KEY` - Private key for license generation (keep secure, not in repo)

---

## Payment Processing

### Provider: Stripe or LemonSqueezy

- Self-serve checkout
- Automatic billing
- Customer portal for management
- No manual invoice processing

### No Enterprise Contracts

- No custom pricing
- No annual invoicing
- No PO numbers
- No contract negotiations

If a company requires these, they can:
- Use credit card like everyone else
- Use a purchasing card
- Expense the subscription

We do not accommodate procurement processes.

---

## Pricing Anti-Patterns to Avoid

### No Usage-Based Pricing
- Unpredictable for users
- Complex to track
- Requires metering infrastructure

### No Per-Seat Pricing for Small Teams
- Friction for adoption
- Complex license management
- Punishes growth

### No Enterprise Tier with SLA
- SLA requires on-call commitment
- Unpredictable support burden
- Forces platform mentality

### No Custom Development
- Distorts roadmap
- Creates maintenance burden
- Violates product focus

### No Consulting Bundles
- Different business model
- Unpredictable time commitment
- Doesn't scale

---

## Revenue Sustainability Model

### Target Metrics

For solo-founder sustainability:

| Metric | Target |
|--------|--------|
| Monthly Active Users (free) | 1,000+ |
| Conversion rate | 2-5% |
| Paying customers | 20-50 |
| MRR | $400-$1,000 |
| Annual revenue | $5,000-$12,000 |

This is a sustainable side-project, not a venture-scale business.

### Cost Structure

| Item | Monthly Cost |
|------|--------------|
| Infrastructure | ~$0 (no backend) |
| Domain | ~$1 |
| Payment processing | ~3% of revenue |
| Support time | Variable |

**Target margin:** >90%

---

## Free vs Paid Feature Decisions

### Feature goes in Free if:
- It's core to the value proposition
- It has near-zero marginal cost
- It drives adoption
- Keeping it paid would cripple usability

### Feature goes in Paid if:
- It's an enhancement, not core functionality
- It requires ongoing development
- It targets professional/CI use cases
- It provides clear additional value

### Gray Area Resolution
When unsure, default to free for adoption, then reconsider based on usage data.

---

## Discounts and Promotions

### Allowed
- Annual discount (built into pricing)
- Open source maintainer discount (50% off)
- Student discount (50% off with verification)
- Launch discount (time-limited)

### Not Allowed
- Negotiated discounts
- Volume discounts (use Team tier instead)
- Perpetual discounts
- Discount codes for influencers

---

## Refund Policy

- 30-day refund, no questions asked
- Refund via original payment method
- License key revoked on refund
- No partial refunds

---

## License Key Security

### Generation
- RSA signature with private key
- Key contains: email, tier, expiration, signature
- Base64 encoded for easy copy/paste

### Storage
- Private key stored securely (not in repo)
- Public key embedded in binary
- Keys are not reversible

### Revocation
- No centralized revocation (no phone-home)
- Rely on expiration dates
- For abuse: don't renew, don't refund

---

## Pricing Changes

Pricing may change, with these constraints:

1. **Existing customers keep their rate** until they cancel
2. **30-day notice** before price increases
3. **No retroactive charges**
4. **Grandfathering** for early adopters if significant increase

---

## What We Will Never Do

1. **Sell user data** — We don't collect it
2. **Add ads** — Not compatible with developer tool
3. **Require online activation** — License works offline
4. **Lock data** — Your policies are yours, always exportable
5. **Paywall core features** — Generation stays free
