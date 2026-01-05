terraform {
  required_version = ">= 1.6.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0" # Use latest 4.x version, allows minor/patch updates
    }
  }

  # Backend configuration: Azure Storage Account
  # State stored in Azure Blob Storage with automatic locking via blob leases
  # Authentication via OIDC (ARM_* environment variables set by GitHub Actions)
  #
  # Storage account name provided via backend config:
  # - GitHub Actions: -backend-config="storage_account_name=${{ secrets.ARM_STORAGE_ACCOUNT_NAME }}"
  # - Local: Create backend.hcl (gitignored) with: storage_account_name = "your-storage-account"
  backend "azurerm" {
    resource_group_name = "rg-tfstate-bootstrap-coruscant-sea"
    container_name      = "bootstraptfstate"
    key                 = "bootstrap/terraform.tfstate"
    use_oidc            = true # Use OIDC authentication (no access keys needed)
    use_azuread_auth    = true # Use OIDC authentication (no access keys needed)
    # storage_account_name provided via -backend-config (from GitHub secret or backend.hcl)
  }
}
