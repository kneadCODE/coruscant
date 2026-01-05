# Root management group (pre-existing, created manually)
data "azurerm_management_group" "root" {
  name = "mg-coruscant-root"
}
