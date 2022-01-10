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

func TestAccAciErrorDisabledRecoveryPolicy_Basic(t *testing.T) {
	var error_disable_recovery_default models.ErrorDisabledRecoveryPolicy
	var error_disable_recovery_updated models.ErrorDisabledRecoveryPolicy
	resourceName := "aci_error_disable_recovery.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciErrorDisabledRecoveryPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateErrorDisabledRecoveryPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccErrorDisabledRecoveryPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryPolicyExists(resourceName, &error_disable_recovery_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "err_dis_recov_intvl", "300"),
				),
			},
			{
				Config: CreateAccErrorDisabledRecoveryPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryPolicyExists(resourceName, &error_disable_recovery_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_error_disable_recovery"),
					resource.TestCheckResourceAttr(resourceName, "err_dis_recov_intvl", "31"),

					testAccCheckAciErrorDisabledRecoveryPolicyIdEqual(&error_disable_recovery_default, &error_disable_recovery_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccErrorDisabledRecoveryPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccErrorDisabledRecoveryPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccErrorDisabledRecoveryPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryPolicyExists(resourceName, &error_disable_recovery_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciErrorDisabledRecoveryPolicyIdNotEqual(&error_disable_recovery_default, &error_disable_recovery_updated),
				),
			},
		},
	})
}

func TestAccAciErrorDisabledRecoveryPolicy_Update(t *testing.T) {
	var error_disable_recovery_default models.ErrorDisabledRecoveryPolicy
	var error_disable_recovery_updated models.ErrorDisabledRecoveryPolicy
	resourceName := "aci_error_disable_recovery.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciErrorDisabledRecoveryPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccErrorDisabledRecoveryPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryPolicyExists(resourceName, &error_disable_recovery_default),
				),
			},
			{
				Config: CreateAccErrorDisabledRecoveryPolicyUpdatedAttr(rName, "err_dis_recov_intvl", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryPolicyExists(resourceName, &error_disable_recovery_updated),
					resource.TestCheckResourceAttr(resourceName, "err_dis_recov_intvl", "65535"),
					testAccCheckAciErrorDisabledRecoveryPolicyIdEqual(&error_disable_recovery_default, &error_disable_recovery_updated),
				),
			},
			{
				Config: CreateAccErrorDisabledRecoveryPolicyUpdatedAttr(rName, "err_dis_recov_intvl", "32752"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryPolicyExists(resourceName, &error_disable_recovery_updated),
					resource.TestCheckResourceAttr(resourceName, "err_dis_recov_intvl", "32752"),
					testAccCheckAciErrorDisabledRecoveryPolicyIdEqual(&error_disable_recovery_default, &error_disable_recovery_updated),
				),
			},

			{
				Config: CreateAccErrorDisabledRecoveryPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciErrorDisabledRecoveryPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciErrorDisabledRecoveryPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccErrorDisabledRecoveryPolicyConfig(rName),
			},

			{
				Config:      CreateAccErrorDisabledRecoveryPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccErrorDisabledRecoveryPolicyUpdatedAttr(rName, "err_dis_recov_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryPolicyUpdatedAttr(rName, "err_dis_recov_intvl", "29"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryPolicyUpdatedAttr(rName, "err_dis_recov_intvl", "65536"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccErrorDisabledRecoveryPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccErrorDisabledRecoveryPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciErrorDisabledRecoveryPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciErrorDisabledRecoveryPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccErrorDisabledRecoveryPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciErrorDisabledRecoveryPolicyExists(name string, error_disable_recovery *models.ErrorDisabledRecoveryPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Error Disable Recovery %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Error Disable Recovery dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		error_disable_recoveryFound := models.ErrorDisabledRecoveryPolicyFromContainer(cont)
		if error_disable_recoveryFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Error Disable Recovery %s not found", rs.Primary.ID)
		}
		*error_disable_recovery = *error_disable_recoveryFound
		return nil
	}
}

func testAccCheckAciErrorDisabledRecoveryPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing error_disable_recovery destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_error_disable_recovery" {
			cont, err := client.Get(rs.Primary.ID)
			error_disable_recovery := models.ErrorDisabledRecoveryPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Error Disable Recovery %s Still exists", error_disable_recovery.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciErrorDisabledRecoveryPolicyIdEqual(m1, m2 *models.ErrorDisabledRecoveryPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("error_disable_recovery DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciErrorDisabledRecoveryPolicyIdNotEqual(m1, m2 *models.ErrorDisabledRecoveryPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("error_disable_recovery DNs are equal")
		}
		return nil
	}
}

func CreateErrorDisabledRecoveryPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing error_disable_recovery creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_error_disable_recovery" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccErrorDisabledRecoveryPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing error_disable_recovery creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccErrorDisabledRecoveryPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing error_disable_recovery creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccErrorDisabledRecoveryPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing error_disable_recovery creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccErrorDisabledRecoveryPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple error_disable_recovery creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccErrorDisabledRecoveryPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing error_disable_recovery creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_error_disable_recovery"
		err_dis_recov_intvl = "31"
		
	}
	`, rName)

	return resource
}

func CreateAccErrorDisabledRecoveryPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing error_disable_recovery updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_error_disable_recovery" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_error_disable_recovery"
		err_dis_recov_intvl = "31"
		
	}
	`)

	return resource
}

func CreateAccErrorDisabledRecoveryPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing error_disable_recovery attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccErrorDisabledRecoveryPolicyUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing error_disable_recovery attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
