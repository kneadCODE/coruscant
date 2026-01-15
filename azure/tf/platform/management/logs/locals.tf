locals {
  # Define regions to deploy platform logs infrastructure
  # Using only active regions from each data boundary
  platform_logs_regions = {
    r1 = local.regions.global.active
    r2 = local.regions.gdpr.active
  }
}
