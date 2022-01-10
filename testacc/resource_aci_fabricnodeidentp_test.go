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

var FabricNodeMemberNodeId = "201"

func TestAccAciFabricNodeMember_Basic(t *testing.T) {
	var fabric_node_member_default models.FabricNodeMember
	var fabric_node_member_updated models.FabricNodeMember
	resourceName := "aci_fabric_node_member.test"

	rName := makeTestVariable(acctest.RandString(5))
	serial := "1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeMemberDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFabricNodeMemberWithoutRequired(serial, rName, "serial"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricNodeMemberConfig(serial, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_default),

					resource.TestCheckResourceAttr(resourceName, "serial", serial),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "ext_pool_id", "['0']"),
					resource.TestCheckResourceAttr(resourceName, "fabric_id", "['1']"),
					resource.TestCheckResourceAttr(resourceName, "node_id", FabricNodeMemberNodeId),
					resource.TestCheckResourceAttr(resourceName, "node_type", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "role", "unspecified"),
				),
			},
			{
				Config: CreateAccFabricNodeMemberConfigWithOptionalValues(serial, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),

					resource.TestCheckResourceAttr(resourceName, "serial", serial),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fabric_node_member"),
					resource.TestCheckResourceAttr(resourceName, "node_id", FabricNodeMemberNodeId),

					resource.TestCheckResourceAttr(resourceName, "node_type", "remote-leaf-wan"),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "2"),

					resource.TestCheckResourceAttr(resourceName, "role", "leaf"),

					testAccCheckAciFabricNodeMemberIdEqual(&fabric_node_member_default, &fabric_node_member_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccFabricNodeMemberRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccFabricNodeMemberConfigWithRequiredParams("20", rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "serial", "20"),
					testAccCheckAciFabricNodeMemberIdNotEqual(&fabric_node_member_default, &fabric_node_member_updated),
				),
			},
		},
	})
}

func TestAccAciFabricNodeMember_Update(t *testing.T) {
	var fabric_node_member_default models.FabricNodeMember
	var fabric_node_member_updated models.FabricNodeMember
	resourceName := "aci_fabric_node_member.test"

	rName := makeTestVariable(acctest.RandString(5))
	serial := "2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricNodeMemberConfig(serial, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_default),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "node_id", "4000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "node_id", "4000"),
					testAccCheckAciFabricNodeMemberIdEqual(&fabric_node_member_default, &fabric_node_member_updated),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "node_id", "1949"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "node_id", "1949"),
					testAccCheckAciFabricNodeMemberIdEqual(&fabric_node_member_default, &fabric_node_member_updated),
				),
			},

			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "node_type", "tier-2-leaf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "node_type", "tier-2-leaf"),
					testAccCheckAciFabricNodeMemberIdEqual(&fabric_node_member_default, &fabric_node_member_updated),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "node_type", "virtual"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "node_type", "virtual"),
					testAccCheckAciFabricNodeMemberIdEqual(&fabric_node_member_default, &fabric_node_member_updated),
				),
			}, {
				Config: CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "pod_id", "254"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "254"),
					testAccCheckAciFabricNodeMemberIdEqual(&fabric_node_member_default, &fabric_node_member_updated),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "pod_id", "126"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "126"),
					testAccCheckAciFabricNodeMemberIdEqual(&fabric_node_member_default, &fabric_node_member_updated),
				),
			},

			{
				Config: CreateAccFabricNodeMemberConfig(serial, rName),
			},
		},
	})
}

func TestAccAciFabricNodeMember_Negative(t *testing.T) {

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	serial := "3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricNodeMemberConfig(serial, rName),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "ext_pool_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "fabric_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "node_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "node_id", "100"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "node_id", "4001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "node_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "pod_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "pod_id", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "pod_id", "255"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, "role", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(serial, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFabricNodeMemberConfig(serial, rName),
			},
		},
	})
}

func TestAccAciFabricNodeMember_MultipleCreateDelete(t *testing.T) {

	rName := makeTestVariable(acctest.RandString(5))
	serial := "5"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricNodeMemberConfigMultiple(serial, rName),
			},
		},
	})
}

func testAccCheckAciFabricNodeMemberExists(name string, fabric_node_member *models.FabricNodeMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric Node Member %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric Node Member dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_node_memberFound := models.FabricNodeMemberFromContainer(cont)
		if fabric_node_memberFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric Node Member %s not found", rs.Primary.ID)
		}
		*fabric_node_member = *fabric_node_memberFound
		return nil
	}
}

func testAccCheckAciFabricNodeMemberDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing fabric_node_member destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fabric_node_member" {
			cont, err := client.Get(rs.Primary.ID)
			fabric_node_member := models.FabricNodeMemberFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric Node Member %s Still exists", fabric_node_member.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFabricNodeMemberIdEqual(m1, m2 *models.FabricNodeMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fabric_node_member DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFabricNodeMemberIdNotEqual(m1, m2 *models.FabricNodeMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("fabric_node_member DNs are equal")
		}
		return nil
	}
}

func CreateFabricNodeMemberWithoutRequired(serial, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_node_member creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "serial":
		rBlock += `
	resource "aci_fabric_node_member" "test" {
	#	serial  = "%s"
		name = "%s"
		node_id = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, serial, rName, FabricNodeMemberNodeId)
}

func CreateAccFabricNodeMemberConfigWithRequiredParams(serial, rName string) string {
	fmt.Println("=== STEP  testing fabric_node_member creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
	
		serial  = "%s"
		name = "%s"
		node_id = "%s"
	}
	`, serial, rName, FabricNodeMemberNodeId)
	return resource
}

func CreateAccFabricNodeMemberConfig(serial, rName string) string {
	fmt.Println("=== STEP  testing fabric_node_member creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
		serial  = "%s"
		name = "%s"
		node_id = "%s"
	}
	`, serial, rName, FabricNodeMemberNodeId)
	return resource
}

func CreateAccFabricNodeMemberConfigMultiple(serial, rName string) string {
	fmt.Println("=== STEP  testing multiple fabric_node_member creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
		serial  = %s+count.index
		name = "%s_${count.index}"
		node_id = "%s"
		count = 5
	}
	`, serial, rName, FabricNodeMemberNodeId)
	return resource
}

func CreateAccFabricNodeMemberConfigWithOptionalValues(serial, rName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_node_member creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
	
		serial  = "%s"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_node_member"
		node_id = "%s"
		node_type = "remote-leaf-wan"
		pod_id = "2"
		role = "leaf"
		
	}
	`, serial, rName, FabricNodeMemberNodeId)

	return resource
}

func CreateAccFabricNodeMemberRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing fabric_node_member updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_member" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_node_member"
		node_id = "102"
		node_type = "remote-leaf-wan"
		pod_id = "2"
		role = "leaf"
		
	}
	`)

	return resource
}

func CreateAccFabricNodeMemberUpdatedAttr(serial, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing fabric_node_member attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
	
		serial  = "%s"
		name  = "%s"
		node_id = "%s"
		%s = "%s"
	}
	`, serial, rName, FabricNodeMemberNodeId, attribute, value)
	return resource
}

func CreateAccFabricNodeMemberUpdatedAttrList(serial, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing fabric_node_member attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
	
		serial  = "%s"
		name  = "%s"
		node_id = "%s"
		%s = %s
	}
	`, serial, rName, FabricNodeMemberNodeId, attribute, value)
	return resource
}
