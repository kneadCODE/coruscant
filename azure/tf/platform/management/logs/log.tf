resource "azurerm_log_analytics_workspace" "log_01" {
  for_each = local.platform_logs_regions

  name                = "log-platform-logs-${each.value.short_name}-01"
  location            = azurerm_resource_group.rg[each.key].location
  resource_group_name = azurerm_resource_group.rg[each.key].name

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

  tags = merge(azurerm_resource_group.rg[each.key].tags, {
    purpose = "platform-logs-hot"
  })
}
