# ============================================================================
# SHARED CONFIGURATION
# ============================================================================
# This file contains both variable declarations and common locals that are
# shared across workspaces. Symlink this file into each workspace to avoid
# duplication.
#
# Variable values provided via:
# - GitHub Actions: TF_VAR_* environment variables (set at job level)
# - Local dev: terraform.tfvars.local (gitignored, per workspace)
#
# To symlink into a workspace:
#   cd infrastructure/azure/tf/management
#   ln -s ../shared-vars.tf .

# ============================================================================
# SHARED VARIABLES
# ============================================================================

variable "subscription_id_foundation" {
  description = "Foundation subscription ID (used for provider and state backend)"
  type        = string
  sensitive   = true

  # Validation: Must be a valid UUID (case-insensitive)
  validation {
    condition     = can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", var.subscription_id_foundation))
    error_message = "subscription_id_foundation must be a valid UUID in format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (lowercase or uppercase)"
  }
}

variable "subscription_id_mapping_json" {
  description = <<-EOT
    JSON map of all subscription IDs in the Azure landing zone hierarchy.

    Expected format:
    {
      "management": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "identity": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "connectivity_prod": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "connectivity_nonprod": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    }
  EOT
  type        = string
  sensitive   = true

  # Validation 1: Must be valid JSON
  validation {
    condition     = can(jsondecode(var.subscription_id_mapping_json))
    error_message = "subscription_id_mapping_json must be valid JSON"
  }

  # Validation 2: Required subscription keys must exist
  validation {
    condition = alltrue([
      for required_key in [
        "management",
        "identity",
        "connectivity_prod",
        "connectivity_nonprod",
        "security_prod",
        "security_nonprod",
        "devops_prod",
        "devops_nonprod",
        "esb_prod",
        "esb_nonprod",
        "observability_prod",
        "observability_nonprod",
        "edge_prod",
        "edge_nonprod",
        "iam_prod",
        "iam_nonprod",
        "payment_prod",
        "payment_nonprod"
      ] : can(lookup(jsondecode(var.subscription_id_mapping_json), required_key))
    ])
    error_message = "subscription_id_mapping_json is missing required keys"
  }

  # Validation 3: All subscription IDs must be valid UUIDs (case-insensitive)
  validation {
    condition = alltrue([
      for key, value in jsondecode(var.subscription_id_mapping_json) :
      can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", value))
    ])
    error_message = "All subscription IDs must be valid UUIDs in format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (lowercase or uppercase)"
  }

  # Validation 4: No duplicate subscription IDs
  validation {
    condition     = length(values(jsondecode(var.subscription_id_mapping_json))) == length(distinct(values(jsondecode(var.subscription_id_mapping_json))))
    error_message = "Duplicate subscription IDs detected. Each subscription must be unique."
  }
}

variable "storage_account_prefix" {
  description = "Prefix for storage account names (must be lowercase alphanumeric, 3-10 chars)"
  type        = string
  sensitive   = true

  validation {
    condition     = can(regex("^[a-z0-9]{3,10}$", var.storage_account_prefix))
    error_message = "storage_account_prefix must be 3-10 characters, lowercase letters and numbers only"
  }
}

# ============================================================================
# SHARED LOCALS
# ============================================================================
# Common parsing logic for JSON-based mappings. These locals are available
# in any workspace that symlinks this file.

locals {
  # Parse subscription mapping from JSON secret
  # Expected format: {"management": "guid", "identity": "guid", ...}
  subscription_ids = jsondecode(var.subscription_id_mapping_json)
}
