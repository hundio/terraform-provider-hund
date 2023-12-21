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
	resource.AddTestSweepers("hund_component", &resource.Sweeper{
		Name: "hund_component",
		F: func(r string) error {
			client, err := sharedClientForDomain(r)
			if err != nil {
				return err
			}

			ctx := context.Background()

			limit := 100
			rsp, err := client.GetAllComponents(ctx, &hundApiV1.GetAllComponentsParams{
				Limit: &limit,
			})
			if err != nil {
				return err
			}

			components, err := hundApiV1.ParseGetAllComponentsResponse(rsp)
			if err != nil {
				return err
			}
			if components.StatusCode() != 200 {
				return errors.New("couldn't retrieve components:\n" + string(components.Body))
			}

			for _, c := range components.HALJSON200.Data {
				rsp, err := client.DeleteAComponent(ctx, c.Id)
				if err != nil {
					return err
				}
				if rsp.StatusCode != 204 {
					component, err := hundApiV1.ParseDeleteAComponentResponse(rsp)
					if err != nil {
						return err
					}

					return errors.New("couldn't delete component:\n" + string(component.Body))
				}
			}

			return nil
		},
	})
}

func TestAccComponentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccComponentResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_component.test", "name_translations.en", "one"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "hund_component.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_event_at"},
			},
			// Update and Read testing
			{
				Config: testAccComponentResourceConfigI18n("two", "zwei"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_component.test", "name", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccComponentResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "hund_group" "test" {
	name = "Test Group"
}

resource "hund_component" "test" {
  name = %[1]q

	group = hund_group.test.id

	watchdog = {service = { manual = {state  = 1} }}
}
`, name)
}

func testAccComponentResourceConfigI18n(name_en string, name_de string) string {
	return providerConfig + fmt.Sprintf(`
resource "hund_group" "test" {
	name = "Test Group"
}

resource "hund_component" "test" {
  name_translations = {
		en = %[1]q,
		de = %[2]q,
		original = "en"
	}

	group = hund_group.test.id

	watchdog = {service = { manual = {state  = 0} }}
}
`, name_en, name_de)
}
