resource "azurerm_role_definition" "locks_manager" {
  name        = "Locks Manager (RG/Resource)"
  scope       = data.azurerm_management_group.root.id
  description = "Can read/create/delete management locks at resource group or resource scope."

  permissions {
    actions = [
      "Microsoft.Authorization/locks/read",
      "Microsoft.Authorization/locks/write",
      "Microsoft.Authorization/locks/delete",
    ]
    not_actions = []
  }
  assignable_scopes = [
    data.azurerm_management_group.platform.id,
    data.azurerm_management_group.landingzone.id,
  ]
}
