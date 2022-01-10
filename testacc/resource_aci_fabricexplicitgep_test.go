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

func TestAccAciVPCExplicitProtectionGroup_Basic(t *testing.T) {
	var vpc_explicit_protection_group_default models.VPCExplicitProtectionGroup
	var vpc_explicit_protection_group_updated models.VPCExplicitProtectionGroup
	resourceName := "aci_vpc_explicit_protection_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCExplicitProtectionGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateVPCExplicitProtectionGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists(resourceName, &vpc_explicit_protection_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_explicit_protection_group_id", ""),
				),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists(resourceName, &vpc_explicit_protection_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vpc_explicit_protection_group"),

					testAccCheckAciVPCExplicitProtectionGroupIdEqual(&vpc_explicit_protection_group_default, &vpc_explicit_protection_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccVPCExplicitProtectionGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists(resourceName, &vpc_explicit_protection_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciVPCExplicitProtectionGroupIdNotEqual(&vpc_explicit_protection_group_default, &vpc_explicit_protection_group_updated),
				),
			},
		},
	})
}

func TestAccAciVPCExplicitProtectionGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCExplicitProtectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVPCExplicitProtectionGroupConfig(rName),
			},

			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(rName, "vpc_explicit_protection_group_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccVPCExplicitProtectionGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupConfig(rName),
			},
		},
	})
}

func TestAccAciVPCExplicitProtectionGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCExplicitProtectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVPCExplicitProtectionGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciVPCExplicitProtectionGroupExists(name string, vpc_explicit_protection_group *models.VPCExplicitProtectionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VPC Explicit Protection Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC Explicit Protection Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vpc_explicit_protection_groupFound := models.VPCExplicitProtectionGroupFromContainer(cont)
		if vpc_explicit_protection_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VPC Explicit Protection Group %s not found", rs.Primary.ID)
		}
		*vpc_explicit_protection_group = *vpc_explicit_protection_groupFound
		return nil
	}
}

func testAccCheckAciVPCExplicitProtectionGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vpc_explicit_protection_group" {
			cont, err := client.Get(rs.Primary.ID)
			vpc_explicit_protection_group := models.VPCExplicitProtectionGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VPC Explicit Protection Group %s Still exists", vpc_explicit_protection_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVPCExplicitProtectionGroupIdEqual(m1, m2 *models.VPCExplicitProtectionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vpc_explicit_protection_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVPCExplicitProtectionGroupIdNotEqual(m1, m2 *models.VPCExplicitProtectionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vpc_explicit_protection_group DNs are equal")
		}
		return nil
	}
}

func CreateVPCExplicitProtectionGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vpc_explicit_protection_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_vpc_explicit_protection_group" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccVPCExplicitProtectionGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccVPCExplicitProtectionGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccVPCExplicitProtectionGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccVPCExplicitProtectionGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple vpc_explicit_protection_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccVPCExplicitProtectionGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing vpc_explicit_protection_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vpc_explicit_protection_group"
		
	}
	`, rName)

	return resource
}

func CreateAccVPCExplicitProtectionGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vpc_explicit_protection_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vpc_explicit_protection_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vpc_explicit_protection_group"
		
	}
	`)

	return resource
}

func CreateAccVPCExplicitProtectionGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vpc_explicit_protection_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccVPCExplicitProtectionGroupUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vpc_explicit_protection_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
