############################################
# variables.tf (fixed / clean / simple)
############################################

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

# ------------------------------------------------------------------------------
# NSG rule object (supports both singular and plural)
# - We default the plural lists to [] so validations can use length() safely.
# ------------------------------------------------------------------------------

variable "subnets" {
  type = map(object({
    address_prefix                                = string
    service_endpoints                             = optional(list(string), [])
    private_endpoint_network_policies_enabled     = optional(bool, true)
    private_link_service_network_policies_enabled = optional(bool, true)
    default_outbound_access_enabled               = optional(bool, true)

    delegations = optional(list(object({
      name = string
      service_delegation = object({
        name    = string
        actions = optional(list(string), [])
      })
    })), [])

    inbound_rules = optional(map(object({
      priority = number
      access   = string
      protocol = string

      # Ports
      source_port_range       = optional(string, null)
      source_port_ranges      = optional(list(string), [])
      destination_port_range  = optional(string, null)
      destination_port_ranges = optional(list(string), [])

      # Addresses
      source_address_prefix        = optional(string, null)
      source_address_prefixes      = optional(list(string), [])
      destination_address_prefix   = optional(string, null)
      destination_address_prefixes = optional(list(string), [])

      description = optional(string, "")
    })), {})

    outbound_rules = optional(map(object({
      priority = number
      access   = string
      protocol = string

      # Ports
      source_port_range       = optional(string, null)
      source_port_ranges      = optional(list(string), [])
      destination_port_range  = optional(string, null)
      destination_port_ranges = optional(list(string), [])

      # Addresses
      source_address_prefix        = optional(string, null)
      source_address_prefixes      = optional(list(string), [])
      destination_address_prefix   = optional(string, null)
      destination_address_prefixes = optional(list(string), [])

      description = optional(string, "")
    })), {})

    route_table_key = optional(string, null)
    nat_gateway_id  = optional(string, null)
  }))

  description = "Map of subnets to create. Key is the subnet name."

  # ---------------------------
  # Subnet CIDR validation
  # ---------------------------
  validation {
    condition     = alltrue([for _, subnet in var.subnets : can(cidrhost(subnet.address_prefix, 0))])
    error_message = "Each subnet address prefix must be valid CIDR notation"
  }

  # ---------------------------
  # Inbound rule validations
  # ---------------------------
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([for rule in values(subnet.inbound_rules) : rule.priority >= 100 && rule.priority <= 4096])
    ])
    error_message = "Inbound rule priorities must be between 100 and 4096"
  }

  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([for rule in values(subnet.inbound_rules) : contains(["Allow", "Deny"], rule.access)])
    ])
    error_message = "Inbound rule access must be either 'Allow' or 'Deny'"
  }

  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([for rule in values(subnet.inbound_rules) : contains(["Tcp", "Udp", "Icmp", "Esp", "Ah", "*"], rule.protocol)])
    ])
    error_message = "Inbound rule protocol must be one of: Tcp, Udp, Icmp, Esp, Ah, *"
  }

  # Require source ports (singular OR plural)
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([
        for rule in values(subnet.inbound_rules) :
        (rule.source_port_range != null || length(rule.source_port_ranges) > 0)
      ])
    ])
    error_message = "Each inbound rule must specify source_port_range or source_port_ranges."
  }

  # Require destination ports
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([
        for rule in values(subnet.inbound_rules) :
        (rule.destination_port_range != null || length(rule.destination_port_ranges) > 0)
      ])
    ])
    error_message = "Each inbound rule must specify destination_port_range or destination_port_ranges."
  }

  # Require source addresses
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([
        for rule in values(subnet.inbound_rules) :
        (rule.source_address_prefix != null || length(rule.source_address_prefixes) > 0)
      ])
    ])
    error_message = "Each inbound rule must specify source_address_prefix or source_address_prefixes."
  }

  # Require destination addresses
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([
        for rule in values(subnet.inbound_rules) :
        (rule.destination_address_prefix != null || length(rule.destination_address_prefixes) > 0)
      ])
    ])
    error_message = "Each inbound rule must specify destination_address_prefix or destination_address_prefixes."
  }

  # Disallow setting both singular + plural for same field
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([
        for rule in values(subnet.inbound_rules) :
        !(
          (rule.source_port_range != null && length(rule.source_port_ranges) > 0) ||
          (rule.destination_port_range != null && length(rule.destination_port_ranges) > 0) ||
          (rule.source_address_prefix != null && length(rule.source_address_prefixes) > 0) ||
          (rule.destination_address_prefix != null && length(rule.destination_address_prefixes) > 0)
        )
      ])
    ])
    error_message = "Inbound rules: set only one of each pair (e.g., source_port_range OR source_port_ranges), not both."
  }

  # ---------------------------
  # Outbound rule validations
  # ---------------------------
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([for rule in values(subnet.outbound_rules) : rule.priority >= 100 && rule.priority <= 4096])
    ])
    error_message = "Outbound rule priorities must be between 100 and 4096"
  }

  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([for rule in values(subnet.outbound_rules) : contains(["Allow", "Deny"], rule.access)])
    ])
    error_message = "Outbound rule access must be either 'Allow' or 'Deny'"
  }

  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([for rule in values(subnet.outbound_rules) : contains(["Tcp", "Udp", "Icmp", "Esp", "Ah", "*"], rule.protocol)])
    ])
    error_message = "Outbound rule protocol must be one of: Tcp, Udp, Icmp, Esp, Ah, *"
  }

  # Require source ports
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([
        for rule in values(subnet.outbound_rules) :
        (rule.source_port_range != null || length(rule.source_port_ranges) > 0)
      ])
    ])
    error_message = "Each outbound rule must specify source_port_range or source_port_ranges."
  }

  # Require destination ports
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([
        for rule in values(subnet.outbound_rules) :
        (rule.destination_port_range != null || length(rule.destination_port_ranges) > 0)
      ])
    ])
    error_message = "Each outbound rule must specify destination_port_range or destination_port_ranges."
  }

  # Require source addresses
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([
        for rule in values(subnet.outbound_rules) :
        (rule.source_address_prefix != null || length(rule.source_address_prefixes) > 0)
      ])
    ])
    error_message = "Each outbound rule must specify source_address_prefix or source_address_prefixes."
  }

  # Require destination addresses
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([
        for rule in values(subnet.outbound_rules) :
        (rule.destination_address_prefix != null || length(rule.destination_address_prefixes) > 0)
      ])
    ])
    error_message = "Each outbound rule must specify destination_address_prefix or destination_address_prefixes."
  }

  # Disallow setting both singular + plural
  validation {
    condition = alltrue([
      for _, subnet in var.subnets :
      alltrue([
        for rule in values(subnet.outbound_rules) :
        !(
          (rule.source_port_range != null && length(rule.source_port_ranges) > 0) ||
          (rule.destination_port_range != null && length(rule.destination_port_ranges) > 0) ||
          (rule.source_address_prefix != null && length(rule.source_address_prefixes) > 0) ||
          (rule.destination_address_prefix != null && length(rule.destination_address_prefixes) > 0)
        )
      ])
    ])
    error_message = "Outbound rules: set only one of each pair (e.g., source_port_range OR source_port_ranges), not both."
  }
}

variable "route_tables" {
  type = map(object({
    bgp_route_propagation_enabled = optional(bool, false)
    routes = map(object({
      address_prefix         = string
      next_hop_type          = string
      next_hop_in_ip_address = optional(string, null)
    }))
  }))

  description = "Map of route tables to create. Key is the route table name suffix (e.g., 'force-firewall')."
  default     = {}

  validation {
    condition = alltrue([
      for _, rt in var.route_tables :
      alltrue([
        for _, route in rt.routes :
        contains(["VirtualNetworkGateway", "VnetLocal", "Internet", "VirtualAppliance", "None"], route.next_hop_type)
      ])
    ])
    error_message = "Route next_hop_type must be one of: VirtualNetworkGateway, VnetLocal, Internet, VirtualAppliance, None"
  }

  validation {
    condition = alltrue([
      for _, rt in var.route_tables :
      alltrue([
        for _, route in rt.routes :
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
