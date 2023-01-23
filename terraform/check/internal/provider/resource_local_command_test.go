package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLocalCommandResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalCommandResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("check_local_command.test", "id", "example-id"),
				),
			},
		},
	})
}

func testAccLocalCommandResourceConfig() string {
	return fmt.Sprintf(`
resource "check_local_command" "test" {
	command = "true"
}`)

}
