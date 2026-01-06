variable "subscription_id" {
  description = "The subscription on which the roles will be applied"
  type        = string
  sensitive   = true
}

variable "sp_tf_apply_obj_id" {
  description = "The object ID for the TF apply service principal"
  type        = string
  sensitive   = true
}

variable "lock_manager_role_name" {
  description = "The name of the custom lock manager role"
  type        = string
}
