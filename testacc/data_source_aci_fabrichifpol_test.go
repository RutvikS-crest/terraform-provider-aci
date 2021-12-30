package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFabricIfPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_fabric_if_pol.test"
	dataSourceName := "data.aci_fabric_if_pol.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFabricIfPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFabricIfPolicyDSWithoutRequired(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricIfPolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auto_neg", resourceName, "auto_neg"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dfe_delay_ms", resourceName, "dfe_delay_ms"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fec_mode", resourceName, "fec_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "link_debounce", resourceName, "link_debounce"),
					resource.TestCheckResourceAttrPair(dataSourceName, "speed", resourceName, "speed"),
				),
			},
			{
				Config:      CreateAccFabricIfPolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config: CreateAccFabricIfPolicyDataSourceUpdate(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFabricIfPolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {	
		name  = "%s"
	}

	data "aci_fabric_if_pol" "test" {
		name  = aci_fabric_if_pol.test.name
		depends_on = [
			aci_fabric_if_pol.test
		]
	}
	`, rName)
	return resource
}

func CreateFabricIfPolicyDSWithoutRequired(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation without required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
		name  = "%s"
	}

	data "aci_fabric_if_pol" "test" {
		depends_on = [
			aci_fabric_if_pol.test
		]
	}
	`, rName)
	return resource
}

func CreateAccFabricIfPolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing fabric_if_pol data source updation")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
		name  = "%s"
	}

	data "aci_fabric_if_pol" "test" {
		name  = aci_fabric_if_pol.test.name
		%s = "%s"
		depends_on = [
			aci_fabric_if_pol.test
		]
	}
	`, rName, key, value)
	return resource
}
