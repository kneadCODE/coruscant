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
# EDGE SUBSCRIPTION (Landing Zone - Online)
# =============================================================================
# Purpose: Internet-facing workloads, edge services, CDN

resource "azurerm_management_group_subscription_association" "edge_prod" {
  provider            = azurerm.edge_prod
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["edge_prod"]}"
}
module "sub_edge_prod_rbac" {
  providers            = { azurerm = azurerm.edge_prod }
  source               = "../../modules/subscription-rbac"
  subscription_id      = local.subscription_ids["edge_prod"]
  sp_tf_apply_obj_id   = var.sp_gha_tf_apply_landingzone_obj_id
  lock_manager_role_id = azurerm_role_definition.locks_manager.role_definition_resource_id
  depends_on           = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.edge_prod]
}

resource "azurerm_management_group_subscription_association" "edge_nonprod" {
  provider            = azurerm.edge_nonprod
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["edge_nonprod"]}"
}
module "sub_edge_nonprod_rbac" {
  providers            = { azurerm = azurerm.edge_nonprod }
  source               = "../../modules/subscription-rbac"
  subscription_id      = local.subscription_ids["edge_nonprod"]
  sp_tf_apply_obj_id   = var.sp_gha_tf_apply_landingzone_obj_id
  lock_manager_role_id = azurerm_role_definition.locks_manager.role_definition_resource_id
  depends_on           = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.edge_nonprod]
}

# =============================================================================
# IAM SUBSCRIPTION (Landing Zone - Online)
# =============================================================================
# Purpose: Application IAM

resource "azurerm_management_group_subscription_association" "iam_prod" {
  provider            = azurerm.iam_prod
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["iam_prod"]}"
}
module "sub_iam_prod_rbac" {
  providers            = { azurerm = azurerm.iam_prod }
  source               = "../../modules/subscription-rbac"
  subscription_id      = local.subscription_ids["iam_prod"]
  sp_tf_apply_obj_id   = var.sp_gha_tf_apply_landingzone_obj_id
  lock_manager_role_id = azurerm_role_definition.locks_manager.role_definition_resource_id
  depends_on           = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.iam_prod]
}

resource "azurerm_management_group_subscription_association" "iam_nonprod" {
  provider            = azurerm.iam_nonprod
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["iam_nonprod"]}"
}
module "sub_iam_nonprod_rbac" {
  providers            = { azurerm = azurerm.iam_nonprod }
  source               = "../../modules/subscription-rbac"
  subscription_id      = local.subscription_ids["iam_nonprod"]
  sp_tf_apply_obj_id   = var.sp_gha_tf_apply_landingzone_obj_id
  lock_manager_role_id = azurerm_role_definition.locks_manager.role_definition_resource_id
  depends_on           = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.iam_nonprod]
}

# =============================================================================
# PAYMENT SUBSCRIPTION (Landing Zone - Online)
# =============================================================================
# Purpose: Payment Gateways

resource "azurerm_management_group_subscription_association" "payment_prod" {
  provider            = azurerm.payment_prod
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["payment_prod"]}"
}
module "sub_payment_prod_rbac" {
  providers            = { azurerm = azurerm.payment_prod }
  source               = "../../modules/subscription-rbac"
  subscription_id      = local.subscription_ids["payment_prod"]
  sp_tf_apply_obj_id   = var.sp_gha_tf_apply_landingzone_obj_id
  lock_manager_role_id = azurerm_role_definition.locks_manager.role_definition_resource_id
  depends_on           = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.payment_prod]
}

resource "azurerm_management_group_subscription_association" "payment_nonprod" {
  provider            = azurerm.payment_nonprod
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["payment_nonprod"]}"
}
module "sub_payment_nonprod_rbac" {
  providers            = { azurerm = azurerm.payment_nonprod }
  source               = "../../modules/subscription-rbac"
  subscription_id      = local.subscription_ids["payment_nonprod"]
  sp_tf_apply_obj_id   = var.sp_gha_tf_apply_landingzone_obj_id
  lock_manager_role_id = azurerm_role_definition.locks_manager.role_definition_resource_id
  depends_on           = [azurerm_role_definition.locks_manager, azurerm_management_group_subscription_association.payment_nonprod]
}
