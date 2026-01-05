resource "azurerm_resource_group" "platform_logs_sea" {
  name       = "rg-platform-logs-sea"
  location   = "southeastasia"
  managed_by = "iac"

  tags = {
    org        = "kneadcode"
    portfolio  = "coruscant"
    workload   = "platform-logs"
    env        = "shared"
    purpose    = "platform-logs"
    owner      = "platform-engineering"
    costcenter = "cc-it-platform"
    managed_by = "iac"
    iac_tool   = "opentofu"
  }
}
resource "azurerm_management_lock" "rg_platform_logs_sea_cannot_delete" {
  name       = "lock-rg-platform-logs-sea-cannot-delete"
  scope      = azurerm_resource_group.platform_logs_sea.id
  lock_level = "CanNotDelete"
  depends_on = [azurerm_resource_group.platform_logs_sea]
}

resource "azurerm_storage_account" "platform_logs_archive_sea_01" {
  name                = local.storage_account_names["st_platform_logs_archive_sea_01"]
  location            = azurerm_resource_group.platform_logs_sea.location
  resource_group_name = azurerm_resource_group.platform_logs_sea.name

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  access_tier              = "Cool"
  account_replication_type = "ZRS" # Since this is for archived logs, we don't need GRS

  cross_tenant_replication_enabled = false
  https_traffic_only_enabled       = true
  min_tls_version                  = "TLS1_2" # Eventually move to TLS1_3 when supported widely
  dns_endpoint_type                = "Standard"
  allow_nested_items_to_be_public  = false
  public_network_access_enabled    = true
  network_rules {
    default_action = "Allow" # TODO: Change to Deny after bootstrap is complete and GHA runners are running inside private network
    bypass         = ["AzureServices"]
    # virtual_network_subnet_ids = var.allowed_subnet_ids # Allow traffic from specific subnets (must have Microsoft.Storage service endpoint enabled)
  }
  routing {
    choice                      = "MicrosoftRouting"
    publish_internet_endpoints  = false
    publish_microsoft_endpoints = false
  }

  identity {
    type = "SystemAssigned"
  }

  default_to_oauth_authentication   = true
  shared_access_key_enabled         = false
  is_hns_enabled                    = false
  nfsv3_enabled                     = false
  large_file_share_enabled          = false
  local_user_enabled                = false
  sftp_enabled                      = false
  infrastructure_encryption_enabled = false # Since this is for logs storage only, we don't need this 2nd layer of encryption

  immutability_policy {
    state                         = "Unlocked" # TODO: Change to Locked after creation
    period_since_creation_in_days = 365
    allow_protected_append_writes = true
  }

  sas_policy {
    expiration_period = "0.00:01:00" // We are not planning to use SAS for this account
    expiration_action = "Block"
  }

  blob_properties {
    versioning_enabled       = false # Off since this st is for archival
    change_feed_enabled      = false # Off since this st is for archival
    last_access_time_enabled = false # Off since this st is for archival

    delete_retention_policy {
      days = 7
    }

    container_delete_retention_policy {
      days = 7
    }

    # Intentionally NOT setting:
    # - restore_policy (can conflict; also implies other features)
    # - cors_rule (not needed)
  }

  tags = merge(azurerm_resource_group.platform_logs_sea.tags, {
    purpose = "platform-logs-archive"
  })

  depends_on = [azurerm_resource_group.platform_logs_sea]
}
resource "azurerm_storage_management_policy" "platform_logs_archive_sea_01" {
  storage_account_id = azurerm_storage_account.platform_logs_archive_sea_01.id

  rule {
    name    = "diag-cool-cold-archive-delete"
    enabled = true

    filters {
      blob_types = ["blockBlob"]

      # Optional: scope to your log blobs only.
      # If you find Azure writes into containers like "insights-logs-XYZ",
      # you can scope by prefix (format is "container/prefix" style).
      # prefix_match = ["insights-logs-"]
    }

    actions {
      base_blob {
        tier_to_cold_after_days_since_modification_greater_than    = 30
        tier_to_archive_after_days_since_modification_greater_than = 90
        delete_after_days_since_modification_greater_than          = 365
      }
      # Optional: if Azure writes lots of snapshots/versions (usually not for diag logs)
      # snapshot {
      #   delete_after_days_since_creation_greater_than = 30
      # }
    }
  }

  depends_on = [azurerm_storage_account.platform_logs_archive_sea_01]
}

resource "azurerm_log_analytics_workspace" "platform_logs_sea_01" {
  name                = "law-platform-logs-sea-01"
  location            = azurerm_resource_group.platform_logs_sea.location
  resource_group_name = azurerm_resource_group.platform_logs_sea.name

  sku                                     = "PerGB2018" # For real world, might make sense to go for CapacityReservation after we have stablizied
  retention_in_days                       = 30
  daily_quota_gb                          = 0.5 # Real world quota would be much higher
  immediate_data_purge_on_30_days_enabled = true

  identity {
    type = "SystemAssigned"
  }

  internet_ingestion_enabled      = true
  internet_query_enabled          = true
  allow_resource_only_permissions = true
  local_authentication_enabled    = false

  tags = merge(azurerm_resource_group.platform_logs_sea.tags, {
    purpose = "platform-logs-hot"
  })

  depends_on = [azurerm_resource_group.platform_logs_sea]
}
