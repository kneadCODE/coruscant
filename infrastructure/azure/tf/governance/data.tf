data "azurerm_management_group" "root" {
  name = "mg-coruscant-root"
}

data "azurerm_subscription" "management" {
  subscription_id = var.subscription_id_management
}

data "azurerm_subscription" "identity" {
  subscription_id = var.subscription_id_identity
}

data "azurerm_subscription" "connectivity-prod" {
  subscription_id = var.subscription_id_connectivity_prod
}

data "azurerm_subscription" "security-prod" {
  subscription_id = var.subscription_id_security_prod
}
