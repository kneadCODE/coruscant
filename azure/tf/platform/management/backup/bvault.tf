resource "azurerm_data_protection_backup_vault" "platform_sea_01" {
  name                = "bvault-platform-sea-01"
  location            = azurerm_resource_group.platform_backup_sea.location
  resource_group_name = azurerm_resource_group.platform_backup_sea.name

  datastore_type               = "OperationalStore"
  redundancy                   = "GeoRedundant"
  cross_region_restore_enabled = true
  retention_duration_in_days   = 30
  soft_delete                  = "On"
  immutability                 = "Unlocked" # Later turn to locked for ransomware protection

  identity {
    type = "SystemAssigned"
  }

  tags = merge(azurerm_resource_group.platform_backup_sea.tags, {
    purpose = "platform-backup-paas"
  })

  depends_on = [azurerm_resource_group.platform_backup_sea]
}
