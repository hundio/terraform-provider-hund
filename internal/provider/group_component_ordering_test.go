package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupComponentOrderingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGroupComponentOrderingResourceConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_group_component_ordering.test", "components.#", "3"),
					resource.TestCheckResourceAttrPair("hund_group_component_ordering.test", "components.0", "hund_component.beta", "id"),
					resource.TestCheckResourceAttrPair("hund_group_component_ordering.test", "components.1", "hund_component.alpha", "id"),
					resource.TestCheckResourceAttrPair("hund_group_component_ordering.test", "components.2", "hund_component.delta", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hund_group_component_ordering.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccGroupComponentOrderingResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_group_component_ordering.test", "components.#", "3"),
					resource.TestCheckResourceAttrPair("hund_group_component_ordering.test", "components.0", "hund_component.alpha", "id"),
					resource.TestCheckResourceAttrPair("hund_group_component_ordering.test", "components.1", "hund_component.beta", "id"),
					resource.TestCheckResourceAttrPair("hund_group_component_ordering.test", "components.2", "hund_component.delta", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccGroupComponentOrderingResourceConfig(sortByName bool) string {
	var ordering string

	if sortByName {
		ordering = `values({for c in [hund_component.beta, hund_component.alpha, hund_component.delta] : c.name => c.id })`
	} else {
		ordering = `[for c in [hund_component.beta, hund_component.alpha, hund_component.delta] : c.id]`
	}

	return providerConfig + fmt.Sprintf(`
resource "hund_group" "test" {
  name = "testing group"
}

resource "hund_component" "beta" {
  name = "Beta"
  group = resource.hund_group.test.id

  watchdog = { service = { manual = {} } }
}

resource "hund_component" "alpha" {
  name = "Alpha"
  group = resource.hund_group.test.id

  watchdog = { service = { manual = {} } }
}

resource "hund_component" "delta" {
  name = "Delta"
  group = resource.hund_group.test.id

  watchdog = { service = { manual = {} } }
}

resource "hund_group_component_ordering" "test" {
	group = resource.hund_group.test.id

  components = %[1]v
}
`, ordering)
}
