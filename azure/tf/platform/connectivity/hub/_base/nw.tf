resource "azurerm_resource_group" "netops" {
  for_each = local.regions[var.region_pair]

  name       = "rg-connectivity-netops-${var.env}-${each.value.short_name}"
  location   = each.value.name
  managed_by = "opentofu"

  tags = merge(local.base_tags, {
    workload      = "network-ops"
    purpose       = "network-ops"
    region        = each.value.name
    region_role   = each.key
    data_boundary = each.value.data_boundary
  })
}

resource "azurerm_management_lock" "netops_rg" {
  for_each = local.regions[var.region_pair]

  name       = "lock-${azurerm_resource_group.netops[each.key].name}-cannot-delete"
  scope      = azurerm_resource_group.netops[each.key].id
  lock_level = "CanNotDelete"
}

resource "azurerm_network_watcher" "nw" {
  for_each = local.regions[var.region_pair]

  name                = "nw-connectivity-${var.env}-${each.value.short_name}"
  resource_group_name = azurerm_resource_group.netops[each.key].name
  location            = azurerm_resource_group.netops[each.key].location

  tags = merge(azurerm_resource_group.netops[each.key].tags, {
    purpose = "network-monitoring"
  })
}
