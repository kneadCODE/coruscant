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
#   ln -s ../shared.tf .

# ============================================================================
# SHARED VARIABLES
# ============================================================================

variable "entra_tenant_id" {
  description = "Entra Tenant ID (Azure AD/Entra ID tenant GUID)"
  type        = string
  sensitive   = true

  # Validation: Must be a valid UUID (case-insensitive)
  validation {
    condition     = can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", var.entra_tenant_id))
    error_message = "entra_tenant_id must be a valid UUID in format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (lowercase or uppercase)"
  }
}

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
        # "observability_prod",
        # "observability_nonprod",
        "edge_prod",
        "edge_nonprod",
        # "iam_prod",
        # "iam_nonprod"
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

variable "storage_account_name_mapping_json" {
  description = <<-EOT
    JSON map of storage account names used across the infrastructure.
    Storage account names alone are not secrets, but since the repo is public,
    we avoid committing any identifiable information.

    Expected format:
    {
      "st_platform_logs_archive_sea_01": "stxxxxxxxx01",
      "st_platform_logs_archive_sea_02": "stxxxxxxxx02"
    }

    Storage account naming constraints enforced:
    - 3-24 characters
    - Lowercase alphanumeric only
    - Must be globally unique across Azure
  EOT
  type        = string
  sensitive   = true

  # Validation 1: Must be valid JSON
  validation {
    condition     = can(jsondecode(var.storage_account_name_mapping_json))
    error_message = "storage_account_name_mapping_json must be valid JSON"
  }

  # Validation 2: All storage account names must meet Azure naming constraints
  # 3-24 chars, lowercase alphanumeric only
  validation {
    condition = alltrue([
      for key, value in jsondecode(var.storage_account_name_mapping_json) :
      can(regex("^[a-z0-9]{3,24}$", value))
    ])
    error_message = "All storage account names must be 3-24 characters, lowercase letters and numbers only"
  }

  # Validation 3: No duplicate storage account names
  validation {
    condition     = length(values(jsondecode(var.storage_account_name_mapping_json))) == length(distinct(values(jsondecode(var.storage_account_name_mapping_json))))
    error_message = "Duplicate storage account names detected. Each storage account name must be unique."
  }

  # Validation 4: Required storage account keys must exist
  validation {
    condition = alltrue([
      for required_key in [
        "st_platform_logs_archive_sea_01",
      ] : can(lookup(jsondecode(var.storage_account_name_mapping_json), required_key))
    ])
    error_message = "storage_account_name_mapping_json is missing required keys"
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

  # Parse storage account name mapping from JSON secret
  # Expected format: {"st_platform_logs_archive_sea_01": "stxxxxx01", ...}
  storage_account_names = jsondecode(var.storage_account_name_mapping_json)
}
