# =============================================================================
# Management Group Hierarchy
# =============================================================================
# This workspace creates the management group hierarchy only.
# Subscription assignments are handled by the subscription-vending workspace.

# Level 1: Top-level management groups (children of root)
resource "azurerm_management_group" "platform" {
  name                       = "mg-coruscant-platform"
  display_name               = "Platform"
  parent_management_group_id = data.azurerm_management_group.root.id
}

resource "azurerm_management_group" "landingzone" {
  name                       = "mg-coruscant-landingzone"
  display_name               = "Landing Zones"
  parent_management_group_id = data.azurerm_management_group.root.id
}

resource "azurerm_management_group" "sandbox" {
  name                       = "mg-coruscant-sandbox"
  display_name               = "Sandbox"
  parent_management_group_id = data.azurerm_management_group.root.id
}

resource "azurerm_management_group" "decommissioned" {
  name                       = "mg-coruscant-decommissioned"
  display_name               = "Decommissioned"
  parent_management_group_id = data.azurerm_management_group.root.id
}

# Level 2: Platform child management groups
resource "azurerm_management_group" "management" {
  name                       = "mg-coruscant-management"
  display_name               = "Management"
  parent_management_group_id = azurerm_management_group.platform.id
}

resource "azurerm_management_group" "identity" {
  name                       = "mg-coruscant-identity"
  display_name               = "Identity"
  parent_management_group_id = azurerm_management_group.platform.id
}

resource "azurerm_management_group" "connectivity" {
  name                       = "mg-coruscant-connectivity"
  display_name               = "Connectivity"
  parent_management_group_id = azurerm_management_group.platform.id
}

resource "azurerm_management_group" "security" {
  name                       = "mg-coruscant-security"
  display_name               = "Security"
  parent_management_group_id = azurerm_management_group.platform.id
}

# Level 2: Landing zone child management groups
resource "azurerm_management_group" "corp" {
  name                       = "mg-coruscant-corp"
  display_name               = "Corp"
  parent_management_group_id = azurerm_management_group.landingzone.id
}

resource "azurerm_management_group" "online" {
  name                       = "mg-coruscant-online"
  display_name               = "Online"
  parent_management_group_id = azurerm_management_group.landingzone.id
}
