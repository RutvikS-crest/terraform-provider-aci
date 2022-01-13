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

func TestAccAciLDAPProvider_Basic(t *testing.T) {
	var ldap_provider_default models.LDAPProvider
	var ldap_provider_updated models.LDAPProvider
	resourceName := "aci_ldap_provider.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPProviderDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLDAPProviderWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLDAPProviderConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "ssl_validation_level", "strict"),
					resource.TestCheckResourceAttr(resourceName, "attribute", ""),
					resource.TestCheckResourceAttr(resourceName, "basedn", ""),
					resource.TestCheckResourceAttr(resourceName, "enable_ssl", "no"),
					resource.TestCheckResourceAttr(resourceName, "filter", ""),
					resource.TestCheckResourceAttr(resourceName, "key", ""),
					resource.TestCheckResourceAttr(resourceName, "monitor_server", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_password", ""),
					resource.TestCheckResourceAttr(resourceName, "monitoring_user", "default"),
					resource.TestCheckResourceAttr(resourceName, "port", "389"),
					resource.TestCheckResourceAttr(resourceName, "retries", "1"),
					resource.TestCheckResourceAttr(resourceName, "rootdn", ""),
					resource.TestCheckResourceAttr(resourceName, "timeout", "30"),
				),
			},
			{
				Config: CreateAccLDAPProviderConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ldap_provider"),

					resource.TestCheckResourceAttr(resourceName, "ssl_validation_level", "permissive"),

					resource.TestCheckResourceAttr(resourceName, "attribute", ""),

					resource.TestCheckResourceAttr(resourceName, "basedn", ""),

					resource.TestCheckResourceAttr(resourceName, "enable_ssl", "yes"),

					resource.TestCheckResourceAttr(resourceName, "filter", ""),

					resource.TestCheckResourceAttr(resourceName, "key", ""),

					resource.TestCheckResourceAttr(resourceName, "monitor_server", "enabled"),

					resource.TestCheckResourceAttr(resourceName, "monitoring_password", ""),

					resource.TestCheckResourceAttr(resourceName, "monitoring_user", ""),
					resource.TestCheckResourceAttr(resourceName, "port", "2"),
					resource.TestCheckResourceAttr(resourceName, "retries", "1"),

					resource.TestCheckResourceAttr(resourceName, "rootdn", ""),
					resource.TestCheckResourceAttr(resourceName, "timeout", "6"),

					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLDAPProviderConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccLDAPProviderRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccLDAPProviderConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLDAPProviderIdNotEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
		},
	})
}

func TestAccAciLDAPProvider_Update(t *testing.T) {
	var ldap_provider_default models.LDAPProvider
	var ldap_provider_updated models.LDAPProvider
	resourceName := "aci_ldap_provider.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPProviderConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_default),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "port", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "port", "65535"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "port", "32767"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "port", "32767"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "retries", "5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "retries", "5"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "retries", "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "retries", "2"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "timeout", "60"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "60"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},
			{
				Config: CreateAccLDAPProviderUpdatedAttr(rName, "timeout", "27"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists(resourceName, &ldap_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "27"),
					testAccCheckAciLDAPProviderIdEqual(&ldap_provider_default, &ldap_provider_updated),
				),
			},

			{
				Config: CreateAccLDAPProviderConfig(rName),
			},
		},
	})
}

func TestAccAciLDAPProvider_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPProviderConfig(rName),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "ssl_validation_level", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "attribute", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "basedn", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "enable_ssl", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "filter", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "key", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "monitor_server", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "monitoring_password", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "monitoring_user", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "port", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "port", "65536"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "retries", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "retries", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "retries", "6"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "rootdn", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "timeout", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "timeout", "4"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, "timeout", "61"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccLDAPProviderUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLDAPProviderConfig(rName),
			},
		},
	})
}

func TestAccAciLDAPProvider_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLDAPProviderConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciLDAPProviderExists(name string, ldap_provider *models.LDAPProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LDAP Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LDAP Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ldap_providerFound := models.LDAPProviderFromContainer(cont)
		if ldap_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LDAP Provider %s not found", rs.Primary.ID)
		}
		*ldap_provider = *ldap_providerFound
		return nil
	}
}

func testAccCheckAciLDAPProviderDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing ldap_provider destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ldap_provider" {
			cont, err := client.Get(rs.Primary.ID)
			ldap_provider := models.LDAPProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LDAP Provider %s Still exists", ldap_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLDAPProviderIdEqual(m1, m2 *models.LDAPProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("ldap_provider DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLDAPProviderIdNotEqual(m1, m2 *models.LDAPProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("ldap_provider DNs are equal")
		}
		return nil
	}
}

func CreateLDAPProviderWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_provider creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_ldap_provider" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLDAPProviderConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing ldap_provider creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccLDAPProviderConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing ldap_provider creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLDAPProviderConfig(rName string) string {
	fmt.Println("=== STEP  testing ldap_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLDAPProviderConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple ldap_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccLDAPProviderConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_provider creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ldap_provider"
		ssl_validation_level = "permissive"
		attribute = ""
		basedn = ""
		enable_ssl = "yes"
		filter = ""
		key = ""
		monitor_server = "enabled"
		monitoring_password = ""
		monitoring_user = ""
		port = "2"
		retries = "1"
		rootdn = ""
		timeout = "6"
		
	}
	`, rName)

	return resource
}

func CreateAccLDAPProviderRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing ldap_provider updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_ldap_provider" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ldap_provider"
		ssl_validation_level = "permissive"
		attribute = ""
		basedn = ""
		enable_ssl = "yes"
		filter = ""
		key = ""
		monitor_server = "enabled"
		monitoring_password = ""
		monitoring_user = ""
		port = "2"
		retries = "1"
		rootdn = ""
		timeout = "6"
		
	}
	`)

	return resource
}

func CreateAccLDAPProviderUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing ldap_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
