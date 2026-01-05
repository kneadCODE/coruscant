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
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["edge_prod"]}"
}

resource "azurerm_resource_provider_registration" "edge_prod" {
  for_each = toset(concat(
    local.azure_default_providers,
    local.base_providers,
    [
      "Microsoft.Cdn",      # Azure Front Door
      "Microsoft.Devices",  # IoT Hub for MQTT/AMQP ingress
      "Microsoft.EventHub", # Event Hubs for Kafka ingress
      "Microsoft.Storage",  # Storage accounts for SFTP ingress
      "Microsoft.Network",  # Spoke VNet, NSGs, Route Tables
    ]
  ))
  provider = azurerm.edge_prod

  name = each.value
}

resource "azurerm_role_assignment" "edge_prod_sp_gha_tf_apply_landingzone" {
  scope                = "/subscriptions/${local.subscription_ids["edge_prod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_landingzone
}

resource "azurerm_management_group_subscription_association" "edge_nonprod" {
  management_group_id = data.azurerm_management_group.online.id
  subscription_id     = "/subscriptions/${local.subscription_ids["edge_nonprod"]}"
}

resource "azurerm_resource_provider_registration" "edge_nonprod" {
  for_each = toset(concat(
    local.azure_default_providers,
    local.base_providers,
    [
      "Microsoft.Cdn",      # Azure Front Door
      "Microsoft.Devices",  # IoT Hub for MQTT/AMQP ingress
      "Microsoft.EventHub", # Event Hubs for Kafka ingress
      "Microsoft.Storage",  # Storage accounts for SFTP ingress
      "Microsoft.Network",  # Spoke VNet, NSGs, Route Tables
    ]
  ))
  provider = azurerm.edge_nonprod

  name = each.value
}

resource "azurerm_role_assignment" "edge_nonprod_sp_gha_tf_apply_landingzone" {
  scope                = "/subscriptions/${local.subscription_ids["edge_nonprod"]}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_gha_tf_apply_landingzone
}
