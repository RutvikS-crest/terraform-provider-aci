package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3outBGPProtocolProfile_Basic(t *testing.T) {
	var l3out_bgp_protocol_profile_default models.L3outBGPProtocolProfile
	var l3out_bgp_protocol_profile_updated models.L3outBGPProtocolProfile
	resourceName := "aci_l3out_bgp_protocol_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBGPProtocolProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outBGPProtocolProfileWithoutRequired(rName, rName, rName, "logical_node_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outBGPProtocolProfileConfig(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists(resourceName, &l3out_bgp_protocol_profile_default),
					resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					// resource.TestCheckResourceAttr(resourceName, "relation_bgp_rs_bgp_node_ctx_pol", ""),
				),
			},
			{
				// in this step all optional attribute expect realational attribute are given for the same resource and then compared
				Config: CreateAccL3outBGPProtocolProfileConfigWithOptionalValues(rName, rName, rName), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists(resourceName, &l3out_bgp_protocol_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_bgp_protocol_profile"),
					// resource.TestCheckResourceAttr(resourceName, "relation_bgp_rs_bgp_node_ctx_pol", ""),
					testAccCheckAciL3outBGPProtocolProfileIdEqual(&l3out_bgp_protocol_profile_default, &l3out_bgp_protocol_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// {
			// 	Config: CreateAccL3outBGPProtocolProfileConfigWithRequiredParams(rName, rName, rNameUpdated),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciL3outBGPProtocolProfileExists(resourceName, &l3out_bgp_protocol_profile_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", GetParentDn(l3out_bgp_protocol_profile_default.DistinguishedName, fmt.Sprintf("/protp"))),
			// 		testAccCheckAciL3outBGPProtocolProfileIdNotEqual(&l3out_bgp_protocol_profile_default, &l3out_bgp_protocol_profile_updated),
			// 	),
			// },
			{
				Config:      CreateAccL3outBGPProtocolProfileConfigUpdateWithoutRequiredParams(rName, rName, rName, "annotation", randomValue),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outBGPProtocolProfileConfig(rName, rName, rName),
			},
		},
	})
}

// func TestAccAciL3outBGPProtocolProfile_Update(t *testing.T) {
// 	var l3out_bgp_protocol_profile_default models.L3outBGPProtocolProfile
// 	// var l3out_bgp_protocol_profile_updated models.L3outBGPProtocolProfile
// 	resourceName := "aci_l3out_bgp_protocol_profile.test"
// 	rName := makeTestVariable(acctest.RandString(5))
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckAciL3outBGPProtocolProfileDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: CreateAccL3outBGPProtocolProfileConfig(rName, rName, rName),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciL3outBGPProtocolProfileExists(resourceName, &l3out_bgp_protocol_profile_default),
// 				),
// 			},

// 			{
// 				Config: CreateAccL3outBGPProtocolProfileConfig(rName, rName, rName),
// 			},
// 		},
// 	})
// }

func TestAccAciL3outBGPProtocolProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBGPProtocolProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outBGPProtocolProfileConfig(rName, rName, rName),
			},
			{
				Config:      CreateAccL3outBGPProtocolProfileWithInValidParentDn(rName, rName, rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outBGPProtocolProfileUpdatedAttr(rName, rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outBGPProtocolProfileUpdatedAttr(rName, rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outBGPProtocolProfileUpdatedAttr(rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccL3outBGPProtocolProfileConfig(rName, rName, rName),
			},
		},
	})
}

func TestAccAciL3outBGPProtocolProfile_RelationParameters(t *testing.T) {
	var l3out_bgp_protocol_profile_default models.L3outBGPProtocolProfile
	var l3out_bgp_protocol_profile_rel models.L3outBGPProtocolProfile
	resourceName := "aci_l3out_bgp_protocol_profile.test"
	rName := acctest.RandString(5)
	randomName1 := acctest.RandString(5)
	randomName2 := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBGPProtocolProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outBGPProtocolProfileConfig(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists(resourceName, &l3out_bgp_protocol_profile_default),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd", ""),
				),
			},
			{
				Config: CreateAccL3outBGPProtocolProfileRelConfig(rName, randomName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists(resourceName, &l3out_bgp_protocol_profile_rel),
					resource.TestCheckResourceAttr(resourceName, "relation_bgp_rs_bgp_node_ctx_pol", fmt.Sprintf("uni/tn-%s/bgpCtxP-%s", rName, randomName1)),
					testAccCheckAciL3outBGPProtocolProfileIdEqual(&l3out_bgp_protocol_profile_default, &l3out_bgp_protocol_profile_rel),
				),
			},
			{
				Config: CreateAccL3outBGPProtocolProfileRelConfig(rName, randomName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists(resourceName, &l3out_bgp_protocol_profile_rel),
					resource.TestCheckResourceAttr(resourceName, "relation_bgp_rs_bgp_node_ctx_pol", fmt.Sprintf("uni/tn-%s/bgpCtxP-%s", rName, randomName2)),
					testAccCheckAciL3outBGPProtocolProfileIdEqual(&l3out_bgp_protocol_profile_default, &l3out_bgp_protocol_profile_rel),
				),
			},
			{
				Config: CreateAccL3outBGPProtocolProfileConfig(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists(resourceName, &l3out_bgp_protocol_profile_default),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd", ""),
				),
			},
		},
	})
}
func testAccCheckAciL3outBGPProtocolProfileExists(name string, l3out_bgp_protocol_profile *models.L3outBGPProtocolProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out BGP Protocol Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out BGP Protocol Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_bgp_protocol_profileFound := models.L3outBGPProtocolProfileFromContainer(cont)
		if l3out_bgp_protocol_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out BGP Protocol Profile %s not found", rs.Primary.ID)
		}
		*l3out_bgp_protocol_profile = *l3out_bgp_protocol_profileFound
		return nil
	}
}

func testAccCheckAciL3outBGPProtocolProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_bgp_protocol_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_bgp_protocol_profile" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_bgp_protocol_profile := models.L3outBGPProtocolProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out BGP Protocol Profile %s Still exists", l3out_bgp_protocol_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outBGPProtocolProfileIdEqual(m1, m2 *models.L3outBGPProtocolProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_bgp_protocol_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outBGPProtocolProfileIdNotEqual(m1, m2 *models.L3outBGPProtocolProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_bgp_protocol_profile DNs are equal")
		}
		return nil
	}
}

func CreateL3outBGPProtocolProfileWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_bgp_protocol_profile creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
		
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	`
	switch attrName {
	case "logical_node_profile_dn":
		rBlock += `
	resource "aci_l3out_bgp_protocol_profile" "test" {
	#	logical_node_profile_dn  = aci_logical_node_profile.test.id
	}
		`

	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName)
}

func CreateAccL3outBGPProtocolProfileConfigWithRequiredParams(fvTenantName, l3extOutName, l3extLNodePName string) string {
	fmt.Println("=== STEP  testing l3out_bgp_protocol_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
	}
	`, fvTenantName, l3extOutName, l3extLNodePName)
	return resource
}

func CreateAccL3outBGPProtocolProfileConfig(fvTenantName, l3extOutName, l3extLNodePName string) string {
	fmt.Println("=== STEP  testing l3out_bgp_protocol_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
	}
	`, fvTenantName, l3extOutName, l3extLNodePName)
	return resource
}

func CreateAccL3outBGPProtocolProfileRelConfig(rName, relName string) string {
	fmt.Println("=== STEP  testing l3out_bgp_protocol_profile relation parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_bgp_timers" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		relation_bgp_rs_bgp_node_ctx_pol = aci_bgp_timers.test.id
	}
	`, rName, rName, rName, relName)
	return resource
}

func CreateAccL3outBGPProtocolProfileWithInValidParentDn(fvTenantName, l3extOutName, l3extLNodePName string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_bgp_protocol_profile creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_l3_outside.test.id	
	}
	`, fvTenantName, l3extOutName, l3extLNodePName)
	return resource
}

func CreateAccL3outBGPProtocolProfileConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_bgp_protocol_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = "${aci_logical_node_profile.test.id}"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_bgp_protocol_profile"	
	}
	`, fvTenantName, l3extOutName, l3extLNodePName)

	return resource
}

func CreateAccL3outBGPProtocolProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_bgp_protocol_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_l3out_bgp_protocol_profile" "test" {
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_bgp_protocol_profile"
	}
	`)

	return resource
}

func CreateAccL3outBGPProtocolProfileConfigUpdateWithoutRequiredParams(fvTenantName, l3extOutName, l3extLNodePName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_bgp_protocol_profile attribute: %s=%s without required arguments \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, attribute, value)
	return resource
}

func CreateAccL3outBGPProtocolProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_bgp_protocol_profile attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, attribute, value)
	return resource
}
