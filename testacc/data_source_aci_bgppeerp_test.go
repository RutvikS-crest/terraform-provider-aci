package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)




  

func TestAccAciPeerConnectivityProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_peer_connectivity_profile.test"
	dataSourceName := "data.aci_peer_connectivity_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		ProviderFactories:    testAccProviders,
		CheckDestroy: testAccCheckAciPeerConnectivityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreatePeerConnectivityProfileDSWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, addr,"logical_node_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreatePeerConnectivityProfileDSWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, addr, "addr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPeerConnectivityProfileConfigDataSource(fvTenantName, l3extOutName, l3extLNodePName, addr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "logical_node_profile_dn", resourceName, "logical_node_profile_dn",),
					resource.TestCheckResourceAttrPair(dataSourceName, "addr", resourceName, "addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "addr_t_ctrl.#", resourceName, "addr_t_ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "addr_t_ctrl.0", resourceName, "addr_t_ctrl.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "allowed_self_as_cnt", resourceName, "allowed_self_as_cnt"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.#", resourceName, "ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.0", resourceName, "ctrl.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "password", resourceName, "password"),
					resource.TestCheckResourceAttrPair(dataSourceName, "peer_ctrl.#", resourceName, "peer_ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "peer_ctrl.0", resourceName, "peer_ctrl.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "private_a_sctrl.#", resourceName, "private_a_sctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "private_a_sctrl.0", resourceName, "private_a_sctrl.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ttl", resourceName, "ttl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "weight", resourceName, "weight"),
					
				),
			},
			{
				Config:      CreateAccPeerConnectivityProfileDataSourceUpdate(fvTenantName, l3extOutName, l3extLNodePName, addr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			
			{
				Config:      CreateAccPeerConnectivityProfileDSWithInvalidParentDn(fvTenantName, l3extOutName, l3extLNodePName, addr),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			
			{
				Config: CreateAccPeerConnectivityProfileDataSourceUpdatedResource(fvTenantName, l3extOutName, l3extLNodePName, addr, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}


func CreateAccPeerConnectivityProfileConfigDataSource(fvTenantName, l3extOutName, l3extLNodePName, addr string) string {
	fmt.Println("=== STEP  testing peer_connectivity_profile Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_peer_connectivity_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
	}

	data "aci_peer_connectivity_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		addr  = aci_peer_connectivity_profile.test.addr
		depends_on = [ aci_peer_connectivity_profile.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, addr)
	return resource
}

func CreatePeerConnectivityProfileDSWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, addr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing peer_connectivity_profile creation without ",attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_peer_connectivity_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
	}
	`
	switch attrName {
	case "logical_node_profile_dn":
		rBlock += `
	data "aci_peer_connectivity_profile" "test" {
	#	logical_node_profile_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
		depends_on = [ aci_peer_connectivity_profile.test ]
	}
		`
	case "addr":
		rBlock += `
	data "aci_peer_connectivity_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
	#	addr  = "%s"
		depends_on = [ aci_peer_connectivity_profile.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock,fvTenantName, l3extOutName, l3extLNodePName, addr)
}

func CreateAccPeerConnectivityProfileDSWithInvalidParentDn(fvTenantName, l3extOutName, l3extLNodePName, addr string) string {
	fmt.Println("=== STEP  testing peer_connectivity_profile Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_peer_connectivity_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
	}

	data "aci_peer_connectivity_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		addr  = "${aci_peer_connectivity_profile.test.addr}_invalid"
		depends_on = [ aci_peer_connectivity_profile.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, addr)
	return resource
}

func CreateAccPeerConnectivityProfileDataSourceUpdate(fvTenantName, l3extOutName, l3extLNodePName, addr, key, value string) string {
	fmt.Println("=== STEP  testing peer_connectivity_profile Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_peer_connectivity_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
	}

	data "aci_peer_connectivity_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		addr  = aci_peer_connectivity_profile.test.addr
		%s = "%s"
		depends_on = [ aci_peer_connectivity_profile.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, addr,key,value)
	return resource
}

func CreateAccPeerConnectivityProfileDataSourceUpdatedResource(fvTenantName, l3extOutName, l3extLNodePName, addr, key, value string) string {
	fmt.Println("=== STEP  testing peer_connectivity_profile Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_peer_connectivity_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
		%s = "%s"
	}

	data "aci_peer_connectivity_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		addr  = aci_peer_connectivity_profile.test.addr
		depends_on = [ aci_peer_connectivity_profile.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, addr,key,value)
	return resource
}