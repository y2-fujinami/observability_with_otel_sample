terraform {
  # terraform自体のバージョン制約
  required_version = "1.12.2"

  # 使用するproviderとバージョン制約
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "6.44.0"
    }
  }
}