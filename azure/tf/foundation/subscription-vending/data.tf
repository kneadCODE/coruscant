# =============================================================================
# Data Sources - Reference Bootstrap Management Groups
# =============================================================================

data "azurerm_management_group" "root" {
  name = "mg-coruscant-root"
}

data "azurerm_management_group" "platform" {
  name = "mg-coruscant-platform"
}

data "azurerm_management_group" "landingzone" {
  name = "mg-coruscant-landingzone"
}

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
