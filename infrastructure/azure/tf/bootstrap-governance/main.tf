locals {
  # Level 1: Top-level MGs (children of root) - never have direct subscriptions
  mg_level1 = {
    platform       = "mg-coruscant-platform"
    landingzone    = "mg-coruscant-landingzone"
    sandbox        = "mg-coruscant-sandbox"
    decommissioned = "mg-coruscant-decommissioned"
  }

  # Level 2: Child MGs with subscription assignments
  mg_level2 = {
    # Platform children
    management = {
      name          = "mg-coruscant-management"
      parent_key    = "platform"
      subscriptions = ["management"]
    }
    identity = {
      name          = "mg-coruscant-identity"
      parent_key    = "platform"
      subscriptions = ["identity"]
    }
    connectivity = {
      name       = "mg-coruscant-connectivity"
      parent_key = "platform"
      subscriptions = [
        "connectivity_prod",
        "connectivity_nonprod"
      ]
    }
    security = {
      name       = "mg-coruscant-security"
      parent_key = "platform"
      subscriptions = [
        "security_prod",
        "security_nonprod"
      ]
    }

    # Landing zone children
    corp = {
      name       = "mg-coruscant-corp"
      parent_key = "landingzone"
      subscriptions = [
        "devops_prod",
        "devops_nonprod",
        # "observability_prod",
        # "observability_nonprod",
        "esb_prod",
        "esb_nonprod"
      ]
    }
    online = {
      name       = "mg-coruscant-online"
      parent_key = "landingzone"
      subscriptions = [
        "edge_prod",
        "edge_nonprod",
        # "iam_prod",
        # "iam_nonprod"
      ]
    }
  }
}


# Level 1: Top-level management groups (children of root)
resource "azurerm_management_group" "level1" {
  for_each = local.mg_level1

  name                       = each.value
  parent_management_group_id = data.azurerm_management_group.root.id
}

# Level 2: Platform and Landing Zone child management groups
resource "azurerm_management_group" "level2" {
  for_each = local.mg_level2

  name                       = each.value.name
  parent_management_group_id = azurerm_management_group.level1[each.value.parent_key].id
}

# Subscription to Management Group associations
# Dynamically creates associations by flattening the mg_level2 subscription lists
resource "azurerm_management_group_subscription_association" "all" {
  for_each = merge([
    for mg_key, mg_config in local.mg_level2 : {
      for sub_key in mg_config.subscriptions : "${mg_key}_${sub_key}" => {
        mg_id  = azurerm_management_group.level2[mg_key].id
        sub_id = local.subscription_ids[sub_key]
      }
    }
  ]...)

  management_group_id = each.value.mg_id
  subscription_id     = "/subscriptions/${each.value.sub_id}"
}
