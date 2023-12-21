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
	resource.AddTestSweepers("hund_group", &resource.Sweeper{
		Name:         "hund_group",
		Dependencies: []string{"hund_component"},
		F: func(r string) error {
			client, err := sharedClientForDomain(r)
			if err != nil {
				return err
			}

			ctx := context.Background()

			limit := 100
			rsp, err := client.GetAllGroups(ctx, &hundApiV1.GetAllGroupsParams{
				Limit: &limit,
			})
			if err != nil {
				return err
			}

			groups, err := hundApiV1.ParseGetAllGroupsResponse(rsp)
			if err != nil {
				return err
			}
			if groups.StatusCode() != 200 {
				return errors.New("couldn't retrieve groups:\n" + string(groups.Body))
			}

			for _, g := range groups.HALJSON200.Data {
				rsp, err := client.DeleteAGroup(ctx, g.Id)
				if err != nil {
					return err
				}
				if rsp.StatusCode != 204 {
					group, err := hundApiV1.ParseDeleteAGroupResponse(rsp)
					if err != nil {
						return err
					}

					return errors.New("couldn't delete group:\n" + string(group.Body))
				}
			}

			return nil
		},
	})
}

func TestAccGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGroupResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_group.test", "name_translations.en", "one"),
					resource.TestCheckResourceAttr("hund_group.test", "components.#", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hund_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccGroupResourceConfigI18n("two", "zwei"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hund_group.test", "name", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccGroupResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "hund_group" "test" {
  name = %[1]q
}
`, name)
}

func testAccGroupResourceConfigI18n(name_en string, name_de string) string {
	return providerConfig + fmt.Sprintf(`
resource "hund_group" "test" {
  name_translations = {
		en = %[1]q,
		de = %[2]q,
		original = "en",
	}
}
`, name_en, name_de)
}
