resource "azurerm_resource_group" "platform_backup_sea" {
  name       = "rg-platform-backup-sea"
  location   = "southeastasia"
  managed_by = "iac"

  tags = {
    org        = "kneadcode"
    portfolio  = "coruscant"
    workload   = "platform-backup"
    env        = "shared"
    purpose    = "platform-backup"
    owner      = "platform-engineering"
    costcenter = "cc-it-platform"
    managed_by = "iac"
    iac_tool   = "opentofu"
  }
}
resource "azurerm_management_lock" "rg_platform_backup_sea_cannot_delete" {
  name       = "lock-rg-platform-backup-sea-cannot-delete"
  scope      = azurerm_resource_group.platform_backup_sea.id
  lock_level = "CanNotDelete"
  depends_on = [azurerm_resource_group.platform_backup_sea]
}
