module "main" {
  source = "../_base"

  env                    = local.envs.nonprod.name
  region_pair            = "global"
  siem_deploy_st         = true
  siem_deploy_law        = true
  storage_account_prefix = var.storage_account_prefix
}
