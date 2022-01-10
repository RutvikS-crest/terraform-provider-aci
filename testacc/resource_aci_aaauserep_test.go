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

func TestAccAciGlobalSecurity_Basic(t *testing.T) {
	var global_security_default models.PasswordChangeExpirationPolicy
	var global_security_updated models.PasswordChangeExpirationPolicy
	resourceName := "aci_global_security.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciGlobalSecurityDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccGlobalSecurityConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciGlobalSecurityExists(resourceName, &global_security_default),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "pwd_strength_check", "yes"),
				),
			},
			{
				Config: CreateAccGlobalSecurityConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciGlobalSecurityExists(resourceName, &global_security_updated),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_global_security"),

					resource.TestCheckResourceAttr(resourceName, "pwd_strength_check", "no"),

					testAccCheckAciGlobalSecurityIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccGlobalSecurityRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestAccAciGlobalSecurity_Negative(t *testing.T) {

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciGlobalSecurityDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccGlobalSecurityConfig(),
			},

			{
				Config:      CreateAccGlobalSecurityUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccGlobalSecurityUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccGlobalSecurityUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccGlobalSecurityUpdatedAttr("pwd_strength_check", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccGlobalSecurityUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccGlobalSecurityConfig(),
			},
		},
	})
}

func testAccCheckAciGlobalSecurityExists(name string, global_security *models.PasswordChangeExpirationPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Global Security %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Global Security dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		global_securityFound := models.PasswordChangeExpirationPolicyFromContainer(cont)
		if global_securityFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Global Security %s not found", rs.Primary.ID)
		}
		*global_security = *global_securityFound
		return nil
	}
}

func testAccCheckAciGlobalSecurityDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing global_security destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_global_security" {
			cont, err := client.Get(rs.Primary.ID)
			global_security := models.PasswordChangeExpirationPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Global Security %s Still exists", global_security.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciGlobalSecurityIdEqual(m1, m2 *models.PasswordChangeExpirationPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("global_security DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciGlobalSecurityIdNotEqual(m1, m2 *models.PasswordChangeExpirationPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("global_security DNs are equal")
		}
		return nil
	}
}

func CreateAccGlobalSecurityConfigWithRequiredParams() string {
	fmt.Println("=== STEP  testing global_security creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {
	
	}
	`)
	return resource
}
func CreateAccGlobalSecurityConfig() string {
	fmt.Println("=== STEP  testing global_security creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {
	
	}
	`)
	return resource
}

func CreateAccGlobalSecurityConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing global_security creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {
	
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_global_security"
		pwd_strength_check = "no"
		
	}
	`)

	return resource
}

func CreateAccGlobalSecurityRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing global_security updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_global_security" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_global_security"
		pwd_strength_check = "no"
		
	}
	`)

	return resource
}

func CreateAccGlobalSecurityUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing global_security attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}

func CreateAccGlobalSecurityUpdatedAttrList(attribute, value string) string {
	fmt.Printf("=== STEP  testing global_security attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {
	
		%s = %s
	}
	`, attribute, value)
	return resource
}
