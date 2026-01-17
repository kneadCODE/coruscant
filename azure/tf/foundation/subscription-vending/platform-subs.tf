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

resource "azurerm_management_group_subscription_association" "management_prod" {
  provider            = azurerm.management_prod
  management_group_id = data.azurerm_management_group.management.id
  subscription_id     = "/subscriptions/${local.subscription_ids["management_prod"]}"
}
module "sub_management_prod_rbac" {
  providers              = { azurerm = azurerm.management_prod }
  source                 = "../../modules/subscription-rbac"
  subscription_id        = local.subscription_ids["management_prod"]
  sp_tf_apply_obj_id     = var.sp_gha_tf_apply_platform_obj_id
  lock_manager_role_name = azurerm_role_definition.locks_manager.name
  depends_on             = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.management_prod]
}
resource "azurerm_management_group_subscription_association" "management_nonprod" {
  provider            = azurerm.management_nonprod
  management_group_id = data.azurerm_management_group.management.id
  subscription_id     = "/subscriptions/${local.subscription_ids["management_nonprod"]}"
}
module "sub_management_nonprod_rbac" {
  providers              = { azurerm = azurerm.management_nonprod }
  source                 = "../../modules/subscription-rbac"
  subscription_id        = local.subscription_ids["management_nonprod"]
  sp_tf_apply_obj_id     = var.sp_gha_tf_apply_platform_obj_id
  lock_manager_role_name = azurerm_role_definition.locks_manager.name
  depends_on             = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.management_nonprod]
}

# =============================================================================
# IDENTITY SUBSCRIPTION
# =============================================================================
# Purpose: Entra ID integration, managed identities, identity governance

resource "azurerm_management_group_subscription_association" "identity_prod" {
  provider            = azurerm.identity_prod
  management_group_id = data.azurerm_management_group.identity.id
  subscription_id     = "/subscriptions/${local.subscription_ids["identity_prod"]}"
}
module "sub_identity_prod_rbac" {
  providers              = { azurerm = azurerm.identity_prod }
  source                 = "../../modules/subscription-rbac"
  subscription_id        = local.subscription_ids["identity_prod"]
  sp_tf_apply_obj_id     = var.sp_gha_tf_apply_platform_obj_id
  lock_manager_role_name = azurerm_role_definition.locks_manager.name
  depends_on             = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.identity_prod]
}
resource "azurerm_management_group_subscription_association" "identity_nonprod" {
  provider            = azurerm.identity_nonprod
  management_group_id = data.azurerm_management_group.identity.id
  subscription_id     = "/subscriptions/${local.subscription_ids["identity_nonprod"]}"
}
module "sub_identity_nonprod_rbac" {
  providers              = { azurerm = azurerm.identity_nonprod }
  source                 = "../../modules/subscription-rbac"
  subscription_id        = local.subscription_ids["identity_nonprod"]
  sp_tf_apply_obj_id     = var.sp_gha_tf_apply_platform_obj_id
  lock_manager_role_name = azurerm_role_definition.locks_manager.name
  depends_on             = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.identity_nonprod]
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
