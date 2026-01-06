resource "azurerm_role_assignment" "tf_apply_contributor" {
  scope                = "/subscriptions/${var.subscription_id}"
  role_definition_name = "Contributor"
  principal_id         = var.sp_tf_apply_obj_id
  principal_type       = "ServicePrincipal"
}

resource "azurerm_role_assignment" "tf_apply_rbac_admin" {
  scope                = "/subscriptions/${var.subscription_id}"
  role_definition_name = "Role Based Access Control Administrator"
  principal_id         = var.sp_tf_apply_obj_id
  principal_type       = "ServicePrincipal"
}

resource "azurerm_role_assignment" "tf_apply_lock_manager" {
  scope                = "/subscriptions/${var.subscription_id}"
  role_definition_name = var.lock_manager_role_name
  principal_id         = var.sp_tf_apply_obj_id
  principal_type       = "ServicePrincipal"
}
