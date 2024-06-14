terraform {
  required_providers {
    pastebin = {
      source = "registry.terraform.io/simonkarman/pastebin"
    }
  }
}

provider "pastebin" {}

resource "pastebin_paste" "example" {}
