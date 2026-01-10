locals {
  # Flatten inbound NSG rules - transforms nested map structure to flat list
  inbound_rules_flat = flatten([
    for subnet_key, subnet in var.subnets : [
      for rule_key, rule in subnet.inbound_rules : {
        subnet_key                   = subnet_key
        rule_key                     = rule_key
        name                         = rule_key # Use map key as rule name (kebab-case)
        priority                     = rule.priority
        direction                    = "Inbound"
        access                       = rule.access
        protocol                     = rule.protocol
        source_port_ranges           = rule.source_port_ranges
        source_address_prefixes      = rule.source_address_prefixes
        destination_port_ranges      = rule.destination_port_ranges
        destination_address_prefixes = rule.destination_address_prefixes
        description                  = rule.description
      }
    ]
  ])

  # Flatten outbound NSG rules - transforms nested map structure to flat list
  outbound_rules_flat = flatten([
    for subnet_key, subnet in var.subnets : [
      for rule_key, rule in subnet.outbound_rules : {
        subnet_key                   = subnet_key
        rule_key                     = rule_key
        name                         = rule_key # Use map key as rule name (kebab-case)
        priority                     = rule.priority
        direction                    = "Outbound"
        access                       = rule.access
        protocol                     = rule.protocol
        source_port_ranges           = rule.source_port_ranges
        source_address_prefixes      = rule.source_address_prefixes
        destination_port_ranges      = rule.destination_port_ranges
        destination_address_prefixes = rule.destination_address_prefixes
        description                  = rule.description
      }
    ]
  ])

  # Combine both inbound and outbound rules
  all_nsg_rules = concat(local.inbound_rules_flat, local.outbound_rules_flat)

  # Create map of NSG rules indexed by unique key for for_each
  nsg_rules_map = {
    for rule in local.all_nsg_rules :
    "${rule.subnet_key}_${rule.direction}_${rule.rule_key}" => rule
  }

  # Flatten delegations for easier iteration
  subnet_delegations = {
    for subnet_key, subnet in var.subnets :
    subnet_key => subnet.delegations
    if length(subnet.delegations) > 0
  }

  # Flatten route table associations (map route_table_key to actual resource ID)
  subnet_route_tables = {
    for subnet_key, subnet in var.subnets :
    subnet_key => azurerm_route_table.route_table[subnet.route_table_key].id
    if subnet.route_table_key != null
  }

  # Flatten NAT Gateway associations
  subnet_nat_gateways = {
    for subnet_key, subnet in var.subnets :
    subnet_key => subnet.nat_gateway_id
    if subnet.nat_gateway_id != null
  }

  # Subnets that have NSG rules (either inbound or outbound)
  subnets_with_nsgs = {
    for k, v in var.subnets : k => v
    if length(v.inbound_rules) > 0 || length(v.outbound_rules) > 0
  }
}
