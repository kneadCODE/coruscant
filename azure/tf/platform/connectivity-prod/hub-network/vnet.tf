resource "azurerm_resource_group" "hub_net_sea" {
  name       = "rg-hub-net-sea"
  location   = "southeastasia"
  managed_by = "iac"

  tags = {
    org        = "kneadcode"
    portfolio  = "coruscant"
    workload   = "hub-network"
    env        = "prod"
    purpose    = "hub-network"
    owner      = "platform-engineering"
    costcenter = "cc-it-platform"
    managed_by = "iac"
    iac_tool   = "opentofu"
  }
}

module "hub_vnet_sea" {
  source = "../../../modules/vnet"

  address_space = local.net_cidr.sea.prod.hub
  env           = "prod"
  purpose       = "hub"
  region        = azurerm_resource_group.hub_net_sea.location
  rg_name       = azurerm_resource_group.hub_net_sea.name

  # ============================================================================
  # Route Tables - Centralized definition
  # ============================================================================
  route_tables = {
    "force-firewall" = {
      bgp_route_propagation_enabled = false # Disable BGP - we control all routes explicitly
      routes = {
        "default-to-firewall" = {
          address_prefix         = "0.0.0.0/0"
          next_hop_type          = "VirtualAppliance"
          next_hop_in_ip_address = "10.0.0.4" # Azure Firewall private IP (first usable IP in /26 is .4)
        }
        "rfc1918-10-to-firewall" = {
          address_prefix         = "10.0.0.0/8"
          next_hop_type          = "VirtualAppliance"
          next_hop_in_ip_address = "10.0.0.4" # Force all internal 10.x traffic through firewall
        }
      }
    }
  }

  # ============================================================================
  # Subnets - Hub Network
  # ============================================================================
  subnets = {
    # ============================================================================
    # Azure Firewall Subnet - NO NSG, NO UDR (Azure Firewall manages its own)
    # ============================================================================
    "AzureFirewallSubnet" = {
      address_prefix = "10.0.0.0/26" # /26 = 64 IPs (Azure requires minimum /26)
    }

    # ============================================================================
    # VPN/ExpressRoute Gateway Subnet - NO NSG, NO UDR (Gateway manages its own)
    # ============================================================================
    "GatewaySubnet" = {
      address_prefix = "10.0.0.64/27" # /27 = 32 IPs (sufficient for gateway HA instances)
    }

    # ============================================================================
    # DNS Private Resolver Inbound Endpoint Subnet
    # ============================================================================
    "snet-dnspr-in" = {
      address_prefix = "10.0.0.96/28" # /28 = 16 IPs
      delegations = [
        {
          name = "dns-resolver-delegation"
          service_delegation = {
            name = "Microsoft.Network/dnsResolvers"
            actions = [
              "Microsoft.Network/virtualNetworks/subnets/join/action",
            ]
          }
        },
      ]
    }

    # ============================================================================
    # DNS Private Resolver Outbound Endpoint Subnet
    # ============================================================================
    "snet-dnspr-out" = {
      address_prefix = "10.0.0.112/28" # /28 = 16 IPs
      delegations = [
        {
          name = "dns-resolver-delegation"
          service_delegation = {
            name = "Microsoft.Network/dnsResolvers"
            actions = [
              "Microsoft.Network/virtualNetworks/subnets/join/action",
            ]
          }
        },
      ]
    }

    # ============================================================================
    # Entra Private Access (EPA) Connectors Subnet
    # For Entra Private Access app proxy connectors (formerly Azure AD App Proxy)
    # ============================================================================
    "snet-epa-connectors" = {
      address_prefix = "10.0.0.128/27" # /27 = 32 IPs for EPA connector VMs

      # Inbound: Only allow breakglass emergency management
      inbound_rules = {
        "allow-breakglass-mgmt" = {
          priority                     = 100
          access                       = "Allow"
          protocol                     = "Tcp"
          source_port_ranges           = ["*"]
          source_address_prefixes      = [local.net_cidr.sea.prod.breakglass]
          destination_port_ranges      = ["22", "3389"]    # SSH/RDP for connector VM management
          destination_address_prefixes = ["10.0.0.128/27"] # This subnet only
          description                  = "Allow breakglass emergency SSH/RDP to EPA connector VMs"
        }
      }

      # Outbound: EPA connectors need to reach Entra (Internet) and internal apps (via firewall)
      outbound_rules = {
        "allow-entra-services" = {
          priority                     = 100
          access                       = "Allow"
          protocol                     = "Tcp"
          source_port_ranges           = ["*"]
          source_address_prefixes      = ["10.0.0.128/27"]
          destination_port_ranges      = ["443", "80"]
          destination_address_prefixes = ["AzureActiveDirectory"] # Service tag for Entra
          description                  = "Allow EPA connectors to reach Entra/Azure AD services"
        }

        "allow-to-firewall" = {
          priority                     = 110
          access                       = "Allow"
          protocol                     = "*"
          source_port_ranges           = ["*"]
          source_address_prefixes      = ["10.0.0.128/27"]
          destination_port_ranges      = ["*"]
          destination_address_prefixes = ["10.0.0.0/26"] # Azure Firewall subnet
          description                  = "Allow all traffic to firewall for inspection/routing to internal apps"
        }
      }

      route_table_key = "force-firewall" # Force all non-Entra traffic through firewall
    }

    # ============================================================================
    # Private Endpoints Subnet
    # For Azure PaaS services (Storage, SQL, Key Vault, etc.)
    # ============================================================================
    "snet-pe-services" = {
      address_prefix                            = "10.0.0.160/27" # /27 = 32 IPs for private endpoints
      private_endpoint_network_policies_enabled = false           # REQUIRED for private endpoints

      # CRITICAL: NO route table on Private Endpoints subnet!
      # Private Endpoints use Azure Private Link (Azure backbone routing)
      # Route tables can break Private Endpoint connectivity

      # NSG rules provide defense in depth
      # Actual enforcement via: (1) Firewall rules, (2) Route tables on spokes, (3) NSG here to prevent hub bypass

      # Inbound: Allow ONLY from Azure Firewall (all spoke traffic flows through FW)
      inbound_rules = {
        "allow-firewall-to-pe" = {
          priority                     = 100
          access                       = "Allow"
          protocol                     = "*"
          source_port_ranges           = ["*"]
          source_address_prefixes      = ["10.0.0.0/26"] # Azure Firewall subnet ONLY
          destination_port_ranges      = ["*"]
          destination_address_prefixes = ["10.0.0.160/27"] # This subnet
          description                  = "Allow traffic from firewall to private endpoints (spoke traffic)"
        }

        "allow-breakglass-to-pe" = {
          priority                     = 110
          access                       = "Allow"
          protocol                     = "*"
          source_port_ranges           = ["*"]
          source_address_prefixes      = [local.net_cidr.sea.prod.breakglass]
          destination_port_ranges      = ["*"]
          destination_address_prefixes = ["10.0.0.160/27"]
          description                  = "Allow breakglass emergency access to private endpoints"
        }

        # DENY all other hub subnets (Gateway, EPA, DNS) from reaching PE directly
        "deny-hub-to-pe" = {
          priority           = 4090
          access             = "Deny"
          protocol           = "*"
          source_port_ranges = ["*"]
          source_address_prefixes = [
            "10.0.0.64/27",  # GatewaySubnet - block direct access
            "10.0.0.96/28",  # DNS Resolver inbound - block direct access
            "10.0.0.112/28", # DNS Resolver outbound - block direct access
            "10.0.0.128/27", # EPA connectors - must go through firewall
          ]
          destination_port_ranges      = ["*"]
          destination_address_prefixes = ["10.0.0.160/27"]
          description                  = "Deny hub subnets from bypassing firewall to reach PE"
        }
      }

      # Outbound: Private endpoints typically don't initiate connections
      outbound_rules = {
        "allow-pe-responses" = {
          priority                     = 100
          access                       = "Allow"
          protocol                     = "*"
          source_port_ranges           = ["*"]
          source_address_prefixes      = ["10.0.0.160/27"]
          destination_port_ranges      = ["*"]
          destination_address_prefixes = ["10.0.0.0/26", local.net_cidr.sea.prod.breakglass]
          description                  = "Allow responses to firewall and breakglass"
        }
      }

      # NO route_table_key - Private Endpoints use Azure backbone, not subnet routing
    }
  }

  tags = azurerm_resource_group.hub_net_sea.tags
}
