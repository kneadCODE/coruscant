# ============================================================================
# SHARED VARIABLE DECLARATIONS
# ============================================================================
# This file contains variable declarations that will be shared across workspaces.
# Symlink this file into each workspace to avoid duplication.
#
# Values provided via:
# - GitHub Actions: TF_VAR_* environment variables (set at job level)
# - Local dev: terraform.tfvars.local (gitignored, per workspace)
#
# To symlink into a workspace:
#   cd infrastructure/azure/tf/governance
#   ln -s ../shared-variables.tf .

variable "entra_tenant_id" {
  description = "Entra Tenant ID"
  type        = string
  sensitive   = true
}

variable "subscription_id_foundation" {
  description = "Foundation subscription ID"
  type        = string
  sensitive   = true
}

variable "subscription_id_management" {
  description = "Management subscription ID"
  type        = string
  sensitive   = true
}

variable "subscription_id_identity" {
  description = "Identity subscription ID"
  type        = string
  sensitive   = true
}

variable "subscription_id_connectivity_prod" {
  description = "Connectivity Prod subscription ID"
  type        = string
  sensitive   = true
}

variable "subscription_id_security_prod" {
  description = "Security Prod subscription ID"
  type        = string
  sensitive   = true
}
