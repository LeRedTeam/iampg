# Non-Goals

This document explicitly defines what IAM Auto-Policy Generator will NOT become.

These are not "maybe later" items. These are hard boundaries that protect the product's focus and sustainability.

---

## This Product Is NOT

### A Cloud Governance Platform
- No policy approval workflows
- No role request systems
- No permission lifecycle management
- No organizational policy inheritance
- No compliance framework mapping

**Why:** Governance requires organizational context, human workflows, and ongoing maintenance. This is a different product category.

### A Centralized IAM Management System
- No central policy repository
- No cross-account policy synchronization
- No IAM role inventory
- No permission analytics dashboard
- No access graph visualization

**Why:** Centralization requires infrastructure, state management, and multi-tenant complexity. Contradicts local-first architecture.

### A Compliance Dashboard
- No SOC2 mapping
- No PCI-DSS controls
- No HIPAA requirement tracking
- No audit report generation
- No compliance scoring

**Why:** Compliance is a different problem domain requiring specialized knowledge and ongoing regulatory tracking.

### A Multi-Cloud Security Suite
- No Azure support
- No GCP support
- No Kubernetes RBAC
- No cloud-agnostic abstractions

**Why:** Multi-cloud multiplies complexity without multiplying value for our target user. AWS-only focus enables depth over breadth.

### A Role Lifecycle Management System
- No role creation workflows
- No role decommissioning
- No permission rotation
- No temporary access grants
- No access reviews

**Why:** Lifecycle management requires state, scheduling, and organizational integration. Different product category.

### A SaaS Platform
- No user accounts
- No team management
- No usage analytics collection
- No cloud-hosted processing
- No data retention

**Why:** SaaS introduces infrastructure cost, compliance burden, and support complexity. Local-first is simpler and more trustworthy.

### A Consulting Tool
- No security assessments
- No architecture reviews
- No best-practice audits
- No remediation playbooks

**Why:** Consulting requires human expertise and context. This is automation, not advisory.

### A Real-Time Monitoring System
- No continuous policy drift detection
- No live permission usage tracking
- No alerting system
- No anomaly detection

**Why:** Real-time monitoring requires always-on infrastructure and ongoing operational cost.

---

## Explicit Feature Exclusions

These specific features will NOT be implemented:

| Feature | Reason |
|---------|--------|
| Web UI | Adds frontend maintenance, hosting cost |
| User accounts | Adds auth complexity, data storage |
| Team management | Adds multi-tenant complexity |
| Policy storage | Adds database dependency |
| Audit logging | Adds compliance burden |
| Email notifications | Adds delivery infrastructure |
| Slack/Teams integration | Adds third-party dependencies |
| Custom policy templates | Adds template engine complexity |
| Policy approval workflows | Adds workflow state management |
| Scheduled scans | Adds scheduler infrastructure |
| Cross-account discovery | Adds assume-role complexity |
| SSO integration | Adds auth provider complexity |

---

## Responses to Common Requests

**"Can you add a dashboard?"**
No. Use your existing AWS console or third-party tools. We generate JSON; visualization is out of scope.

**"Can you support Azure/GCP?"**
No. Multi-cloud multiplies maintenance burden without clear revenue justification. AWS-only focus enables better depth.

**"Can you store policies centrally?"**
No. Use git, S3, or your existing artifact storage. We don't store anything.

**"Can you integrate with our SIEM?"**
No. Export to SARIF/JSON and integrate yourself. We don't maintain integration code.

**"Can you add approval workflows?"**
No. Use your existing change management process. We generate policies; approval is your organizational concern.

**"Can you detect overly-permissive existing policies?"**
Limited. The `refine` command (paid) analyzes policies you provide. We don't scan your AWS account.

**"Can you run as a service in our VPC?"**
No. Run the CLI directly. No hosted version exists or will exist.

---

## Why These Boundaries Matter

1. **Focus enables quality**: Doing one thing well beats doing many things poorly.

2. **Simplicity enables sustainability**: A solo founder cannot maintain a platform.

3. **Constraints enable trust**: Users can trust what we don't do (store credentials, phone home).

4. **Boundaries enable speed**: Saying no to scope creep enables shipping.

5. **Limits enable pricing**: Clear scope enables clear value proposition.

---

## Revisiting Non-Goals

These boundaries are not permanent, but changing them requires:

1. Explicit founder approval
2. Clear revenue justification
3. Maintenance burden assessment
4. Support impact analysis
5. Documentation update across all `.claude/` files

The default answer to scope expansion is NO.
