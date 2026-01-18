resource "azurerm_resource_group" "backup" {
  count = (var.backup_deploy_rsv || var.backup_deploy_bvault) ? 1 : 0

  name       = "rg-platform-backup-${var.env}-${local.active_region_short}"
  location   = local.active_region
  managed_by = "opentofu"

  tags = merge(local.base_tags, {
    region     = local.active_region
    regionrole = "active"
    workload   = "platform-backup"
    purpose    = "platform-backup"
  })
}

resource "azurerm_management_lock" "backup_rg" {
  count = (var.backup_deploy_rsv || var.backup_deploy_bvault) ? 1 : 0

  name       = "lock-${azurerm_resource_group.backup[0].name}-cannot-delete"
  scope      = azurerm_resource_group.backup[0].id
  lock_level = "CanNotDelete"
}

resource "azurerm_recovery_services_vault" "platform_1" {
  count = var.backup_deploy_rsv ? 1 : 0

  name                = "rsv-platform-${var.env}-${local.active_region_short}-01"
  location            = local.active_region
  resource_group_name = azurerm_resource_group.backup[0].name

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

  tags = merge(azurerm_resource_group.backup[0].tags, {
    purpose = "platform-backup-iaas"
  })
}

resource "azurerm_backup_policy_vm" "data_disk_daily_30d" {
  count = var.backup_deploy_rsv ? 1 : 0

  name                = "rsvbkpol-data-disk-daily-30d"
  resource_group_name = azurerm_resource_group.backup[0].name
  recovery_vault_name = azurerm_recovery_services_vault.platform_1[0].name

  policy_type = "V1" # TODO: Upgrade to v2 after the TF provider is more stable

  backup {
    frequency = "Daily"
    time      = "00:00"
  }
  timezone = "UTC"

  instant_restore_retention_days = 5 # Note: V1 supports 1-5 days typically

  retention_daily {
    count = 30
  }
}

resource "azurerm_data_protection_backup_vault" "platform_compliance_1" {
  count = var.backup_deploy_bvault ? 1 : 0

  name                = "bvault-platform-compliance-${var.env}-${local.active_region_short}-01"
  location            = azurerm_resource_group.backup[0].location
  resource_group_name = azurerm_resource_group.backup[0].name

  datastore_type               = "VaultStore"
  redundancy                   = "GeoRedundant"
  cross_region_restore_enabled = false
  soft_delete                  = "AlwaysOn"
  retention_duration_in_days   = 14
  immutability                 = "Unlocked" # Later turn to locked for ransomware protection

  identity {
    type = "SystemAssigned"
  }

  tags = merge(azurerm_resource_group.backup[0].tags, {
    purpose = "platform-backup-paas-compliance"
  })
}

resource "azurerm_data_protection_backup_policy_postgresql_flexible_server" "weekly_1y" {
  count = var.backup_deploy_bvault ? 1 : 0

  name     = "bkpol-psql-weekly-1y"
  vault_id = azurerm_data_protection_backup_vault.platform_compliance_1[0].id

  backup_repeating_time_intervals = [
    "R/2025-01-05T00:00:00Z/P1W", # Weekly backup
  ]
  time_zone = "UTC"

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
