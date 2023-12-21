resource "hund_issue_template" "template" {
  name  = "a terraform template"
  kind  = "issue"
  title = "lol {{vars.ordinal}}"
  body  = "yeah {{vars.summary}}"

  label = "assessed"

  variables = {
    summary = { required = true }
    ordinal = { type = "number" }
    unused  = { type = "i18n-string" }
  }
}

resource "hund_issue_template" "update_template" {
  name  = "a terraform template (updates)"
  kind  = "update"
  body  = "update: {{vars.summary}}"
  label = "monitoring"

  variables = {
    summary = { required = true, type = "i18n-string" }
  }
}

# Example application of Issue Template
resource "hund_issue" "example" {
  component_ids = ["5d72d51f8fbb65b5d3a587e1"]

  template = {
    issue_template_id = resource.hund_issue_template.template.id

    body_translations = {
      en       = "different summary: {{vars.unused}}"
      de       = "andere Zusammenfassung: {{vars.unused}}"
      pl       = "inne podsumowanie: {{vars.unused}}"
      original = "en"
    }

    variables = {
      summary = { string = "test summary" }
      ordinal = { number = 27 }
      unused = { i18n_string = {
        en = "lol"
        de = "lül"
        pl = "łuł"
      } }
    }
  }
}

# Example application of Issue Update Template
resource "hund_issue_update" "example0" {
  archive_on_destroy = true

  issue_id = resource.hund_issue.example.id

  label = "investigating"

  template = {
    issue_template_id = resource.hund_issue_template.update_template.id

    variables = {
      summary = { i18n_string = {
        en = "lol"
      } }
    }
  }
}
