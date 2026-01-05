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
resource "azurerm_role_assignment" "edge_prod_sp_gha_tf_apply_landingzone" {
  provider             = azurerm.edge_prod
  scope                = "/subscriptions/${local.subscription_ids["edge_prod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_landingzone_obj_id
  principal_type       = "ServicePrincipal"
}

resource "azurerm_management_group_subscription_association" "edge_nonprod" {
  provider            = azurerm.edge_nonprod
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["edge_nonprod"]}"
}
resource "azurerm_role_assignment" "edge_nonprod_sp_gha_tf_apply_landingzone" {
  provider             = azurerm.edge_nonprod
  scope                = "/subscriptions/${local.subscription_ids["edge_nonprod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_landingzone_obj_id
  principal_type       = "ServicePrincipal"
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
resource "azurerm_role_assignment" "iam_prod_sp_gha_tf_apply_landingzone" {
  provider             = azurerm.iam_prod
  scope                = "/subscriptions/${local.subscription_ids["iam_prod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_landingzone_obj_id
  principal_type       = "ServicePrincipal"
}

resource "azurerm_management_group_subscription_association" "iam_nonprod" {
  provider            = azurerm.iam_nonprod
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["iam_nonprod"]}"
}
resource "azurerm_role_assignment" "iam_nonprod_sp_gha_tf_apply_landingzone" {
  provider             = azurerm.iam_nonprod
  scope                = "/subscriptions/${local.subscription_ids["iam_nonprod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_landingzone_obj_id
  principal_type       = "ServicePrincipal"
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
resource "azurerm_role_assignment" "payment_prod_sp_gha_tf_apply_landingzone" {
  provider             = azurerm.payment_prod
  scope                = "/subscriptions/${local.subscription_ids["payment_prod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_landingzone_obj_id
  principal_type       = "ServicePrincipal"
}

resource "azurerm_management_group_subscription_association" "payment_nonprod" {
  provider            = azurerm.payment_nonprod
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["payment_nonprod"]}"
}
resource "azurerm_role_assignment" "payment_nonprod_sp_gha_tf_apply_landingzone" {
  provider             = azurerm.payment_nonprod
  scope                = "/subscriptions/${local.subscription_ids["payment_nonprod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_landingzone_obj_id
  principal_type       = "ServicePrincipal"
}
