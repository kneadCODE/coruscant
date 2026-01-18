# =============================================================================
# LOGS OUTPUTS
# =============================================================================

output "logs_rg_name" {
  description = "Name of the logs resource group"
  value       = (var.logs_deploy_st || var.logs_deploy_law) ? azurerm_resource_group.logs[0].name : null
}

output "logs_st_archive_name" {
  description = "Name of the logs archive storage account"
  value       = var.logs_deploy_st ? azurerm_storage_account.logs_archive_1[0].name : null
  sensitive   = true
}

output "logs_law_name" {
  description = "Name of the Log Analytics Workspace"
  value       = var.logs_deploy_law ? azurerm_log_analytics_workspace.logs_1[0].name : null
}

# =============================================================================
# BACKUP OUTPUTS
# =============================================================================

output "backup_rg_name" {
  description = "Name of the backup resource group"
  value       = (var.backup_deploy_rsv || var.backup_deploy_bvault) ? azurerm_resource_group.backup[0].name : null
}

output "backup_rsv_name" {
  description = "Name of the Recovery Services Vault"
  value       = var.backup_deploy_rsv ? azurerm_recovery_services_vault.platform_1[0].name : null
}

output "backup_bvault_name" {
  description = "Name of the Backup Vault"
  value       = var.backup_deploy_bvault ? azurerm_data_protection_backup_vault.platform_compliance_1[0].name : null
}
