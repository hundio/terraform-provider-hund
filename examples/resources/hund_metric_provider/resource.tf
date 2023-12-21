resource "hund_component" "component0" {
  name        = "Terraform Component"
  description = "Created by Terraform"

  group = "63f6c4938fbb652be74a9ae6"

  watchdog = {
    service = {
      webhook = { webhook_key = "KEY" }
    }
  }
}

resource "hund_metric_provider" "default" {
  watchdog = hund_component.component0.watchdog.id
  service  = { webhook = hund_component.component0.watchdog.service.webhook }

  # Managing a Watchdog's default MetricProvider requires importation.
  default = true

  instances = {}
}

resource "hund_metric_provider" "builtin" {
  watchdog = hund_component.component0.watchdog.id

  service = { builtin = {} }

  instances = {
    percent_uptime     = { enabled = true },
    incidents_reported = { enabled = false }
  }
}
