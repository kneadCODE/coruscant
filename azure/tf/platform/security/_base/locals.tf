locals {
  # Derive active region from the region pair (all _base resources deploy to active region only)
  active_region       = local.regions[var.region_pair].active.name
  active_region_short = local.regions[var.region_pair].active.short_name

  base_tags = merge(local.root_tags, {
    env              = var.env
    compliance_scope = var.region_pair # region_pair maps directly to compliance scope (global/gdpr)
    owner            = "platform-engineering"
    costcenter       = "cc-it-platform"
    managed_by       = "opentofu"
  })
}
