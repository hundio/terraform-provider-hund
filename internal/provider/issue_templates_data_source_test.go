package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIssueTemplatesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIssueTemplatesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_issue_templates.test", "issue_templates.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issue_templates.test", "issue_templates.*.id", "hund_issue_template.issue", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issue_templates.test", "issue_templates.*.id", "hund_issue_template.update", "id"),
				),
			},
		},
	})
}

func TestAccIssueTemplatesDataSource_issue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIssueTemplatesDataSourceConfig_issue(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_issue_templates.test", "issue_templates.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issue_templates.test", "issue_templates.*.id", "hund_issue_template.issue", "id"),
				),
			},
		},
	})
}

func testAccIssueTemplatesDataSourceConfig_base() string {
	return providerConfig + `
resource "hund_issue_template" "issue" {
  name = "Issue Template Test"

  kind  = "issue"
  title = "Test Template {{vars.ordinal}}"

	body_translations = {
		en = "different summary: {{vars.unused}}"
		de = "andere Zusammenfassung: {{vars.unused}}"
		pl = "inne podsumowanie: {{vars.unused}}"
		original = "en"
	}

  label = "assessed"

  variables = {
    summary = { required = true }
    ordinal = { type = "number" }
		date = { type = "datetime" }
    unused  = { type = "i18n-string" }
  }
}

resource "hund_issue_template" "update" {
  name = "Issue Update Template Test"

  kind  = "update"

	body_translations = {
		en = "different summary: {{vars.unused}}"
		de = "andere Zusammenfassung: {{vars.unused}}"
		pl = "inne podsumowanie: {{vars.unused}}"
		original = "en"
	}

  label = "monitoring"

  variables = {
    summary = { required = true }
    ordinal = { type = "number" }
		date = { type = "datetime" }
    unused  = { type = "i18n-string" }
  }
}
`
}

func testAccIssueTemplatesDataSourceConfig() string {
	return testAccIssueTemplatesDataSourceConfig_base() + `
	data "hund_issue_templates" "test" {
		depends_on = [hund_issue_template.issue, hund_issue_template.update]
	}
	`
}

func testAccIssueTemplatesDataSourceConfig_issue() string {
	return testAccIssueTemplatesDataSourceConfig_base() + `
	data "hund_issue_templates" "test" {
		depends_on = [hund_issue_template.issue, hund_issue_template.update]

		kind = "issue"
	}
	`
}
