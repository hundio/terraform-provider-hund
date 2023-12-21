data "hund_issue_templates" "example" {
  kind = "update"
}

output "template_names" {
  value = [for t in data.hund_issue_templates.example.issue_templates : t.name]
}
