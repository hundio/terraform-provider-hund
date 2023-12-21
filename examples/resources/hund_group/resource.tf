resource "hund_group" "group" {
  name = "Terraform Group"
}

resource "hund_component" "components" {
  name        = "Terraform Component"
  description = "Created by Terraform"

  group = resource.hund_group.group.id

  watchdog = { service = { manual = {} } }
}
