terraform {
  required_providers {
    hund = {
      source = "registry.terraform.io/hundio/hund"
    }
  }
}

provider "hund" {
  domain = "porbo.hund.localhost"
}
