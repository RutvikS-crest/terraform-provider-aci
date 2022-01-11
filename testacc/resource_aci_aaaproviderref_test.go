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

func TestAccAciLoginDomainProvider_Basic(t *testing.T) {
	var login_domain_provider_default models.ProviderGroupMember
	var login_domain_provider_updated models.ProviderGroupMember
	resourceName := "aci_login_domain_provider.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	aaaDuoProviderGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoginDomainProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLoginDomainProviderWithoutRequired(aaaDuoProviderGroupName, rName, "duo_provider_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLoginDomainProviderWithoutRequired(aaaDuoProviderGroupName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLoginDomainProviderConfig(aaaDuoProviderGroupName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainProviderExists(resourceName, &login_domain_provider_default),
					//resource.TestCheckResourceAttr(resourceName, "duo_provider_group_dn", GetParentDn(login_domain_provider_default.DistinguishedName, fmt.Sprintf("/providerref-%s", name))),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "order", ""),
				),
			},
			{
				Config: CreateAccLoginDomainProviderConfigWithOptionalValues(aaaDuoProviderGroupName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainProviderExists(resourceName, &login_domain_provider_updated),
					//resource.TestCheckResourceAttr(resourceName, "duo_provider_group_dn", GetParentDn(login_domain_provider_updated.DistinguishedName, fmt.Sprintf("/providerref-%s", name))),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_login_domain_provider"),
					resource.TestCheckResourceAttr(resourceName, "order", "1"),

					testAccCheckAciLoginDomainProviderIdEqual(&login_domain_provider_default, &login_domain_provider_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLoginDomainProviderConfigUpdatedName(aaaDuoProviderGroupName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccLoginDomainProviderRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLoginDomainProviderConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainProviderExists(resourceName, &login_domain_provider_updated),
					//resource.TestCheckResourceAttr(resourceName, "duo_provider_group_dn", GetParentDn(login_domain_provider_updated.DistinguishedName, fmt.Sprintf("/providerref-%s", name))),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciLoginDomainProviderIdNotEqual(&login_domain_provider_default, &login_domain_provider_updated),
				),
			},
			{
				Config: CreateAccLoginDomainProviderConfig(aaaDuoProviderGroupName, rName),
			},
			{
				Config: CreateAccLoginDomainProviderConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainProviderExists(resourceName, &login_domain_provider_updated),
					//resource.TestCheckResourceAttr(resourceName, "duo_provider_group_dn", GetParentDn(login_domain_provider_updated.DistinguishedName, fmt.Sprintf("/providerref-%s", name))),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLoginDomainProviderIdNotEqual(&login_domain_provider_default, &login_domain_provider_updated),
				),
			},
		},
	})
}

func TestAccAciLoginDomainProvider_Update(t *testing.T) {
	var login_domain_provider_default models.ProviderGroupMember
	var login_domain_provider_updated models.ProviderGroupMember
	resourceName := "aci_login_domain_provider.test"
	rName := makeTestVariable(acctest.RandString(5))

	aaaDuoProviderGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoginDomainProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLoginDomainProviderConfig(aaaDuoProviderGroupName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainProviderExists(resourceName, &login_domain_provider_default),
				),
			},
			{
				Config: CreateAccLoginDomainProviderUpdatedAttr(aaaDuoProviderGroupName, rName, "order", "16"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainProviderExists(resourceName, &login_domain_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "order", "16"),
					testAccCheckAciLoginDomainProviderIdEqual(&login_domain_provider_default, &login_domain_provider_updated),
				),
			},
			{
				Config: CreateAccLoginDomainProviderUpdatedAttr(aaaDuoProviderGroupName, rName, "order", "8"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainProviderExists(resourceName, &login_domain_provider_updated),
					resource.TestCheckResourceAttr(resourceName, "order", "8"),
					testAccCheckAciLoginDomainProviderIdEqual(&login_domain_provider_default, &login_domain_provider_updated),
				),
			},

			{
				Config: CreateAccLoginDomainProviderConfig(aaaDuoProviderGroupName, rName),
			},
		},
	})
}

func TestAccAciLoginDomainProvider_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaDuoProviderGroupName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoginDomainProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLoginDomainProviderConfig(aaaDuoProviderGroupName, rName),
			},
			{
				Config:      CreateAccLoginDomainProviderWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLoginDomainProviderUpdatedAttr(aaaDuoProviderGroupName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLoginDomainProviderUpdatedAttr(aaaDuoProviderGroupName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLoginDomainProviderUpdatedAttr(aaaDuoProviderGroupName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLoginDomainProviderUpdatedAttr(aaaDuoProviderGroupName, rName, "order", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLoginDomainProviderUpdatedAttr(aaaDuoProviderGroupName, rName, "order", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccLoginDomainProviderUpdatedAttr(aaaDuoProviderGroupName, rName, "order", "17"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccLoginDomainProviderUpdatedAttr(aaaDuoProviderGroupName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLoginDomainProviderConfig(aaaDuoProviderGroupName, rName),
			},
		},
	})
}

func TestAccAciLoginDomainProvider_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaDuoProviderGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoginDomainProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLoginDomainProviderConfigMultiple(aaaDuoProviderGroupName, rName),
			},
		},
	})
}

func testAccCheckAciLoginDomainProviderExists(name string, login_domain_provider *models.ProviderGroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Login Domain Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Login Domain Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		login_domain_providerFound := models.ProviderGroupMemberFromContainer(cont)
		if login_domain_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Login Domain Provider %s not found", rs.Primary.ID)
		}
		*login_domain_provider = *login_domain_providerFound
		return nil
	}
}

func testAccCheckAciLoginDomainProviderDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing login_domain_provider destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_login_domain_provider" {
			cont, err := client.Get(rs.Primary.ID)
			login_domain_provider := models.ProviderGroupMemberFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Login Domain Provider %s Still exists", login_domain_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLoginDomainProviderIdEqual(m1, m2 *models.ProviderGroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("login_domain_provider DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLoginDomainProviderIdNotEqual(m1, m2 *models.ProviderGroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("login_domain_provider DNs are equal")
		}
		return nil
	}
}

func CreateLoginDomainProviderWithoutRequired(aaaDuoProviderGroupName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing login_domain_provider creation without ", attrName)
	rBlock := `
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "duo_provider_group_dn":
		rBlock += `
	resource "aci_login_domain_provider" "test" {
	#	duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaDuoProviderGroupName, rName)
}

func CreateAccLoginDomainProviderConfigWithRequiredParams(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing login_domain_provider creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}
func CreateAccLoginDomainProviderConfigUpdatedName(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing login_domain_provider creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}

func CreateAccLoginDomainProviderConfig(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing login_domain_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}

func CreateAccLoginDomainProviderConfigMultiple(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing multiple login_domain_provider creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}

func CreateAccLoginDomainProviderWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing login_domain_provider creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccLoginDomainProviderConfigWithOptionalValues(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  Basic: testing login_domain_provider creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = "${aci_duo_provider_group.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_login_domain_provider"
		order = "1"
		
	}
	`, aaaDuoProviderGroupName, rName)

	return resource
}

func CreateAccLoginDomainProviderRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing login_domain_provider updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_login_domain_provider" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_login_domain_provider"
		order = "1"
		
	}
	`)

	return resource
}

func CreateAccLoginDomainProviderUpdatedAttr(aaaDuoProviderGroupName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing login_domain_provider attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s"
		%s = "%s"
	}
	`, aaaDuoProviderGroupName, rName, attribute, value)
	return resource
}
