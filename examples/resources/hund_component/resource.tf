resource "hund_group" "group" {
  name = "Terraform Group"
}

resource "hund_component" "component0" {
  name        = "Terraform Component"
  description = "Created by Terraform"

  group = hund_group.group.id

  watchdog = {
    service = {
      icmp = {
        regions = ["wa-us-1"]
        target  = "example.com"
      }
    }
  }
}

resource "hund_component" "utr" {
  name  = "Terraform Component (utr)"
  group = hund_group.group.id

  watchdog = {
    service = {
      uptimerobot = { monitor_api_key = "m1234-NOTAKEY" }
    }
  }
}
