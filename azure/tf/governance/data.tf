# Data sources for management group references
data "azurerm_management_group" "platform" {
  name = "mg-coruscant-platform"
}

data "azurerm_management_group" "landing_zone" {
  name = "mg-coruscant-landingzone"
}

data "azurerm_management_group" "corp" {
  name = "mg-coruscant-corp"
}

data "azurerm_management_group" "online" {
  name = "mg-coruscant-online"
}

data "azurerm_management_group" "decommissioned" {
  name = "mg-coruscant-decommissioned"
}

data "azurerm_management_group" "sandbox" {
  name = "mg-coruscant-decommissioned"
}
