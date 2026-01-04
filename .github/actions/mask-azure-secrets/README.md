# Mask Azure Secrets Action

Comprehensive GitHub Actions composite action that masks Azure globally unique identifiers to prevent sensitive data leaks in public repositories.

## Problem

When running Terraform/OpenTofu plans in GitHub Actions, Azure resource names and IDs appear in logs and PR comments. Even though you mark variables as `sensitive = true`, **computed attributes** like DNS endpoints are not automatically masked.

### Example Leaks

```hcl
# This appears in plan output even with sensitive = true:
primary_blob_host = "stknccrsntplgasea01.blob.core.windows.net"
id = "/subscriptions/a998a203-5e24-4c8d-b2d2-536f14f7bf87/resourceGroups/..."
```

## Solution

This action masks:
- ✅ Tenant IDs
- ✅ Subscription IDs (from JSON mapping)
- ✅ Storage account names + all derived endpoints
- ✅ Container Registry endpoints (`*.azurecr.io`)
- ✅ Key Vault endpoints (`*.vault.azure.net`)
- ✅ App Service / Functions endpoints (`*.azurewebsites.net`)
- ✅ API Management endpoints (`*.azure-api.net`)
- ✅ Cosmos DB endpoints (`*.documents.azure.com`)
- ✅ Azure SQL endpoints (`*.database.windows.net`)
- ✅ Redis Cache endpoints (`*.redis.cache.windows.net`)
- ✅ Service Bus / Event Hubs endpoints (`*.servicebus.windows.net`)
- ✅ And 10+ other Azure service endpoints

## Usage

```yaml
steps:
  - name: Checkout repository
    uses: actions/checkout@v4

  - name: Mask Azure secrets
    uses: ./.github/actions/mask-azure-secrets
    with:
      tenant-id: ${{ secrets.ARM_TENANT_ID }}
      subscription-id-foundation: ${{ secrets.ARM_SUBSCRIPTION_ID_FOUNDATION }}
      tfstate-storage-account-name: ${{ secrets.ARM_TFSTATE_STORAGE_ACCOUNT_NAME }}
      subscription-id-mapping-json: ${{ secrets.ARM_SUBSCRIPTION_ID_MAPPING_JSON }}
      storage-account-name-mapping-json: ${{ secrets.ARM_STORAGE_ACCOUNT_NAME_MAPPING_JSON }}

  # Now run your OpenTofu/Terraform commands
  - name: OpenTofu Plan
    run: tofu plan
```

## Inputs

| Input | Description | Required |
|-------|-------------|----------|
| `tenant-id` | Azure Entra Tenant ID | Yes |
| `subscription-id-foundation` | Foundation subscription ID | Yes |
| `tfstate-storage-account-name` | Terraform state storage account name | Yes |
| `subscription-id-mapping-json` | JSON object mapping subscription names to IDs | Yes |
| `storage-account-name-mapping-json` | JSON object mapping logical names to Azure resource names | Yes |

## Input Format

### `subscription-id-mapping-json`

```json
{
  "management": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "identity": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "connectivity_prod": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

### `storage-account-name-mapping-json`

This JSON mapping can contain **any Azure resource with a globally unique DNS name**, not just storage accounts:

```json
{
  "st_platform_logs_archive_sea_01": "stknccrsntplgasea01",
  "cr_platform_prod": "crknccrsntprod01",
  "kv_platform_prod": "kv-knccrsnt-prod-01",
  "app_kyber_prod": "app-kyber-prod"
}
```

The action will automatically mask all common Azure DNS patterns for each name.

## How It Works

1. **Parses JSON secrets** using `jq` to extract individual values
2. **Masks each value** using GitHub Actions `::add-mask::` command
3. **Generates derived endpoints** for each resource name and masks them too
4. **Works proactively** - masks before they appear in logs

## Supported Azure Services

The action masks endpoints for these globally unique Azure services:

- Storage Accounts (blob, dfs, file, queue, table, web)
- Container Registry
- Key Vault
- App Service / Azure Functions
- API Management
- Cosmos DB
- Azure SQL Server
- Redis Cache
- Service Bus / Event Hubs
- SignalR Service
- Static Web Apps
- Azure Front Door
- Azure CDN
- Cognitive Services
- Azure Machine Learning

## Adding New Resources

To add support for new Azure resources, edit [`mask-secrets.sh`](./mask-secrets.sh) and add the DNS pattern:

```bash
# In the mask_azure_resource function
echo "::add-mask::${name}.YOUR-SERVICE.azure.net"
echo "::add-mask::https://${name}.YOUR-SERVICE.azure.net/"
```

Then add the resource name to your `ARM_STORAGE_ACCOUNT_NAME_MAPPING_JSON` secret.

## Security Best Practices

1. **Use this action FIRST** in your job, before any Azure authentication or Terraform commands
2. **Keep JSON mappings updated** whenever you create new Azure resources
3. **Verify masking works** by checking workflow logs - sensitive values should show as `***`
4. **Don't commit actual values** - all identifiers should come from GitHub secrets

## Example Workflow

```yaml
jobs:
  plan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      # CRITICAL: Mask secrets BEFORE any Azure operations
      - uses: ./.github/actions/mask-azure-secrets
        with:
          tenant-id: ${{ secrets.ARM_TENANT_ID }}
          subscription-id-foundation: ${{ secrets.ARM_SUBSCRIPTION_ID_FOUNDATION }}
          tfstate-storage-account-name: ${{ secrets.ARM_TFSTATE_STORAGE_ACCOUNT_NAME }}
          subscription-id-mapping-json: ${{ secrets.ARM_SUBSCRIPTION_ID_MAPPING_JSON }}
          storage-account-name-mapping-json: ${{ secrets.ARM_STORAGE_ACCOUNT_NAME_MAPPING_JSON }}

      - uses: azure/login@v2
        with:
          client-id: ${{ secrets.ARM_CLIENT_ID }}
          tenant-id: ${{ secrets.ARM_TENANT_ID }}
          subscription-id: ${{ secrets.ARM_SUBSCRIPTION_ID_FOUNDATION }}

      - name: OpenTofu Plan
        run: |
          tofu init
          tofu plan -no-color
```

## License

Part of the Coruscant project. See root LICENSE file.
