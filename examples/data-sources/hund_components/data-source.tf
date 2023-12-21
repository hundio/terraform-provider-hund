data "hund_components" "example" {
  group = "654c04c48fbb65547077cb2f"
}

output "component_names" {
  value = [for c in data.hund_components.example.components : c.name]
}
