package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

func TestAccIssueResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIssueResourceConfig("one", "**body**"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue.test", "title_translations.en", "one"),
					resource.TestMatchResourceAttr("hund_issue.test", "body_html", regexp.MustCompile("<strong>body</strong>")),
				),
			},
			// ImportState testing
			{
				ResourceName:            "hund_issue.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration"},
			},
			// Update and Read testing
			{
				Config: testAccIssueResourceConfig("two", "*body*"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue.test", "title_translations.en", "two"),
					resource.TestMatchResourceAttr("hund_issue.test", "body_html", regexp.MustCompile("<em>body</em>")),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIssueResource_archive(t *testing.T) {
	var apiIssue hundApiV1.Issue

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIssueResourceConfig_archive(false),
				Check:  testAccIssueResourceCheckExistence("hund_issue.test", &apiIssue),
			},
			// Delete Issue from State only
			{
				Config: testAccIssueResourceConfig_archive(true),
				Check:  testAccIssueResourceCheckArchival(&apiIssue),
			},
		},
	})
}

func TestAccIssueResource_scheduled(t *testing.T) {
	datum := time.Now().AddDate(0, 0, 1)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIssueResourceConfig_scheduled(datum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue.test", "scheduled", "true"),
					resource.TestCheckNoResourceAttr("hund_issue.test", "schedule.notify_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "hund_issue.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration"},
			},
			// Update and Read testing
			{
				Config: testAccIssueResourceConfig_scheduled(datum.AddDate(0, 0, 1)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue.test", "scheduled", "true"),
					resource.TestCheckNoResourceAttr("hund_issue.test", "schedule.notify_at"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIssueResource_retrospective(t *testing.T) {
	datum := time.Now().AddDate(0, 0, -1)
	datumEnd := datum.AddDate(0, 0, 1)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIssueResourceConfig_retrospective(datum, datumEnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue.test", "retrospective", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "hund_issue.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration"},
			},
			// Update and Read testing
			{
				Config: testAccIssueResourceConfig_retrospective(datum.AddDate(0, 0, -1), datumEnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue.test", "retrospective", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIssueResource_template(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIssueResourceConfig_template(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue.test", "title", "Test Title 4"),
					resource.TestCheckResourceAttr("hund_issue.test", "body", "my text - 2023-11-15"),
					resource.TestCheckResourceAttr("hund_issue.test", "label", "assessed"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "hund_issue.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration"},
			},
			// Update and Read testing
			{
				Config: testAccIssueResourceConfig_template(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_issue.test", "title", "Test Title 4"),
					resource.TestCheckResourceAttr("hund_issue.test", "body", "overridden body 4"),
					resource.TestCheckResourceAttr("hund_issue.test", "label", "assessed"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIssueResourceConfig(title string, body string) string {
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

  title = %[1]q
	body = %[2]q
}
`, title, body)
}

func testAccIssueResourceConfig_archive(archive bool) string {
	var issueResource string

	if !archive {
		issueResource = `
		resource "hund_issue" "test" {
			archive_on_destroy = true

			component_ids = [hund_component.test.id]

			title = "Test Issue"
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

	%[1]v
	`, issueResource)
}

func testAccIssueResourceConfig_scheduled(datum time.Time) string {
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

		title = "Test Scheduled Issue"
		body = "Test Body"

		schedule = {
			starts_at = %[1]q
			ends_at = timeadd(%[1]q, "2h")
		}
	}
	`, testToTfTimestamp(datum))
}

func testAccIssueResourceConfig_retrospective(began time.Time, ended time.Time) string {
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

		title = "Test Scheduled Issue"
		body = "Test Body"

		began_at = %[1]q
		ended_at = %[2]q
	}
	`, testToTfTimestamp(began), testToTfTimestamp(ended))
}

func testAccIssueResourceConfig_template(overrideBody bool) string {
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
		kind = "issue"
		name = "Test Template"

		title = "Test Title {{vars.ordinal}}"
		body = "{{vars.summary}} - {{vars.timestamp | date: '%%F' }}"
		label = "assessed"

		variables = {
			summary = { required = true }
			ordinal = { type = "number" }
			timestamp = { type = "datetime" }
			unused  = { type = "i18n-string" }
		}
	}

	resource "hund_issue" "test" {
		component_ids = [hund_component.test.id]

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

func testAccIssueResourceCheckExistence(resourceName string, apiIssue *hundApiV1.Issue) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		issue, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource name not found: %s", resourceName)
		}

		id, ok := issue.Primary.Attributes["id"]
		if !ok {
			return fmt.Errorf("could not find id attribute: %s", resourceName)
		}

		api, err := testAccIssueResourceRequestApiIssue(id)
		if err != nil {
			return err
		}

		*apiIssue = *api

		return nil
	}
}

func testAccIssueResourceCheckArchival(apiIssue *hundApiV1.Issue) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := testAccIssueResourceRequestApiIssue(apiIssue.Id)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccIssueResourceRequestApiIssue(id string) (*hundApiV1.Issue, error) {
	client, err := sharedClientForDomain("default")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	rsp, err := client.RetrieveAIssue(ctx, id)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("could not find issue id %s (status code: %v)", id, rsp.StatusCode)
	}

	issue, err := hundApiV1.ParseRetrieveAIssueResponse(rsp)
	if err != nil {
		return nil, err
	}

	return issue.HALJSON200, nil
}
