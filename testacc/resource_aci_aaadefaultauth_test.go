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

func TestAccAciDefaultAuthentication_Basic(t *testing.T) {
	var default_authentication_default models.DefaultAuthenticationMethodforallLogins
	var default_authentication_updated models.DefaultAuthenticationMethodforallLogins
	resourceName := "aci_default_authentication.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDefaultAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDefaultAuthenticationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_default),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "fallback_check", "no"),
					resource.TestCheckResourceAttr(resourceName, "provider_group", ""),
					resource.TestCheckResourceAttr(resourceName, "realm", "local"),
					resource.TestCheckResourceAttr(resourceName, "realm_sub_type", "default"),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_default_authentication"),

					resource.TestCheckResourceAttr(resourceName, "fallback_check", "false"),

					resource.TestCheckResourceAttr(resourceName, "provider_group", ""),

					resource.TestCheckResourceAttr(resourceName, "realm", "ldap"),

					resource.TestCheckResourceAttr(resourceName, "realm_sub_type", "duo"),

					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccDefaultAuthenticationRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestAccAciDefaultAuthentication_Update(t *testing.T) {
	var default_authentication_default models.DefaultAuthenticationMethodforallLogins
	var default_authentication_updated models.DefaultAuthenticationMethodforallLogins
	resourceName := "aci_default_authentication.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDefaultAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDefaultAuthenticationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_default),
				),
			},

			{
				Config: CreateAccDefaultAuthenticationUpdatedAttr("realm", "radius"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "radius"),
					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationUpdatedAttr("realm", "rsa"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "rsa"),
					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationUpdatedAttr("realm", "saml"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "saml"),
					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationUpdatedAttr("realm", "tacacs"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDefaultAuthenticationExists(resourceName, &default_authentication_updated),
					resource.TestCheckResourceAttr(resourceName, "realm", "tacacs"),
					testAccCheckAciDefaultAuthenticationIdEqual(&default_authentication_default, &default_authentication_updated),
				),
			},
			{
				Config: CreateAccDefaultAuthenticationConfig(),
			},
		},
	})
}

func TestAccAciDefaultAuthentication_Negative(t *testing.T) {

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDefaultAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDefaultAuthenticationConfig(),
			},

			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("fallback_check", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("provider_group", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("realm", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr("realm_sub_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccDefaultAuthenticationUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccDefaultAuthenticationConfig(),
			},
		},
	})
}

func TestAccAciDefaultAuthentication_MultipleCreateDelete(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDefaultAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDefaultAuthenticationConfigMultiple(),
			},
		},
	})
}

func testAccCheckAciDefaultAuthenticationExists(name string, default_authentication *models.DefaultAuthenticationMethodforallLogins) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Default Authentication %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Default Authentication dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		default_authenticationFound := models.DefaultAuthenticationMethodforallLoginsFromContainer(cont)
		if default_authenticationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Default Authentication %s not found", rs.Primary.ID)
		}
		*default_authentication = *default_authenticationFound
		return nil
	}
}

func testAccCheckAciDefaultAuthenticationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing default_authentication destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_default_authentication" {
			cont, err := client.Get(rs.Primary.ID)
			default_authentication := models.DefaultAuthenticationMethodforallLoginsFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Default Authentication %s Still exists", default_authentication.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciDefaultAuthenticationIdEqual(m1, m2 *models.DefaultAuthenticationMethodforallLogins) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("default_authentication DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciDefaultAuthenticationIdNotEqual(m1, m2 *models.DefaultAuthenticationMethodforallLogins) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("default_authentication DNs are equal")
		}
		return nil
	}
}

func CreateAccDefaultAuthenticationConfigWithRequiredParams() string {
	fmt.Println("=== STEP  testing default_authentication creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
	}
	`)
	return resource
}
func CreateAccDefaultAuthenticationConfig() string {
	fmt.Println("=== STEP  testing default_authentication creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
	}
	`)
	return resource
}

func CreateAccDefaultAuthenticationConfigMultiple() string {
	fmt.Println("=== STEP  testing multiple default_authentication creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
		count = 5
	}
	`)
	return resource
}

func CreateAccDefaultAuthenticationConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing default_authentication creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_default_authentication"
		fallback_check = "false"
		provider_group = ""
		realm = "ldap"
		realm_sub_type = "duo"
		
	}
	`)

	return resource
}

func CreateAccDefaultAuthenticationRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing default_authentication updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_default_authentication" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_default_authentication"
		fallback_check = "false"
		provider_group = ""
		realm = "ldap"
		realm_sub_type = "duo"
		
	}
	`)

	return resource
}

func CreateAccDefaultAuthenticationUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing default_authentication attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
