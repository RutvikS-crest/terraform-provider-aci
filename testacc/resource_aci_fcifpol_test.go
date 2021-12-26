package acctest

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

func TestAccAciInterfaceFCPolicy_Basic(t *testing.T) {
	var interface_fc_policy_default models.InterfaceFCPolicy
	var interface_fc_policy_updated models.InterfaceFCPolicy
	resourceName := "aci_interface_fc_policy.testacc"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciInterfaceFCPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateInterfaceFCPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccInterfaceFCPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "automaxspeed", "32G"),
					resource.TestCheckResourceAttr(resourceName, "fill_pattern", "IDLE"),
					resource.TestCheckResourceAttr(resourceName, "port_mode", "f"),
					resource.TestCheckResourceAttr(resourceName, "rx_bb_credit", "6"),
					resource.TestCheckResourceAttr(resourceName, "speed", "auto"),
					resource.TestCheckResourceAttr(resourceName, "trunk_mode", "1"),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_interface_fc_policy"),
					resource.TestCheckResourceAttr(resourceName, "automaxspeed", "16G"),
					resource.TestCheckResourceAttr(resourceName, "fill_pattern", "ARBFF"),
					resource.TestCheckResourceAttr(resourceName, "port_mode", "np"),
					resource.TestCheckResourceAttr(resourceName, "rx_bb_credit", "17"), resource.TestCheckResourceAttr(resourceName, "rx_bb_credit", ""),
					resource.TestCheckResourceAttr(resourceName, "speed", "16G"),
					resource.TestCheckResourceAttr(resourceName, "trunk_mode", "auto"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccInterfaceFCPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)* failed validation`),
			},

			{
				Config: CreateAccInterfaceFCPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciInterfaceFCPolicyIdNotEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
		},
	})
}

func TestAccAciInterfaceFCPolicy_Update(t *testing.T) {
	var interface_fc_policy_default models.InterfaceFCPolicy
	var interface_fc_policy_updated models.InterfaceFCPolicy
	resourceName := "aci_interface_fc_policy.testacc"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciInterfaceFCPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccInterfaceFCPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_default),
				),
			},

			{
				Config: CreateAccInterfaceFCPolicyUpdatedAttr(rName, "automaxspeed", "2G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "automaxspeed", "2G"),
					testAccCheckAciInterfaceFCPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyUpdatedAttr(rName, "automaxspeed", "4G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "automaxspeed", "4G"),
					testAccCheckAciInterfaceFCPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyUpdatedAttr(rName, "automaxspeed", "8G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "automaxspeed", "8G"),
					testAccCheckAciInterfaceFCPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyUpdatedAttr(rName, "speed", "32G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "32G"),
					testAccCheckAciInterfaceFCPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyUpdatedAttr(rName, "speed", "4G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "4G"),
					testAccCheckAciInterfaceFCPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyUpdatedAttr(rName, "speed", "8G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "8G"),
					testAccCheckAciInterfaceFCPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyUpdatedAttr(rName, "speed", "unknown"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "unknown"),
					testAccCheckAciInterfaceFCPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyUpdatedAttr(rName, "trunkMode", "trunk-off"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "trunkMode", "trunk-off"),
					testAccCheckAciInterfaceFCPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyUpdatedAttr(rName, "trunkMode", "trunk-on"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "trunkMode", "trunk-on"),
					testAccCheckAciInterfaceFCPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyUpdatedAttr(rName, "trunkMode", "un-init"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "trunkMode", "un-init"),
					testAccCheckAciInterfaceFCPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFCPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciInterfaceFCPolicy_Negative(t *testing.T) {

	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciInterfaceFCPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccInterfaceFCPolicyConfig(rName),
			},

			{
				Config:      CreateAccInterfaceFCPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccInterfaceFCPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccInterfaceFCPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccInterfaceFCPolicyUpdatedAttr(rName, "automaxspeed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccInterfaceFCPolicyUpdatedAttr(rName, "fill_pattern", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccInterfaceFCPolicyUpdatedAttr(rName, "port_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccInterfaceFCPolicyUpdatedAttr(rName, "rx_bb_credit", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccInterfaceFCPolicyUpdatedAttr(rName, "speed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccInterfaceFCPolicyUpdatedAttr(rName, "trunk_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},

			{
				Config:      CreateAccInterfaceFCPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccInterfaceFCPolicyConfig(rName),
			},
		},
	})
}

func testAccCheckAciInterfaceFCPolicyExists(name string, interface_fc_policy *models.InterfaceFCPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface FC Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface FC Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		interface_fc_policyFound := models.InterfaceFCPolicyFromContainer(cont)
		if interface_fc_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface FC Policy %s not found", rs.Primary.ID)
		}
		*interface_fc_policy = *interface_fc_policyFound
		return nil
	}
}

func testAccCheckAciInterfaceFCPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing interface_fc_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_interface_fc_policy" {
			cont, err := client.Get(rs.Primary.ID)
			interface_fc_policy := models.InterfaceFCPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface FC Policy %s Still exists", interface_fc_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciInterfaceFCPolicyIdEqual(m1, m2 *models.InterfaceFCPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("interface_fc_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciInterfaceFCPolicyIdNotEqual(m1, m2 *models.InterfaceFCPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("interface_fc_policy DNs are equal")
		}
		return nil
	}
}

func CreateInterfaceFCPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing interface_fc_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_interface_fc_policy" "test" {
	
	#	name  = "%s"
		description = "created while acceptance testing"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccInterfaceFCPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing interface_fc_policy creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccInterfaceFCPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing interface_fc_policy creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccInterfaceFCPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing interface_fc_policy creation with Invalid Name: %s", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccInterfaceFCPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing interface_fc_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_interface_fc_policy"
		automaxspeed = "16G"fill_pattern = "ARBFF"port_mode = "np"rx_bb_credit = "17"rx_bb_credit = ""speed = "16G"trunk_mode = "auto"
	}
	`, rName)

	return resource
}

func CreateAccInterfaceFCPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing interface_fc_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_interface_fc_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_interface_fc_policy"
		automaxspeed = "16G"fill_pattern = "ARBFF"port_mode = "np"rx_bb_credit = "17"rx_bb_credit = ""speed = "16G"trunk_mode = "auto"
	}
	`)

	return resource
}

func CreateAccInterfaceFCPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing interface_fc_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
