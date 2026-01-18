terraform {
  required_version = ">= 1.11.0" # opentofu

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0" # Use latest 4.x version, allows minor/patch updates
    }
  }
}
