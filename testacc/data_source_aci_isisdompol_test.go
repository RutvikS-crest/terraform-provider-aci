package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciISISDomainPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_isis_domain_policy.test"
	dataSourceName := "data.aci_isis_domain_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciISISDomainPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateISISDomainPolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccISISDomainPolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mtu", resourceName, "mtu"),
					resource.TestCheckResourceAttrPair(dataSourceName, "redistrib_metric", resourceName, "redistrib_metric"),
				),
			},
			{
				Config:      CreateAccISISDomainPolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccISISDomainPolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccISISDomainPolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccISISDomainPolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing isis_domain_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
	}

	data "aci_isis_domain_policy" "test" {
	
		name  = aci_isis_domain_policy.test.name
		depends_on = [ aci_isis_domain_policy.test ]
	}
	`, rName)
	return resource
}

func CreateISISDomainPolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing isis_domain_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_isis_domain_policy" "test" {
	
	#	name  = aci_isis_domain_policy.test.name
		depends_on = [ aci_isis_domain_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccISISDomainPolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing isis_domain_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
	}

	data "aci_isis_domain_policy" "test" {
	
		name  = "${aci_isis_domain_policy.test.name}_invalid"
		name  = aci_isis_domain_policy.test.name
		depends_on = [ aci_isis_domain_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccISISDomainPolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing isis_domain_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
	}

	data "aci_isis_domain_policy" "test" {
	
		name  = aci_isis_domain_policy.test.name
		%s = "%s"
		depends_on = [ aci_isis_domain_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccISISDomainPolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing isis_domain_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_isis_domain_policy" "test" {
	
		name  = aci_isis_domain_policy.test.name
		depends_on = [ aci_isis_domain_policy.test ]
	}
	`, rName, key, value)
	return resource
}
