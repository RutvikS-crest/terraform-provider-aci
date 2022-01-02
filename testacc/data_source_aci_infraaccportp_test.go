package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLeafInterfaceProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_leaf_interface_profile.test"
	dataSourceName := "data.aci_leaf_interface_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafInterfaceProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLeafInterfaceProfileDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafInterfaceProfileConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccLeafInterfaceProfileDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config: CreateAccLeafInterfaceProfileDataSourceUpdate(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLeafInterfaceProfileConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing leaf_interface_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
	
		name  = "%s"
	}

	data "aci_leaf_interface_profile" "test" {
	
		name  = aci_leaf_interface_profile.test.name
		depends_on = [
			aci_leaf_interface_profile.test
		]
	}
	`, rName)
	return resource
}

func CreateLeafInterfaceProfileDSWithoutRequired(rName, attr string) string {
	fmt.Println("=== STEP  testing leaf_interface_profile creation without required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name  = "%s"
	}

	data "aci_leaf_interface_profile" "test" {
	
		depends_on = [
			aci_leaf_interface_profile.test
		]
	}
	`, rName)
	return resource
}

func CreateAccLeafInterfaceProfileDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing leaf_interface_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
	
		name  = "%s"
	}

	data "aci_leaf_interface_profile" "test" {
	
		name  = aci_leaf_interface_profile.test.name
		%s = "%s"
		depends_on = [
			aci_leaf_interface_profile.test
		]
	}
	`, rName, key, value)
	return resource
}
