terraform {
  required_providers {
    pastebin = {
      source = "registry.terraform.io/simonkarman/pastebin"
    }
  }
}

provider "pastebin" {}

data "pastebin_noop" "example" {}
