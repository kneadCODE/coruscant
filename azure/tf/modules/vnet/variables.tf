variable "rg_name" {
  type        = string
  description = "The name of the resource group"
}

variable "env" {
  type        = string
  description = "The environment for the VNET"

  validation {
    condition     = contains(local.allowed_envs, var.env)
    error_message = "Environment must be one of allowed environments"
  }
}

variable "region" {
  type        = string
  description = "Azure region to deploy the VNET to"

  validation {
    condition     = contains(local.allowed_regions, var.region)
    error_message = "Region must be one of the allowed regions"
  }
}

variable "purpose" {
  type        = string
  description = "Purpose/name of the VNET"

  validation {
    condition     = length(var.purpose) > 0 && length(var.purpose) <= 15
    error_message = "Purpose must be between 1 and 15 characters"
  }

  validation {
    condition     = can(regex("^[a-z0-9-]+$", var.purpose))
    error_message = "Purpose must contain only lowercase letters, numbers, and hyphens"
  }
}

variable "instance_number" {
  type        = number
  description = "The instance number for this VNET"
  default     = 1

  validation {
    condition     = var.instance_number >= 1 && var.instance_number <= 9
    error_message = "Instance number must be between 1 and 9"
  }
}

variable "address_space" {
  type        = string
  description = "The address space for the VNET in CIDR notation (e.g., '10.0.0.0/16')"

  validation {
    condition     = can(cidrhost(var.address_space, 0))
    error_message = "Address space must be valid CIDR notation"
  }
}

variable "subnets" {
  type = map(object({
    address_prefix                                = string                     # CIDR for the subnet (e.g., "10.0.0.0/24")
    service_endpoints                             = optional(list(string), []) # Service endpoints (Microsoft.Storage, etc.)
    private_endpoint_network_policies_enabled     = optional(bool, true)       # Enable/disable network policies for private endpoints
    private_link_service_network_policies_enabled = optional(bool, true)       # Enable/disable network policies for private link service
    default_outbound_access_enabled               = optional(bool, true)       # Default outbound internet access

    # Subnet delegation configuration
    delegations = optional(list(object({
      name = string # Delegation name
      service_delegation = object({
        name    = string                     # Service name (e.g., Microsoft.Web/serverFarms)
        actions = optional(list(string), []) # Delegated actions
      })
    })), [])

    # Inbound NSG rules for this subnet (map key is used as rule name)
    inbound_rules = optional(map(object({
      priority                     = number       # 100-4096
      access                       = string       # Allow or Deny
      protocol                     = string       # Tcp, Udp, Icmp, Esp, Ah, or * (any)
      source_port_ranges           = list(string) # Source ports/ranges (e.g., ["*"], ["80", "443"], ["1024-65535"])
      source_address_prefixes      = list(string) # Source CIDRs/service tags (e.g., ["Internet"], ["10.0.0.0/16"])
      destination_port_ranges      = list(string) # Destination ports/ranges
      destination_address_prefixes = list(string) # Destination CIDRs/service tags
      description                  = optional(string, "")
    })), {})

    # Outbound NSG rules for this subnet (map key is used as rule name)
    outbound_rules = optional(map(object({
      priority                     = number       # 100-4096
      access                       = string       # Allow or Deny
      protocol                     = string       # Tcp, Udp, Icmp, Esp, Ah, or * (any)
      source_port_ranges           = list(string) # Source ports/ranges
      source_address_prefixes      = list(string) # Source CIDRs/service tags
      destination_port_ranges      = list(string) # Destination ports/ranges
      destination_address_prefixes = list(string) # Destination CIDRs/service tags
      description                  = optional(string, "")
    })), {})

    # Optional associations
    route_table_key = optional(string, null) # Route table map key (references var.route_tables)
    nat_gateway_id  = optional(string, null) # NAT gateway resource ID
  }))

  description = "Map of subnets to create. Key is the subnet name."

  validation {
    condition     = alltrue([for subnet in var.subnets : can(cidrhost(subnet.address_prefix, 0))])
    error_message = "Each subnet address prefix must be valid CIDR notation"
  }

  # Inbound rule validations
  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.inbound_rules) : rule.priority >= 100 && rule.priority <= 4096])
    ])
    error_message = "Inbound rule priorities must be between 100 and 4096"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.inbound_rules) : contains(["Allow", "Deny"], rule.access)])
    ])
    error_message = "Inbound rule access must be either 'Allow' or 'Deny'"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.inbound_rules) : contains(["Tcp", "Udp", "Icmp", "Esp", "Ah", "*"], rule.protocol)])
    ])
    error_message = "Inbound rule protocol must be one of: Tcp, Udp, Icmp, Esp, Ah, *"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.inbound_rules) : length(rule.source_port_ranges) > 0])
    ])
    error_message = "Each inbound rule must specify at least one source_port_range"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.inbound_rules) : length(rule.source_address_prefixes) > 0])
    ])
    error_message = "Each inbound rule must specify at least one source_address_prefix"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.inbound_rules) : length(rule.destination_port_ranges) > 0])
    ])
    error_message = "Each inbound rule must specify at least one destination_port_range"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.inbound_rules) : length(rule.destination_address_prefixes) > 0])
    ])
    error_message = "Each inbound rule must specify at least one destination_address_prefix"
  }

  # Outbound rule validations
  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.outbound_rules) : rule.priority >= 100 && rule.priority <= 4096])
    ])
    error_message = "Outbound rule priorities must be between 100 and 4096"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.outbound_rules) : contains(["Allow", "Deny"], rule.access)])
    ])
    error_message = "Outbound rule access must be either 'Allow' or 'Deny'"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.outbound_rules) : contains(["Tcp", "Udp", "Icmp", "Esp", "Ah", "*"], rule.protocol)])
    ])
    error_message = "Outbound rule protocol must be one of: Tcp, Udp, Icmp, Esp, Ah, *"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.outbound_rules) : length(rule.source_port_ranges) > 0])
    ])
    error_message = "Each outbound rule must specify at least one source_port_range"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.outbound_rules) : length(rule.source_address_prefixes) > 0])
    ])
    error_message = "Each outbound rule must specify at least one source_address_prefix"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.outbound_rules) : length(rule.destination_port_ranges) > 0])
    ])
    error_message = "Each outbound rule must specify at least one destination_port_range"
  }

  validation {
    condition = alltrue([
      for subnet in var.subnets :
      alltrue([for rule in values(subnet.outbound_rules) : length(rule.destination_address_prefixes) > 0])
    ])
    error_message = "Each outbound rule must specify at least one destination_address_prefix"
  }
}

variable "route_tables" {
  type = map(object({
    bgp_route_propagation_enabled = optional(bool, false) # Default: disable BGP route propagation
    routes = map(object({                                 # Map of routes (key is route name)
      address_prefix         = string                     # Destination CIDR (e.g., "0.0.0.0/0", "10.0.0.0/8")
      next_hop_type          = string                     # VirtualNetworkGateway, VnetLocal, Internet, VirtualAppliance, None
      next_hop_in_ip_address = optional(string, null)     # Required when next_hop_type = VirtualAppliance
    }))
  }))

  description = "Map of route tables to create. Key is the route table name suffix (e.g., 'force-firewall')."
  default     = {}

  validation {
    condition = alltrue([
      for rt_key, rt in var.route_tables :
      alltrue([
        for route_key, route in rt.routes :
        contains(["VirtualNetworkGateway", "VnetLocal", "Internet", "VirtualAppliance", "None"], route.next_hop_type)
      ])
    ])
    error_message = "Route next_hop_type must be one of: VirtualNetworkGateway, VnetLocal, Internet, VirtualAppliance, None"
  }

  validation {
    condition = alltrue([
      for rt_key, rt in var.route_tables :
      alltrue([
        for route_key, route in rt.routes :
        route.next_hop_type != "VirtualAppliance" || route.next_hop_in_ip_address != null
      ])
    ])
    error_message = "next_hop_in_ip_address is required when next_hop_type is VirtualAppliance"
  }
}

variable "ddos_protection_plan_id" {
  type        = string
  description = "DDoS Protection Plan ID (leave empty to disable)"
  default     = ""
}

variable "tags" {
  type        = map(any)
  description = "Tags to apply to all resources"
}
