locals {
  # Create NSGs only for subnets that actually have any rules
  subnets_with_nsgs = {
    for subnet_key, subnet in var.subnets :
    subnet_key => subnet
    if length(keys(try(subnet.inbound_rules, {}))) > 0 || length(keys(try(subnet.outbound_rules, {}))) > 0
  }

  # Subnets that should be associated with a route table
  subnet_route_tables = {
    for subnet_key, subnet in var.subnets :
    subnet_key => azurerm_route_table.route_table[subnet.route_table_key].id
    if try(subnet.route_table_key, null) != null
  }

  # Subnets that should be associated with a NAT Gateway
  subnet_nat_gateways = {
    for subnet_key, subnet in var.subnets :
    subnet_key => subnet.nat_gateway_id
    if try(subnet.nat_gateway_id, null) != null
  }

  # Flatten inbound+outbound rules across all subnets
  # Produces map keyed by "<subnet>_<dir>_<rulekey>"
  nsg_rules_map = merge(
    # --------------------------
    # Inbound rules
    # --------------------------
    {
      for r in flatten([
        for subnet_key, subnet in local.subnets_with_nsgs : [
          for rule_key, rule in try(subnet.inbound_rules, {}) : {
            map_key    = "${subnet_key}_Inbound_${rule_key}"
            subnet_key = subnet_key
            name       = rule_key
            direction  = "Inbound"

            priority    = rule.priority
            access      = rule.access
            protocol    = rule.protocol
            description = try(rule.description, "")

            # --------------------------
            # Normalize ports
            # - If caller provided singular, use it
            # - Else if caller provided ranges == ["*"], convert to singular "*" (Azure rejects SourcePortRanges:["*"])
            # - Else use ranges list
            # --------------------------
            source_port_range = (
              try(rule.source_port_range, null) != null ? rule.source_port_range :
              (length(try(rule.source_port_ranges, [])) == 1 && try(rule.source_port_ranges[0], "") == "*") ? "*" :
              null
            )

            source_port_ranges = (
              try(rule.source_port_range, null) != null ? null :
              (length(try(rule.source_port_ranges, [])) == 1 && try(rule.source_port_ranges[0], "") == "*") ? null :
              (length(try(rule.source_port_ranges, [])) > 0 ? rule.source_port_ranges : null)
            )

            destination_port_range = (
              try(rule.destination_port_range, null) != null ? rule.destination_port_range :
              (length(try(rule.destination_port_ranges, [])) == 1 && try(rule.destination_port_ranges[0], "") == "*") ? "*" :
              null
            )

            destination_port_ranges = (
              try(rule.destination_port_range, null) != null ? null :
              (length(try(rule.destination_port_ranges, [])) == 1 && try(rule.destination_port_ranges[0], "") == "*") ? null :
              (length(try(rule.destination_port_ranges, [])) > 0 ? rule.destination_port_ranges : null)
            )

            # --------------------------
            # Normalize addresses
            # - If singular provided, use it
            # - Else if list has exactly 1 item, use singular (works for service tags too)
            # - Else use prefixes list
            # --------------------------
            source_address_prefix = (
              try(rule.source_address_prefix, null) != null ? rule.source_address_prefix :
              (length(try(rule.source_address_prefixes, [])) == 1 ? rule.source_address_prefixes[0] : null)
            )

            source_address_prefixes = (
              try(rule.source_address_prefix, null) != null ? null :
              (length(try(rule.source_address_prefixes, [])) > 1 ? rule.source_address_prefixes : null)
            )

            destination_address_prefix = (
              try(rule.destination_address_prefix, null) != null ? rule.destination_address_prefix :
              (length(try(rule.destination_address_prefixes, [])) == 1 ? rule.destination_address_prefixes[0] : null)
            )

            destination_address_prefixes = (
              try(rule.destination_address_prefix, null) != null ? null :
              (length(try(rule.destination_address_prefixes, [])) > 1 ? rule.destination_address_prefixes : null)
            )
          }
        ]
      ]) : r.map_key => r
    },

    # --------------------------
    # Outbound rules
    # --------------------------
    {
      for r in flatten([
        for subnet_key, subnet in local.subnets_with_nsgs : [
          for rule_key, rule in try(subnet.outbound_rules, {}) : {
            map_key    = "${subnet_key}_Outbound_${rule_key}"
            subnet_key = subnet_key
            name       = rule_key
            direction  = "Outbound"

            priority    = rule.priority
            access      = rule.access
            protocol    = rule.protocol
            description = try(rule.description, "")

            source_port_range = (
              try(rule.source_port_range, null) != null ? rule.source_port_range :
              (length(try(rule.source_port_ranges, [])) == 1 && try(rule.source_port_ranges[0], "") == "*") ? "*" :
              null
            )

            source_port_ranges = (
              try(rule.source_port_range, null) != null ? null :
              (length(try(rule.source_port_ranges, [])) == 1 && try(rule.source_port_ranges[0], "") == "*") ? null :
              (length(try(rule.source_port_ranges, [])) > 0 ? rule.source_port_ranges : null)
            )

            destination_port_range = (
              try(rule.destination_port_range, null) != null ? rule.destination_port_range :
              (length(try(rule.destination_port_ranges, [])) == 1 && try(rule.destination_port_ranges[0], "") == "*") ? "*" :
              null
            )

            destination_port_ranges = (
              try(rule.destination_port_range, null) != null ? null :
              (length(try(rule.destination_port_ranges, [])) == 1 && try(rule.destination_port_ranges[0], "") == "*") ? null :
              (length(try(rule.destination_port_ranges, [])) > 0 ? rule.destination_port_ranges : null)
            )

            source_address_prefix = (
              try(rule.source_address_prefix, null) != null ? rule.source_address_prefix :
              (length(try(rule.source_address_prefixes, [])) == 1 ? rule.source_address_prefixes[0] : null)
            )

            source_address_prefixes = (
              try(rule.source_address_prefix, null) != null ? null :
              (length(try(rule.source_address_prefixes, [])) > 1 ? rule.source_address_prefixes : null)
            )

            destination_address_prefix = (
              try(rule.destination_address_prefix, null) != null ? rule.destination_address_prefix :
              (length(try(rule.destination_address_prefixes, [])) == 1 ? rule.destination_address_prefixes[0] : null)
            )

            destination_address_prefixes = (
              try(rule.destination_address_prefix, null) != null ? null :
              (length(try(rule.destination_address_prefixes, [])) > 1 ? rule.destination_address_prefixes : null)
            )
          }
        ]
      ]) : r.map_key => r
    }
  )
}
