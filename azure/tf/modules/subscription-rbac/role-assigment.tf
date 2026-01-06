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
  scope              = "/subscriptions/${var.subscription_id}"
  role_definition_id = azurerm_role_definition.locks_manager.role_definition_resource_id
  principal_id       = var.sp_tf_apply_obj_id
  principal_type     = "ServicePrincipal"
  depends_on         = [azurerm_role_definition.locks_manager]
}
