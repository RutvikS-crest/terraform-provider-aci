package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCoopPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_coop_policy.test"
	dataSourceName := "data.aci_coop_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCoopPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateCoopPolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCoopPolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "coop_group_policy_type", resourceName, "coop_group_policy_type"),
				),
			},
			{
				Config:      CreateAccCoopPolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccCoopPolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccCoopPolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCoopPolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing coop_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
	}

	data "aci_coop_policy" "test" {
	
		name  = aci_coop_policy.test.name
		depends_on = [ aci_coop_policy.test ]
	}
	`, rName)
	return resource
}

func CreateCoopPolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing coop_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_coop_policy" "test" {
	
	#	name  = aci_coop_policy.test.name
		depends_on = [ aci_coop_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccCoopPolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing coop_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
	}

	data "aci_coop_policy" "test" {
	
		name  = "${aci_coop_policy.test.name}_invalid"
		depends_on = [ aci_coop_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccCoopPolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing coop_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
	}

	data "aci_coop_policy" "test" {
	
		name  = aci_coop_policy.test.name
		%s = "%s"
		depends_on = [ aci_coop_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccCoopPolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing coop_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_coop_policy" "test" {
	
		name  = aci_coop_policy.test.name
		depends_on = [ aci_coop_policy.test ]
	}
	`, rName, key, value)
	return resource
}
