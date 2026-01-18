# =============================================================================
# SIEM OUTPUTS
# =============================================================================

output "siem_rg_name" {
  description = "Name of the siem resource group"
  value       = module.main.siem_rg_name
}

output "siem_st_archive_name" {
  description = "Name of the siem archive storage account"
  value       = module.main.siem_st_archive_name
  sensitive   = true
}

output "siem_law_name" {
  description = "Name of the Log Analytics Workspace"
  value       = module.main.siem_law_name
}
