# =============================================================================
# SIEM OUTPUTS
# =============================================================================

output "siem_rg_name" {
  description = "Name of the siem resource group"
  value       = (var.siem_deploy_st || var.siem_deploy_law) ? azurerm_resource_group.siem[0].name : null
}

output "siem_st_archive_name" {
  description = "Name of the siem archive storage account"
  value       = var.siem_deploy_st ? azurerm_storage_account.siem_archive_1[0].name : null
  sensitive   = true
}

output "siem_law_name" {
  description = "Name of the Log Analytics Workspace"
  value       = var.siem_deploy_law ? azurerm_log_analytics_workspace.siem_1[0].name : null
}
