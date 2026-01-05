# =============================================================================
# Subscription Vending - Manages subscriptions
# =============================================================================
# This workspace handles subscription lifecycle operations:
# 1. Subscription-to-MG associations
# 2. Resource provider registration (least privilege per subscription)
# 3. Service principal RBAC assignments
#
# Organization: Resources grouped by subscription for clarity

# =============================================================================
# MANAGEMENT SUBSCRIPTION
# =============================================================================
# Purpose: Management and governance (policy, cost management, monitoring)

resource "azurerm_management_group_subscription_association" "management" {
  management_group_id = data.azurerm_management_group.management.id
  subscription_id     = "/subscriptions/${local.subscription_ids["management"]}"
}
resource "azurerm_resource_provider_registration" "management" {
  for_each = toset(concat(
    local.azure_default_providers,
    local.base_providers,
    [
      "Microsoft.OperationalInsights", # Log Analytics Workspace for platform logs
      "Microsoft.Storage",             # Storage accounts for logs
      "Microsoft.RecoveryServices",    # Recovery Services Vault and Backup Vault
    ]
  ))
  provider = azurerm.management

  name = each.value
}
resource "azurerm_role_assignment" "management_sp_gha_tf_apply_platform" {
  scope                = "/subscriptions/${local.subscription_ids["management"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_platform
}

# =============================================================================
# IDENTITY SUBSCRIPTION
# =============================================================================
# Purpose: Entra ID integration, managed identities, identity governance

resource "azurerm_management_group_subscription_association" "identity" {
  management_group_id = data.azurerm_management_group.identity.id
  subscription_id     = "/subscriptions/${local.subscription_ids["identity"]}"
}
resource "azurerm_resource_provider_registration" "identity" {
  for_each = toset(concat(
    local.azure_default_providers,
    local.base_providers,
  ))
  provider = azurerm.identity

  name = each.value
}
resource "azurerm_role_assignment" "identity_sp_gha_tf_apply_platform" {
  scope                = "/subscriptions/${local.subscription_ids["identity"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_platform
}

# =============================================================================
# CONNECTIVITY SUBSCRIPTION
# =============================================================================
# Purpose: Hub VNets, VPN/ExpressRoute gateways, Azure Firewall, DNS

resource "azurerm_management_group_subscription_association" "connectivity_prod" {
  management_group_id = data.azurerm_management_group.connectivity.id
  subscription_id     = "/subscriptions/${local.subscription_ids["connectivity_prod"]}"
}
resource "azurerm_resource_provider_registration" "connectivity_prod" {
  for_each = toset(concat(
    local.azure_default_providers,
    local.base_providers,
    [
      "Microsoft.Network", # VNets, NSGs, Route Tables, Firewalls, DNS, Gateways, Bastion
    ]
  ))
  provider = azurerm.connectivity_prod

  name = each.value
}
resource "azurerm_role_assignment" "connectivity_prod_sp_gha_tf_apply_platform" {
  scope                = "/subscriptions/${local.subscription_ids["connectivity_prod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_platform
}

resource "azurerm_management_group_subscription_association" "connectivity_nonprod" {
  management_group_id = data.azurerm_management_group.connectivity.id
  subscription_id     = "/subscriptions/${local.subscription_ids["connectivity_nonprod"]}"
}
resource "azurerm_resource_provider_registration" "connectivity_nonprod" {
  for_each = toset(concat(
    local.azure_default_providers,
    local.base_providers,
    [
      "Microsoft.Network", # VNets, NSGs, Route Tables, Firewalls, DNS, Gateways, Bastion
    ]
  ))
  provider = azurerm.connectivity_nonprod

  name = each.value
}
resource "azurerm_role_assignment" "connectivity_nonprod_sp_gha_tf_apply_platform" {
  scope                = "/subscriptions/${local.subscription_ids["connectivity_nonprod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_platform
}

# =============================================================================
# SECURITY SUBSCRIPTION
# =============================================================================
# Purpose: Security tooling (Sentinel, Defender, HashiCorp Vault VMs, DDoS Protection)

resource "azurerm_management_group_subscription_association" "security_prod" {
  management_group_id = data.azurerm_management_group.security.id
  subscription_id     = "/subscriptions/${local.subscription_ids["security_prod"]}"
}
resource "azurerm_resource_provider_registration" "security_prod" {
  for_each = toset(concat(
    local.azure_default_providers,
    local.base_providers,
    [
      "Microsoft.OperationalInsights", # Log Analytics Workspace with Sentinel
      "Microsoft.Storage",             # Storage accounts for security logs
      "Microsoft.Compute",             # VMs for HashiCorp Vault
      "Microsoft.Network",             # Spoke VNet, NSGs, Route Tables, DDoS Protection Plan
    ]
  ))
  provider = azurerm.security_prod

  name = each.value
}
resource "azurerm_role_assignment" "security_prod_sp_gha_tf_apply_platform" {
  scope                = "/subscriptions/${local.subscription_ids["security_prod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_platform
}

resource "azurerm_management_group_subscription_association" "security_nonprod" {
  management_group_id = data.azurerm_management_group.security.id
  subscription_id     = "/subscriptions/${local.subscription_ids["security_nonprod"]}"
}
resource "azurerm_resource_provider_registration" "security_nonprod" {
  for_each = toset(concat(
    local.azure_default_providers,
    local.base_providers,
    [
      "Microsoft.OperationalInsights", # Log Analytics Workspace with Sentinel
      "Microsoft.Storage",             # Storage accounts for security logs
      "Microsoft.Compute",             # VMs for HashiCorp Vault
      "Microsoft.Network",             # Spoke VNet, NSGs, Route Tables, DDoS Protection Plan
    ]
  ))
  provider = azurerm.security_nonprod

  name = each.value
}
resource "azurerm_role_assignment" "security_nonprod_sp_gha_tf_apply_platform" {
  scope                = "/subscriptions/${local.subscription_ids["security_nonprod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_platform
}
