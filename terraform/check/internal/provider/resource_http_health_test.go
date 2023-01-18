package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	testUrl := "http://example.com"
	testUrl2 := "https://httpbin.org/get"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig(testUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("check_http_health.test", "url", testUrl),
					resource.TestCheckResourceAttr("check_http_health.test", "id", "example-id"),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      "check_http_health.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// This is not normally necessary, but is here because this
			// 	// example code does not have an actual upstream service.
			// 	// Once the Read method is able to refresh information from
			// 	// the upstream service, this can be removed.
			// 	ImportStateVerifyIgnore: []string{"url"},
			// },
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig(testUrl2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("check_http_health.test", "url", testUrl2),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig(url string) string {
	return fmt.Sprintf(`
resource "check_http_health" "test" {
  url = %[1]q
  consecutive_successes = 5
  headers = {
	hello = "world"
  }
}
`, url)
}
