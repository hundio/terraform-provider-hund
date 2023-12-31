---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "hund_issue_template Resource - terraform-provider-hund"
subcategory: ""
description: |-
  Resource representing templates for creating Issues, as well as Issue Updates.
---

# hund_issue_template (Resource)

Resource representing templates for creating Issues, as well as Issue Updates.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `kind` (String) The "kind" of this IssueTemplate. This field can be either `issue` or `update`, depending on whether this IssueTemplate can be applied to an Issue or Issue Update, respectively.
- `name` (String) An internal name for identifying this IssueTemplate.

### Optional

- `body` (String) The body to use for an Issue/Update applied against this template. This field supports [Liquid templating](https://shopify.github.io/liquid/), in the default translation.
- `body_translations` (Map of String) The body to use for an Issue/Update applied against this template. This field supports [Liquid templating](https://shopify.github.io/liquid/), translated into multiple languages. Map keys express the language each string value is to be interpreted in. The `original` field of this map denotes the language used for the non-`_translations` version of this attribute.
- `label` (String) The label to use for an Issue/Update applied against this template.
- `title` (String) When `kind` is `issue`, then the applied Issue will take on this title. This field supports [Liquid templating](https://shopify.github.io/liquid/), in the default translation.
- `title_translations` (Map of String) When `kind` is `issue`, then the applied Issue will take on this title. This field supports [Liquid templating](https://shopify.github.io/liquid/), translated into multiple languages. Map keys express the language each string value is to be interpreted in. The `original` field of this map denotes the language used for the non-`_translations` version of this attribute.
- `variables` (Attributes Map) An object defining a set of typed variables that can be provided in an application of this IssueTemplate. The variables can be accessed from any field in the IssueTemplate supporting Liquid. (see [below for nested schema](#nestedatt--variables))

### Read-Only

- `created_at` (String) The timestamp at which this IssueTemplate was created.
- `id` (String) The ObjectId of this IssueTemplate.
- `updated_at` (String) The timestamp at which this IssueTemplate was last updated.

<a id="nestedatt--variables"></a>
### Nested Schema for `variables`

Optional:

- `required` (Boolean) Whether this variable is required when applying the template to an Issue/Update.
- `type` (String) The expected type of this variable. One of `datetime`, `i18n-string`, `number`, or `string`.
