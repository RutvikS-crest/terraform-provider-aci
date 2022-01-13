package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciErrorDisableRecoveryDataSource_Basic(t *testing.T) {
	resourceName := "aci_error_disable_recovery.test"
	dataSourceName := "data.aci_error_disable_recovery.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciErrorDisabledRecoveryPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateErrorDisableRecoveryDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccErrorDisableRecoveryConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "err_dis_recov_intvl", resourceName, "err_dis_recov_intvl"),
				),
			},
			{
				Config:      CreateAccErrorDisableRecoveryDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccErrorDisableRecoveryDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccErrorDisableRecoveryDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccErrorDisableRecoveryConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing error_disable_recovery Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
	}

	data "aci_error_disable_recovery" "test" {
	
		name  = aci_error_disable_recovery.test.name
		depends_on = [ aci_error_disable_recovery.test ]
	}
	`, rName)
	return resource
}

func CreateErrorDisableRecoveryDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing error_disable_recovery Data Source without ", attrName)
	rBlock := `
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_error_disable_recovery" "test" {
	
	#	name  = aci_error_disable_recovery.test.name
		depends_on = [ aci_error_disable_recovery.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccErrorDisableRecoveryDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing error_disable_recovery Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
	}

	data "aci_error_disable_recovery" "test" {
	
		name  = "${aci_error_disable_recovery.test.name}_invalid"
		depends_on = [ aci_error_disable_recovery.test ]
	}
	`, rName)
	return resource
}

func CreateAccErrorDisableRecoveryDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing error_disable_recovery Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
	}

	data "aci_error_disable_recovery" "test" {
	
		name  = aci_error_disable_recovery.test.name
		%s = "%s"
		depends_on = [ aci_error_disable_recovery.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccErrorDisableRecoveryDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing error_disable_recovery Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_error_disable_recovery" "test" {
	
		name  = aci_error_disable_recovery.test.name
		depends_on = [ aci_error_disable_recovery.test ]
	}
	`, rName, key, value)
	return resource
}
