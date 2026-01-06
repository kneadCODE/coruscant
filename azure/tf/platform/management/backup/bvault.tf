resource "azurerm_data_protection_backup_vault" "platform_compliance_sea_01" {
  name                = "bvault-platform-compliance-sea-01"
  location            = azurerm_resource_group.platform_backup_sea.location
  resource_group_name = azurerm_resource_group.platform_backup_sea.name

  datastore_type               = "VaultStore"
  redundancy                   = "GeoRedundant"
  cross_region_restore_enabled = false
  soft_delete                  = "AlwaysOn"
  retention_duration_in_days   = 14
  immutability                 = "Unlocked" # Later turn to locked for ransomware protection

  identity {
    type = "SystemAssigned"
  }

  tags = merge(azurerm_resource_group.platform_backup_sea.tags, {
    purpose = "platform-backup-paas-compliance"
  })

  depends_on = [azurerm_resource_group.platform_backup_sea]
}

# Later if we need an operational vault, create one here
