package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciMCPInstancePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_mcp_instance_policy.test"
	dataSourceName := "data.aci_mcp_instance_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMiscablingProtocolInstancePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateMCPInstancePolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMCPInstancePolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.#", resourceName, "ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "init_delay_time", resourceName, "init_delay_time"),
					resource.TestCheckResourceAttrPair(dataSourceName, "key", resourceName, "key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "loop_detect_mult", resourceName, "loop_detect_mult"),
					resource.TestCheckResourceAttrPair(dataSourceName, "loop_protect_act.#", resourceName, "loop_protect_act.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "loop_protect_act.0", resourceName, "loop_protect_act.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tx_freq", resourceName, "tx_freq"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tx_freq_msec", resourceName, "tx_freq_msec"),
				),
			},
			{
				Config:      CreateAccMCPInstancePolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccMCPInstancePolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccMCPInstancePolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccMCPInstancePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
	}

	data "aci_mcp_instance_policy" "test" {
	
		name  = aci_mcp_instance_policy.test.name
		depends_on = [ aci_mcp_instance_policy.test ]
	}
	`, rName)
	return resource
}

func CreateMCPInstancePolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing mcp_instance_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_mcp_instance_policy" "test" {
	
	#	name  = aci_mcp_instance_policy.test.name
		depends_on = [ aci_mcp_instance_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccMCPInstancePolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
	}

	data "aci_mcp_instance_policy" "test" {
	
		name  = "${aci_mcp_instance_policy.test.name}_invalid"
		name  = aci_mcp_instance_policy.test.name
		depends_on = [ aci_mcp_instance_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccMCPInstancePolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
	}

	data "aci_mcp_instance_policy" "test" {
	
		name  = aci_mcp_instance_policy.test.name
		%s = "%s"
		depends_on = [ aci_mcp_instance_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccMCPInstancePolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_mcp_instance_policy" "test" {
	
		name  = aci_mcp_instance_policy.test.name
		depends_on = [ aci_mcp_instance_policy.test ]
	}
	`, rName, key, value)
	return resource
}
