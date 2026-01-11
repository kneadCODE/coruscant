locals {
  allowed_regions = ["southeastasia"]

  region_shorthand = {
    "southeastasia" : "sea",
  }

  allowed_envs = ["prod", "nonprod", "staging", "qa", "dev"]
}
