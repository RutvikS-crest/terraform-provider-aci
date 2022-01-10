package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSystemDataSource_Basic(t *testing.T) {
	resourceName := "aci_system.test"
	dataSourceName := "data.aci_system.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	fabricNodeName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateSystemDSWithoutRequired(fabricNodeName, "fabric_node_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			}, {
				Config: CreateAccSystemConfigDataSource(fabricNodeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_node_dn", resourceName, "fabric_node_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "address", resourceName, "address"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etep_addr", resourceName, "etep_addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "system_id", resourceName, "system_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "node_type", resourceName, "node_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "remote_network_id", resourceName, "remote_network_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "remote_node", resourceName, "remote_node"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rldirect_mode", resourceName, "rldirect_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "role", resourceName, "role"),
					resource.TestCheckResourceAttrPair(dataSourceName, "server_type", resourceName, "server_type"),
				),
			},
			{
				Config:      CreateAccSystemDataSourceUpdate(fabricNodeName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSystemDSWithInvalidParentDn(fabricNodeName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccSystemDataSourceUpdatedResource(fabricNodeName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSystemConfigDataSource(fabricNodeName string) string {
	fmt.Println("=== STEP  testing system Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_system" "test" {
		fabric_node_dn  = aci_fabric_node.test.id
	}

	data "aci_system" "test" {
		fabric_node_dn  = aci_fabric_node.test.id
		depends_on = [ aci_system.test ]
	}
	`, fabricNodeName)
	return resource
}

func CreateSystemDSWithoutRequired(fabricNodeName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing system Data Source without ", attrName)
	rBlock := `
	
	resource "aci_fabric_node" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_system" "test" {
		fabric_node_dn  = aci_fabric_node.test.id
	}
	`
	switch attrName {
	case "fabric_node_dn":
		rBlock += `
	data "aci_system" "test" {
	#	fabric_node_dn  = aci_fabric_node.test.id
	
		depends_on = [ aci_system.test ]
	}
		`

	}
	return fmt.Sprintf(rBlock, fabricNodeName)
}

func CreateAccSystemDSWithInvalidParentDn(fabricNodeName string) string {
	fmt.Println("=== STEP  testing system Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_system" "test" {
		fabric_node_dn  = aci_fabric_node.test.id
	}

	data "aci_system" "test" {
		fabric_node_dn  = aci_fabric_node.test.id
		depends_on = [ aci_system.test ]
	}
	`, fabricNodeName)
	return resource
}

func CreateAccSystemDataSourceUpdate(fabricNodeName, key, value string) string {
	fmt.Println("=== STEP  testing system Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_system" "test" {
		fabric_node_dn  = aci_fabric_node.test.id
	}

	data "aci_system" "test" {
		fabric_node_dn  = aci_fabric_node.test.id
		%s = "%s"
		depends_on = [ aci_system.test ]
	}
	`, fabricNodeName, key, value)
	return resource
}

func CreateAccSystemDataSourceUpdatedResource(fabricNodeName, key, value string) string {
	fmt.Println("=== STEP  testing system Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_system" "test" {
		fabric_node_dn  = aci_fabric_node.test.id
		%s = "%s"
	}

	data "aci_system" "test" {
		fabric_node_dn  = aci_fabric_node.test.id
		depends_on = [ aci_system.test ]
	}
	`, fabricNodeName, key, value)
	return resource
}
