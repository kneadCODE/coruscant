variable "env" {
  type        = string
  description = "The environment (prod or nonprod)"
  validation {
    condition     = contains(local.allowed_envs, var.env)
    error_message = "Invalid environment. Must be one of: ${join(", ", local.allowed_envs)}"
  }
}

variable "region_pair" {
  type        = string
  description = "The region pair to deploy resources in (global or gdpr)"
  validation {
    condition     = contains(keys(local.regions), var.region_pair)
    error_message = "Invalid region pair. Must be one of: ${join(", ", keys(local.regions))}"
  }
}
