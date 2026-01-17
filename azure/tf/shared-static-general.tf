# ============================================================================
# SHARED CONFIGURATION
# ============================================================================
# To be symlinked into a workspace

locals {
  # Two region groups based on regulatory intent:
  # - global: non-GDPR / general workloads
  # - gdpr: GDPR / EU-style residency workloads
  regions = {
    global = {
      active = {
        name          = "eastus2" # Virginia
        display_name  = "East US 2"
        short_name    = "eus2"
        data_boundary = "global"
        region_role   = "active"
      }
      standby = {
        name          = "centralus" # Iowa
        display_name  = "Central US"
        short_name    = "cus"
        data_boundary = "global"
        region_role   = "standby"
      }
    }

    gdpr = {
      active = {
        name          = "westeurope" # Netherlands
        display_name  = "West Europe"
        short_name    = "weu"
        data_boundary = "gdpr"
        region_role   = "active"
      }
      standby = {
        name          = "northeurope" # Ireland
        display_name  = "North Europe"
        short_name    = "neu"
        data_boundary = "gdpr"
        region_role   = "standby"
      }
    }
  }

  # Flatten to a list so we can derive everything else without repeating paths.
  region_objects = flatten([
    for scope_key, scope in local.regions : [
      for role_key, r in scope : merge(r, {
        regulatory_scope = scope_key # "global" | "gdpr"
        role             = role_key  # "active" | "standby"
      })
    ]
  ])

  # All allowed Azure locations
  allowed_regions = [for r in local.region_objects : r.name]

  # Map: location -> shorthand (e.g., "eastus2" -> "eus2")
  region_shorthand_lookup = {
    for r in local.region_objects : r.name => r.short_name
  }
}
