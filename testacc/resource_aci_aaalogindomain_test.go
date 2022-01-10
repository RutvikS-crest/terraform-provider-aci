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

func TestAccAciLoginDomain_Basic(t *testing.T) {
	var login_domain_default models.LoginDomain
	var login_domain_updated models.LoginDomain
	resourceName := "aci_login_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoginDomainDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLoginDomainWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLoginDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainExists(resourceName, &login_domain_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccLoginDomainConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainExists(resourceName, &login_domain_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_login_domain"),

					testAccCheckAciLoginDomainIdEqual(&login_domain_default, &login_domain_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLoginDomainConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccLoginDomainRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccLoginDomainConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainExists(resourceName, &login_domain_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLoginDomainIdNotEqual(&login_domain_default, &login_domain_updated),
				),
			},
		},
	})
}

func TestAccAciLoginDomain_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoginDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLoginDomainConfig(rName),
			},

			{
				Config:      CreateAccLoginDomainUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLoginDomainUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLoginDomainUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLoginDomainUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLoginDomainConfig(rName),
			},
		},
	})
}

func TestAccAciLoginDomain_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoginDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLoginDomainConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciLoginDomainExists(name string, login_domain *models.LoginDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Login Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Login Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		login_domainFound := models.LoginDomainFromContainer(cont)
		if login_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Login Domain %s not found", rs.Primary.ID)
		}
		*login_domain = *login_domainFound
		return nil
	}
}

func testAccCheckAciLoginDomainDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing login_domain destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_login_domain" {
			cont, err := client.Get(rs.Primary.ID)
			login_domain := models.LoginDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Login Domain %s Still exists", login_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLoginDomainIdEqual(m1, m2 *models.LoginDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("login_domain DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLoginDomainIdNotEqual(m1, m2 *models.LoginDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("login_domain DNs are equal")
		}
		return nil
	}
}

func CreateLoginDomainWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing login_domain creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_login_domain" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLoginDomainConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing login_domain creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccLoginDomainConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing login_domain creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLoginDomainConfig(rName string) string {
	fmt.Println("=== STEP  testing login_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLoginDomainConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple login_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccLoginDomainConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing login_domain creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_login_domain"
		
	}
	`, rName)

	return resource
}

func CreateAccLoginDomainRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing login_domain updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_login_domain" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_login_domain"
		
	}
	`)

	return resource
}

func CreateAccLoginDomainUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing login_domain attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccLoginDomainUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing login_domain attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_login_domain" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
