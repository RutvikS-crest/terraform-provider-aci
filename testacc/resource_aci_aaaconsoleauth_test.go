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

func TestAccAciConsoleAuthentication_Basic(t *testing.T) {
	var console_authentication_default models.ConsoleAuthenticationMethod
	var console_authentication_updated models.ConsoleAuthenticationMethod
	resourceName := "aci_console_authentication.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConsoleAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccConsoleAuthenticationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConsoleAuthenticationExists(resourceName, &console_authentication_default),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "provider_group", ""),
					resource.TestCheckResourceAttr(resourceName, "realm", "local"),
					resource.TestCheckResourceAttr(resourceName, "realm_sub_type", "default"),
				),
			},
			{
				Config: CreateAccConsoleAuthenticationConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_console_authentication"),

					resource.TestCheckResourceAttr(resourceName, "provider_group", ""),

					resource.TestCheckResourceAttr(resourceName, "realm", "ldap"),

					resource.TestCheckResourceAttr(resourceName, "realm_sub_type", "duo"),

					testAccCheckAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccConsoleAuthenticationRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestAccAciConsoleAuthentication_Update(t *testing.T) {
	var console_authentication_default models.ConsoleAuthenticationMethod
	var console_authentication_updated models.ConsoleAuthenticationMethod
	resourceName := "aci_console_authentication.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConsoleAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccConsoleAuthenticationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConsoleAuthenticationExists(resourceName, &console_authentication_default),
				),
			},

			{
				Config: CreateAccConsoleAuthenticationUpdatedAttr("realm", "radius"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "radius"),
					testAccCheckAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccConsoleAuthenticationUpdatedAttr("realm", "rsa"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "rsa"),
					testAccCheckAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccConsoleAuthenticationUpdatedAttr("realm", "saml"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "saml"),
					testAccCheckAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccConsoleAuthenticationUpdatedAttr("realm", "tacacs"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConsoleAuthenticationExists(resourceName, &console_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "tacacs"),
					testAccCheckAciConsoleAuthenticationIdEqual(&console_authentication_default, &console_authentication_updated),
				),
			},
			{
				Config: CreateAccConsoleAuthenticationConfig(),
			},
		},
	})
}

func TestAccAciConsoleAuthentication_Negative(t *testing.T) {

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConsoleAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccConsoleAuthenticationConfig(),
			},

			{
				Config:      CreateAccConsoleAuthenticationUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccConsoleAuthenticationUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccConsoleAuthenticationUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccConsoleAuthenticationUpdatedAttr("provider_group", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConsoleAuthenticationUpdatedAttr("realm", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConsoleAuthenticationUpdatedAttr("realm_sub_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccConsoleAuthenticationUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccConsoleAuthenticationConfig(),
			},
		},
	})
}

func testAccCheckAciConsoleAuthenticationExists(name string, console_authentication *models.ConsoleAuthenticationMethod) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Console Authentication %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Console Authentication dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		console_authenticationFound := models.ConsoleAuthenticationMethodFromContainer(cont)
		if console_authenticationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Console Authentication %s not found", rs.Primary.ID)
		}
		*console_authentication = *console_authenticationFound
		return nil
	}
}

func testAccCheckAciConsoleAuthenticationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing console_authentication destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_console_authentication" {
			cont, err := client.Get(rs.Primary.ID)
			console_authentication := models.ConsoleAuthenticationMethodFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Console Authentication %s Still exists", console_authentication.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciConsoleAuthenticationIdEqual(m1, m2 *models.ConsoleAuthenticationMethod) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("console_authentication DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciConsoleAuthenticationIdNotEqual(m1, m2 *models.ConsoleAuthenticationMethod) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("console_authentication DNs are equal")
		}
		return nil
	}
}

func CreateAccConsoleAuthenticationConfigWithRequiredParams() string {
	fmt.Println("=== STEP  testing console_authentication creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_console_authentication" "test" {
	
	}
	`)
	return resource
}
func CreateAccConsoleAuthenticationConfig() string {
	fmt.Println("=== STEP  testing console_authentication creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_console_authentication" "test" {
	
	}
	`)
	return resource
}

func CreateAccConsoleAuthenticationConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing console_authentication creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_console_authentication" "test" {
	
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_console_authentication"
		provider_group = ""
		realm = "ldap"
		realm_sub_type = "duo"
		
	}
	`)

	return resource
}

func CreateAccConsoleAuthenticationRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing console_authentication updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_console_authentication" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_console_authentication"
		provider_group = ""
		realm = "ldap"
		realm_sub_type = "duo"
		
	}
	`)

	return resource
}

func CreateAccConsoleAuthenticationUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing console_authentication attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_console_authentication" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
