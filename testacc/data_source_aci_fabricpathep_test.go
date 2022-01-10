package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFabricPathEpDataSource_Basic(t *testing.T) {
	resourceName := "aci_fabric_path_ep.test"
	dataSourceName := "data.aci_fabric_path_ep.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	name := makeTestVariable(acctest.RandString(5))
	fabricPodName := makeTestVariable(acctest.RandString(5))
	fabricPathEpContName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateFabricPathEpDSWithoutRequired(fabricPodName, fabricPathEpContName, name, "fabric_path_end-point_container_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateFabricPathEpDSWithoutRequired(fabricPodName, fabricPathEpContName, name, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricPathEpConfigDataSource(fabricPodName, fabricPathEpContName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_path_end-point_container_dn", resourceName, "fabric_path_end-point_container_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccFabricPathEpDataSourceUpdate(fabricPodName, fabricPathEpContName, name, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFabricPathEpDSWithInvalidParentDn(fabricPodName, fabricPathEpContName, name),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccFabricPathEpDataSourceUpdatedResource(fabricPodName, fabricPathEpContName, name, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFabricPathEpConfigDataSource(fabricPodName, fabricPathEpContName, name string) string {
	fmt.Println("=== STEP  testing fabric_path_ep Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_pod" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fabric_path_end-point_container" "test" {
		name 		= "%s"
		fabric_pod_dn = aci_fabric_pod.test.id
	}
	
	resource "aci_fabric_path_ep" "test" {
		fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
		name  = "%s"
	}

	data "aci_fabric_path_ep" "test" {
		fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
		name  = aci_fabric_path_ep.test.name
		depends_on = [ aci_fabric_path_ep.test ]
	}
	`, fabricPodName, fabricPathEpContName, name)
	return resource
}

func CreateFabricPathEpDSWithoutRequired(fabricPodName, fabricPathEpContName, name, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_path_ep Data Source without ", attrName)
	rBlock := `
	
	resource "aci_fabric_pod" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fabric_path_end-point_container" "test" {
		name 		= "%s"
		fabric_pod_dn = aci_fabric_pod.test.id
	}
	
	resource "aci_fabric_path_ep" "test" {
		fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "fabric_path_end-point_container_dn":
		rBlock += `
	data "aci_fabric_path_ep" "test" {
	#	fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
		name  = aci_fabric_path_ep.test.name
		depends_on = [ aci_fabric_path_ep.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_fabric_path_ep" "test" {
		fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
	#	name  = aci_fabric_path_ep.test.name
		depends_on = [ aci_fabric_path_ep.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fabricPodName, fabricPathEpContName, name)
}

func CreateAccFabricPathEpDSWithInvalidParentDn(fabricPodName, fabricPathEpContName, name string) string {
	fmt.Println("=== STEP  testing fabric_path_ep Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_pod" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fabric_path_end-point_container" "test" {
		name 		= "%s"
		fabric_pod_dn = aci_fabric_pod.test.id
	}
	
	resource "aci_fabric_path_ep" "test" {
		fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
		name  = "%s"
	}

	data "aci_fabric_path_ep" "test" {
		fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
		name  = "${aci_fabric_path_ep.test.name}_invalid"
		depends_on = [ aci_fabric_path_ep.test ]
	}
	`, fabricPodName, fabricPathEpContName, name)
	return resource
}

func CreateAccFabricPathEpDataSourceUpdate(fabricPodName, fabricPathEpContName, name, key, value string) string {
	fmt.Println("=== STEP  testing fabric_path_ep Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_pod" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fabric_path_end-point_container" "test" {
		name 		= "%s"
		fabric_pod_dn = aci_fabric_pod.test.id
	}
	
	resource "aci_fabric_path_ep" "test" {
		fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
		name  = "%s"
	}

	data "aci_fabric_path_ep" "test" {
		fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
		name  = aci_fabric_path_ep.test.name
		%s = "%s"
		depends_on = [ aci_fabric_path_ep.test ]
	}
	`, fabricPodName, fabricPathEpContName, name, key, value)
	return resource
}

func CreateAccFabricPathEpDataSourceUpdatedResource(fabricPodName, fabricPathEpContName, name, key, value string) string {
	fmt.Println("=== STEP  testing fabric_path_ep Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_pod" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_fabric_path_end-point_container" "test" {
		name 		= "%s"
		fabric_pod_dn = aci_fabric_pod.test.id
	}
	
	resource "aci_fabric_path_ep" "test" {
		fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_fabric_path_ep" "test" {
		fabric_path_end-point_container_dn  = aci_fabric_path_end-point_container.test.id
		name  = aci_fabric_path_ep.test.name
		depends_on = [ aci_fabric_path_ep.test ]
	}
	`, fabricPodName, fabricPathEpContName, name, key, value)
	return resource
}
