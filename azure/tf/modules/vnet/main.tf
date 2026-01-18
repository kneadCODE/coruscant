# ====================================
# Virtual Network
# ====================================

resource "azurerm_virtual_network" "vnet" {
  name                = "vnet-${var.purpose}-${var.env}-${local.region_shorthand_lookup[var.region]}-0${var.instance_number}"
  resource_group_name = var.rg_name
  location            = var.region

  address_space = [var.address_space]

  dynamic "ddos_protection_plan" {
    for_each = var.ddos_protection_plan_id != "" ? [1] : []
    content {
      enable = true
      id     = var.ddos_protection_plan_id
    }
  }

  tags = var.tags
}

# ====================================
# Route Tables
# ====================================

resource "azurerm_route_table" "route_table" {
  for_each = var.route_tables

  name                          = "rt-${var.purpose}-${each.key}-${var.env}-${local.region_shorthand_lookup[var.region]}-0${var.instance_number}"
  resource_group_name           = var.rg_name
  location                      = var.region
  bgp_route_propagation_enabled = each.value.bgp_route_propagation_enabled

  # Inline routes for cleaner configuration
  dynamic "route" {
    for_each = each.value.routes
    content {
      name                   = route.key
      address_prefix         = route.value.address_prefix
      next_hop_type          = route.value.next_hop_type
      next_hop_in_ip_address = route.value.next_hop_in_ip_address
    }
  }

  tags = var.tags
}

# ====================================
# Network Security Groups
# ====================================

resource "azurerm_network_security_group" "subnet" {
  for_each = local.subnets_with_nsgs

  name                = "nsg-${var.purpose}-${each.key}-${var.env}-${local.region_shorthand_lookup[var.region]}-0${var.instance_number}"
  resource_group_name = var.rg_name
  location            = var.region
  tags                = var.tags
}

# ====================================
# Network Security Rules
# ====================================

resource "azurerm_network_security_rule" "rules" {
  for_each = local.nsg_rules_map

  name                         = each.value.name
  resource_group_name          = var.rg_name
  network_security_group_name  = azurerm_network_security_group.subnet[each.value.subnet_key].name
  priority                     = each.value.priority
  direction                    = each.value.direction
  access                       = each.value.access
  protocol                     = each.value.protocol
  description                  = each.value.description != "" ? each.value.description : null
  source_port_range            = each.value.source_port_range
  source_port_ranges           = each.value.source_port_ranges
  source_address_prefix        = each.value.source_address_prefix
  source_address_prefixes      = each.value.source_address_prefixes
  destination_port_range       = each.value.destination_port_range
  destination_port_ranges      = each.value.destination_port_ranges
  destination_address_prefix   = each.value.destination_address_prefix
  destination_address_prefixes = each.value.destination_address_prefixes

  depends_on = [azurerm_network_security_group.subnet]
}

# ====================================
# Subnets
# ====================================

resource "azurerm_subnet" "subnet" {
  for_each = var.subnets

  name                 = each.key
  resource_group_name  = var.rg_name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = [each.value.address_prefix]

  # Service endpoints
  service_endpoints = length(each.value.service_endpoints) > 0 ? each.value.service_endpoints : null

  # Private endpoint and private link service policies
  private_endpoint_network_policies             = each.value.private_endpoint_network_policies_enabled ? "Enabled" : "Disabled"
  private_link_service_network_policies_enabled = each.value.private_link_service_network_policies_enabled

  # Default outbound access (Azure default is true, but can be explicitly disabled)
  default_outbound_access_enabled = each.value.default_outbound_access_enabled

  # Subnet delegations
  dynamic "delegation" {
    for_each = each.value.delegations
    content {
      name = delegation.value.name
      service_delegation {
        name    = delegation.value.service_delegation.name
        actions = delegation.value.service_delegation.actions
      }
    }
  }

  depends_on = [azurerm_virtual_network.vnet]
}

# ====================================
# Subnet - NSG Associations
# ====================================

resource "azurerm_subnet_network_security_group_association" "subnet_nsg" {
  for_each = local.subnets_with_nsgs

  subnet_id                 = azurerm_subnet.subnet[each.key].id
  network_security_group_id = azurerm_network_security_group.subnet[each.key].id

  depends_on = [
    azurerm_subnet.subnet,
    azurerm_network_security_group.subnet,
    azurerm_network_security_rule.rules
  ]
}

# ====================================
# Subnet - Route Table Associations
# ====================================

resource "azurerm_subnet_route_table_association" "subnet_rt" {
  for_each = local.subnet_route_tables

  subnet_id      = azurerm_subnet.subnet[each.key].id
  route_table_id = each.value

  depends_on = [azurerm_subnet.subnet]
}

# ====================================
# Subnet - NAT Gateway Associations
# ====================================

resource "azurerm_subnet_nat_gateway_association" "subnet_nat" {
  for_each = local.subnet_nat_gateways

  subnet_id      = azurerm_subnet.subnet[each.key].id
  nat_gateway_id = each.value

  depends_on = [azurerm_subnet.subnet]
}
