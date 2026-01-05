provider "azurerm" {
  subscription_id = var.subscription_id_foundation

  resource_provider_registrations = "none" # Disable auto-registration of resource providers
  use_oidc                        = true   # Use OIDC authentication (no access keys needed)

  features {
    resource_group {
      prevent_deletion_if_contains_resources = true # Safety: Prevent accidental RG deletion
    }
    key_vault {
      purge_soft_delete_on_destroy    = false # Safety: Keep soft-delete enabled
      recover_soft_deleted_key_vaults = true  # Recover instead of error on conflicts
    }
  }
}
