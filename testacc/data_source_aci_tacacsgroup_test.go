package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTACACSAccountingDataSource_Basic(t *testing.T) {
	resourceName := "aci_tacacs_accounting.test"
	dataSourceName := "data.aci_tacacs_accounting.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSAccountingDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateTACACSAccountingDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSAccountingConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccTACACSAccountingDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTACACSAccountingDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccTACACSAccountingDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccTACACSAccountingConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing tacacs_accounting Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_accounting" "test" {
	
		name  = aci_tacacs_accounting.test.name
		depends_on = [ aci_tacacs_accounting.test ]
	}
	`, rName)
	return resource
}

func CreateTACACSAccountingDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_tacacs_accounting" "test" {
	
	#	name  = aci_tacacs_accounting.test.name
		depends_on = [ aci_tacacs_accounting.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccTACACSAccountingDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing tacacs_accounting Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_accounting" "test" {
	
		name  = "${aci_tacacs_accounting.test.name}_invalid"
		name  = aci_tacacs_accounting.test.name
		depends_on = [ aci_tacacs_accounting.test ]
	}
	`, rName)
	return resource
}

func CreateAccTACACSAccountingDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_accounting Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_accounting" "test" {
	
		name  = aci_tacacs_accounting.test.name
		%s = "%s"
		depends_on = [ aci_tacacs_accounting.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccTACACSAccountingDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_accounting Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_tacacs_accounting" "test" {
	
		name  = aci_tacacs_accounting.test.name
		depends_on = [ aci_tacacs_accounting.test ]
	}
	`, rName, key, value)
	return resource
}
