package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciQosInstancePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_qos_instance_policy.test"
	dataSourceName := "data.aci_qos_instance_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciQOSInstancePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateQosInstancePolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccQosInstancePolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etrap_age_timer", resourceName, "etrap_age_timer"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etrap_bw_thresh", resourceName, "etrap_bw_thresh"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etrap_byte_ct", resourceName, "etrap_byte_ct"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etrap_st", resourceName, "etrap_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_flush_interval", resourceName, "fabric_flush_interval"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_flush_st", resourceName, "fabric_flush_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.#", resourceName, "ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.0", resourceName, "ctrl.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "uburst_spine_queues", resourceName, "uburst_spine_queues"),
					resource.TestCheckResourceAttrPair(dataSourceName, "uburst_tor_queues", resourceName, "uburst_tor_queues"),
				),
			},
			{
				Config:      CreateAccQosInstancePolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccQosInstancePolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccQosInstancePolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccQosInstancePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing qos_instance_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
	}

	data "aci_qos_instance_policy" "test" {
	
		name  = aci_qos_instance_policy.test.name
		depends_on = [ aci_qos_instance_policy.test ]
	}
	`, rName)
	return resource
}

func CreateQosInstancePolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing qos_instance_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_qos_instance_policy" "test" {
	
	#	name  = aci_qos_instance_policy.test.name
		depends_on = [ aci_qos_instance_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccQosInstancePolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing qos_instance_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
	}

	data "aci_qos_instance_policy" "test" {
	
		name  = "${aci_qos_instance_policy.test.name}_invalid"
		name  = aci_qos_instance_policy.test.name
		depends_on = [ aci_qos_instance_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccQosInstancePolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing qos_instance_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
	}

	data "aci_qos_instance_policy" "test" {
	
		name  = aci_qos_instance_policy.test.name
		%s = "%s"
		depends_on = [ aci_qos_instance_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccQosInstancePolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing qos_instance_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_qos_instance_policy" "test" {
	
		name  = aci_qos_instance_policy.test.name
		depends_on = [ aci_qos_instance_policy.test ]
	}
	`, rName, key, value)
	return resource
}
