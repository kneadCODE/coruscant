resource "azurerm_resource_group" "netops_sea" {
  name       = "rg-connectivity-netops-prod-sea"
  location   = "southeastasia"
  managed_by = "iac"

  tags = {
    org        = "kneadcode"
    portfolio  = "coruscant"
    workload   = "hub-network"
    env        = "prod"
    purpose    = "hub-network"
    owner      = "platform-engineering"
    costcenter = "cc-it-platform"
    managed_by = "iac"
    iac_tool   = "opentofu"
  }
}

resource "azurerm_network_watcher" "sea_01" {
  name                = "nw-connectivity-prod-sea-01"
  resource_group_name = azurerm_resource_group.netops_sea.name
  location            = azurerm_resource_group.netops_sea.location
  tags = merge(azurerm_resource_group.netops_sea.tags, {
    "purpose" = "network-monitoring"
  })

  depends_on = [azurerm_resource_group.netops_sea]
}
