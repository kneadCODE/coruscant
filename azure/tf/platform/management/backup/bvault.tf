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

resource "azurerm_data_protection_backup_policy_postgresql_flexible_server" "weekly_1y_sea" { # For prod instances
  name     = "bkpol-psql-weekly-1y-sea"
  vault_id = azurerm_data_protection_backup_vault.platform_compliance_sea_01.id

  backup_repeating_time_intervals = [
    "R/2025-01-05T00:00:00Z/P1W", # Weekly backup
  ]
  time_zone = "Asia/Singapore"

  default_retention_rule {
    life_cycle {
      data_store_type = "VaultStore"
      duration        = "P1Y"
    }
  }

  retention_rule {
    name     = "weekly_1y"
    priority = 1

    criteria {
      absolute_criteria = "FirstOfWeek"
    }

    life_cycle {
      data_store_type = "VaultStore"
      duration        = "P1Y"
    }
  }
}

resource "azurerm_data_protection_backup_policy_postgresql_flexible_server" "weekly_30d_sea" { # For non-prod instances
  name     = "bkpol-psql-weekly-30d-sea"
  vault_id = azurerm_data_protection_backup_vault.platform_compliance_sea_01.id

  backup_repeating_time_intervals = [
    "R/2025-01-05T00:00:00Z/P1W", # Weekly backup
  ]
  time_zone = "Asia/Singapore"

  default_retention_rule {
    life_cycle {
      data_store_type = "VaultStore"
      duration        = "P30D"
    }
  }

  retention_rule {
    name     = "weekly-30d"
    priority = 1

    criteria {
      absolute_criteria = "FirstOfWeek"
    }

    life_cycle {
      data_store_type = "VaultStore"
      duration        = "P30D"
    }
  }
}
