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

func TestAccAciTACACSAccounting_Basic(t *testing.T) {
	var tacacs_accounting_default models.TACACSMonitoringDestinationGroup
	var tacacs_accounting_updated models.TACACSMonitoringDestinationGroup
	resourceName := "aci_tacacs_accounting.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSAccountingDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateTACACSAccountingWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSAccountingConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSAccountingExists(resourceName, &tacacs_accounting_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccTACACSAccountingConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSAccountingExists(resourceName, &tacacs_accounting_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_tacacs_accounting"),

					testAccCheckAciTACACSAccountingIdEqual(&tacacs_accounting_default, &tacacs_accounting_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccTACACSAccountingConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccTACACSAccountingRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccTACACSAccountingConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSAccountingExists(resourceName, &tacacs_accounting_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciTACACSAccountingIdNotEqual(&tacacs_accounting_default, &tacacs_accounting_updated),
				),
			},
		},
	})
}

func TestAccAciTACACSAccounting_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSAccountingDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSAccountingConfig(rName),
			},

			{
				Config:      CreateAccTACACSAccountingUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSAccountingUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSAccountingUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccTACACSAccountingUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccTACACSAccountingConfig(rName),
			},
		},
	})
}

func TestAccAciTACACSAccounting_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSAccountingDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSAccountingConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciTACACSAccountingExists(name string, tacacs_accounting *models.TACACSMonitoringDestinationGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACS Accounting %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACS Accounting dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tacacs_accountingFound := models.TACACSMonitoringDestinationGroupFromContainer(cont)
		if tacacs_accountingFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACS Accounting %s not found", rs.Primary.ID)
		}
		*tacacs_accounting = *tacacs_accountingFound
		return nil
	}
}

func testAccCheckAciTACACSAccountingDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing tacacs_accounting destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_accounting" {
			cont, err := client.Get(rs.Primary.ID)
			tacacs_accounting := models.TACACSMonitoringDestinationGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACS Accounting %s Still exists", tacacs_accounting.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSAccountingIdEqual(m1, m2 *models.TACACSMonitoringDestinationGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("tacacs_accounting DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTACACSAccountingIdNotEqual(m1, m2 *models.TACACSMonitoringDestinationGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("tacacs_accounting DNs are equal")
		}
		return nil
	}
}

func CreateTACACSAccountingWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_tacacs_accounting" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccTACACSAccountingConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing tacacs_accounting creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccTACACSAccountingConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing tacacs_accounting creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTACACSAccountingConfig(rName string) string {
	fmt.Println("=== STEP  testing tacacs_accounting creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTACACSAccountingConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple tacacs_accounting creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccTACACSAccountingConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_accounting"
		
	}
	`, rName)

	return resource
}

func CreateAccTACACSAccountingRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tacacs_accounting" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_accounting"
		
	}
	`)

	return resource
}

func CreateAccTACACSAccountingUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing tacacs_accounting attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
