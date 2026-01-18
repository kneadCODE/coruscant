resource "azurerm_resource_group" "hub_net" {
  for_each = local.regions[var.region_pair]

  name       = "rg-hub-net-${var.env}-${each.value.short_name}"
  location   = each.value.name
  managed_by = "opentofu"

  tags = merge(local.base_tags, {
    workload      = "hub-network"
    purpose       = "hub-network"
    region        = each.value.name
    region_role   = each.key
    data_boundary = each.value.data_boundary
  })
}

resource "azurerm_management_lock" "hub_vnet_rg" {
  for_each = local.regions[var.region_pair]

  name       = "lock-${azurerm_resource_group.hub_net[each.key].name}-cannot-delete"
  scope      = azurerm_resource_group.hub_net[each.key].id
  lock_level = "CanNotDelete"
}

module "hub_vnet" {
  for_each = local.regions[var.region_pair]
  source   = "../../../../modules/vnet"

  address_space = local.net_cidr[var.env][var.region_pair][each.key].hub
  env           = var.env
  purpose       = "hub"
  region        = each.value.name
  rg_name       = azurerm_resource_group.hub_net[each.key].name

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
          next_hop_in_ip_address = local.hub_subnets[each.key].firewall_ip # Azure Firewall private IP
        }
        "rfc1918-10-to-firewall" = {
          address_prefix         = "10.0.0.0/8"
          next_hop_type          = "VirtualAppliance"
          next_hop_in_ip_address = local.hub_subnets[each.key].firewall_ip # Force all internal 10.x traffic through firewall
        }
      }
    }
  }

  # ============================================================================
  # Subnets - Hub Network
  # Dynamically calculated from hub CIDR using cidrsubnet()
  # Hub CIDR: /19 (8,192 IPs) - subnet layout below
  # ============================================================================
  subnets = {
    # ============================================================================
    # Azure Firewall Subnet - NO NSG, NO UDR (Azure Firewall manages its own)
    # cidrsubnet(hub/19, 7, 0) = first /26 (64 IPs) - Azure requires minimum /26
    # ============================================================================
    "AzureFirewallSubnet" = {
      address_prefix = local.hub_subnets[each.key].firewall
    }

    # ============================================================================
    # VPN/ExpressRoute Gateway Subnet - NO NSG, NO UDR (Gateway manages its own)
    # cidrsubnet(hub/19, 8, 2) = /27 (32 IPs) - sufficient for gateway HA instances
    # ============================================================================
    "GatewaySubnet" = {
      address_prefix = local.hub_subnets[each.key].gateway
    }

    # ============================================================================
    # DNS Private Resolver Inbound Endpoint Subnet
    # cidrsubnet(hub/19, 9, 6) = /28 (16 IPs)
    # ============================================================================
    "snet-dnspr-in" = {
      address_prefix = local.hub_subnets[each.key].dns_in
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
    # cidrsubnet(hub/19, 9, 7) = /28 (16 IPs)
    # ============================================================================
    "snet-dnspr-out" = {
      address_prefix = local.hub_subnets[each.key].dns_out
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
    # cidrsubnet(hub/19, 8, 4) = /27 (32 IPs)
    # ============================================================================
    "snet-epa-connectors" = {
      address_prefix = local.hub_subnets[each.key].epa

      # Inbound: Only allow breakglass emergency management
      inbound_rules = {
        "allow-breakglass-mgmt" = {
          priority                     = 100
          access                       = "Allow"
          protocol                     = "Tcp"
          source_port_range            = "*"
          source_address_prefixes      = [local.net_cidr[var.env][var.region_pair][each.key].breakglass]
          destination_port_ranges      = ["22", "3389"]
          destination_address_prefixes = [local.hub_subnets[each.key].epa]
          description                  = "Allow breakglass emergency SSH/RDP to EPA connector VMs"
        }
      }

      # Outbound: EPA connectors need to reach Entra (Internet) and internal apps (via firewall)
      outbound_rules = {
        "allow-entra-services" = {
          priority                   = 100
          access                     = "Allow"
          protocol                   = "Tcp"
          source_port_range          = "*"
          source_address_prefixes    = [local.hub_subnets[each.key].epa]
          destination_port_ranges    = ["443", "80"]
          destination_address_prefix = "AzureActiveDirectory"
          description                = "Allow EPA connectors to reach Entra/Azure AD services"
        }

        "allow-to-firewall" = {
          priority                     = 110
          access                       = "Allow"
          protocol                     = "*"
          source_port_range            = "*"
          source_address_prefixes      = [local.hub_subnets[each.key].epa]
          destination_port_range       = "*"
          destination_address_prefixes = [local.hub_subnets[each.key].firewall]
          description                  = "Allow all traffic to firewall for inspection/routing to internal apps"
        }
      }

      route_table_key = "force-firewall"
    }

    # ============================================================================
    # Private Endpoints Subnet
    # For Azure PaaS services (Storage, SQL, Key Vault, etc.)
    # cidrsubnet(hub/19, 8, 5) = /27 (32 IPs)
    # ============================================================================
    "snet-pep" = {
      address_prefix                            = local.hub_subnets[each.key].pep
      private_endpoint_network_policies_enabled = false

      # CRITICAL: NO route table on Private Endpoints subnet!
      # Private Endpoints use Azure Private Link (Azure backbone routing)

      # Inbound: Allow ONLY from Azure Firewall (all spoke traffic flows through FW)
      inbound_rules = {
        "allow-firewall-to-pep" = {
          priority                     = 100
          access                       = "Allow"
          protocol                     = "*"
          source_port_range            = "*"
          source_address_prefixes      = [local.hub_subnets[each.key].firewall]
          destination_port_range       = "*"
          destination_address_prefixes = [local.hub_subnets[each.key].pep]
          description                  = "Allow traffic from firewall to private endpoints (spoke traffic)"
        }

        "allow-breakglass-to-pep" = {
          priority                     = 110
          access                       = "Allow"
          protocol                     = "*"
          source_port_range            = "*"
          source_address_prefixes      = [local.net_cidr[var.env][var.region_pair][each.key].breakglass]
          destination_port_range       = "*"
          destination_address_prefixes = [local.hub_subnets[each.key].pep]
          description                  = "Allow breakglass emergency access to private endpoints"
        }

        # DENY all other hub subnets (Gateway, EPA, DNS) from reaching PEP directly
        "deny-hub-to-pep" = {
          priority          = 4090
          access            = "Deny"
          protocol          = "*"
          source_port_range = "*"
          source_address_prefixes = [
            local.hub_subnets[each.key].gateway,
            local.hub_subnets[each.key].dns_in,
            local.hub_subnets[each.key].dns_out,
            local.hub_subnets[each.key].epa,
          ]
          destination_port_range       = "*"
          destination_address_prefixes = [local.hub_subnets[each.key].pep]
          description                  = "Deny hub subnets from bypassing firewall to reach PEP"
        }
      }

      # Outbound: Private endpoints typically don't initiate connections
      outbound_rules = {
        "allow-pep-responses" = {
          priority                     = 100
          access                       = "Allow"
          protocol                     = "*"
          source_port_range            = "*"
          source_address_prefixes      = [local.hub_subnets[each.key].pep]
          destination_port_range       = "*"
          destination_address_prefixes = [local.hub_subnets[each.key].firewall, local.net_cidr[var.env][var.region_pair][each.key].breakglass]
          description                  = "Allow responses to firewall and breakglass"
        }
      }
    }
  }

  tags = azurerm_resource_group.hub_net[each.key].tags

  depends_on = [azurerm_network_watcher.nw] # Ensure network watchers exist before VNet
}
