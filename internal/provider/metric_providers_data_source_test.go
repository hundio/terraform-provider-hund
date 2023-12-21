package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetricProvidersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccMetricProvidersDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_metric_providers.test", "metric_providers.#", "3"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_metric_providers.test", "metric_providers.*.id", "hund_metric_provider.test0", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_metric_providers.test", "metric_providers.*.id", "hund_metric_provider.test1", "id"),
				),
			},
		},
	})
}

func TestAccMetricProvidersDataSource_watchdog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccMetricProvidersDataSourceConfig_watchdog(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_metric_providers.test", "metric_providers.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("data.hund_metric_providers.test", "metric_providers.*.id", "hund_metric_provider.test0", "id"),
				),
			},
		},
	})
}

func TestAccMetricProvidersDataSource_default(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccMetricProvidersDataSourceConfig_default(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hund_metric_providers.test", "metric_providers.#", "1"),
					resource.TestCheckResourceAttrPair("data.hund_metric_providers.test", "metric_providers.0.watchdog", "hund_component.test1", "watchdog.id"),
					resource.TestCheckResourceAttr("data.hund_metric_providers.test", "metric_providers.0.default", "true"),
				),
			},
		},
	})
}

func testAccMetricProvidersDataSourceConfig_base() string {
	return providerConfig + `
		resource "hund_group" "test" {
			name = "Test Group 1"
		}

		resource "hund_component" "test0" {
			group = hund_group.test.id
			name = "Test Component"

			watchdog = {service = {manual = {}}}
		}

		resource "hund_component" "test1" {
			group = hund_group.test.id
			name = "Test Component"

			watchdog = {service = {webhook = {}}}
		}

		resource "hund_metric_provider" "test0" {
			watchdog = hund_component.test0.watchdog.id

			service = { builtin = {} }
		}

		resource "hund_metric_provider" "test1" {
			watchdog = hund_component.test1.watchdog.id

			service = { builtin = {} }
		}
	`
}

func testAccMetricProvidersDataSourceConfig() string {
	return testAccMetricProvidersDataSourceConfig_base() + `
		data "hund_metric_providers" "test" {
			depends_on = [
				hund_metric_provider.test0,
				hund_metric_provider.test1
			]
		}
	`
}

func testAccMetricProvidersDataSourceConfig_watchdog() string {
	return testAccMetricProvidersDataSourceConfig_base() + `
		data "hund_metric_providers" "test" {
			depends_on = [
				hund_metric_provider.test0,
				hund_metric_provider.test1
			]

			watchdog = hund_component.test0.watchdog.id
		}
	`
}

func testAccMetricProvidersDataSourceConfig_default() string {
	return testAccMetricProvidersDataSourceConfig_base() + `
		data "hund_metric_providers" "test" {
			depends_on = [
				hund_metric_provider.test0,
				hund_metric_provider.test1
			]

			default = true
		}
	`
}
