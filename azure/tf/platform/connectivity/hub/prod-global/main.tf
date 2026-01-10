module "main" {
  source = "../_base"

  env         = local.envs.prod.name
  region_pair = "global"
}
