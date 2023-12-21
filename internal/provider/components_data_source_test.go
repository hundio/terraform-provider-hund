package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccComponentsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccComponentsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_components.test", "components.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_components.test", "components.*.id", "hund_component.test0", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_components.test", "components.*.id", "hund_component.test1", "id"),
				),
			},
		},
	})
}

func TestAccComponentsDataSource_group(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccComponentsDataSourceConfig_group(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_components.test", "components.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_components.test", "components.*.id", "hund_component.test0", "id"),
				),
			},
		},
	})
}

func TestAccComponentsDataSource_issue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccComponentsDataSourceConfig_issue(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_components.test", "components.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_components.test", "components.*.id", "hund_component.test1", "id"),
				),
			},
		},
	})
}

func testAccComponentsDataSourceConfig_base() string {
	return providerConfig + `
		resource "hund_group" "test0" {
			name = "Test Group 1"
		}

		resource "hund_group" "test1" {
			name = "Test Group 2"
		}

		resource "hund_component" "test0" {
			group = hund_group.test0.id
			name = "Test Component"

			watchdog = {service = {manual = {}}}
		}

		resource "hund_component" "test1" {
			group = hund_group.test1.id
			name = "Test Component"

			watchdog = {service = {manual = {}}}
		}
	`
}

func testAccComponentsDataSourceConfig() string {
	return testAccComponentsDataSourceConfig_base() + `
		data "hund_components" "test" {
			depends_on = [
				hund_component.test0,
				hund_component.test1
			]
		}
	`
}

func testAccComponentsDataSourceConfig_group() string {
	return testAccComponentsDataSourceConfig_base() + `
	data "hund_components" "test" {
		depends_on = [
			hund_component.test0,
			hund_component.test1
		]

		group = hund_group.test0.id
	}
	`
}

func testAccComponentsDataSourceConfig_issue() string {
	return testAccComponentsDataSourceConfig_base() + `
	resource "hund_issue" "test" {
		title = "a terraform issue"
		body = "**lol** ok.\nnah"

		component_ids = [hund_component.test1.id]
	}

	data "hund_components" "test" {
		depends_on = [
			hund_component.test0,
			hund_component.test1
		]

		issue = hund_issue.test.id
	}
	`
}
