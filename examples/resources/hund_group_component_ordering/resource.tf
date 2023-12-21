locals {
  component_names = toset([
    "Delta",
    "Gamma",
    "Beta",
    "Eta",
  ])
}

resource "hund_group" "group0" {
  name = "Terraform Group"
}

resource "hund_group_component_ordering" "group0" {
  group = resource.hund_group.group0.id

  # Sort components alphabetically by name
  components = values({ for c in resource.hund_component.components : c.name => c.id })
}

resource "hund_component" "components" {
  for_each = local.component_names

  name        = each.key
  description = "Created by Terraform"

  group = resource.hund_group.group0.id

  watchdog = { service = { manual = {} } }
}
