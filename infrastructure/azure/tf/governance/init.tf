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
    container_name      = "containertfstate"
    key                 = "governance/terraform.tfstate"
    use_oidc            = true # Use OIDC authentication (no access keys needed)
    use_azuread_auth    = true # Use OIDC authentication (no access keys needed)
    # storage_account_name provided via -backend-config (from GitHub secret or backend.hcl)
  }
}

# Azure Resource Manager Provider Configuration
provider "azurerm" {
  resource_provider_registrations = "none" # Disable auto-registration of resource providers
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true # Safety: Prevent accidental RG deletion
    }
    key_vault {
      purge_soft_delete_on_destroy    = false # Safety: Keep soft-delete enabled
      recover_soft_deleted_key_vaults = true  # Recover instead of error on conflicts
    }
  }

  # When managing Management Groups and moving subscriptions:
  # - Set ARM_SUBSCRIPTION_ID environment variable (for auth context)
  # - Do NOT hardcode subscription_id here (allows operating on all subscriptions)
  # The provider uses ARM_SUBSCRIPTION_ID for authentication but can manage any subscription
  # where the service principal has permissions
}

# # Provider Alias: Security Subscription
# provider "azurerm" {
#   alias = "security_prod"
#   features {
#     resource_group {
#       prevent_deletion_if_contains_resources = true
#     }
#     key_vault {
#       purge_soft_delete_on_destroy    = false
#       recover_soft_deleted_key_vaults = true
#     }
#   }

#   subscription_id = var.subscription_id_security_prod
# }

# # Provider Alias: Connectivity Production Subscription
# provider "azurerm" {
#   alias = "connectivity_prod"
#   features {
#     resource_group {
#       prevent_deletion_if_contains_resources = true
#     }
#   }

#   subscription_id = var.subscription_id_connectivity_prod
# }

# # Provider Alias: Connectivity Development Subscription
# provider "azurerm" {
#   alias = "connectivity_dev"
#   features {
#     resource_group {
#       prevent_deletion_if_contains_resources = true
#     }
#   }

#   subscription_id = var.subscription_id_connectivity_dev
# }

# # Provider Alias: Identity Subscription
# provider "azurerm" {
#   alias = "identity"
#   features {
#     resource_group {
#       prevent_deletion_if_contains_resources = true
#     }
#   }

#   subscription_id = var.subscription_id_identity
# }

# Provider: Random (for unique resource naming)
provider "random" {}
