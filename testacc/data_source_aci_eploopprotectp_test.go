package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciEndpointLoopProtectionDataSource_Basic(t *testing.T) {
	resourceName := "aci_endpoint_loop_protection.test"
	dataSourceName := "data.aci_endpoint_loop_protection.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPLoopProtectionPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateEndpointLoopProtectionDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndpointLoopProtectionConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "action.#", resourceName, "action.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "action.0", resourceName, "action.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "loop_detect_intvl", resourceName, "loop_detect_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "loop_detect_mult", resourceName, "loop_detect_mult"),
				),
			},
			{
				Config:      CreateAccEndpointLoopProtectionDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccEndpointLoopProtectionDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccEndpointLoopProtectionDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccEndpointLoopProtectionConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing endpoint_loop_protection Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
	}

	data "aci_endpoint_loop_protection" "test" {
	
		name  = aci_endpoint_loop_protection.test.name
		depends_on = [ aci_endpoint_loop_protection.test ]
	}
	`, rName)
	return resource
}

func CreateEndpointLoopProtectionDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_loop_protection Data Source without ", attrName)
	rBlock := `
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_endpoint_loop_protection" "test" {
	
	#	name  = aci_endpoint_loop_protection.test.name
		depends_on = [ aci_endpoint_loop_protection.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccEndpointLoopProtectionDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing endpoint_loop_protection Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
	}

	data "aci_endpoint_loop_protection" "test" {
	
		name  = "${aci_endpoint_loop_protection.test.name}_invalid"
		name  = aci_endpoint_loop_protection.test.name
		depends_on = [ aci_endpoint_loop_protection.test ]
	}
	`, rName)
	return resource
}

func CreateAccEndpointLoopProtectionDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing endpoint_loop_protection Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
	}

	data "aci_endpoint_loop_protection" "test" {
	
		name  = aci_endpoint_loop_protection.test.name
		%s = "%s"
		depends_on = [ aci_endpoint_loop_protection.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccEndpointLoopProtectionDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing endpoint_loop_protection Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_endpoint_loop_protection" "test" {
	
		name  = aci_endpoint_loop_protection.test.name
		depends_on = [ aci_endpoint_loop_protection.test ]
	}
	`, rName, key, value)
	return resource
}
