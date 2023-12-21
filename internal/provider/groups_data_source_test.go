package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccGroupsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_groups.test", "groups.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_groups.test", "groups.*.id", "hund_group.test0", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_groups.test", "groups.*.id", "hund_group.test1", "id"),
					resource.TestCheckResourceAttr("data.hund_groups.test", "groups.0.components.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_groups.test", "groups.*.components.*", "hund_component.test0", "id"),
					resource.TestCheckResourceAttr("data.hund_groups.test", "groups.1.components.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_groups.test", "groups.*.components.*", "hund_component.test1", "id"),
				),
			},
		},
	})
}

func testAccGroupsDataSourceConfig_base() string {
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

func testAccGroupsDataSourceConfig() string {
	return testAccGroupsDataSourceConfig_base() + `
		data "hund_groups" "test" {
			depends_on = [
				hund_component.test0,
				hund_component.test1
			]
		}
	`
}
