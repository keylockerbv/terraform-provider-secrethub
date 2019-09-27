package secrethub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceService_create(t *testing.T) {
	repoPath := testAcc.namespace + "/" + testAcc.repository

	config := fmt.Sprintf(`
		resource "secrethub_service" "test" {
			repo = "%s"
		}
	`, repoPath)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  testAccPreCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkServiceExistsRemotely(repoPath, ""),
				),
			},
		},
	})
}

func checkServiceExistsRemotely(path string, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := *testAccProvider.Meta().(providerMeta).client

		services, err := client.Services().List(path)
		if err != nil {
			return fmt.Errorf("cannot list services: %s", err)
		}

		for _, service := range services {
			if service.Description == description {
				return nil
			}
		}

		return fmt.Errorf("expected service on repo %s with description \"%s\"", path, description)
	}
}