variable "env" {
  type        = string
  description = "The environment (prod or nonprod)"
  validation {
    condition     = contains(local.allowed_envs, var.env)
    error_message = "Invalid environment. Must be one of: ${join(", ", local.allowed_envs)}"
  }
}

variable "region_pair" {
  type        = string
  description = "The region pair to deploy resources in (global or gdpr)"
  validation {
    condition     = contains(keys(local.regions), var.region_pair)
    error_message = "Invalid region pair. Must be one of: ${join(", ", keys(local.regions))}"
  }
}

variable "region_role" {
  type        = string
  description = "Which region within the pair to deploy to (active or standby)"
  default     = "active"
  validation {
    condition     = contains(["active", "standby"], var.region_role)
    error_message = "Invalid region role. Must be 'active' or 'standby'"
  }
}

variable "siem_deploy_st" {
  type        = bool
  description = "Deploy SIEM storage account for archival?"
  default     = false
}

variable "siem_deploy_law" {
  type        = bool
  description = "Deploy Log Analytics Workspace for hot SIEM logs?"
  default     = false
}

variable "storage_account_prefix" {
  type        = string
  description = "Prefix for storage account names (must be lowercase alphanumeric, 3-10 chars)"
  sensitive   = true
  validation {
    condition     = can(regex("^[a-z0-9]{3,10}$", var.storage_account_prefix))
    error_message = "storage_account_prefix must be 3-10 characters, lowercase letters and numbers only"
  }
}
