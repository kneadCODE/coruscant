# ============================================================================
# Hub Outputs - Only what's needed for cross-state references
# ============================================================================

output "hub_vnet_ids" {
  description = "Map of region role (active/standby) to hub VNet IDs"
  value = {
    for role, vnet in module.hub_vnet : role => vnet.vnet_id
  }
}

output "hub_vnet_names" {
  description = "Map of region role (active/standby) to hub VNet names"
  value = {
    for role, vnet in module.hub_vnet : role => vnet.vnet_name
  }
}

output "hub_firewall_ips" {
  description = "Map of region role (active/standby) to Azure Firewall private IPs"
  value = {
    for role in keys(local.hub_subnets) : role => local.hub_subnets[role].firewall_ip
  }
}
