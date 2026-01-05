# =============================================================================
# Azure Provider Configuration
# =============================================================================

# Default provider - uses foundation subscription for state/control plane operations
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

  subscription_id = var.subscription_id_foundation
}

# =============================================================================
# Provider Aliases - One per subscription for provider registration
# =============================================================================

provider "azurerm" {
  alias                           = "management"
  subscription_id                 = local.subscription_ids["management"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.OperationalInsights", # Log Analytics Workspace for platform logs
    "Microsoft.Storage",             # Storage accounts for logs
    "Microsoft.RecoveryServices",    # Recovery Services Vault and Backup Vault
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "identity"
  subscription_id                 = local.subscription_ids["identity"]
  resource_provider_registrations = "none"
  resource_providers_to_register  = local.base_resource_providers
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "connectivity_prod"
  subscription_id                 = local.subscription_ids["connectivity_prod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Network", # VNets, NSGs, Route Tables, Firewalls, DNS, Gateways, Bastion
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "connectivity_nonprod"
  subscription_id                 = local.subscription_ids["connectivity_nonprod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Network", # VNets, NSGs, Route Tables, Firewalls, DNS, Gateways, Bastion
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "security_prod"
  subscription_id                 = local.subscription_ids["security_prod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.OperationalInsights", # Log Analytics Workspace with Sentinel
    "Microsoft.Storage",             # Storage accounts for security logs
    "Microsoft.Compute",             # VMs for HashiCorp Vault
    "Microsoft.Network",             # Spoke VNet, NSGs, Route Tables, DDoS Protection Plan
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "security_nonprod"
  subscription_id                 = local.subscription_ids["security_nonprod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.OperationalInsights", # Log Analytics Workspace with Sentinel
    "Microsoft.Storage",             # Storage accounts for security logs
    "Microsoft.Compute",             # VMs for HashiCorp Vault
    "Microsoft.Network",             # Spoke VNet, NSGs, Route Tables, DDoS Protection Plan
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "devops_prod"
  subscription_id                 = local.subscription_ids["devops_prod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Compute",           # VMSS for AKS self-hosted GitHub Actions runners
    "Microsoft.ContainerService",  # AKS for self-hosted GitHub Actions runners
    "Microsoft.ContainerRegistry", # Azure Container Registries
    "Microsoft.KeyVault",          # Key Vault for GHA
    "Microsoft.Network",           # Spoke VNet, NSGs, Route Tables
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "devops_nonprod"
  subscription_id                 = local.subscription_ids["devops_nonprod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Compute",           # VMSS for AKS self-hosted GitHub Actions runners
    "Microsoft.ContainerService",  # AKS for self-hosted GitHub Actions runners
    "Microsoft.ContainerRegistry", # Azure Container Registries
    "Microsoft.KeyVault",          # Key Vault for GHA
    "Microsoft.Network",           # Spoke VNet, NSGs, Route Tables
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "esb_prod"
  subscription_id                 = local.subscription_ids["esb_prod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Compute", # VMs for self-hosted Kafka
    "Microsoft.Network", # Spoke VNet, NSGs, Route Tables, ILB
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "esb_nonprod"
  subscription_id                 = local.subscription_ids["esb_nonprod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Compute", # VMs for self-hosted Kafka
    "Microsoft.Network", # Spoke VNet, NSGs, Route Tables, ILB
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "observability_prod"
  subscription_id                 = local.subscription_ids["observability_prod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Compute",          # VMSS for AKS
    "Microsoft.ContainerService", # AKS for hosting OTEL collectors
    "Microsoft.Network",          # Spoke VNet, NSGs, Route Tables
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "observability_nonprod"
  subscription_id                 = local.subscription_ids["observability_nonprod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Compute",          # VMSS for AKS
    "Microsoft.ContainerService", # AKS for hosting OTEL collectors
    "Microsoft.Network",          # Spoke VNet, NSGs, Route Tables
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "edge_prod"
  subscription_id                 = local.subscription_ids["edge_prod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Compute",          # VMSS for AKS
    "Microsoft.Cdn",              # Azure Front Door
    "Microsoft.ContainerService", # AKS for hosting API Gw
    "Microsoft.Devices",          # IoT Hub for MQTT/AMQP ingress
    "Microsoft.EventHub",         # Event Hubs for Kafka ingress
    "Microsoft.Storage",          # Storage accounts for SFTP ingress
    "Microsoft.Network",          # Spoke VNet, NSGs, Route Tables
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "edge_nonprod"
  subscription_id                 = local.subscription_ids["edge_nonprod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Compute",          # VMSS for AKS
    "Microsoft.Cdn",              # Azure Front Door
    "Microsoft.ContainerService", # AKS for hosting API Gw
    "Microsoft.Devices",          # IoT Hub for MQTT/AMQP ingress
    "Microsoft.EventHub",         # Event Hubs for Kafka ingress
    "Microsoft.Storage",          # Storage accounts for SFTP ingress
    "Microsoft.Network",          # Spoke VNet, NSGs, Route Tables
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "iam_prod"
  subscription_id                 = local.subscription_ids["iam_prod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Cache",            # Redis
    "Microsoft.Compute",          # VMSS for AKS
    "Microsoft.ContainerService", # AKS
    "Microsoft.DBforPostgreSQL",  # PG Flexible servers
    "Microsoft.Storage",          # Storage accounts
    "Microsoft.Network",          # Spoke VNet, NSGs, Route Tables
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "iam_nonprod"
  subscription_id                 = local.subscription_ids["iam_nonprod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Cache",            # Redis
    "Microsoft.Compute",          # VMSS for AKS
    "Microsoft.ContainerService", # AKS
    "Microsoft.DBforPostgreSQL",  # PG Flexible servers
    "Microsoft.Storage",          # Storage accounts
    "Microsoft.Network",          # Spoke VNet, NSGs, Route Tables
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "payment_prod"
  subscription_id                 = local.subscription_ids["payment_prod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Cache",            # Redis
    "Microsoft.Compute",          # VMSS for AKS
    "Microsoft.ContainerService", # AKS
    "Microsoft.DBforPostgreSQL",  # PG Flexible servers
    "Microsoft.Storage",          # Storage accounts
    "Microsoft.Network",          # Spoke VNet, NSGs, Route Tables
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

provider "azurerm" {
  alias                           = "payment_nonprod"
  subscription_id                 = local.subscription_ids["payment_nonprod"]
  resource_provider_registrations = "none"
  resource_providers_to_register = toset(concat(local.base_resource_providers, [
    "Microsoft.Cache",            # Redis
    "Microsoft.Compute",          # VMSS for AKS
    "Microsoft.ContainerService", # AKS
    "Microsoft.DBforPostgreSQL",  # PG Flexible servers
    "Microsoft.Storage",          # Storage accounts
    "Microsoft.Network",          # Spoke VNet, NSGs, Route Tables
  ]))
  features {
    resource_group {
      prevent_deletion_if_contains_resources = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}
