# =============================================================================
# Subscription-Vending Specific Variables
# =============================================================================
# Note: Shared variables (subscription mappings, tenant ID, etc.) are in shared.tf

variable "sp_gha_tf_apply_platform_obj_id" {
  description = "Service Principal ID for platform subscriptions (management, identity, connectivity, security)"
  type        = string
  sensitive   = true

  validation {
    condition     = can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", var.sp_gha_tf_apply_platform_obj_id))
    error_message = "sp_gha_tf_apply_platform_obj_id must be a valid UUID (Service Principal Object ID)"
  }
}

variable "sp_gha_tf_apply_landingzone_obj_id" {
  description = "Service Principal Object ID for landing zone subscriptions (corp, online)"
  type        = string
  sensitive   = true

  validation {
    condition     = can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", var.sp_gha_tf_apply_landingzone_obj_id))
    error_message = "sp_gha_tf_apply_landingzone_obj_id must be a valid UUID (Service Principal Object ID)"
  }
}
