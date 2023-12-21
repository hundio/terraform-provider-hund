data "hund_groups" "example" {
}

data "hund_components" "components" {
  group = one([for g in data.hund_groups.example.groups : g.id if g.name == "Terraform Group"])
}

output "component_names" {
  value = data.hund_components.components.components[*].name
}
