# ============================================================================
# Pass-through outputs from _base module
# ============================================================================
output "hub_vnet_names" {
  description = "Map of region role (active/standby) to hub VNet names"
  value       = module.main.hub_vnet_names
}

output "hub_firewall_private_ips" {
  description = "Map of region role (active/standby) to Azure Firewall private IPs"
  value       = module.main.hub_firewall_private_ips
}
