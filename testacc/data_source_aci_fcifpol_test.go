package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciInterfaceFCPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_interface_fc_policy.test"
	dataSourceName := "data.aci_interface_fc_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciInterfaceFCPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateInterfaceFCPolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccInterfaceFCPolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "automaxspeed", resourceName, "automaxspeed"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fill_pattern", resourceName, "fill_pattern"),
					resource.TestCheckResourceAttrPair(dataSourceName, "port_mode", resourceName, "port_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rx_bb_credit", resourceName, "rx_bb_credit"),
					resource.TestCheckResourceAttrPair(dataSourceName, "speed", resourceName, "speed"),
					resource.TestCheckResourceAttrPair(dataSourceName, "trunk_mode", resourceName, "trunk_mode"),
				),
			},
			{
				Config:      CreateAccInterfaceFCPolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config: CreateAccInterfaceFCPolicyDataSourceUpdate(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccInterfaceFCPolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing interface_fc_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
	
		name  = "%s"
	}

	data "aci_interface_fc_policy" "test" {
	
		name  = aci_interface_fc_policy.test.name
		depends_on = [
			aci_interface_fc_policy.test
		]
	}
	`, rName)
	return resource
}

func CreateInterfaceFCPolicyDSWithoutRequired(rName, attr string) string {
	fmt.Println("=== STEP  testing interface_fc_policy creation without required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
	
		name  = "%s"
	}

	data "aci_interface_fc_policy" "test" {
	
		depends_on = [
			aci_interface_fc_policy.test
		]
	}
	`, rName)
	return resource
}

func CreateAccInterfaceFCPolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing interface_fc_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
	
		name  = "%s"
	}

	data "aci_interface_fc_policy" "test" {
	
		name  = aci_interface_fc_policy.test.name
		%s = "%s"
		depends_on = [
			aci_interface_fc_policy.test
		]
	}
	`, rName, key, value)
	return resource
}
