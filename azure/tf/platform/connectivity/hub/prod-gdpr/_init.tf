terraform {
  backend "azurerm" {
    resource_group_name = "rg-tfstate-bootstrap-coruscant-sea"
    container_name      = "platformtfstate"
    key                 = "connectivity/prod/hub/gdpr.tfstate"
    use_oidc            = true # Use OIDC authentication (no access keys needed)
    use_azuread_auth    = true # Use OIDC authentication (no access keys needed)
    # storage_account_name provided via -backend-config (from GitHub secret or local.backend.hcl)
  }
}

provider "azurerm" {
  subscription_id = local.subscription_ids["connectivity_${local.envs.prod.name}"]

  resource_provider_registrations = "none" # Disable auto-registration of resource providers
  use_oidc                        = true   # Use OIDC authentication (no access keys needed)
  storage_use_azuread             = true   # Use OIDC authentication (no access keys needed)

  features {
    resource_group {
      prevent_deletion_if_contains_resources = true # Safety: Prevent accidental RG deletion
    }
    key_vault {
      purge_soft_delete_on_destroy    = false # Safety: Keep soft-delete enabled
      recover_soft_deleted_key_vaults = true  # Recover instead of error on conflicts
    }
    storage {
      data_plane_available = false
    }
  }
}
