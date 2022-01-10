package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFabricNodeDataSource_Basic(t *testing.T) {
	resourceName := "aci_fabric_node.test"
	dataSourceName := "data.aci_fabric_node.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	fabricNodeId := "101"

	fabricPodName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateFabricNodeDSWithoutRequired(fabricPodName, fabricNodeId, "fabric_pod_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricNodeConfigDataSource(fabricPodName, fabricNodeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_pod_dn", resourceName, "fabric_pod_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_node_id", resourceName, "fabric_node_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ad_st", resourceName, "ad_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "apic_type", resourceName, "apic_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_st", resourceName, "fabric_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "last_state_mod_ts", resourceName, "last_state_mod_ts"),
				),
			},
			{
				Config:      CreateAccFabricNodeDataSourceUpdate(fabricPodName, fabricNodeId, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFabricNodeDSWithInvalidParentDn(fabricPodName, fabricNodeId),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccFabricNodeDataSourceUpdatedResource(fabricPodName, fabricNodeId, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFabricNodeConfigDataSource(fabricPodName, fabricNodeId string) string {
	fmt.Println("=== STEP  testing fabric_node Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_pod" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fabric_node" "test" {
		fabric_pod_dn  = aci_fabric_pod.test.id
		fabric_node_id  = "%s"
	}

	data "aci_fabric_node" "test" {
		fabric_pod_dn  = aci_fabric_pod.test.id
		fabric_node_id  = aci_fabric_node.test.fabric_node_id
		depends_on = [ aci_fabric_node.test ]
	}
	`, fabricPodName, fabricNodeId)
	return resource
}

func CreateFabricNodeDSWithoutRequired(fabricPodName, fabricNodeId, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_node Data Source without ", attrName)
	rBlock := `
	
	resource "aci_fabric_pod" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fabric_node" "test" {
		fabric_pod_dn  = aci_fabric_pod.test.id
		fabric_node_id  = "%s"
	}
	`
	switch attrName {
	case "fabric_pod_dn":
		rBlock += `
	data "aci_fabric_node" "test" {
	#	fabric_pod_dn  = aci_fabric_pod.test.id
		fabric_node_id  = aci_fabric_node.test.fabric_node_id
		depends_on = [ aci_fabric_node.test ]
	}
		`
	case "fabric_node_id":
		rBlock += `
	data "aci_fabric_node" "test" {
		fabric_pod_dn  = aci_fabric_pod.test.id
	#	fabric_node_id  = aci_fabric_node.test.fabric_node_id
		depends_on = [ aci_fabric_node.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fabricPodName, fabricNodeId)
}

func CreateAccFabricNodeDSWithInvalidParentDn(fabricPodName, fabricNodeId string) string {
	fmt.Println("=== STEP  testing fabric_node Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_pod" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fabric_node" "test" {
		fabric_pod_dn  = aci_fabric_pod.test.id
		fabric_node_id  = "%s"
	}

	data "aci_fabric_node" "test" {
		fabric_pod_dn  = aci_fabric_pod.test.id
		fabric_node_id  = "${aci_fabric_node.test.fabric_node_id}_invalid"
		depends_on = [ aci_fabric_node.test ]
	}
	`, fabricPodName, fabricNodeId)
	return resource
}

func CreateAccFabricNodeDataSourceUpdate(fabricPodName, fabricNodeId, key, value string) string {
	fmt.Println("=== STEP  testing fabric_node Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_pod" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fabric_node" "test" {
		fabric_pod_dn  = aci_fabric_pod.test.id
		fabric_node_id  = "%s"
	}

	data "aci_fabric_node" "test" {
		fabric_pod_dn  = aci_fabric_pod.test.id
		fabric_node_id  = aci_fabric_node.test.fabric_node_id
		%s = "%s"
		depends_on = [ aci_fabric_node.test ]
	}
	`, fabricPodName, fabricNodeId, key, value)
	return resource
}

func CreateAccFabricNodeDataSourceUpdatedResource(fabricPodName, fabricNodeId, key, value string) string {
	fmt.Println("=== STEP  testing fabric_node Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_pod" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fabric_node" "test" {
		fabric_pod_dn  = aci_fabric_pod.test.id
		fabric_node_id  = "%s"
		%s = "%s"
	}

	data "aci_fabric_node" "test" {
		fabric_pod_dn  = aci_fabric_pod.test.id
		fabric_node_id  = aci_fabric_node.test.fabric_node_id
		depends_on = [ aci_fabric_node.test ]
	}
	`, fabricPodName, fabricNodeId, key, value)
	return resource
}
