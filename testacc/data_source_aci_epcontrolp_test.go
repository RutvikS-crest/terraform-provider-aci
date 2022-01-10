package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciEndpointControlsDataSource_Basic(t *testing.T) {
	resourceName := "aci_endpoint_controls.test"
	dataSourceName := "data.aci_endpoint_controls.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEndpointControlPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateEndpointControlsDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndpointControlsConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hold_intvl", resourceName, "hold_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rogue_ep_detect_intvl", resourceName, "rogue_ep_detect_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rogue_ep_detect_mult", resourceName, "rogue_ep_detect_mult"),
				),
			},
			{
				Config:      CreateAccEndpointControlsDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccEndpointControlsDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccEndpointControlsDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccEndpointControlsConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing endpoint_controls Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s"
	}

	data "aci_endpoint_controls" "test" {
	
		name  = aci_endpoint_controls.test.name
		depends_on = [ aci_endpoint_controls.test ]
	}
	`, rName)
	return resource
}

func CreateEndpointControlsDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_controls Data Source without ", attrName)
	rBlock := `
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_endpoint_controls" "test" {
	
	#	name  = aci_endpoint_controls.test.name
		depends_on = [ aci_endpoint_controls.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccEndpointControlsDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing endpoint_controls Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s"
	}

	data "aci_endpoint_controls" "test" {
	
		name  = "${aci_endpoint_controls.test.name}_invalid"
		name  = aci_endpoint_controls.test.name
		depends_on = [ aci_endpoint_controls.test ]
	}
	`, rName)
	return resource
}

func CreateAccEndpointControlsDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing endpoint_controls Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s"
	}

	data "aci_endpoint_controls" "test" {
	
		name  = aci_endpoint_controls.test.name
		%s = "%s"
		depends_on = [ aci_endpoint_controls.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccEndpointControlsDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing endpoint_controls Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_endpoint_controls" "test" {
	
		name  = aci_endpoint_controls.test.name
		depends_on = [ aci_endpoint_controls.test ]
	}
	`, rName, key, value)
	return resource
}
