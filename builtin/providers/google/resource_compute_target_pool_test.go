package google

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeTargetPool_basic(t *testing.T) {
	var targetPool compute.TargetPool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetPoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetPool_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetPoolExists(
						"google_compute_target_pool.foobar", &targetPool),
				),
			},
		},
	})
}

func TestAccComputeTargetPool_health_checks(t *testing.T) {
	var targetPool compute.TargetPool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetPoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetPool_health_checks,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetPoolExists(
						"google_compute_target_pool.foobar", &targetPool),
					testAccCheckComputeTargetPoolHealthChecks(
						[]string{"google_compute_http_health_check.foobar"}, &targetPool),
				),
			},
		},
	})
}

func testAccCheckComputeTargetPoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_target_pool" {
			continue
		}

		_, err := config.clientCompute.TargetPools.Get(
			config.Project, config.Region, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("TargetPool still exists")
		}
	}

	return nil
}

func testAccCheckComputeTargetPoolExists(n string, targetPool *compute.TargetPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.TargetPools.Get(
			config.Project, config.Region, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("TargetPool not found")
		}

		*targetPool = *found

		return nil
	}
}

func testAccCheckComputeTargetPoolHealthChecks(names []string, targetPool *compute.TargetPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		urls := make([]string, len(names))

		for i, name := range names {
			rs, ok := s.RootModule().Resources[name]
			if !ok {
				return fmt.Errorf("Not found: %s", name)
			}

			url, ok := rs.Primary.Attributes["self_link"]
			if !ok {
				return fmt.Errorf("Not found: %s", name)
			}

			urls[i] = url
		}

		if !reflect.DeepEqual(targetPool.HealthChecks, urls) {
			return fmt.Errorf("Health checks don't match: expected %q, got %q",
				urls, targetPool.HealthChecks)
		}

		return nil
	}
}

const testAccComputeTargetPool_basic = `
resource "google_compute_target_pool" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	instances = ["us-central1-a/foo", "us-central1-b/bar"]
	name = "terraform-test"
	session_affinity = "CLIENT_IP_PROTO"
}`

const testAccComputeTargetPool_health_checks = `
resource "google_compute_http_health_check" "foobar" {
	name = "foobar"
	description = "Resource created for Terraform acceptance testing"
}
resource "google_compute_target_pool" "foobar" {
	name = "terraform-test"
	description = "Resource created for Terraform acceptance testing"
	health_checks = [ "${google_compute_http_health_check.foobar.self_link}" ]
}`
