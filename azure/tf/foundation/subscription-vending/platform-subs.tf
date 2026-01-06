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
  provider            = azurerm.management
  management_group_id = data.azurerm_management_group.management.id
  subscription_id     = "/subscriptions/${local.subscription_ids["management"]}"
}
module "sub_management_rbac" {
  providers              = { azurerm = azurerm.management }
  source                 = "../../modules/subscription-rbac"
  subscription_id        = local.subscription_ids["management"]
  sp_tf_apply_obj_id     = var.sp_gha_tf_apply_platform_obj_id
  lock_manager_role_name = azurerm_role_definition.locks_manager.name
  depends_on             = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.management]
}

# =============================================================================
# IDENTITY SUBSCRIPTION
# =============================================================================
# Purpose: Entra ID integration, managed identities, identity governance

resource "azurerm_management_group_subscription_association" "identity" {
  provider            = azurerm.identity
  management_group_id = data.azurerm_management_group.identity.id
  subscription_id     = "/subscriptions/${local.subscription_ids["identity"]}"
}
module "sub_identity_rbac" {
  providers              = { azurerm = azurerm.identity }
  source                 = "../../modules/subscription-rbac"
  subscription_id        = local.subscription_ids["identity"]
  sp_tf_apply_obj_id     = var.sp_gha_tf_apply_platform_obj_id
  lock_manager_role_name = azurerm_role_definition.locks_manager.name
  depends_on             = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.identity]
}

# =============================================================================
# CONNECTIVITY SUBSCRIPTION
# =============================================================================
# Purpose: Hub VNets, VPN/ExpressRoute gateways, Azure Firewall, DNS

resource "azurerm_management_group_subscription_association" "connectivity_prod" {
  provider            = azurerm.connectivity_prod
  management_group_id = data.azurerm_management_group.connectivity.id
  subscription_id     = "/subscriptions/${local.subscription_ids["connectivity_prod"]}"
}
module "sub_connectivity_prod_rbac" {
  providers              = { azurerm = azurerm.connectivity_prod }
  source                 = "../../modules/subscription-rbac"
  subscription_id        = local.subscription_ids["connectivity_prod"]
  sp_tf_apply_obj_id     = var.sp_gha_tf_apply_platform_obj_id
  lock_manager_role_name = azurerm_role_definition.locks_manager.name
  depends_on             = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.connectivity_prod]
}

resource "azurerm_management_group_subscription_association" "connectivity_nonprod" {
  provider            = azurerm.connectivity_nonprod
  management_group_id = data.azurerm_management_group.connectivity.id
  subscription_id     = "/subscriptions/${local.subscription_ids["connectivity_nonprod"]}"
}
module "sub_connectivity_nonprod_rbac" {
  providers              = { azurerm = azurerm.connectivity_nonprod }
  source                 = "../../modules/subscription-rbac"
  subscription_id        = local.subscription_ids["connectivity_nonprod"]
  sp_tf_apply_obj_id     = var.sp_gha_tf_apply_platform_obj_id
  lock_manager_role_name = azurerm_role_definition.locks_manager.name
  depends_on             = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.connectivity_nonprod]
}

# =============================================================================
# SECURITY SUBSCRIPTION
# =============================================================================
# Purpose: Security tooling (Sentinel, Defender, HashiCorp Vault VMs, DDoS Protection)

resource "azurerm_management_group_subscription_association" "security_prod" {
  provider            = azurerm.security_prod
  management_group_id = data.azurerm_management_group.security.id
  subscription_id     = "/subscriptions/${local.subscription_ids["security_prod"]}"
}
module "sub_security_prod_rbac" {
  providers              = { azurerm = azurerm.security_prod }
  source                 = "../../modules/subscription-rbac"
  subscription_id        = local.subscription_ids["security_prod"]
  sp_tf_apply_obj_id     = var.sp_gha_tf_apply_platform_obj_id
  lock_manager_role_name = azurerm_role_definition.locks_manager.name
  depends_on             = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.security_prod]
}

resource "azurerm_management_group_subscription_association" "security_nonprod" {
  provider            = azurerm.security_nonprod
  management_group_id = data.azurerm_management_group.security.id
  subscription_id     = "/subscriptions/${local.subscription_ids["security_nonprod"]}"
}
module "sub_security_nonprod_rbac" {
  providers              = { azurerm = azurerm.security_nonprod }
  source                 = "../../modules/subscription-rbac"
  subscription_id        = local.subscription_ids["security_nonprod"]
  sp_tf_apply_obj_id     = var.sp_gha_tf_apply_platform_obj_id
  lock_manager_role_name = azurerm_role_definition.locks_manager.name
  depends_on             = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.security_nonprod]
}
