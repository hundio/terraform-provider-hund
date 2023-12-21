resource "hund_issue" "example" {
  archive_on_destroy = true

  title = "a terraform issue"
  body_translations = {
    en       = "**lol** ok.\nnah"
    de       = "**lul** ok.\nnein"
    original = "en"
  }
  component_ids = ["5d72d51f8fbb65b5d3a587e1"]
  began_at      = "2023-09-23T13:00:00Z"
}

resource "hund_issue_update" "example0" {
  archive_on_destroy = true

  issue_id        = resource.hund_issue.example.id
  effective_after = timeadd(resource.hund_issue.example.began_at, "2h")
  body            = "heads up"
}

resource "hund_issue_update" "example1" {
  archive_on_destroy = true

  depends_on = [hund_issue_update.example0]

  issue_id = resource.hund_issue.example.id
  body     = "keeping watch (shh)"
  label    = "monitoring"
}

resource "hund_issue_update" "example2" {
  archive_on_destroy = true

  depends_on = [hund_issue_update.example1]

  issue_id = resource.hund_issue.example.id
  label    = "resolved"
}
