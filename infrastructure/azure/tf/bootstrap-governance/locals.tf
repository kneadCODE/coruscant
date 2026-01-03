locals {
  # Parse subscription mapping from JSON secret
  # Expected format: {"management": "guid", "identity": "guid", ...}
  subscription_ids = jsondecode(var.subscription_id_mapping_json)
}
