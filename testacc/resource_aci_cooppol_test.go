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

func TestAccAciCoopPolicy_Basic(t *testing.T) {
	var coop_policy_default models.COOPGroupPolicy
	var coop_policy_updated models.COOPGroupPolicy
	resourceName := "aci_coop_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCoopPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateCoopPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCoopPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCoopPolicyExists(resourceName, &coop_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "coop_group_policy_type", "compatible"),
				),
			},
			{
				Config: CreateAccCoopPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCoopPolicyExists(resourceName, &coop_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_coop_policy"),

					resource.TestCheckResourceAttr(resourceName, "coop_group_policy_type", "strict"),

					testAccCheckAciCoopPolicyIdEqual(&coop_policy_default, &coop_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccCoopPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccCoopPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccCoopPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCoopPolicyExists(resourceName, &coop_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciCoopPolicyIdNotEqual(&coop_policy_default, &coop_policy_updated),
				),
			},
		},
	})
}

func TestAccAciCoopPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCoopPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCoopPolicyConfig(rName),
			},

			{
				Config:      CreateAccCoopPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCoopPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCoopPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccCoopPolicyUpdatedAttr(rName, "coop_group_policy_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCoopPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCoopPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciCoopPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCoopPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCoopPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciCoopPolicyExists(name string, coop_policy *models.COOPGroupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Coop Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Coop Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		coop_policyFound := models.COOPGroupPolicyFromContainer(cont)
		if coop_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Coop Policy %s not found", rs.Primary.ID)
		}
		*coop_policy = *coop_policyFound
		return nil
	}
}

func testAccCheckAciCoopPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing coop_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_coop_policy" {
			cont, err := client.Get(rs.Primary.ID)
			coop_policy := models.COOPGroupPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Coop Policy %s Still exists", coop_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCoopPolicyIdEqual(m1, m2 *models.COOPGroupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("coop_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCoopPolicyIdNotEqual(m1, m2 *models.COOPGroupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("coop_policy DNs are equal")
		}
		return nil
	}
}

func CreateCoopPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing coop_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_coop_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccCoopPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing coop_policy creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccCoopPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing coop_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccCoopPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing coop_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccCoopPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple coop_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccCoopPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing coop_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_coop_policy"
		coop_group_policy_type = "strict"
		
	}
	`, rName)

	return resource
}

func CreateAccCoopPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing coop_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_coop_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_coop_policy"
		coop_group_policy_type = "strict"
		
	}
	`)

	return resource
}

func CreateAccCoopPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing coop_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccCoopPolicyUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing coop_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
