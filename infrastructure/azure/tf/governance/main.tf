# Force refresh

data "azurerm_management_group" "root" {
  name = "mg-coruscant-root"
}

resource "azurerm_management_group" "platform" {
  name                       = "mg-coruscant-platform"
  parent_management_group_id = data.azurerm_management_group.root.id
}

resource "azurerm_management_group" "landingzone" {
  name                       = "mg-coruscant-landingzone"
  parent_management_group_id = data.azurerm_management_group.root.id
}

resource "azurerm_management_group" "sandbox" {
  name                       = "mg-coruscant-sandbox"
  parent_management_group_id = data.azurerm_management_group.root.id
}

resource "azurerm_management_group" "decommissioned" {
  name                       = "mg-coruscant-decommissioned"
  parent_management_group_id = data.azurerm_management_group.root.id
}

resource "azurerm_management_group" "management" {
  name                       = "mg-coruscant-management"
  parent_management_group_id = azurerm_management_group.platform.id
  subscription_ids = [
    var.subscription_id_management,
  ]
}

resource "azurerm_management_group" "identity" {
  name                       = "mg-coruscant-identity"
  parent_management_group_id = azurerm_management_group.platform.id
  subscription_ids = [
    var.subscription_id_identity,
  ]
}

resource "azurerm_management_group" "connectivity-prod" {
  name                       = "mg-coruscant-connectivity-prod"
  parent_management_group_id = azurerm_management_group.platform.id
  subscription_ids = [
    var.subscription_id_connectivity_prod,
  ]
}

resource "azurerm_management_group" "security-prod" {
  name                       = "mg-coruscant-security-prod"
  parent_management_group_id = azurerm_management_group.platform.id
  subscription_ids = [
    var.subscription_id_security_prod,
  ]
}


# module "avm-ptn-alz-management" {
#   source  = "Azure/avm-ptn-alz-management/azurerm"
#   version = "0.9.0"
#
#   automation_account_name      = "aa-prod-eus-001"
#   location                     = "eastus"
#   log_analytics_workspace_name = "law-prod-eus-001"
#   resource_group_name          = "rg-management-eus-001"
# }
#
# module "avm-ptn-alz" {
#   source  = "Azure/avm-ptn-alz/azurerm"
#   version = "0.15.0"
#
#   automation_account_name      = "aa-prod-eus-001"
#   location                     = "eastus"
#   log_analytics_workspace_name = "law-prod-eus-001"
#   resource_group_name          = "rg-management-eus-001"
# }
#
#
# module "avm-ptn-alz-connectivity-hub-and-spoke-vnet" {
#   source  = "Azure/avm-ptn-alz-connectivity-hub-and-spoke-vnet/azurerm"
#   version = "0.16.7"
#
#   default_naming_convention = {}
#
#   automation_account_name      = "aa-prod-eus-001"
#   location                     = "eastus"
#   log_analytics_workspace_name = "law-prod-eus-001"
#   resource_group_name          = "rg-management-eus-001"
# }
#
#
# module "avm-ptn-alz-connectivity-virtual-wan" {
#   source  = "Azure/avm-ptn-alz-connectivity-virtual-wan/azurerm"
#   version = "0.13.4"
#
#   automation_account_name      = "aa-prod-eus-001"
#   location                     = "eastus"
#   log_analytics_workspace_name = "law-prod-eus-001"
#   resource_group_name          = "rg-management-eus-001"
# }
