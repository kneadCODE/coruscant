locals {
  subscription_keys = keys(
    nonsensitive(local.subscription_ids)
  )

  subscription_map = {
    for key in local.subscription_keys :
    key => local.subscription_ids[key]
  }
}
