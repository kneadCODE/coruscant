# To symlink into a workspace:
#   cd <dir>
#   ln -s ../../shared-static-network.tf .

###### DESIGN ######
# Overall Address Space: "10.0.0.0/12" (10.0.0.0 - 10.15.255.255)
#
# Per Region: /14 (262,144 IPs each)
# SEA (Primary Region): 10.0.0.0/14   (10.0.0.0 - 10.3.255.255)
# Region 2: 10.4.0.0/14                (10.4.0.0 - 10.7.255.255)
# Region 3: 10.8.0.0/14                (10.8.0.0 - 10.11.255.255)
# Region 4: 10.12.0.0/14               (10.12.0.0 - 10.15.255.255)
#
# Per Env per region: /16 (65,536 IPs each)
# Hub VNET per Env per Region: /19 (8,192 IPs)
# Spoke VNET per Env per Region: /21 (2,048 IPs each)

locals {
  net_cidr = {
    sea = {                            # 10.0.0.0/14 (Southeast Asia region)
      prod = {                         # 10.0.0.0/16 (65,536 IPs for prod environment)
        hub           = "10.0.0.0/19"  # 8,192 IPs  - Hub VNet (Firewall, Gateway, DNS, etc.)
        breakglass    = "10.0.32.0/24" # 256 IPs    - Breakglass/Emergency access VNet (SEA prod only)
        security      = "10.0.40.0/21" # 2,048 IPs  - Security spoke (SIEM, vulnerability scanners, etc.)
        devops        = "10.0.48.0/21" # 2,048 IPs  - DevOps spoke (CI/CD, artifact repos, build agents) (SEA only)
        esb           = "10.0.56.0/21" # 2,048 IPs  - Enterprise Service Bus spoke (integration layer)
        observability = "10.0.64.0/21" # 2,048 IPs  - Observability spoke (monitoring, logging, metrics)
        edge          = "10.0.72.0/21" # 2,048 IPs  - Edge spoke (API gateways, load balancers, CDN)
        iam           = "10.0.80.0/21" # 2,048 IPs  - Identity & Access Management spoke (auth services)
        payment       = "10.0.88.0/21" # 2,048 IPs  - Payment spoke (payment processing, PCI-DSS workloads)
        # Reserved: 10.0.96.0/19 through 10.0.255.255 for future prod spokes
      }
      nonprod = {           # 10.1.0.0/16 (65,536 IPs for nonprod environment)
        hub = "10.1.0.0/19" # 8,192 IPs  - Hub VNet (Firewall, Gateway, DNS, etc.)
        # No breakglass in nonprod - use prod breakglass if needed
        security      = "10.1.40.0/21" # 2,048 IPs  - Security spoke
        devops        = "10.1.48.0/21" # 2,048 IPs  - DevOps spoke (SEA only)
        esb           = "10.1.56.0/21" # 2,048 IPs  - Enterprise Service Bus spoke
        observability = "10.1.64.0/21" # 2,048 IPs  - Observability spoke
        edge          = "10.1.72.0/21" # 2,048 IPs  - Edge spoke
        iam           = "10.1.80.0/21" # 2,048 IPs  - Identity & Access Management spoke
        payment       = "10.1.88.0/21" # 2,048 IPs  - Payment spoke
        # Reserved: 10.1.96.0/19 through 10.1.255.255 for future nonprod spokes
      }
    }
    # Future regions can be added here:
    # region2 = { prod = { hub = "10.4.0.0/19", ... }, nonprod = { ... } }
    # region3 = { prod = { hub = "10.8.0.0/19", ... }, nonprod = { ... } }
    # region4 = { prod = { hub = "10.12.0.0/19", ... }, nonprod = { ... } }
  }
}
