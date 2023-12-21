package provider

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMetricProviderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetricProviderResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchTypeSetElemNestedAttrs("hund_metric_provider.test", "instances.*", map[string]*regexp.Regexp{
						"slug": regexp.MustCompile("^percent_uptime"),
					}),
					resource.TestMatchTypeSetElemNestedAttrs("hund_metric_provider.test", "instances.*", map[string]*regexp.Regexp{
						"slug": regexp.MustCompile("^incidents_reported"),
					}),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hund_metric_provider.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccMetricProviderResourceConfig(false),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccMetricProviderResource_generalProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetricProviderResourceConfig_generalProvider(true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchTypeSetElemNestedAttrs("hund_metric_provider.test", "instances.*", map[string]*regexp.Regexp{
						"slug": regexp.MustCompile("^metric0"),
					}),
					resource.TestMatchTypeSetElemNestedAttrs("hund_metric_provider.test", "instances.*", map[string]*regexp.Regexp{
						"slug": regexp.MustCompile("^metric1"),
					}),
					resource.TestMatchTypeSetElemNestedAttrs("hund_metric_provider.test", "instances.*", map[string]*regexp.Regexp{
						"slug": regexp.MustCompile("^metric2"),
					}),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hund_metric_provider.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccMetricProviderResourceConfig_generalProvider(false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchTypeSetElemNestedAttrs("hund_metric_provider.test", "instances.*", map[string]*regexp.Regexp{
						"slug": regexp.MustCompile("^metric0"),
					}),
					resource.TestMatchTypeSetElemNestedAttrs("hund_metric_provider.test", "instances.*", map[string]*regexp.Regexp{
						"slug": regexp.MustCompile("^metric1"),
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccMetricProviderResource_serviceChange(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetricProviderResourceConfig_serviceChange(false),
			},
			// Update and Read testing
			{
				Config: testAccMetricProviderResourceConfig_serviceChange(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hund_metric_provider.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func TestAccMetricProviderResource_defaultProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricProviderResourceBaseConfig_defaultProvider,
			},
			{
				Config:      testAccMetricProviderResourceConfig_defaultProvider(true),
				ExpectError: regexp.MustCompile("Cannot create Default MetricProvider"),
			},
			{
				Config:             testAccMetricProviderResourceConfig_defaultProvider(true),
				ResourceName:       "hund_metric_provider.test",
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateId:      "default",
				ImportStateIdFunc:  testAccMetricProviderResourceImportStateIdFunc_defaultProvider,
			},
			{
				Config: testAccMetricProviderResourceConfig_defaultProvider(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_metric_provider.test", "service.icmp.target", "example.com"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hund_metric_provider.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccMetricProviderResourceConfig_defaultProvider(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hund_metric_provider.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func TestAccMetricProviderResource_incompleteInstances(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMetricProviderResourceConfig_incompleteInstances(false),
				ExpectError: regexp.MustCompile("Missing instance from MetricProvider"),
			},
			{
				Config:      testAccMetricProviderResourceConfig_incompleteInstances(true),
				ExpectError: regexp.MustCompile("Extraneous instance in MetricProvider"),
			},
		},
	})
}

func TestAccMetricProviderResource_autoInstances(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricProviderResourceConfig_autoInstances(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_metric_provider.test", "instances.%", "2"),
					resource.TestMatchTypeSetElemNestedAttrs("hund_metric_provider.test", "instances.*", map[string]*regexp.Regexp{
						"slug": regexp.MustCompile("^incidents_reported"),
					}),
					resource.TestMatchTypeSetElemNestedAttrs("hund_metric_provider.test", "instances.*", map[string]*regexp.Regexp{
						"slug": regexp.MustCompile("^percent_uptime"),
					}),
					resource.TestCheckResourceAttr("hund_metric_provider.test_webhook", "instances.%", "0"),
				),
			},
		},
	})
}

const testAccMetricProviderResourceBaseConfig = providerConfig + `
resource "hund_group" "test" {
	name = "Test Group"
}

resource "hund_component" "test" {
  name = "Test Component"

	group = hund_group.test.id

	watchdog = {service = { manual = {state  = 1} }}
}
`

func testAccMetricProviderResourceConfig(enabled bool) string {
	return testAccMetricProviderResourceBaseConfig + fmt.Sprintf(`
resource "hund_metric_provider" "test" {
	watchdog = hund_component.test.watchdog.id

	service = { builtin = {} }

	instances = {
		percent_uptime = { enabled = %[1]v },
		incidents_reported = { enabled = %[1]v }
	}
}
`, enabled)
}

func testAccMetricProviderResourceConfig_generalProvider(enabled bool, includeMetric2 bool) string {
	var metric2Config string

	if includeMetric2 {
		metric2Config = `metric2 = { enabled = true, title = "Test Metric 2" }`
	}

	return testAccMetricProviderResourceBaseConfig + fmt.Sprintf(`
resource "hund_metric_provider" "test" {
	watchdog = hund_component.test.watchdog.id

	service = { webhook = {} }

	instances = {
		metric0 = { enabled = %[1]v, title = "Test Metric 0" },
		metric1 = { enabled = %[1]v, title = "Test Metric 1" },
		%[2]v
	}
}
`, enabled, metric2Config)
}

func testAccMetricProviderResourceConfig_serviceChange(changeService bool) string {
	var service, instances string

	if changeService {
		service = `{ webhook = {} }`
		instances = `
metric0 = { enabled = true, title = "Test Metric 0" }
		`
	} else {
		service = `{ builtin = {} }`
		instances = `
percent_uptime = { enabled = true },
incidents_reported = { enabled = true }
		`
	}

	return testAccMetricProviderResourceBaseConfig + fmt.Sprintf(`
resource "hund_metric_provider" "test" {
	watchdog = hund_component.test.watchdog.id

	service = %[1]v

	instances = {
		%[2]v
	}
}
`, service, instances)
}

const testAccMetricProviderResourceBaseConfig_defaultProvider = providerConfig + `
resource "hund_group" "test" {
	name = "Test Group"
}

resource "hund_component" "test" {
	name = "Test Component"

	group = hund_group.test.id

	watchdog = {service = { icmp = { target = "example.com", regions = ["wa-us-1"] } }}
}
`

func testAccMetricProviderResourceConfig_defaultProvider(useDefault bool) string {
	var service string

	if useDefault {
		service = `
service = { icmp = hund_component.test.watchdog.service.icmp }
`
	} else {
		service = `
service = { icmp = {
	target = "other.example.com",
	regions = ["wa-us-1"]
} }
`
	}

	return testAccMetricProviderResourceBaseConfig_defaultProvider + fmt.Sprintf(`

	resource "hund_metric_provider" "test" {
		watchdog = hund_component.test.watchdog.id
		default = %[1]v

		%[2]v

		instances = {
			res = {},
			"icmp.total_addresses" = { enabled = false },
			"icmp.passed_addresses" = { enabled = false },
		}
	}
`, useDefault, service)
}

func testAccMetricProviderResourceConfig_incompleteInstances(extraneous bool) string {
	var instances string

	if extraneous {
		instances = `
incidents_reported = {},
extra_metric = {},
		`
	}

	return testAccMetricProviderResourceBaseConfig + fmt.Sprintf(`
resource "hund_metric_provider" "test" {
	watchdog = hund_component.test.watchdog.id

	service = { builtin = {} }

	instances = {
		percent_uptime = {},
		%[1]v
	}
}
`, instances)
}

func testAccMetricProviderResourceConfig_autoInstances() string {
	return testAccMetricProviderResourceBaseConfig + `
resource "hund_metric_provider" "test" {
	watchdog = hund_component.test.watchdog.id

	service = { builtin = {} }
}

resource "hund_metric_provider" "test_webhook" {
	watchdog = hund_component.test.watchdog.id

	service = { webhook = {} }
}
`
}

func testAccMetricProviderResourceImportStateIdFunc_defaultProvider(s *terraform.State) (string, error) {
	component, ok := s.RootModule().Resources["hund_component.test"]

	if !ok {
		return "", errors.New("could not find hund_component.test resource")
	}

	return "default/" + component.Primary.Attributes["watchdog.id"], nil
}
