package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFabricIFPolDataSource_Basic(t *testing.T) {
	resourceName := "aci_fabric_if_pol.test"
	dataSourceName := "data.aci_fabric_if_pol.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIFPolDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFabricIFPolDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricIFPolConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auto_neg", resourceName, "auto_neg"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fec_mode", resourceName, "fec_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "link_debounce", resourceName, "link_debounce"),
					resource.TestCheckResourceAttrPair(dataSourceName, "speed", resourceName, "speed"),
				),
			},
			{
				Config:      CreateAccFabricIFPolDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccFabricIFPolDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},

			{
				Config: CreateAccFabricIFPolDataSourceUpdate(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFabricIFPolConfigDataSource(rName string) string {
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

func CreateAccFabricIFPolConfigDataSourceUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with updated name")
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

func CreateAccFabricIFPolDataSourceUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with Invalid Name")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
	
		name  = "%s"
	}

	data "aci_fabric_if_pol" "test" {
	
		name  = "${aci_fabric_if_pol.test.name}_invaid"
		depends_on = [
			aci_fabric_if_pol.test
		]
	}
	`, rName)
	return resource
}

func CreateFabricIFPolDSWithoutRequired(rName, attribute string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with required arguments only")
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

func CreateAccFabricIFPolDataSourceUpdate(rName, key, value string) string {
	fmt.Printf("=== STEP  testing fabric_if_pol creation with %s = %s", key, value)
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
