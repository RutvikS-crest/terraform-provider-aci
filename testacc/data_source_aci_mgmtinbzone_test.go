package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciMgmtZoneDataSource_Basic(t *testing.T) {
	resourceName := "aci_mgmt_zone.test"
	dataSourceName := "data.aci_mgmt_zone.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	mgmtGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMgmtZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateMgmtZoneDSWithoutRequired(mgmtGrpName, "managed_node_connectivity_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			}, {
				Config: CreateAccMgmtZoneConfigDataSource(mgmtGrpName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "managed_node_connectivity_group_dn", resourceName, "managed_node_connectivity_group_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccMgmtZoneDataSourceUpdate(mgmtGrpName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccMgmtZoneDSWithInvalidParentDn(mgmtGrpName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccMgmtZoneDataSourceUpdatedResource(mgmtGrpName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccMgmtZoneConfigDataSource(mgmtGrpName string) string {
	fmt.Println("=== STEP  testing mgmt_zone Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
	}

	data "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		depends_on = [ aci_mgmt_zone.test ]
	}
	`, mgmtGrpName)
	return resource
}

func CreateMgmtZoneDSWithoutRequired(mgmtGrpName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing mgmt_zone Data Source without ", attrName)
	rBlock := `
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
	}
	`
	switch attrName {
	case "managed_node_connectivity_group_dn":
		rBlock += `
	data "aci_mgmt_zone" "test" {
	#	managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
	
		depends_on = [ aci_mgmt_zone.test ]
	}
		`

	}
	return fmt.Sprintf(rBlock, mgmtGrpName)
}

func CreateAccMgmtZoneDSWithInvalidParentDn(mgmtGrpName string) string {
	fmt.Println("=== STEP  testing mgmt_zone Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
	}

	data "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		depends_on = [ aci_mgmt_zone.test ]
	}
	`, mgmtGrpName)
	return resource
}

func CreateAccMgmtZoneDataSourceUpdate(mgmtGrpName, key, value string) string {
	fmt.Println("=== STEP  testing mgmt_zone Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
	}

	data "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		%s = "%s"
		depends_on = [ aci_mgmt_zone.test ]
	}
	`, mgmtGrpName, key, value)
	return resource
}

func CreateAccMgmtZoneDataSourceUpdatedResource(mgmtGrpName, key, value string) string {
	fmt.Println("=== STEP  testing mgmt_zone Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		%s = "%s"
	}

	data "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		depends_on = [ aci_mgmt_zone.test ]
	}
	`, mgmtGrpName, key, value)
	return resource
}
