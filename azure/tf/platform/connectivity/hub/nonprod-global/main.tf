
module "main" {
  source = "../_base"

  env         = local.envs.nonprod.name
  region_pair = "global"
}
