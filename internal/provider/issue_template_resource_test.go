package provider

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

func init() {
	resource.AddTestSweepers("hund_issue_template", &resource.Sweeper{
		Name: "hund_issue_template",
		F: func(r string) error {
			client, err := sharedClientForDomain(r)
			if err != nil {
				return err
			}

			ctx := context.Background()

			limit := 100
			rsp, err := client.GetAllIssueTemplates(ctx, &hundApiV1.GetAllIssueTemplatesParams{
				Limit: &limit,
			})
			if err != nil {
				return err
			}

			templates, err := hundApiV1.ParseGetAllIssueTemplatesResponse(rsp)
			if err != nil {
				return err
			}
			if templates.StatusCode() != 200 {
				return errors.New("couldn't retrieve issue_templates:\n" + string(templates.Body))
			}

			for _, c := range templates.HALJSON200.Data {
				rsp, err := client.DeleteAIssueTemplate(ctx, c.Id)
				if err != nil {
					return err
				}
				if rsp.StatusCode != 204 {
					template, err := hundApiV1.ParseDeleteAIssueTemplateResponse(rsp)
					if err != nil {
						return err
					}

					return errors.New("couldn't delete issue_template:\n" + string(template.Body))
				}
			}

			return nil
		},
	})
}

func TestAccIssueTemplateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIssueTemplateResourceConfig("Template One"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue_template.test", "body_translations.en", "Test Template Body {{vars.summary}}"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hund_issue_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIssueTemplateResourceConfigI18n("Template Two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue_template.test", "body", "different summary: {{vars.unused}}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIssueTemplateResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "hund_issue_template" "test" {
  name = %[1]q

  kind  = "issue"
  title = "Test Template {{vars.ordinal}}"
  body  = "Test Template Body {{vars.summary}}"

  label = "assessed"

  variables = {
    summary = { required = true }
    ordinal = { type = "number" }
		date = { type = "datetime" }
    unused  = { type = "i18n-string" }
  }
}
`, name)
}

func testAccIssueTemplateResourceConfigI18n(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "hund_issue_template" "test" {
  name = %[1]q

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
`, name)
}
