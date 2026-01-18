resource "azurerm_resource_group" "siem" {
  count = (var.siem_deploy_st || var.siem_deploy_law) ? 1 : 0

  name       = "rg-siem-${var.env}-${local.active_region_short}"
  location   = local.active_region
  managed_by = "opentofu"

  tags = merge(local.base_tags, {
    region     = local.active_region
    regionrole = "active"
    workload   = "siem"
    purpose    = "siem"
  })
}

resource "azurerm_management_lock" "siem_rg" {
  count = (var.siem_deploy_st || var.siem_deploy_law) ? 1 : 0

  name       = "lock-${azurerm_resource_group.siem[0].name}-cannot-delete"
  scope      = azurerm_resource_group.siem[0].id
  lock_level = "CanNotDelete"
}

resource "azurerm_storage_account" "siem_archive_1" {
  count = var.siem_deploy_st ? 1 : 0

  name                = "${var.storage_account_prefix}siema${local.env_shorthand_lookup[var.env]}${local.active_region_short}1"
  location            = azurerm_resource_group.siem[0].location
  resource_group_name = azurerm_resource_group.siem[0].name

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  access_tier              = "Cool"
  account_replication_type = "LRS" # Since this is for archived siem, we don't need GRS. Can't use ZRS as it does not support archive tier

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
  infrastructure_encryption_enabled = false # Since this is for siem storage only, we don't need this 2nd layer of encryption

  sas_policy {
    expiration_period = "0.00:01:00" # We are not planning to use SAS for this account
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

  tags = merge(azurerm_resource_group.siem[0].tags, {
    purpose = "platform-siem-archive"
  })
}

resource "azurerm_storage_management_policy" "siem_archive_1" {
  count = var.siem_deploy_st ? 1 : 0

  storage_account_id = azurerm_storage_account.siem_archive_1[0].id

  rule {
    name    = "diag-cool-cold-archive-delete"
    enabled = true

    filters {
      blob_types = ["blockBlob"]

      # Optional: scope to your log blobs only.
      # If you find Azure writes into containers like "insights-siem-XYZ",
      # you can scope by prefix (format is "container/prefix" style).
      # prefix_match = ["insights-siem-"]
    }

    actions {
      base_blob {
        tier_to_cold_after_days_since_modification_greater_than    = 30
        tier_to_archive_after_days_since_modification_greater_than = 90
        delete_after_days_since_modification_greater_than          = 365
      }
      # Optional: if Azure writes lots of snapshots/versions (usually not for diag siem)
      # snapshot {
      #   delete_after_days_since_creation_greater_than = 30
      # }
    }
  }
}

resource "azurerm_log_analytics_workspace" "siem_1" {
  count = var.siem_deploy_law ? 1 : 0

  name                = "log-siem-${var.env}-${local.active_region_short}-01"
  location            = azurerm_resource_group.siem[0].location
  resource_group_name = azurerm_resource_group.siem[0].name

  sku                                     = "PerGB2018" # For real world, might make sense to go for CapacityReservation after we have stabilized
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

  tags = merge(azurerm_resource_group.siem[0].tags, {
    purpose = "siem-hot"
  })
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "siem_1" {
  count = var.siem_deploy_law ? 1 : 0

  workspace_id = azurerm_log_analytics_workspace.siem_1[0].id
}
