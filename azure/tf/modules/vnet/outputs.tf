# ====================================
# Virtual Network Outputs
# ====================================

output "vnet_id" {
  description = "The ID of the virtual network"
  value       = azurerm_virtual_network.vnet.id
}

output "vnet_name" {
  description = "The name of the virtual network"
  value       = azurerm_virtual_network.vnet.name
}

output "vnet_location" {
  description = "The location of the virtual network"
  value       = azurerm_virtual_network.vnet.location
}

output "vnet_address_space" {
  description = "The address space of the virtual network"
  value       = azurerm_virtual_network.vnet.address_space
}

output "vnet_guid" {
  description = "The GUID of the virtual network"
  value       = azurerm_virtual_network.vnet.guid
}

# ====================================
# Subnet Outputs
# ====================================

output "subnet_ids" {
  description = "Map of subnet names to their IDs"
  value = {
    for k, v in azurerm_subnet.subnet : k => v.id
  }
}

output "subnet_address_prefixes" {
  description = "Map of subnet names to their address prefix"
  value = {
    for k, v in azurerm_subnet.subnet : k => v.address_prefixes[0]
  }
}

output "subnets" {
  description = "Map of all subnet details including id, name, address_prefix"
  value = {
    for k, v in azurerm_subnet.subnet : k => {
      id             = v.id
      name           = v.name
      address_prefix = v.address_prefixes[0]
    }
  }
}

# ====================================
# Network Security Group Outputs
# ====================================

output "nsg_ids" {
  description = "Map of subnet names to their NSG IDs"
  value = {
    for k, v in azurerm_network_security_group.subnet : k => v.id
  }
}

output "nsg_names" {
  description = "Map of subnet names to their NSG names"
  value = {
    for k, v in azurerm_network_security_group.subnet : k => v.name
  }
}

output "nsgs" {
  description = "Map of all NSG details including id, name, and associated subnet"
  value = {
    for k, v in azurerm_network_security_group.subnet : k => {
      id     = v.id
      name   = v.name
      subnet = k
    }
  }
}

# ====================================
# Convenient Outputs for Peering
# ====================================

output "vnet_peering_info" {
  description = "Information needed for VNet peering operations"
  value = {
    id             = azurerm_virtual_network.vnet.id
    name           = azurerm_virtual_network.vnet.name
    resource_group = var.rg_name
    address_space  = azurerm_virtual_network.vnet.address_space
  }
}

# ====================================
# Private Endpoint Configuration Outputs
# ====================================

output "private_endpoint_subnets" {
  description = "Map of subnets configured for private endpoints (network policies disabled)"
  value = {
    for k, v in var.subnets : k => {
      id                       = azurerm_subnet.subnet[k].id
      name                     = azurerm_subnet.subnet[k].name
      address_prefix           = azurerm_subnet.subnet[k].address_prefixes[0]
      network_policies_enabled = v.private_endpoint_network_policies_enabled
    } if !v.private_endpoint_network_policies_enabled
  }
}

output "private_link_service_subnets" {
  description = "Map of subnets configured for private link service (network policies disabled)"
  value = {
    for k, v in var.subnets : k => {
      id                       = azurerm_subnet.subnet[k].id
      name                     = azurerm_subnet.subnet[k].name
      address_prefix           = azurerm_subnet.subnet[k].address_prefixes[0]
      network_policies_enabled = v.private_link_service_network_policies_enabled
    } if !v.private_link_service_network_policies_enabled
  }
}

# ====================================
# Delegation Outputs
# ====================================

output "delegated_subnets" {
  description = "Map of subnets that have delegations configured"
  value = {
    for k, v in var.subnets : k => {
      id             = azurerm_subnet.subnet[k].id
      name           = azurerm_subnet.subnet[k].name
      address_prefix = azurerm_subnet.subnet[k].address_prefixes[0]
      delegations    = v.delegations
    } if length(v.delegations) > 0
  }
}

# ====================================
# Route Table Outputs
# ====================================

output "route_table_ids" {
  description = "Map of route table keys to their resource IDs"
  value = {
    for k, v in azurerm_route_table.route_table : k => v.id
  }
}

output "route_table_names" {
  description = "Map of route table keys to their names"
  value = {
    for k, v in azurerm_route_table.route_table : k => v.name
  }
}

# ====================================
# Association Outputs
# ====================================

output "route_table_associations" {
  description = "Map of subnets to their associated route table IDs"
  value       = local.subnet_route_tables
}

output "nat_gateway_associations" {
  description = "Map of subnets to their associated NAT gateway IDs"
  value       = local.subnet_nat_gateways
}
