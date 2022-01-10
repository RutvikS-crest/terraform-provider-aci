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

func TestAccAciEPLoopProtectionPolicy_Basic(t *testing.T) {
	var endpoint_loop_protection_default models.EPLoopProtectionPolicy
	var endpoint_loop_protection_updated models.EPLoopProtectionPolicy
	resourceName := "aci_endpoint_loop_protection.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPLoopProtectionPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateEPLoopProtectionPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEPLoopProtectionPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0", "port-disable"),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_intvl", "60"),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "4"),
				),
			},
			{
				Config: CreateAccEPLoopProtectionPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_endpoint_loop_protection"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0", "bd-learn-disable"),

					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_intvl", "31"),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "2"),

					testAccCheckAciEPLoopProtectionPolicyIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccEPLoopProtectionPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccEPLoopProtectionPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccEPLoopProtectionPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciEPLoopProtectionPolicyIdNotEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
		},
	})
}

func TestAccAciEPLoopProtectionPolicy_Update(t *testing.T) {
	var endpoint_loop_protection_default models.EPLoopProtectionPolicy
	var endpoint_loop_protection_updated models.EPLoopProtectionPolicy
	resourceName := "aci_endpoint_loop_protection.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPLoopProtectionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEPLoopProtectionPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_default),
				),
			},

			{

				Config: CreateAccEPLoopProtectionPolicyUpdatedAttrList(rName, "action", StringListtoString([]string{"bd-learn-disable"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "action.0", "bd-learn-disable"),
				),
			},
			{

				Config: CreateAccEPLoopProtectionPolicyUpdatedAttrList(rName, "action", StringListtoString([]string{"bd-learn-disable", "port-disable"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.0", "bd-learn-disable"),
					resource.TestCheckResourceAttr(resourceName, "action.1", "port-disable"),
				),
			},
			{
				Config: CreateAccEPLoopProtectionPolicyUpdatedAttrList(rName, "action", StringListtoString([]string{"port-disable", "bd-learn-disable"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "action.0", "port-disable"),
					resource.TestCheckResourceAttr(resourceName, "action.1", "bd-learn-disable"),
				),
			}, {
				Config: CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "loop_detect_intvl", "300"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_intvl", "300"),
					testAccCheckAciEPLoopProtectionPolicyIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
			{
				Config: CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "loop_detect_intvl", "135"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_intvl", "135"),
					testAccCheckAciEPLoopProtectionPolicyIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
			{
				Config: CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "loop_detect_mult", "255"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "255"),
					testAccCheckAciEPLoopProtectionPolicyIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},
			{
				Config: CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "loop_detect_mult", "127"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPLoopProtectionPolicyExists(resourceName, &endpoint_loop_protection_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "127"),
					testAccCheckAciEPLoopProtectionPolicyIdEqual(&endpoint_loop_protection_default, &endpoint_loop_protection_updated),
				),
			},

			{
				Config: CreateAccEPLoopProtectionPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciEPLoopProtectionPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPLoopProtectionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEPLoopProtectionPolicyConfig(rName),
			},

			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttrList(rName, "action", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttrList(rName, "action", StringListtoString([]string{"bd-learn-disable", "bd-learn-disable"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			// TODO: add unspecified case for "action" if applicable

			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "loop_detect_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "loop_detect_intvl", "29"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "loop_detect_intvl", "301"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "loop_detect_mult", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "loop_detect_mult", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, "loop_detect_mult", "256"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEPLoopProtectionPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciEPLoopProtectionPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPLoopProtectionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEPLoopProtectionPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciEPLoopProtectionPolicyExists(name string, endpoint_loop_protection *models.EPLoopProtectionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Loop Protection %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Loop Protection dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_loop_protectionFound := models.EPLoopProtectionPolicyFromContainer(cont)
		if endpoint_loop_protectionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Loop Protection %s not found", rs.Primary.ID)
		}
		*endpoint_loop_protection = *endpoint_loop_protectionFound
		return nil
	}
}

func testAccCheckAciEPLoopProtectionPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing endpoint_loop_protection destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_endpoint_loop_protection" {
			cont, err := client.Get(rs.Primary.ID)
			endpoint_loop_protection := models.EPLoopProtectionPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Endpoint Loop Protection %s Still exists", endpoint_loop_protection.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciEPLoopProtectionPolicyIdEqual(m1, m2 *models.EPLoopProtectionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("endpoint_loop_protection DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciEPLoopProtectionPolicyIdNotEqual(m1, m2 *models.EPLoopProtectionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("endpoint_loop_protection DNs are equal")
		}
		return nil
	}
}

func CreateEPLoopProtectionPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_loop_protection creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_endpoint_loop_protection" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccEPLoopProtectionPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing endpoint_loop_protection creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccEPLoopProtectionPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing endpoint_loop_protection creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccEPLoopProtectionPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing endpoint_loop_protection creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccEPLoopProtectionPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple endpoint_loop_protection creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccEPLoopProtectionPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_loop_protection creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_endpoint_loop_protection"
		action = ["bd-learn-disable"]
		admin_st = "enabled"
		loop_detect_intvl = "31"
		loop_detect_mult = "2"
		
	}
	`, rName)

	return resource
}

func CreateAccEPLoopProtectionPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing endpoint_loop_protection updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_endpoint_loop_protection" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_endpoint_loop_protection"
		action = ["bd-learn-disable"]
		admin_st = "enabled"
		loop_detect_intvl = "31"
		loop_detect_mult = "2"
		
	}
	`)

	return resource
}

func CreateAccEPLoopProtectionPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing endpoint_loop_protection attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccEPLoopProtectionPolicyUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing endpoint_loop_protection attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
