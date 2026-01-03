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
  description = "Foundation subscription ID (used for provider and state backend)"
  type        = string
  sensitive   = true
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
