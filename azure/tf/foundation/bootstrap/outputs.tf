# =============================================================================
# Outputs
# =============================================================================

# Level 1 MG IDs
output "mg_platform_id" {
  description = "Management group ID for Platform"
  value       = azurerm_management_group.platform.id
}

output "mg_landingzone_id" {
  description = "Management group ID for Landing Zones"
  value       = azurerm_management_group.landingzone.id
}

output "mg_sandbox_id" {
  description = "Management group ID for Sandbox"
  value       = azurerm_management_group.sandbox.id
}

output "mg_decommissioned_id" {
  description = "Management group ID for Decommissioned"
  value       = azurerm_management_group.decommissioned.id
}

# Level 2 MG IDs - Platform
output "mg_management_id" {
  description = "Management group ID for Management"
  value       = azurerm_management_group.management.id
}

output "mg_identity_id" {
  description = "Management group ID for Identity"
  value       = azurerm_management_group.identity.id
}

output "mg_connectivity_id" {
  description = "Management group ID for Connectivity"
  value       = azurerm_management_group.connectivity.id
}

output "mg_security_id" {
  description = "Management group ID for Security"
  value       = azurerm_management_group.security.id
}

# Level 2 MG IDs - Landing Zones
output "mg_corp_id" {
  description = "Management group ID for Corp"
  value       = azurerm_management_group.corp.id
}

output "mg_online_id" {
  description = "Management group ID for Online"
  value       = azurerm_management_group.online.id
}
