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
# DEVOPS SUBSCRIPTION
# =============================================================================
# Purpose: AKS for GHA runners, ACR, Key Vault, Spoke VNet

resource "azurerm_management_group_subscription_association" "devops_prod" {
  provider            = azurerm.devops_prod
  management_group_id = data.azurerm_management_group.corp.id
  subscription_id     = "/subscriptions/${local.subscription_ids["devops_prod"]}"
}
module "sub_devops_prod_rbac" {
  providers          = { azurerm = azurerm.devops_prod }
  source             = "../../modules/subscription-rbac"
  subscription_id    = local.subscription_ids["devops_prod"]
  sp_tf_apply_obj_id = var.sp_gha_tf_apply_landingzone_obj_id
  depends_on         = [azurerm_management_group_subscription_association.devops_prod]
}

resource "azurerm_management_group_subscription_association" "devops_nonprod" {
  provider            = azurerm.devops_nonprod
  management_group_id = data.azurerm_management_group.corp.id
  subscription_id     = "/subscriptions/${local.subscription_ids["devops_nonprod"]}"
}
module "sub_devops_nonprod_rbac" {
  providers          = { azurerm = azurerm.devops_nonprod }
  source             = "../../modules/subscription-rbac"
  subscription_id    = local.subscription_ids["devops_nonprod"]
  sp_tf_apply_obj_id = var.sp_gha_tf_apply_landingzone_obj_id
  depends_on         = [azurerm_management_group_subscription_association.devops_nonprod]
}

# =============================================================================
# ESB SUBSCRIPTION (Landing Zone - Corp)
# =============================================================================
# Purpose: Enterprise service bus, integration services, API management

resource "azurerm_management_group_subscription_association" "esb_prod" {
  provider            = azurerm.esb_prod
  management_group_id = data.azurerm_management_group.corp.id
  subscription_id     = "/subscriptions/${local.subscription_ids["esb_prod"]}"
}
module "sub_esb_prod_rbac" {
  providers          = { azurerm = azurerm.esb_prod }
  source             = "../../modules/subscription-rbac"
  subscription_id    = local.subscription_ids["esb_prod"]
  sp_tf_apply_obj_id = var.sp_gha_tf_apply_landingzone_obj_id
  depends_on         = [azurerm_management_group_subscription_association.esb_prod]
}

resource "azurerm_management_group_subscription_association" "esb_nonprod" {
  provider            = azurerm.esb_nonprod
  management_group_id = data.azurerm_management_group.corp.id
  subscription_id     = "/subscriptions/${local.subscription_ids["esb_nonprod"]}"
}
module "sub_esb_nonprod_rbac" {
  providers          = { azurerm = azurerm.esb_nonprod }
  source             = "../../modules/subscription-rbac"
  subscription_id    = local.subscription_ids["esb_nonprod"]
  sp_tf_apply_obj_id = var.sp_gha_tf_apply_landingzone_obj_id
  depends_on         = [azurerm_management_group_subscription_association.esb_nonprod]
}

# =============================================================================
# OBSERVABILITY SUBSCRIPTION (Landing Zone - Corp)
# =============================================================================
# Purpose: Application Monitoring (MELT) and OTEL

resource "azurerm_management_group_subscription_association" "observability_prod" {
  provider            = azurerm.observability_prod
  management_group_id = data.azurerm_management_group.corp.id
  subscription_id     = "/subscriptions/${local.subscription_ids["observability_prod"]}"
}
module "sub_observability_prod_rbac" {
  providers          = { azurerm = azurerm.observability_prod }
  source             = "../../modules/subscription-rbac"
  subscription_id    = local.subscription_ids["observability_prod"]
  sp_tf_apply_obj_id = var.sp_gha_tf_apply_landingzone_obj_id
  depends_on         = [azurerm_management_group_subscription_association.observability_prod]
}

resource "azurerm_management_group_subscription_association" "observability_nonprod" {
  provider            = azurerm.observability_nonprod
  management_group_id = data.azurerm_management_group.corp.id
  subscription_id     = "/subscriptions/${local.subscription_ids["observability_nonprod"]}"
}
module "sub_observability_nonprod_rbac" {
  providers          = { azurerm = azurerm.observability_nonprod }
  source             = "../../modules/subscription-rbac"
  subscription_id    = local.subscription_ids["observability_nonprod"]
  sp_tf_apply_obj_id = var.sp_gha_tf_apply_landingzone_obj_id
  depends_on         = [azurerm_management_group_subscription_association.observability_nonprod]
}
