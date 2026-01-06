resource "azurerm_consumption_budget_subscription" "cheapo" {
  for_each = local.subscription_map

  name            = "cheapo"
  subscription_id = "/subscriptions/${each.value}"

  amount     = 1 # Yes 1 USD. Cause too poor
  time_grain = "Monthly"

  time_period {
    start_date = "2026-01-01T00:00:00Z"
  }

  # 25% alert
  notification {
    enabled        = true
    threshold      = 25.0
    operator       = "GreaterThanOrEqualTo"
    threshold_type = "Actual"
    contact_roles  = ["Owner", "Cost Management Contributor"]
  }

  # 50% alert
  notification {
    enabled        = true
    threshold      = 50.0
    operator       = "GreaterThanOrEqualTo"
    threshold_type = "Actual"
    contact_roles  = ["Owner", "Cost Management Contributor"]
  }

  # 75% alert
  notification {
    enabled        = true
    threshold      = 75.0
    operator       = "GreaterThanOrEqualTo"
    threshold_type = "Actual"
    contact_roles  = ["Owner", "Cost Management Contributor"]
  }

  # 100% alert
  notification {
    enabled        = true
    threshold      = 100.0
    operator       = "GreaterThanOrEqualTo"
    threshold_type = "Actual"
    contact_roles  = ["Owner", "Cost Management Contributor"]
  }
}
