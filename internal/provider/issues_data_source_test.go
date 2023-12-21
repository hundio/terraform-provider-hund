package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIssuesDataSource(t *testing.T) {
	datum := time.Now()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIssuesDataSourceConfig(datum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_issues.test", "issues.#", "3"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issues.test", "issues.*.id", "hund_issue.standing", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issues.test", "issues.*.id", "hund_issue.scheduled", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issues.test", "issues.*.id", "hund_issue.resolved", "id"),
				),
			},
		},
	})
}

func TestAccIssuesDataSource_standing(t *testing.T) {
	datum := time.Now()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIssuesDataSourceConfig_standing(datum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_issues.test", "issues.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issues.test", "issues.*.id", "hund_issue.standing", "id"),
				),
			},
		},
	})
}

func TestAccIssuesDataSource_upcoming(t *testing.T) {
	datum := time.Now()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIssuesDataSourceConfig_upcoming(datum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_issues.test", "issues.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issues.test", "issues.*.id", "hund_issue.scheduled", "id"),
				),
			},
		},
	})
}

func TestAccIssuesDataSource_resolved(t *testing.T) {
	datum := time.Now()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIssuesDataSourceConfig_resolved(datum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_issues.test", "issues.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issues.test", "issues.*.id", "hund_issue.resolved", "id"),
				),
			},
		},
	})
}

func TestAccIssuesDataSource_components(t *testing.T) {
	datum := time.Now()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIssuesDataSourceConfig_components(datum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_issues.test", "issues.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issues.test", "issues.*.id", "hund_issue.standing", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_issues.test", "issues.*.id", "hund_issue.resolved", "id"),
				),
			},
		},
	})
}

func testAccIssuesDataSourceConfig_base(datum time.Time) string {
	return providerConfig + fmt.Sprintf(`
		resource "hund_group" "test" {
			name = "Test Group"
		}

		resource "hund_component" "test" {
			group = hund_group.test.id
			name = "Test Component"

			watchdog = {service = {manual = {}}}
		}

		resource "hund_component" "test1" {
			group = hund_group.test.id
			name = "Test Component"

			watchdog = {service = {manual = {}}}
		}

		resource "hund_issue" "standing" {
			component_ids = [hund_component.test.id, hund_component.test1.id]

			title = "Test Issue Standing"
			body = "Test Body"
		}

		resource "hund_issue" "scheduled" {
			component_ids = [hund_component.test.id]

			title = "Test Issue Scheduled"
			body = "Test Body"

			schedule = {
				starts_at = %[1]q
				ends_at = %[2]q
			}
		}

		resource "hund_issue" "resolved" {
			component_ids = [hund_component.test.id, hund_component.test1.id]

			title = "Test Issue Resolved"
			body = "Test Body"

			began_at = %[3]q
			ended_at = %[4]q
		}
	`,
		testToTfTimestamp(datum.AddDate(0, 0, 1)), testToTfTimestamp(datum.AddDate(0, 0, 2)),
		testToTfTimestamp(datum.AddDate(0, 0, -3)), testToTfTimestamp(datum.AddDate(0, 0, -1)),
	)
}

func testAccIssuesDataSourceConfig(datum time.Time) string {
	return testAccIssuesDataSourceConfig_base(datum) + `
		data "hund_issues" "test" {
			depends_on = [
				hund_issue.standing,
				hund_issue.scheduled,
				hund_issue.resolved
			]
		}
	`
}

func testAccIssuesDataSourceConfig_standing(datum time.Time) string {
	return testAccIssuesDataSourceConfig_base(datum) + `
		data "hund_issues" "test" {
			depends_on = [
				hund_issue.standing,
				hund_issue.scheduled,
				hund_issue.resolved
			]

			standing = true
		}
	`
}

func testAccIssuesDataSourceConfig_upcoming(datum time.Time) string {
	return testAccIssuesDataSourceConfig_base(datum) + `
		data "hund_issues" "test" {
			depends_on = [
				hund_issue.standing,
				hund_issue.scheduled,
				hund_issue.resolved
			]

			upcoming = true
		}
	`
}

func testAccIssuesDataSourceConfig_resolved(datum time.Time) string {
	return testAccIssuesDataSourceConfig_base(datum) + `
		data "hund_issues" "test" {
			depends_on = [
				hund_issue.standing,
				hund_issue.scheduled,
				hund_issue.resolved
			]

			resolved = true
		}
	`
}

func testAccIssuesDataSourceConfig_components(datum time.Time) string {
	return testAccIssuesDataSourceConfig_base(datum) + `
		resource "hund_component" "test2" {
			group = hund_group.test.id
			name = "Test Component"

			watchdog = {service = {manual = {}}}
		}

		data "hund_issues" "test" {
			depends_on = [
				hund_issue.standing,
				hund_issue.scheduled,
				hund_issue.resolved
			]

			components = [
				hund_component.test1.id,
				hund_component.test2.id
			]
		}
	`
}
