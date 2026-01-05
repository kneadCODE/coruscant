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

# resource "azurerm_management_group_subscription_association" "devops_prod" {
#   provider            = azurerm.devops_prod
#   management_group_id = data.azurerm_management_group.corp.id
#   subscription_id     = "/subscriptions/${local.subscription_ids["devops_prod"]}"
# }
# resource "azurerm_resource_provider_registration" "devops_prod" {
#   for_each = toset(concat(
#     local.azure_default_providers,
#     local.base_providers,
#     [
#       "Microsoft.ContainerService",  # AKS for self-hosted GitHub Actions runners
#       "Microsoft.ContainerRegistry", # Azure Container Registries
#       "Microsoft.KeyVault",          # Key Vault for GHA
#       "Microsoft.Network",           # Spoke VNet, NSGs, Route Tables
#     ]
#   ))
#   provider = azurerm.devops_prod

#   name = each.value
# }
# resource "azurerm_role_assignment" "devops_prod_sp_gha_tf_apply_landingzone" {
#   provider             = azurerm.devops_prod
#   scope                = "/subscriptions/${local.subscription_ids["devops_prod"]}"
#   role_definition_name = "Contributor"
#   principal_id         = var.sp_gha_tf_apply_landingzone_obj_id
#   principal_type       = "ServicePrincipal"
# }

# resource "azurerm_management_group_subscription_association" "devops_nonprod" {
#   provider            = azurerm.devops_nonprod
#   management_group_id = data.azurerm_management_group.corp.id
#   subscription_id     = "/subscriptions/${local.subscription_ids["devops_nonprod"]}"
# }
# resource "azurerm_resource_provider_registration" "devops_nonprod" {
#   for_each = toset(concat(
#     local.azure_default_providers,
#     local.base_providers,
#     [
#       "Microsoft.ContainerService",  # AKS for self-hosted GitHub Actions runners
#       "Microsoft.ContainerRegistry", # Azure Container Registries
#       "Microsoft.KeyVault",          # Key Vault for GHA
#       "Microsoft.Network",           # Spoke VNet, NSGs, Route Tables
#     ]
#   ))
#   provider = azurerm.devops_nonprod

#   name = each.value
# }
# resource "azurerm_role_assignment" "devops_nonprod_sp_gha_tf_apply_landingzone" {
#   provider             = azurerm.devops_nonprod
#   scope                = "/subscriptions/${local.subscription_ids["devops_nonprod"]}"
#   role_definition_name = "Contributor"
#   principal_id         = var.sp_gha_tf_apply_landingzone_obj_id
#   principal_type       = "ServicePrincipal"
# }

# # =============================================================================
# # ESB SUBSCRIPTION (Landing Zone - Corp)
# # =============================================================================
# # Purpose: Enterprise service bus, integration services, API management

# resource "azurerm_management_group_subscription_association" "esb_prod" {
#   provider            = azurerm.esb_prod
#   management_group_id = data.azurerm_management_group.corp.id
#   subscription_id     = "/subscriptions/${local.subscription_ids["esb_prod"]}"
# }
# resource "azurerm_resource_provider_registration" "esb_prod" {
#   for_each = toset(concat(
#     local.azure_default_providers,
#     local.base_providers,
#     [
#       "Microsoft.Compute", # VMs for self-hosted Kafka
#       "Microsoft.Network", # Spoke VNet, NSGs, Route Tables, ILB
#     ]
#   ))
#   provider = azurerm.esb_prod

#   name = each.value
# }
# resource "azurerm_role_assignment" "esb_prod_sp_gha_tf_apply_landingzone" {
#   provider             = azurerm.esb_prod
#   scope                = "/subscriptions/${local.subscription_ids["esb_prod"]}"
#   role_definition_name = "Contributor"
#   principal_id         = var.sp_gha_tf_apply_landingzone_obj_id
#   principal_type       = "ServicePrincipal"
# }

# resource "azurerm_management_group_subscription_association" "esb_nonprod" {
#   provider            = azurerm.esb_nonprod
#   management_group_id = data.azurerm_management_group.corp.id
#   subscription_id     = "/subscriptions/${local.subscription_ids["esb_nonprod"]}"
# }
# resource "azurerm_resource_provider_registration" "esb_nonprod" {
#   for_each = toset(concat(
#     local.azure_default_providers,
#     local.base_providers,
#     [
#       "Microsoft.Compute", # VMs for self-hosted Kafka
#       "Microsoft.Network", # Spoke VNet, NSGs, Route Tables, ILB
#     ]
#   ))
#   provider = azurerm.esb_nonprod

#   name = each.value
# }
# resource "azurerm_role_assignment" "esb_nonprod_sp_gha_tf_apply_landingzone" {
#   provider             = azurerm.esb_nonprod
#   scope                = "/subscriptions/${local.subscription_ids["esb_nonprod"]}"
#   role_definition_name = "Contributor"
#   principal_id         = var.sp_gha_tf_apply_landingzone_obj_id
#   principal_type       = "ServicePrincipal"
# }
