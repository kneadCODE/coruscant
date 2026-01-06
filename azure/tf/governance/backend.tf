terraform {
  # Backend configuration: Azure Storage Account
  # State stored in Azure Blob Storage with automatic locking via blob leases
  # Authentication via OIDC (ARM_* environment variables set by GitHub Actions)
  backend "azurerm" {
    resource_group_name = "rg-tfstate-bootstrap-coruscant-sea"
    container_name      = "governancetfstate"
    key                 = "governance/terraform.tfstate"
    use_oidc            = true
    use_azuread_auth    = true
    # storage_account_name provided via -backend-config (from GitHub secret or backend.hcl)
  }
}
