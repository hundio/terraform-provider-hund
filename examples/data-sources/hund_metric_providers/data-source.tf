data "hund_metric_providers" "example" {
  watchdog = "656e624f8fbb65049112ea7f"
}

output "provider_services" {
  value = nonsensitive(distinct([for mp in data.hund_metric_providers.example.metric_providers : one([for k, v in mp.service : k if v != null])]))
}

output "provided_slugs" {
  value = flatten([for mp in data.hund_metric_providers.example.metric_providers : keys(mp.instances)])
}
