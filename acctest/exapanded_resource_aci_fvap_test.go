package aci

import (
	"fmt"
	"log"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAciApplicationProfile_Basic(t *testing.T) {
	var application_profile models.ApplicationProfile
	resourceName := "aci_application_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationProfileConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "name", "expanded_test_for_terraform"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", GetParentDn(application_profile.DistinguishedName, fmt.Sprintf("/ap-%s", application_profile.Name))),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplicationProfile_description(t *testing.T) {
	var application_profile_default_description models.ApplicationProfile
	var application_profile_updated_description models.ApplicationProfile
	resourceName := "aci_application_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationProfileConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default_description),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApplicationProfileUpdatedDescription,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated_description),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description for terraform test"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default_description, &application_profile_updated_description),
				),
			},
		},
	})
}

func testAccCheckAciApplicationProfileIdEqual(ap1, ap2 *models.ApplicationProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ap1.DistinguishedName != ap2.DistinguishedName {
			return fmt.Errorf("ApplicationProfile DNs are not equal")
		}
		return nil
	}
}

testAccApplicationProfileConfig :=fmt.Sprintf( `
resource "aci_tenant" "test" {
	name = "expanded_test_for_terraform"
}

resource "aci_application_profile" "test" {
	tenant_dn = aci_tenant.test.id
	name = "%s"
}
`,acctest.RandString(5))

const testAccApplicationProfileUpdatedDescription = `
resource "aci_tenant" "test" {
	name = "expanded_test_for_terraform"
}

resource "aci_application_profile" "test" {
	tenant_dn = aci_tenant.test.id
	name = "expanded_test_for_terraform"
	description = "updated description for terraform test"
}
`

func testAccCheckAciApplicationProfileExists(name string, application_profile *models.ApplicationProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Application Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Application Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		application_profileFound := models.ApplicationProfileFromContainer(cont)
		if application_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Application Profile %s not found", rs.Primary.ID)
		}
		*application_profile = *application_profileFound
		log.Printf("application_profile.Name %s", application_profile.Name)
		log.Printf("application_profile.DistinguishedName %s", application_profile.DistinguishedName)
		return nil
	}
}

func testAccCheckAciApplicationProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_application_profile" {
			cont, err := client.Get(rs.Primary.ID)
			application_profile := models.ApplicationProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Application Profile %s Still exists", application_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}
