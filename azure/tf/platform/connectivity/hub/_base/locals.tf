locals {
  base_tags = merge(local.root_tags, {
    env        = var.env
    owner      = "platform-engineering"
    costcenter = "cc-it-platform"
    managed_by = "opentofu"
  })

  # ============================================================================
  # Hub Subnet Calculation
  # Dynamically calculate subnet CIDRs from hub VNet CIDR using cidrsubnet()
  # Hub CIDR: /19 (8,192 IPs)
  #
  # Layout (example for 10.0.0.0/19):
  #   10.0.0.0/26   - AzureFirewallSubnet (64 IPs) - cidrsubnet(hub, 7, 0)
  #   10.0.0.64/27  - GatewaySubnet (32 IPs)       - cidrsubnet(hub, 8, 2)
  #   10.0.0.96/28  - DNS Resolver Inbound (16 IPs) - cidrsubnet(hub, 9, 6)
  #   10.0.0.112/28 - DNS Resolver Outbound (16 IPs) - cidrsubnet(hub, 9, 7)
  #   10.0.0.128/27 - EPA Connectors (32 IPs)      - cidrsubnet(hub, 8, 4)
  #   10.0.0.160/27 - Private Endpoints (32 IPs)   - cidrsubnet(hub, 8, 5)
  #   10.0.0.192/26 - Reserved for future use
  # ============================================================================
  hub_subnets = {
    for role in ["active", "standby"] : role => {
      hub_cidr    = local.net_cidr[var.env][var.region_pair][role].hub
      firewall    = cidrsubnet(local.net_cidr[var.env][var.region_pair][role].hub, 7, 0)              # /26 (64 IPs)
      firewall_ip = cidrhost(cidrsubnet(local.net_cidr[var.env][var.region_pair][role].hub, 7, 0), 4) # .4 is first usable
      gateway     = cidrsubnet(local.net_cidr[var.env][var.region_pair][role].hub, 8, 2)              # /27 (32 IPs)
      dns_in      = cidrsubnet(local.net_cidr[var.env][var.region_pair][role].hub, 9, 6)              # /28 (16 IPs)
      dns_out     = cidrsubnet(local.net_cidr[var.env][var.region_pair][role].hub, 9, 7)              # /28 (16 IPs)
      epa         = cidrsubnet(local.net_cidr[var.env][var.region_pair][role].hub, 8, 4)              # /27 (32 IPs)
      pep         = cidrsubnet(local.net_cidr[var.env][var.region_pair][role].hub, 8, 5)              # /27 (32 IPs)
    }
  }
}
