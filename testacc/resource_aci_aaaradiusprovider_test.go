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

func TestAccAciRadiusProvider_Basic(t *testing.T) {
	var radius_provider_default models.RADIUSProvider
	var radius_provider_updated models.RADIUSProvider
	resourceName := "aci_radius_provider.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRadiusProviderDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateRadiusProviderWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRadiusProviderConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderExists(resourceName, &radius_provider_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auth_port", "1812"),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "pap"),
					resource.TestCheckResourceAttr(resourceName, "key", ""),
					resource.TestCheckResourceAttr(resourceName, "monitor_server", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "monitoring_password", ""),
					resource.TestCheckResourceAttr(resourceName, "monitoring_user", "default"),
					resource.TestCheckResourceAttr(resourceName, "retries", "1"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "5"),
				),
			},
			{
				Config: CreateAccRadiusProviderConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderExists(resourceName, &radius_provider_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_radius_provider"),
					resource.TestCheckResourceAttr(resourceName, "auth_port", "2"),

					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "chap"),

					resource.TestCheckResourceAttr(resourceName, "key", ""),

					resource.TestCheckResourceAttr(resourceName, "monitor_server", "enabled"),

					resource.TestCheckResourceAttr(resourceName, "monitoring_password", ""),

					resource.TestCheckResourceAttr(resourceName, "monitoring_user", ""),
					resource.TestCheckResourceAttr(resourceName, "retries", "1"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "1"),

					testAccCheckAciRadiusProviderIdEqual(&radius_provider_default, &radius_provider_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccRadiusProviderConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccRadiusProviderRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccRadiusProviderConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderExists(resourceName, &radius_provider_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciRadiusProviderIdNotEqual(&radius_provider_default, &radius_provider_updated),
				),
			},
		},
	})
}

func TestAccAciRadiusProvider_Update(t *testing.T) {
	var radius_provider_default models.RADIUSProvider
	var radius_provider_updated models.RADIUSProvider
	resourceName := "aci_radius_provider.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRadiusProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRadiusProviderConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderExists(resourceName, &radius_provider_default),
				),
			},
			{
				Config: CreateAccRadiusProviderUpdatedAttr(rName, "auth_port", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderExists(resourceName, &radius_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "auth_port", "65535"),
					testAccCheckAciRadiusProviderIdEqual(&radius_provider_default, &radius_provider_updated),
				),
			},
			{
				Config: CreateAccRadiusProviderUpdatedAttr(rName, "auth_port", "32767"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderExists(resourceName, &radius_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "auth_port", "32767"),
					testAccCheckAciRadiusProviderIdEqual(&radius_provider_default, &radius_provider_updated),
				),
			},
			{
				Config: CreateAccRadiusProviderUpdatedAttr(rName, "retries", "5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderExists(resourceName, &radius_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "retries", "5"),
					testAccCheckAciRadiusProviderIdEqual(&radius_provider_default, &radius_provider_updated),
				),
			},
			{
				Config: CreateAccRadiusProviderUpdatedAttr(rName, "retries", "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderExists(resourceName, &radius_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "retries", "2"),
					testAccCheckAciRadiusProviderIdEqual(&radius_provider_default, &radius_provider_updated),
				),
			},
			{
				Config: CreateAccRadiusProviderUpdatedAttr(rName, "timeout", "60"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderExists(resourceName, &radius_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "60"),
					testAccCheckAciRadiusProviderIdEqual(&radius_provider_default, &radius_provider_updated),
				),
			},
			{
				Config: CreateAccRadiusProviderUpdatedAttr(rName, "timeout", "30"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderExists(resourceName, &radius_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "30"),
					testAccCheckAciRadiusProviderIdEqual(&radius_provider_default, &radius_provider_updated),
				),
			},

			{
				Config: CreateAccRadiusProviderConfig(rName),
			},
		},
	})
}

func TestAccAciRadiusProvider_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRadiusProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRadiusProviderConfig(rName),
			},

			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "auth_port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "auth_port", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "auth_port", "65536"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "auth_protocol", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "key", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "monitor_server", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "monitoring_password", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "monitoring_user", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "retries", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "retries", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "retries", "6"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "timeout", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "timeout", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, "timeout", "61"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccRadiusProviderUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccRadiusProviderConfig(rName),
			},
		},
	})
}

func TestAccAciRadiusProvider_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRadiusProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRadiusProviderConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciRadiusProviderExists(name string, radius_provider *models.RADIUSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Radius Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Radius Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		radius_providerFound := models.RADIUSProviderFromContainer(cont)
		if radius_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Radius Provider %s not found", rs.Primary.ID)
		}
		*radius_provider = *radius_providerFound
		return nil
	}
}

func testAccCheckAciRadiusProviderDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing radius_provider destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_radius_provider" {
			cont, err := client.Get(rs.Primary.ID)
			radius_provider := models.RADIUSProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Radius Provider %s Still exists", radius_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRadiusProviderIdEqual(m1, m2 *models.RADIUSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("radius_provider DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciRadiusProviderIdNotEqual(m1, m2 *models.RADIUSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("radius_provider DNs are equal")
		}
		return nil
	}
}

func CreateRadiusProviderWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing radius_provider creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_radius_provider" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccRadiusProviderConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing radius_provider creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccRadiusProviderConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing radius_provider creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccRadiusProviderConfig(rName string) string {
	fmt.Println("=== STEP  testing radius_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccRadiusProviderConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple radius_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccRadiusProviderConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing radius_provider creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_radius_provider"
		auth_port = "2"
		auth_protocol = "chap"
		key = ""
		monitor_server = "enabled"
		monitoring_password = ""
		monitoring_user = ""
		retries = "1"
		timeout = "1"
		
	}
	`, rName)

	return resource
}

func CreateAccRadiusProviderRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing radius_provider updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_radius_provider" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_radius_provider"
		auth_port = "2"
		auth_protocol = "chap"
		key = ""
		monitor_server = "enabled"
		monitoring_password = ""
		monitoring_user = ""
		retries = "1"
		timeout = "1"
		
	}
	`)

	return resource
}

func CreateAccRadiusProviderUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing radius_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
