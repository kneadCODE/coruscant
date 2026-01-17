resource "azurerm_storage_account" "archive_01" {
  for_each = local.platform_logs_regions

  name                = "${var.storage_account_prefix}ploga${each.value.short_name}u1"
  location            = azurerm_resource_group.rg[each.key].location
  resource_group_name = azurerm_resource_group.rg[each.key].name

  account_kind             = "StorageV2"
  account_tier             = "Standard"
  access_tier              = "Cool"
  account_replication_type = "LRS" # Since this is for archived logs, we don't need GRS. Can't use ZRS as it does not support archive tier

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

  tags = merge(azurerm_resource_group.rg[each.key].tags, {
    purpose = "platform-logs-archive"
  })
}

resource "azurerm_storage_management_policy" "archive_01" {
  for_each = local.platform_logs_regions

  storage_account_id = azurerm_storage_account.archive_01[each.key].id

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
}
