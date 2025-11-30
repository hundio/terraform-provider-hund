package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

func TestAccIssueUpdateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIssueUpdateResourceConfig("**body**"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue_update.test", "body_translations.en", "**body**"),
					resource.TestMatchResourceAttr("hund_issue_update.test", "body_html", regexp.MustCompile("<strong>body</strong>")),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hund_issue_update.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIssueUpdateResourceImportStateIdFunc,
			},
			// Update and Read testing
			{
				Config: testAccIssueUpdateResourceConfig("*body*"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue_update.test", "body_translations.en", "*body*"),
					resource.TestMatchResourceAttr("hund_issue_update.test", "body_html", regexp.MustCompile("<em>body</em>")),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIssueUpdateResource_archive(t *testing.T) {
	var apiUpdate hundApiV1.UpdateExpansionary

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIssueUpdateResourceConfig_archive(false),
				Check:  testAccIssueUpdateResourceCheckExistence("hund_issue_update.test", &apiUpdate),
			},
			// Delete Issue from State only
			{
				Config: testAccIssueUpdateResourceConfig_archive(true),
				Check:  testAccIssueUpdateResourceCheckArchival(&apiUpdate),
			},
		},
	})
}
func TestAccIssueUpdateResource_retrospective(t *testing.T) {
	datum := time.Now().AddDate(0, -1, 0)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIssueUpdateResourceConfig_retrospective(datum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("hund_issue_update.test", "body"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hund_issue_update.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIssueUpdateResourceImportStateIdFunc,
			},
			// Update and Read testing
			{
				Config: testAccIssueUpdateResourceConfig_retrospective(datum.AddDate(0, 0, -1)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("hund_issue_update.test", "body"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIssueUpdateResource_template(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIssueUpdateResourceConfig_template(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue_update.test", "body", "my text - 2023-11-15"),
					resource.TestCheckResourceAttr("hund_issue_update.test", "label", "monitoring"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hund_issue_update.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIssueUpdateResourceImportStateIdFunc,
			},
			// Update and Read testing
			{
				Config: testAccIssueUpdateResourceConfig_template(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue_update.test", "body", "overridden body 4"),
					resource.TestCheckResourceAttr("hund_issue_update.test", "label", "monitoring"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIssueUpdateResourceImportStateIdFunc(s *terraform.State) (string, error) {
	update, ok := s.RootModule().Resources["hund_issue_update.test"]

	if !ok {
		return "", errors.New("could not find hund_issue_update.test resource")
	}

	return update.Primary.Attributes["issue_id"] + "/" + update.Primary.Attributes["id"], nil
}

func testAccIssueUpdateResourceConfig(body string) string {
	return providerConfig + fmt.Sprintf(`
resource "hund_group" "test" {
	name = "Test Group"
}

resource "hund_component" "test" {
	group = hund_group.test.id
	name = "Test Component"

	watchdog = {service = {manual = {}}}
}

resource "hund_issue" "test" {
	component_ids = [hund_component.test.id]

  title = "Test Issue"
	body = "Test Body"
}

resource "hund_issue_update" "test" {
	issue_id = hund_issue.test.id

	body = %[1]q
}
`, body)
}

func testAccIssueUpdateResourceConfig_archive(archive bool) string {
	var updateResource string

	if !archive {
		updateResource = `
		resource "hund_issue_update" "test" {
			archive_on_destroy = true

			issue_id = hund_issue.test.id

			body = "Test Body"
		}
		`
	}

	return providerConfig + fmt.Sprintf(`
	resource "hund_group" "test" {
		name = "Test Group"
	}

	resource "hund_component" "test" {
		group = hund_group.test.id
		name = "Test Component"

		watchdog = {service = {manual = {}}}
	}

	resource "hund_issue" "test" {
		component_ids = [hund_component.test.id]

		title = "Test Issue"
		body = "Test Body"
	}

	%[1]v
	`, updateResource)
}

func testAccIssueUpdateResourceConfig_retrospective(datum time.Time) string {
	return providerConfig + fmt.Sprintf(`
resource "hund_group" "test" {
	name = "Test Group"
}

resource "hund_component" "test" {
	group = hund_group.test.id
	name = "Test Component"

	watchdog = {service = {manual = {}}}
}

resource "hund_issue" "test" {
	component_ids = [hund_component.test.id]

  title = "Test Issue"
	body = "Test Body"

	began_at = %[1]q
}

resource "hund_issue_update" "test" {
	issue_id = hund_issue.test.id

	label = "resolved"

	effective_after = %[2]q
}
`, testToTfTimestamp(datum), testToTfTimestamp(datum.AddDate(0, 0, 1)))
}

func testAccIssueUpdateResourceConfig_template(overrideBody bool) string {
	var bodyOverride string

	if overrideBody {
		bodyOverride = `body = "overridden body {{vars.ordinal}}"`
	}

	return providerConfig + fmt.Sprintf(`
	resource "hund_group" "test" {
		name = "Test Group"
	}

	resource "hund_component" "test" {
		group = hund_group.test.id
		name = "Test Component"

		watchdog = {service = {manual = {}}}
	}

	resource "hund_issue_template" "test" {
		kind = "update"
		name = "Test Update Template"

		body = "{{vars.summary}} - {{vars.timestamp | date: '%%F' }}"
		label = "monitoring"

		variables = {
			summary = { required = true }
			ordinal = { type = "number" }
			timestamp = { type = "datetime" }
			unused  = { type = "i18n-string" }
		}
	}

	resource "hund_issue" "test" {
		component_ids = [hund_component.test.id]

		title = "Test Issue"
		body = "Test Body"
	}

	resource "hund_issue_update" "test" {
		issue_id = hund_issue.test.id

		template = {
			issue_template_id = hund_issue_template.test.id

			%[1]v

			variables = {
				summary = { string = "my text" }
				ordinal = { number = 4 }
				timestamp = { datetime = "2023-11-15T19:26:24Z" }
				unused = { i18n_string = { en = "one", de = "ein" } }
			}
		}
	}
	`, bodyOverride)
}

func testAccIssueUpdateResourceCheckExistence(resourceName string, apiUpdate *hundApiV1.UpdateExpansionary) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		update, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource name not found: %s", resourceName)
		}

		id, ok := update.Primary.Attributes["id"]
		if !ok {
			return fmt.Errorf("could not find id attribute: %s", resourceName)
		}

		issueId, ok := update.Primary.Attributes["issue_id"]
		if !ok {
			return fmt.Errorf("could not find issue_id attribute: %s", resourceName)
		}

		api, err := testAccIssueUpdateResourceRequestApiIssue(issueId, id)
		if err != nil {
			return err
		}

		*apiUpdate = *api

		return nil
	}
}

func testAccIssueUpdateResourceCheckArchival(apiUpdate *hundApiV1.UpdateExpansionary) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		issueId, err := apiUpdate.Issue.AsUpdateExpansionaryIssue0()
		if err != nil {
			return err
		}

		_, err = testAccIssueUpdateResourceRequestApiIssue(issueId, apiUpdate.Id)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccIssueUpdateResourceRequestApiIssue(issueId string, id string) (*hundApiV1.UpdateExpansionary, error) {
	client, err := sharedClientForDomain("default")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	rsp, err := client.RetrieveAUpdate(ctx, issueId, id)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("could not find issue id %s (status code: %v)", id, rsp.StatusCode)
	}

	update, err := hundApiV1.ParseRetrieveAUpdateResponse(rsp)
	if err != nil {
		return nil, err
	}

	return update.HALJSON200, nil
}
