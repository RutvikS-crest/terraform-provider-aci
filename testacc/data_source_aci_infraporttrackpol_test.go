package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciPortTrackingDataSource_Basic(t *testing.T) {
	resourceName := "aci_port_tracking.test"
	dataSourceName := "data.aci_port_tracking.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPortTrackingDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreatePortTrackingDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPortTrackingConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "delay", resourceName, "delay"),
					resource.TestCheckResourceAttrPair(dataSourceName, "include_apic_ports", resourceName, "include_apic_ports"),
					resource.TestCheckResourceAttrPair(dataSourceName, "minlinks", resourceName, "minlinks"),
				),
			},
			{
				Config:      CreateAccPortTrackingDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccPortTrackingDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccPortTrackingDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccPortTrackingConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing port_tracking Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
	}

	data "aci_port_tracking" "test" {
	
		name  = aci_port_tracking.test.name
		depends_on = [ aci_port_tracking.test ]
	}
	`, rName)
	return resource
}

func CreatePortTrackingDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing port_tracking Data Source without ", attrName)
	rBlock := `
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_port_tracking" "test" {
	
	#	name  = aci_port_tracking.test.name
		depends_on = [ aci_port_tracking.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccPortTrackingDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing port_tracking Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
	}

	data "aci_port_tracking" "test" {
	
		name  = "${aci_port_tracking.test.name}_invalid"
		depends_on = [ aci_port_tracking.test ]
	}
	`, rName)
	return resource
}

func CreateAccPortTrackingDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing port_tracking Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
	}

	data "aci_port_tracking" "test" {
	
		name  = aci_port_tracking.test.name
		%s = "%s"
		depends_on = [ aci_port_tracking.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccPortTrackingDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing port_tracking Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_port_tracking" "test" {
	
		name  = aci_port_tracking.test.name
		depends_on = [ aci_port_tracking.test ]
	}
	`, rName, key, value)
	return resource
}
