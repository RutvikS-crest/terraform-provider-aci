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

func TestAccAciMiscablingProtocolInstancePolicy_Basic(t *testing.T) {
	var mcp_instance_policy_default models.MiscablingProtocolInstancePolicy
	var mcp_instance_policy_updated models.MiscablingProtocolInstancePolicy
	resourceName := "aci_mcp_instance_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMiscablingProtocolInstancePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateMiscablingProtocolInstancePolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMiscablingProtocolInstancePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "init_delay_time", "180"),
					resource.TestCheckResourceAttr(resourceName, "key", ""),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "3"),
					resource.TestCheckResourceAttr(resourceName, "loop_protect_act.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "loop_protect_act.0", "port-disable"),
					resource.TestCheckResourceAttr(resourceName, "tx_freq", "2"),
					resource.TestCheckResourceAttr(resourceName, "tx_freq_msec", "0"),
				),
			},
			{
				Config: CreateAccMiscablingProtocolInstancePolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_mcp_instance_policy"),

					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "pdu-per-vlan"),
					resource.TestCheckResourceAttr(resourceName, "init_delay_time", "1"),

					resource.TestCheckResourceAttr(resourceName, "key", ""),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "2"),
					resource.TestCheckResourceAttr(resourceName, "loop_protect_act", "port-disable"),
					resource.TestCheckResourceAttr(resourceName, "tx_freq", "1"),
					resource.TestCheckResourceAttr(resourceName, "tx_freq_msec", "1"),

					testAccCheckAciMiscablingProtocolInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccMiscablingProtocolInstancePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccMiscablingProtocolInstancePolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciMiscablingProtocolInstancePolicyIdNotEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
		},
	})
}

func TestAccAciMiscablingProtocolInstancePolicy_Update(t *testing.T) {
	var mcp_instance_policy_default models.MiscablingProtocolInstancePolicy
	var mcp_instance_policy_updated models.MiscablingProtocolInstancePolicy
	resourceName := "aci_mcp_instance_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMiscablingProtocolInstancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMiscablingProtocolInstancePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_default),
				),
			},

			{

				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"pdu-per-vlan"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "pdu-per-vlan"),
				),
			},
			{

				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"pdu-per-vlan", "stateful-ha"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "pdu-per-vlan"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "stateful-ha"),
				),
			},
			{

				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"stateful-ha"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "stateful-ha"),
				),
			},
			{
				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"stateful-ha", "pdu-per-vlan"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "stateful-ha"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "pdu-per-vlan"),
				),
			}, {
				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "init_delay_time", "1800"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "init_delay_time", "1800"),
					testAccCheckAciMiscablingProtocolInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "init_delay_time", "900"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "init_delay_time", "900"),
					testAccCheckAciMiscablingProtocolInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "loop_detect_mult", "255"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "255"),
					testAccCheckAciMiscablingProtocolInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "loop_detect_mult", "127"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_detect_mult", "127"),
					testAccCheckAciMiscablingProtocolInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},

			{
				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttrList(rName, "loop_protect_act", StringListtoString([]string{"port-disable"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "loop_protect_act.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "loop_protect_act.0", "port-disable"),
				),
			}, {
				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "tx_freq", "300"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tx_freq", "300"),
					testAccCheckAciMiscablingProtocolInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "tx_freq", "150"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tx_freq", "150"),
					testAccCheckAciMiscablingProtocolInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "tx_freq_msec", "999"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tx_freq_msec", "999"),
					testAccCheckAciMiscablingProtocolInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},
			{
				Config: CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "tx_freq_msec", "499"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInstancePolicyExists(resourceName, &mcp_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tx_freq_msec", "499"),
					testAccCheckAciMiscablingProtocolInstancePolicyIdEqual(&mcp_instance_policy_default, &mcp_instance_policy_updated),
				),
			},

			{
				Config: CreateAccMiscablingProtocolInstancePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciMiscablingProtocolInstancePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMiscablingProtocolInstancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMiscablingProtocolInstancePolicyConfig(rName),
			},

			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"pdu-per-vlan", "pdu-per-vlan"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			// TODO: add unspecified case for "ctrl" if applicable

			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "init_delay_time", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "init_delay_time", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "init_delay_time", "1801"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "key", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "loop_detect_mult", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "loop_detect_mult", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "loop_detect_mult", "256"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttrList(rName, "loop_protect_act", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttrList(rName, "loop_protect_act", StringListtoString([]string{"port-disable", "port-disable"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			// TODO: add unspecified case for "loop_protect_act" if applicable

			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "tx_freq", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "tx_freq", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "tx_freq", "301"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "tx_freq_msec", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "tx_freq_msec", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, "tx_freq_msec", "1000"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccMiscablingProtocolInstancePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciMiscablingProtocolInstancePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMiscablingProtocolInstancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMiscablingProtocolInstancePolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciMiscablingProtocolInstancePolicyExists(name string, mcp_instance_policy *models.MiscablingProtocolInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("MCP Instance Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MCP Instance Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		mcp_instance_policyFound := models.MiscablingProtocolInstancePolicyFromContainer(cont)
		if mcp_instance_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("MCP Instance Policy %s not found", rs.Primary.ID)
		}
		*mcp_instance_policy = *mcp_instance_policyFound
		return nil
	}
}

func testAccCheckAciMiscablingProtocolInstancePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing mcp_instance_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_mcp_instance_policy" {
			cont, err := client.Get(rs.Primary.ID)
			mcp_instance_policy := models.MiscablingProtocolInstancePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("MCP Instance Policy %s Still exists", mcp_instance_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMiscablingProtocolInstancePolicyIdEqual(m1, m2 *models.MiscablingProtocolInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("mcp_instance_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciMiscablingProtocolInstancePolicyIdNotEqual(m1, m2 *models.MiscablingProtocolInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("mcp_instance_policy DNs are equal")
		}
		return nil
	}
}

func CreateMiscablingProtocolInstancePolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing mcp_instance_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_mcp_instance_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccMiscablingProtocolInstancePolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccMiscablingProtocolInstancePolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccMiscablingProtocolInstancePolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccMiscablingProtocolInstancePolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple mcp_instance_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccMiscablingProtocolInstancePolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing mcp_instance_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_mcp_instance_policy"
		admin_st = "enabled"
		ctrl = ["pdu-per-vlan"]
		init_delay_time = "1"
		key = ""
		loop_detect_mult = "2"
		loop_protect_act = "port-disable"
		tx_freq = "1"
		tx_freq_msec = "1"
		
	}
	`, rName)

	return resource
}

func CreateAccMiscablingProtocolInstancePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing mcp_instance_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_mcp_instance_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_mcp_instance_policy"
		admin_st = "enabled"
		ctrl = ["pdu-per-vlan"]
		init_delay_time = "1"
		key = ""
		loop_detect_mult = "2"
		loop_protect_act = "port-disable"
		tx_freq = "1"
		tx_freq_msec = "1"
		
	}
	`)

	return resource
}

func CreateAccMiscablingProtocolInstancePolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing mcp_instance_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccMiscablingProtocolInstancePolicyUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing mcp_instance_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
