package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciApplicationProfileDataSource_Basic(t *testing.T) {
	var application_profile models.ApplicationProfile
	resourceName := "aci_application_profile.test"
	dataSourceName := "data.aci_application_profile.test"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccApplicationProfileConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
				),
			},
		},
	})
}

func CreateAccApplicationProfileConfigDataSource(rName string) string {
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = aci_application_profile.test.name
	}
	`, rName, rName)
	return resource
}
