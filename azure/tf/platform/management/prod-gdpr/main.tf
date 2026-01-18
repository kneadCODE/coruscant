module "management" {
  source = "../_base"

  env                    = local.envs.prod.name
  region_pair            = "gdpr"
  logs_deploy_st         = true
  logs_deploy_law        = true
  backup_deploy_rsv      = true
  backup_deploy_bvault   = true
  storage_account_prefix = var.storage_account_prefix
}
