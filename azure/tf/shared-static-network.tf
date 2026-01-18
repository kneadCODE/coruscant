# To symlink into a workspace:
#   cd <dir>
#   ln -s ../../shared-static-network.tf .

###### DESIGN ######
# Overall Address Space: "10.0.0.0/12" (10.0.0.0 - 10.15.255.255)
#
# Per Env per Region: /16 (65,536 IPs each)
# Hub VNET per Env per Region: /19 (8,192 IPs)
# Spoke VNET per Env per Region: /21 (2,048 IPs each)
#
# Allocation:
#   10.0.0.0/16  - prod/global/active   (East US 2)
#   10.1.0.0/16  - prod/global/standby  (Central US)
#   10.2.0.0/16  - prod/gdpr/active     (West Europe)
#   10.3.0.0/16  - prod/gdpr/standby    (North Europe)
#   10.4.0.0/16  - nonprod/global/active   (East US 2)
#   10.5.0.0/16  - nonprod/global/standby  (Central US)
#   10.6.0.0/16 - 10.15.0.0/16 - Reserved for future expansion

locals {
  net_cidr = {
    prod = {
      global = {
        active = {                       # 10.0.0.0/16 (65,536 IPs)
          hub           = "10.0.0.0/19"  # 8,192 IPs  - Hub VNet (Firewall, Gateway, DNS, etc.)
          breakglass    = "10.0.32.0/24" # 256 IPs    - Breakglass/Emergency access VNet
          security      = "10.0.40.0/21" # 2,048 IPs  - Security spoke (SIEM, vulnerability scanners, etc.)
          devops        = "10.0.48.0/21" # 2,048 IPs  - DevOps spoke (CI/CD, artifact repos, build agents)
          esb           = "10.0.56.0/21" # 2,048 IPs  - Enterprise Service Bus spoke (integration layer)
          observability = "10.0.64.0/21" # 2,048 IPs  - Observability spoke (monitoring, logging, metrics)
          edge          = "10.0.72.0/21" # 2,048 IPs  - Edge spoke (API gateways, load balancers, CDN)
          iam           = "10.0.80.0/21" # 2,048 IPs  - Identity & Access Management spoke (auth services)
          payment       = "10.0.88.0/21" # 2,048 IPs  - Payment spoke (payment processing, PCI-DSS workloads)
          # Reserved: 10.0.96.0 - 10.0.255.255 for future spokes
        }
        standby = {                      # 10.1.0.0/16 (65,536 IPs)
          hub           = "10.1.0.0/19"  # 8,192 IPs  - Hub VNet (Firewall, Gateway, DNS, etc.)
          breakglass    = "10.1.32.0/24" # 256 IPs    - Breakglass/Emergency access VNet (needed for DR)
          security      = "10.1.40.0/21" # 2,048 IPs  - Security spoke (SIEM, vulnerability scanners, etc.)
          devops        = "10.1.48.0/21" # 2,048 IPs  - DevOps spoke (warm standby for DR deployments)
          esb           = "10.1.56.0/21" # 2,048 IPs  - Enterprise Service Bus spoke (integration layer)
          observability = "10.1.64.0/21" # 2,048 IPs  - Observability spoke (monitoring, logging, metrics)
          edge          = "10.1.72.0/21" # 2,048 IPs  - Edge spoke (API gateways, load balancers, CDN)
          iam           = "10.1.80.0/21" # 2,048 IPs  - Identity & Access Management spoke (auth services)
          payment       = "10.1.88.0/21" # 2,048 IPs  - Payment spoke (payment processing, PCI-DSS workloads)
          # Reserved: 10.1.96.0 - 10.1.255.255 for future spokes
        }
      }
      gdpr = {
        active = {                       # 10.2.0.0/16 (65,536 IPs)
          hub           = "10.2.0.0/19"  # 8,192 IPs  - Hub VNet (Firewall, Gateway, DNS, etc.)
          breakglass    = "10.2.32.0/24" # 256 IPs    - Breakglass/Emergency access VNet
          security      = "10.2.40.0/21" # 2,048 IPs  - Security spoke (SIEM, vulnerability scanners, etc.)
          devops        = "10.2.48.0/21" # 2,048 IPs  - DevOps spoke (CI/CD, artifact repos, build agents)
          esb           = "10.2.56.0/21" # 2,048 IPs  - Enterprise Service Bus spoke (integration layer)
          observability = "10.2.64.0/21" # 2,048 IPs  - Observability spoke (monitoring, logging, metrics)
          edge          = "10.2.72.0/21" # 2,048 IPs  - Edge spoke (API gateways, load balancers, CDN)
          iam           = "10.2.80.0/21" # 2,048 IPs  - Identity & Access Management spoke (auth services)
          payment       = "10.2.88.0/21" # 2,048 IPs  - Payment spoke (payment processing, PCI-DSS workloads)
          # Reserved: 10.2.96.0 - 10.2.255.255 for future spokes
        }
        standby = {                      # 10.3.0.0/16 (65,536 IPs)
          hub           = "10.3.0.0/19"  # 8,192 IPs  - Hub VNet (Firewall, Gateway, DNS, etc.)
          breakglass    = "10.3.32.0/24" # 256 IPs    - Breakglass/Emergency access VNet (needed for DR)
          security      = "10.3.40.0/21" # 2,048 IPs  - Security spoke (SIEM, vulnerability scanners, etc.)
          devops        = "10.3.48.0/21" # 2,048 IPs  - DevOps spoke (warm standby for DR deployments)
          esb           = "10.3.56.0/21" # 2,048 IPs  - Enterprise Service Bus spoke (integration layer)
          observability = "10.3.64.0/21" # 2,048 IPs  - Observability spoke (monitoring, logging, metrics)
          edge          = "10.3.72.0/21" # 2,048 IPs  - Edge spoke (API gateways, load balancers, CDN)
          iam           = "10.3.80.0/21" # 2,048 IPs  - Identity & Access Management spoke (auth services)
          payment       = "10.3.88.0/21" # 2,048 IPs  - Payment spoke (payment processing, PCI-DSS workloads)
          # Reserved: 10.3.96.0 - 10.3.255.255 for future spokes
        }
      }
    }
    nonprod = {
      global = {
        active = {                       # 10.4.0.0/16 (65,536 IPs)
          hub           = "10.4.0.0/19"  # 8,192 IPs  - Hub VNet (Firewall, Gateway, DNS, etc.)
          breakglass    = "10.4.32.0/24" # 256 IPs    - Breakglass/Emergency access VNet
          security      = "10.4.40.0/21" # 2,048 IPs  - Security spoke (SIEM, vulnerability scanners, etc.)
          devops        = "10.4.48.0/21" # 2,048 IPs  - DevOps spoke (CI/CD, artifact repos, build agents)
          esb           = "10.4.56.0/21" # 2,048 IPs  - Enterprise Service Bus spoke (integration layer)
          observability = "10.4.64.0/21" # 2,048 IPs  - Observability spoke (monitoring, logging, metrics)
          edge          = "10.4.72.0/21" # 2,048 IPs  - Edge spoke (API gateways, load balancers, CDN)
          iam           = "10.4.80.0/21" # 2,048 IPs  - Identity & Access Management spoke (auth services)
          payment       = "10.4.88.0/21" # 2,048 IPs  - Payment spoke (payment processing, PCI-DSS workloads)
          # Reserved: 10.4.96.0 - 10.4.255.255 for future spokes
        }
        standby = {                      # 10.5.0.0/16 (65,536 IPs)
          hub           = "10.5.0.0/19"  # 8,192 IPs  - Hub VNet (Firewall, Gateway, DNS, etc.)
          breakglass    = "10.5.32.0/24" # 256 IPs    - Breakglass/Emergency access VNet (needed for DR)
          security      = "10.5.40.0/21" # 2,048 IPs  - Security spoke (SIEM, vulnerability scanners, etc.)
          devops        = "10.5.48.0/21" # 2,048 IPs  - DevOps spoke (warm standby for DR deployments)
          esb           = "10.5.56.0/21" # 2,048 IPs  - Enterprise Service Bus spoke (integration layer)
          observability = "10.5.64.0/21" # 2,048 IPs  - Observability spoke (monitoring, logging, metrics)
          edge          = "10.5.72.0/21" # 2,048 IPs  - Edge spoke (API gateways, load balancers, CDN)
          iam           = "10.5.80.0/21" # 2,048 IPs  - Identity & Access Management spoke (auth services)
          payment       = "10.5.88.0/21" # 2,048 IPs  - Payment spoke (payment processing, PCI-DSS workloads)
          # Reserved: 10.5.96.0 - 10.5.255.255 for future spokes
        }
      }
    }
  }
}
