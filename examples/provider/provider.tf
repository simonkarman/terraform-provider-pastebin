terraform {
  required_providers {
    pastebin = {
      source = "registry.terraform.io/simonkarman/pastebin"
    }
  }
}

provider "pastebin" {}
