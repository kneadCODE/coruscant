# Azure Infrastructure

OpenTofu infrastructure for Coruscant using Azure with OIDC authentication and multi-subscription support.

## Quick Start

### Local Development

1. **Login & configure workspace:**

   ```bash
   az login
   cd infrastructure/azure/tf/governance

   cp backend.hcl.example backend.hcl  # Edit with your storage account name

   cat > terraform.tfvars.local <<EOF
   subscription_id_management        = "YOUR-MGMT-SUB-ID"
   subscription_id_security          = "YOUR-SEC-SUB-ID"
   subscription_id_connectivity_prod = "YOUR-CONN-PROD-SUB-ID"
   subscription_id_connectivity_dev  = "YOUR-CONN-DEV-SUB-ID"
   subscription_id_identity          = "YOUR-IDENTITY-SUB-ID"
   EOF
   ```

2. **Run OpenTofu:**

   ```bash
   tofu init -backend-config=backend.hcl
   tofu plan
   tofu apply
   ```

## GitHub Actions Setup

### Required Secrets

**Environment Secrets (Common)**:

- `ARM_TENANT_ID` - Azure tenant ID
- `ARM_SUBSCRIPTION_ID_FOUNDATION` - Any subscription with Reader access
- `ARM_TFSTATE_STORAGE_ACCOUNT_NAME` - State storage account name
- `ARM_CLIENT_ID_ADMIN` - Service principal app ID (Contributor role)
- `ARM_SUBSCRIPTION_ID` - Default subscription
- `ARM_SUBSCRIPTION_ID_MANAGEMENT` - Management subscription
- `ARM_SUBSCRIPTION_ID_SECURITY` - Security subscription
- `ARM_SUBSCRIPTION_ID_CONNECTIVITY_PROD` - Connectivity prod subscription
- `ARM_SUBSCRIPTION_ID_CONNECTIVITY_DEV` - Connectivity dev subscription
- `ARM_SUBSCRIPTION_ID_IDENTITY` - Identity subscription

**Environment Secrets (`azure-tf-plan` environment, PRs only)**

- `ARM_CLIENT_ID_READONLY` - Service principal app ID (Reader role)

**Environment Secrets (`azure-tf-apply` environment, main branch only)**

- `ARM_CLIENT_ID_ADMIN` - Service principal app ID (Admin role)

### Environment Configuration

Create `azure-infra` environment in GitHub:

1. Settings → Environments → New environment: `azure-infra`
2. Deployment branches → Selected branches and tags → `main`
3. Add environment secrets above

### Azure App Setup

Use **Application (client) ID** (not Object ID) for `ARM_CLIENT_ID_*` secrets.

**Grant RBAC roles:**

```bash
# Read-only app: Reader on all subscriptions
for SUB in <mgmt> <security> <conn-prod> <conn-dev> <identity>; do
  az role assignment create --assignee <READONLY_APP_ID> --role "Reader" --scope "/subscriptions/$SUB"
done

# Admin app: Contributor on all subscriptions
for SUB in <mgmt> <security> <conn-prod> <conn-dev> <identity>; do
  az role assignment create --assignee <ADMIN_APP_ID> --role "Contributor" --scope "/subscriptions/$SUB"
done
```

**Configure OIDC federated credentials:**

```bash
# Read-only (all branches can plan)
az ad app federated-credential create --id <READONLY_APP_ID> --parameters '{
  "name": "github-coruscant-readonly",
  "issuer": "https://token.actions.githubusercontent.com",
  "subject": "repo:YOUR_ORG/coruscant:ref:refs/heads/main",
  "audiences": ["api://AzureADTokenExchange"]
}'

# Admin (only azure-infra environment can apply)
az ad app federated-credential create --id <ADMIN_APP_ID> --parameters '{
  "name": "github-coruscant-admin",
  "issuer": "https://token.actions.githubusercontent.com",
  "subject": "repo:YOUR_ORG/coruscant:environment:azure-infra",
  "audiences": ["api://AzureADTokenExchange"]
}'
```

## Architecture

### Multi-Subscription Pattern

Workspaces use provider aliases to operate across subscriptions:

```hcl
provider "azurerm" {
  subscription_id = var.subscription_id_management  # Default
}

provider "azurerm" {
  alias           = "security"
  subscription_id = var.subscription_id_security
}
```

### Security Model

- **Feature branches**: Read-only credentials (can plan, can't apply)
- **Main branch**: Admin credentials (can plan and apply)
- **Environment restriction**: `azure-infra` environment locked to main branch

### Workflows

- **PR**: Security scan + plan (read-only)
- **Merge to main**: Auto-apply (admin)
- **Manual**: workflow_dispatch (main branch only)
- **Drift detection**: Mondays 8 AM SGT

## GitOps Workflow

1. Create feature branch
2. Modify `.tf` files
3. Push → Plan runs (read-only)
4. Create PR → Review plan
5. Merge → Apply runs automatically

**Deletions**: Remove from `.tf` files, never use `tofu destroy`.

## Troubleshooting

**Subscription not found**: Grant RBAC access
**Authentication failed**: Verify OIDC federated credentials
**No value for variable**: Create `terraform.tfvars.local` with subscription IDs
