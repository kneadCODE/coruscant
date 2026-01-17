# Azure Infrastructure as Code (OpenTofu)

This directory contains the Infrastructure as Code (IaC) definitions for the Coruscant project's Azure infrastructure using OpenTofu.

## Table of Contents

- [Overview](#overview)
- [Directory Structure](#directory-structure)
- [State Management](#state-management)
- [Service Principals & Permissions](#service-principals--permissions)
- [Workflows](#workflows)
- [Security Model](#security-model)
- [Development Guidelines](#development-guidelines)

## Overview

The infrastructure is organized into three main categories based on lifecycle, change frequency, and security requirements:

1. **Foundation**: One-time bootstrap and occasional subscription vending (high privilege)
2. **Platform**: Management subscriptions infrastructure (regular changes, scoped permissions)
3. **Landing Zone**: Application landing zones (regular changes, scoped permissions)

### Why OpenTofu?

- True open-source (MPL 2.0 license vs Terraform's restrictive BSL)
- Native state encryption
- Linux Foundation backed, community-driven
- 100% compatible drop-in replacement for Terraform
- No vendor lock-in

## Directory Structure

```
azure/tf/
├── foundation/
│   ├── bootstrap/                      # ⚠️ COMPLETE: One-time setup (workflow disabled)
│   └── subscription-vending/           # Occasional: Onboard and configure new Azure subscriptions
├── governance/                         # Policy assignments and governance controls (high privilege)
├── platform/                           # Platform subscriptions
└── landingzone/                        # Landing zones subscriptions (future)
```

### Workspace Organization

Each subdirectory with a `backend.tf` file is an **independent workspace** with its own:

- Terraform state file
- Backend configuration
- Lifecycle and change frequency
- Security permissions and service principal

**Examples**:

- `foundation/bootstrap` → Single workspace (⚠️ Bootstrap complete, workflow disabled)
- `foundation/subscription-vending` → Single workspace
- `governance` → Single workspace
- `platform/management/logs` → Single workspace
- `platform/management/backup` → Single workspace

## State Management

### Storage Account Architecture

All Terraform state files are stored in a single Foundation Storage Account (name not disclosed for security) with five isolated containers:

| Container Name              | Purpose                          | Service Principal Access                          |
|-----------------------------|----------------------------------|---------------------------------------------------|
| `bootstraptfstate`          | Bootstrap workspace state        | `sp-gha-tf-apply-bootstrap` (RW) - ⚠️ SP deleted  |
| `subscriptionvendingtfstate`| Subscription vending state       | `sp-gha-tf-apply-subscriptionvending` (RW)        |
| `governancetfstate`         | Governance workspace state       | `sp-gha-tf-apply-governance` (RW)                 |
| `platformtfstate`           | All platform workspaces          | `sp-gha-tf-apply-platform` (RW)                   |
| `landingzonetfstate`        | All landing zone workspaces      | `sp-gha-tf-apply-landingzone` (RW)                |

**Global Plan Access**: `sp-gha-tf-plan-global` has Storage Blob Data Contributor on all 5 containers (needs write access for state locking during plan operations).

### Backend Configuration

Each workspace's `backend.tf` specifies:

```hcl
backend "azurerm" {
  resource_group_name = "rg-tfstate-bootstrap-coruscant-sea"
  container_name      = "platformtfstate"              # Container based on workspace category
  key                 = "management/logs.tfstate"       # Unique path within container
  use_oidc            = true
  use_azuread_auth    = true
  # storage_account_name provided via -backend-config at runtime
}
```

## Service Principals & Permissions

The project uses **6 dedicated service principals** following the principle of least privilege (bootstrap SP deleted after completion).

> **Note**: Service Principals are created manually via Azure Portal App Registrations. Federated Credentials (for OIDC) are also configured manually in the Azure Portal. Subscriptions are created manually through Azure Portal or Azure CLI.

| Service Principal | Purpose | Azure Permissions | Storage Access | GitHub Secret | Federated Credential | Usage |
|-------------------|---------|-------------------|----------------|---------------|---------------------|-------|
| `sp-gha-tf-plan-global` | Read-only plan operations (PRs only) | Reader at `mg-coruscant-root` | Storage Blob Data Contributor on all 5 containers (for state locking) | `ARM_CLIENT_ID_SP_GHA_TF_PLAN_GLOBAL` | Repo: `kneadCODE/coruscant`<br>Env: `azure-tf-plan`<br>Branch: All | PR plan operations across all workflows |
| `sp-gha-tf-apply-bootstrap` | One-time bootstrap apply | ~~Management Group Contributor at `mg-coruscant-root`~~ | ~~Storage Blob Data Contributor on `bootstraptfstate`~~ | ~~`ARM_CLIENT_ID_SP_GHA_TF_APPLY_BOOTSTRAP`~~ | ⚠️ **DELETED** - Bootstrap complete | ⚠️ **DSIABLED** - Bootstrap complete, workflow disabled |
| `sp-gha-tf-apply-subscriptionvending` | Create/configure Azure subscriptions | Management Group Contributor, User Access Administrator, Contributor at `mg-coruscant-root` | Storage Blob Data Contributor on `subscriptionvendingtfstate` | `ARM_CLIENT_ID_SP_GHA_TF_APPLY_SUBSCRIPTIONVENDING` | Repo: `kneadCODE/coruscant`<br>Workflow: `opentofu-subscription-vending.yml`<br>Env: `azure-tf-apply-subscriptionvending`<br>Branch: `main` | Manual dispatch only for subscription vending |
| `sp-gha-tf-apply-governance` | Apply governance policies | Resource Policy Contributor & Cost Management Contributor at `mg-coruscant-root` | Storage Blob Data Contributor on `governancetfstate` | `ARM_CLIENT_ID_SP_GHA_TF_APPLY_GOVERNANCE` | Repo: `kneadCODE/coruscant`<br>Workflow: `opentofu.yml`<br>Env: `azure-tf-apply-governance`<br>Branch: `main` | Auto-apply after merge for `governance` workspace |
| `sp-gha-tf-apply-platform` | Apply platform infrastructure | Contributor scoped to platform subscriptions | Storage Blob Data Contributor on `platformtfstate` | `ARM_CLIENT_ID_SP_GHA_TF_APPLY_PLATFORM` | Repo: `kneadCODE/coruscant`<br>Workflow: `opentofu.yml`<br>Env: `azure-tf-apply-platform`<br>Branch: `main` | Auto-apply after merge for `platform/**` workspaces |
| `sp-gha-tf-apply-landingzone` | Apply landing zone infrastructure | Contributor scoped to landing zone subscriptions | Storage Blob Data Contributor on `landingzonetfstate` | `ARM_CLIENT_ID_SP_GHA_TF_APPLY_LANDINGZONE` | Repo: `kneadCODE/coruscant`<br>Workflow: `opentofu.yml`<br>Env: `azure-tf-apply-landingzone`<br>Branch: `main` | Auto-apply after merge for `landingzone/**` workspaces |

**All environments require manual approval before job execution. All apply environments are restricted to `main` branch only.**

## Workflows

| Workflow | Trigger | Workspaces | Service Principal (PR) | Service Principal (Apply) | Apply Behavior | Notes |
|----------|---------|------------|------------------------|---------------------------|----------------|-------|
| ~~`opentofu-bootstrap.yml`~~ | ~~PR to `main`<br>Manual dispatch on `main`~~ | ~~`foundation/bootstrap`~~ | N/A | N/A | N/A | ⚠️ **DISABLED** - Bootstrap complete, workflow disabled |
| `opentofu-subscription-vending.yml` | PR to `main`<br>Manual dispatch on `main` | `foundation/subscription-vending` | `sp-gha-tf-plan-global`<br>(env: `azure-tf-plan`) | `sp-gha-tf-apply-subscriptionvending`<br>(env: `azure-tf-apply-subscriptionvending`) | **Manual dispatch only**<br>After PR merged, manually trigger workflow | High privilege operations |
| `opentofu.yml` | PR to `main`<br>Push to `main` | `governance`<br>`platform/**`<br>`landingzone/**` | `sp-gha-tf-plan-global`<br>(env: `azure-tf-plan`) | **Dynamic selection:**<br>`governance` → `sp-gha-tf-apply-governance`<br>(env: `azure-tf-apply-governance`)<br>`platform/**` → `sp-gha-tf-apply-platform`<br>(env: `azure-tf-apply-platform`)<br>`landingzone/**` → `sp-gha-tf-apply-landingzone`<br>(env: `azure-tf-apply-landingzone`) | **Automatic after merge**<br>Environment approval required | PRs can only modify one workspace at a time (enforced) |
| `opentofu-scheduled.yml` | Weekly Mon 8AM SGT<br>(0:00 UTC Mon) | `governance`<br>`platform/**`<br>`landingzone/**` | `sp-gha-tf-plan-global`<br>(no environment) | N/A (plan only) | N/A (drift detection only) | Non-blocking<br>Creates GitHub issues when drift detected |

## Security Model

### Multi-Layered Security

1. **Branch Protection**: All changes must go through PR review on `main` branch
2. **Environment Protection**: All environments require manual approval before job execution
3. **Service Principal Isolation**: Each category uses dedicated SPs with scoped permissions
4. **Federated Credentials**: OIDC-based authentication (no long-lived secrets)
5. **State Locking**: Prevents concurrent modifications via Azure Blob lease mechanism
6. **Plan Artifacts**: What's reviewed in PR is exactly what gets applied (plan artifact reuse)
7. **CODEOWNERS**: Directory ownership enforced via GitHub (team/individual approvals required)

### GitOps Approval Model

```
PR Created → Plan Executes (read-only SP) → Plan Comment Posted
     ↓
CODEOWNER Review & Approval → Infrastructure Approval
     ↓
Merge to Main → Apply Job Triggers → Environment Approval Required
     ↓
Manual Approval → Apply Executes (admin SP) → Infrastructure Updated
```

### High-Privilege Operations Protection

Bootstrap and subscription vending use **manual workflow dispatch** instead of automatic apply:

```
PR Created → Plan Executes → CODEOWNER Review → Merge to Main
     ↓
Manual Workflow Dispatch on Main → Environment Approval → Apply
```

This adds an extra layer of control for operations with high-privilege SPs.

## Development Guidelines

### Making Infrastructure Changes

#### Platform or Landing Zone Changes

1. **Create a feature branch**:

   ```bash
   git checkout -b infra/update-platform-logging
   ```

2. **Modify infrastructure** (only one workspace per PR):

   ```bash
   cd azure/tf/platform/management/logs
   # Edit .tf files
   ```

3. **Format and validate locally**:

   ```bash
   tofu fmt -recursive
   tofu init -backend=false
   tofu validate
   ```

4. **Create pull request**:
   - Workflow runs `tofu plan` using read-only SP
   - Plan output posted as PR comment
   - CODEOWNER review required

5. **Review and merge**:
   - CODEOWNER approves infrastructure changes
   - Merge to `main`

6. **Automatic apply**:
   - Workflow detects merged changes
   - Requests environment approval (manual gate)
   - Applies changes using scoped admin SP

#### Subscription Vending Changes

1. **Follow steps 1-5 above** for `foundation/subscription-vending`

2. **Manual trigger apply**:

   ```bash
   # After PR merged to main
   # Go to Actions → OpenTofu (Subscription Vending) → Run workflow
   # Select branch: main → Run workflow
   # Approve environment when prompted
   ```

#### Bootstrap Changes

**WARNING**: Bootstrap should rarely/never change after initial setup. If needed:

1. **Create feature branch and modify** `foundation/bootstrap`
2. **Format and validate locally**
3. **Create pull request**:
   - Workflow runs `tofu plan` using plan SP
   - Plan output posted as PR comment
   - CODEOWNER review required
4. **Review and merge to `main`**
5. **Manual trigger apply**:
   - Go to Actions → OpenTofu (Bootstrap) → Run workflow
   - Select branch: `main` → Run workflow
   - Approve environment when prompted
6. **Consider deleting the bootstrap SP after completion** (high privilege, temporary use only)

### Workspace Rules

1. **One workspace per PR**: PRs modifying multiple workspaces will fail (enforced by workflow)
2. **Format before commit**: Run `tofu fmt -recursive` before committing
3. **No backend.hcl commits**: This file is gitignored (contains environment-specific config)
4. **State file security**: Never commit `.tfstate` files (gitignored)

### CODEOWNERS

Infrastructure directories are protected by CODEOWNERS:

| Directory | Owner Team | Notes |
|-----------|------------|-------|
| `azure/tf/foundation/bootstrap/` | `@kneadCODE/infra-foundation-team` | ⚠️ Bootstrap complete, archived |
| `azure/tf/foundation/subscription-vending/` | `@kneadCODE/infra-foundation-team` | High privilege, occasional |
| `azure/tf/governance/` | `@kneadCODE/infra-governance-team` | Policy assignments & budgets, higher privilege but scoped access |
| `azure/tf/platform/**` | `@kneadCODE/infra-platform-team` | Regular changes, scoped access |
| `azure/tf/landingzone/**` | `@kneadCODE/infra-landingzone-team` | Regular changes, scoped access |

**Current**: All use `@<team>` as placeholder. Replace with team handles when teams are established.

## References

- [OpenTofu Documentation](https://opentofu.org/docs/)
- [Azure Provider Documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs)
- [Azure Landing Zones](https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/landing-zone/)
- [GitHub Actions OIDC with Azure](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-azure)
