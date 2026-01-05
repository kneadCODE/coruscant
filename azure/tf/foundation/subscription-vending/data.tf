# =============================================================================
# Data Sources - Reference Bootstrap Management Groups
# =============================================================================

data "azurerm_management_group" "management" {
  name = "mg-coruscant-management"
}

data "azurerm_management_group" "identity" {
  name = "mg-coruscant-identity"
}

data "azurerm_management_group" "connectivity" {
  name = "mg-coruscant-connectivity"
}

data "azurerm_management_group" "security" {
  name = "mg-coruscant-security"
}

data "azurerm_management_group" "corp" {
  name = "mg-coruscant-corp"
}

data "azurerm_management_group" "online" {
  name = "mg-coruscant-online"
}
