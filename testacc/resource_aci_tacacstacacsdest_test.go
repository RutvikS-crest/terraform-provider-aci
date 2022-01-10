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

func TestAccAciTACACSAccountingDestination_Basic(t *testing.T) {
	var tacacs_accounting_destination_default models.TACACSDestination
	var tacacs_accounting_destination_updated models.TACACSDestination
	resourceName := "aci_tacacs_accounting_destination.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	host := makeTestVariable(acctest.RandString(5))
	hostUpdated := makeTestVariable(acctest.RandString(5))

	port := makeTestVariable(acctest.RandString(5))
	portUpdated := makeTestVariable(acctest.RandString(5))
	tacacsGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSAccountingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTACACSAccountingDestinationWithoutRequired(tacacsGroupName, host, port, "tacacs_monitoring_destination_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTACACSAccountingDestinationWithoutRequired(tacacsGroupName, host, port, "host"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateTACACSAccountingDestinationWithoutRequired(tacacsGroupName, host, port, "port"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSAccountingDestinationConfig(tacacsGroupName, host, port),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSAccountingDestinationExists(resourceName, &tacacs_accounting_destination_default),
					resource.TestCheckResourceAttr(resourceName, "tacacs_monitoring_destination_group_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", tacacsGroupName)),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "port", port),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "pap"),
					resource.TestCheckResourceAttr(resourceName, "key", ""),
				),
			},
			{
				Config: CreateAccTACACSAccountingDestinationConfigWithOptionalValues(tacacsGroupName, host, port),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSAccountingDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "tacacs_monitoring_destination_group_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", tacacsGroupName)),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "port", port),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_tacacs_accounting_destination"),

					resource.TestCheckResourceAttr(resourceName, "auth_protocol", "chap"),

					resource.TestCheckResourceAttr(resourceName, "key", ""),

					testAccCheckAciTACACSAccountingDestinationIdEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccTACACSAccountingDestinationRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSAccountingDestinationConfigWithRequiredParams(rNameUpdated, host, port),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSAccountingDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "tacacs_monitoring_destination_group_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "port", port),
					testAccCheckAciTACACSAccountingDestinationIdNotEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
			{
				Config: CreateAccTACACSAccountingDestinationConfig(tacacsGroupName, host, port),
			},
			{
				Config: CreateAccTACACSAccountingDestinationConfigWithRequiredParams(rName, hostUpdated, port),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSAccountingDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "tacacs_monitoring_destination_group_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "port", port),
					resource.TestCheckResourceAttr(resourceName, "host", hostUpdated),
					testAccCheckAciTACACSAccountingDestinationIdNotEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
			{
				Config: CreateAccTACACSAccountingDestinationConfigWithRequiredParams(rName, host, portUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSAccountingDestinationExists(resourceName, &tacacs_accounting_destination_updated),
					resource.TestCheckResourceAttr(resourceName, "tacacs_monitoring_destination_group_dn", fmt.Sprintf("uni/fabric/tacacsgroup-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "host", host),
					resource.TestCheckResourceAttr(resourceName, "port", portUpdated),
					testAccCheckAciTACACSAccountingDestinationIdNotEqual(&tacacs_accounting_destination_default, &tacacs_accounting_destination_updated),
				),
			},
		},
	})
}

func TestAccAciTACACSAccountingDestination_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	host := makeTestVariable(acctest.RandString(5))
	port := makeTestVariable(acctest.RandString(5))

	tacacsGroupName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSAccountingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSAccountingDestinationConfig(tacacsGroupName, host, port),
			},
			{
				Config:      CreateAccTACACSAccountingDestinationWithInValidParentDn(rName, host, port),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccTACACSAccountingDestinationUpdatedAttr(tacacsGroupName, host, port, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSAccountingDestinationUpdatedAttr(tacacsGroupName, host, port, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTACACSAccountingDestinationUpdatedAttr(tacacsGroupName, host, port, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccTACACSAccountingDestinationUpdatedAttr(tacacsGroupName, host, port, "auth_protocol", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccTACACSAccountingDestinationUpdatedAttr(tacacsGroupName, host, port, "key", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccTACACSAccountingDestinationUpdatedAttr(tacacsGroupName, host, port, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccTACACSAccountingDestinationConfig(tacacsGroupName, host, port),
			},
		},
	})
}

func TestAccAciTACACSAccountingDestination_MultipleCreateDelete(t *testing.T) {
	host := makeTestVariable(acctest.RandString(5))

	port := makeTestVariable(acctest.RandString(5))

	tacacsGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSAccountingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTACACSAccountingDestinationConfigMultiple(tacacsGroupName, host, port),
			},
		},
	})
}

func testAccCheckAciTACACSAccountingDestinationExists(name string, tacacs_accounting_destination *models.TACACSDestination) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACS Accounting Destination %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACS Accounting Destination dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tacacs_accounting_destinationFound := models.TACACSDestinationFromContainer(cont)
		if tacacs_accounting_destinationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACS Accounting Destination %s not found", rs.Primary.ID)
		}
		*tacacs_accounting_destination = *tacacs_accounting_destinationFound
		return nil
	}
}

func testAccCheckAciTACACSAccountingDestinationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing tacacs_accounting_destination destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_accounting_destination" {
			cont, err := client.Get(rs.Primary.ID)
			tacacs_accounting_destination := models.TACACSDestinationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACS Accounting Destination %s Still exists", tacacs_accounting_destination.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSAccountingDestinationIdEqual(m1, m2 *models.TACACSDestination) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("tacacs_accounting_destination DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTACACSAccountingDestinationIdNotEqual(m1, m2 *models.TACACSDestination) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("tacacs_accounting_destination DNs are equal")
		}
		return nil
	}
}

func CreateTACACSAccountingDestinationWithoutRequired(tacacsGroupName, host, port, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting_destination creation without ", attrName)
	rBlock := `
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tacacs_monitoring_destination_group_dn":
		rBlock += `
	resource "aci_tacacs_accounting_destination" "test" {
	#	tacacs_accounting_dn   = aci_tacacs_accounting.test.id
		host  = "%s"	port  = "%s"
	}
		`
	case "host":
		rBlock += `
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn   = aci_tacacs_accounting.test.id
	#	host  = "%s"
		port  = "%s"
	}
		`
	case "port":
		rBlock += `
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn   = aci_tacacs_accounting.test.id
		host  = "%s"
	#	port  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, tacacsGroupName, host, port)
}

func CreateAccTACACSAccountingDestinationConfigWithRequiredParams(tacacsGroupName, host, port string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn   = aci_tacacs_accounting.test.id
		host  = "%s"
		port  = "%s"
	}
	`, tacacsGroupName, host, port)
	return resource
}

func CreateAccTACACSAccountingDestinationConfig(tacacsGroupName, host, port string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn   = aci_tacacs_accounting.test.id
		host  = "%s"
		port  = "%s"
	}
	`, tacacsGroupName, host, port)
	return resource
}

func CreateAccTACACSAccountingDestinationConfigMultiple(tacacsGroupName, host, port string) string {
	fmt.Println("=== STEP  testing multiple tacacs_accounting_destination creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn   = aci_tacacs_accounting.test.id
		host  = "%s_${count.index}"
		port  = "%s_${count.index}"
		count = 5
	}
	`, tacacsGroupName, host, port)
	return resource
}

func CreateAccTACACSAccountingDestinationWithInValidParentDn(rName, host, port string) string {
	fmt.Println("=== STEP  Negative Case: testing tacacs_accounting_destination creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = aci_tenant.test.id
		host  = "%s"
		port  = "%s"	
	}
	`, rName, host, port)
	return resource
}

func CreateAccTACACSAccountingDestinationConfigWithOptionalValues(tacacsGroupName, host, port string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting_destination creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn  = "${aci_tacacs_accounting.test.id}"
		host  = "%s"
		port  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_accounting_destination"
		auth_protocol = "chap"
		key = ""
		
	}
	`, tacacsGroupName, host, port)

	return resource
}

func CreateAccTACACSAccountingDestinationRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting_destination updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tacacs_accounting_destination" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tacacs_accounting_destination"
		auth_protocol = "chap"
		key = ""
		
	}
	`)

	return resource
}

func CreateAccTACACSAccountingDestinationUpdatedAttr(tacacsGroupName, host, port, attribute, value string) string {
	fmt.Printf("=== STEP  testing tacacs_accounting_destination attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn   = aci_tacacs_accounting.test.id
		host  = "%s"
		port  = "%s"
		%s = "%s"
	}
	`, tacacsGroupName, host, port, attribute, value)
	return resource
}

func CreateAccTACACSAccountingDestinationUpdatedAttrList(tacacsGroupName, host, port, attribute, value string) string {
	fmt.Printf("=== STEP  testing tacacs_accounting_destination attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_accounting_dn   = aci_tacacs_accounting.test.id
		host  = "%s"
		port  = "%s"
		%s = %s
	}
	`, tacacsGroupName, host, port, attribute, value)
	return resource
}
