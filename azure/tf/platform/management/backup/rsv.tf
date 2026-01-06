resource "azurerm_recovery_services_vault" "platform_sea_01" {
  name                = "rsv-platform-sea-01"
  location            = azurerm_resource_group.platform_backup_sea.location
  resource_group_name = azurerm_resource_group.platform_backup_sea.name

  sku                           = "Standard"
  storage_mode_type             = "GeoRedundant"
  cross_region_restore_enabled  = false
  soft_delete_enabled           = true
  immutability                  = "Unlocked" # Later turn to locked for ransomware protection
  public_network_access_enabled = true

  identity {
    type = "SystemAssigned"
  }

  monitoring {
    alerts_for_all_job_failures_enabled            = true
    alerts_for_critical_operation_failures_enabled = true
  }

  tags = merge(azurerm_resource_group.platform_backup_sea.tags, {
    purpose = "platform-backup-iaas"
  })

  depends_on = [azurerm_resource_group.platform_backup_sea]
}
