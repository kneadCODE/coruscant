resource "azurerm_resource_group" "rg" {
  for_each = local.platform_logs_regions

  name       = "rg-platform-logs-${each.value.short_name}"
  location   = each.value.name
  managed_by = "iac"

  tags = {
    org        = "kneadcode"
    portfolio  = "coruscant"
    workload   = "platform-logs"
    env        = "universal"
    owner      = "platform-engineering"
    costcenter = "cc-it-platform"
    managed_by = "opentofu"
    purpose    = "platform-logs"
  }
}

resource "azurerm_management_lock" "rg" {
  for_each = azurerm_resource_group.rg

  name       = "lock-${each.value.name}-cannot-delete"
  scope      = each.value.id
  lock_level = "CanNotDelete"
}
