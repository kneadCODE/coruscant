# =============================================================================
# LOGS OUTPUTS
# =============================================================================

output "logs_rg_name" {
  description = "Name of the logs resource group"
  value       = module.management.logs_rg_name
}

output "logs_st_archive_name" {
  description = "Name of the logs archive storage account"
  value       = module.management.logs_st_archive_name
  sensitive   = true
}

output "logs_law_name" {
  description = "Name of the Log Analytics Workspace"
  value       = module.management.logs_law_name
}

# =============================================================================
# BACKUP OUTPUTS
# =============================================================================

output "backup_rg_name" {
  description = "Name of the backup resource group"
  value       = module.management.backup_rg_name
}

output "backup_rsv_name" {
  description = "Name of the Recovery Services Vault"
  value       = module.management.backup_rsv_name
}

output "backup_bvault_name" {
  description = "Name of the Backup Vault"
  value       = module.management.backup_bvault_name
}
