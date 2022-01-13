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

func TestAccAciQOSInstancePolicy_Basic(t *testing.T) {
	var qos_instance_policy_default models.QOSInstancePolicy
	var qos_instance_policy_updated models.QOSInstancePolicy
	resourceName := "aci_qos_instance_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciQOSInstancePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateQOSInstancePolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccQOSInstancePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "etrap_age_timer", ""),
					resource.TestCheckResourceAttr(resourceName, "etrap_bw_thresh", ""),
					resource.TestCheckResourceAttr(resourceName, "etrap_byte_ct", ""),
					resource.TestCheckResourceAttr(resourceName, "etrap_st", "no"),
					resource.TestCheckResourceAttr(resourceName, "fabric_flush_interval", "500"),
					resource.TestCheckResourceAttr(resourceName, "fabric_flush_st", "no"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "none"),
					resource.TestCheckResourceAttr(resourceName, "uburst_spine_queues", "10"),
					resource.TestCheckResourceAttr(resourceName, "uburst_tor_queues", "10"),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_qos_instance_policy"),

					resource.TestCheckResourceAttr(resourceName, "etrap_st", "yes"),
					resource.TestCheckResourceAttr(resourceName, "fabric_flush_interval", "101"),

					resource.TestCheckResourceAttr(resourceName, "fabric_flush_st", "yes"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "dot1p-preserve"),
					resource.TestCheckResourceAttr(resourceName, "uburst_spine_queues", "1"),
					resource.TestCheckResourceAttr(resourceName, "uburst_tor_queues", "1"),

					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccQOSInstancePolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccQOSInstancePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccQOSInstancePolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciQOSInstancePolicyIdNotEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
		},
	})
}

func TestAccAciQOSInstancePolicy_Update(t *testing.T) {
	var qos_instance_policy_default models.QOSInstancePolicy
	var qos_instance_policy_updated models.QOSInstancePolicy
	resourceName := "aci_qos_instance_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciQOSInstancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccQOSInstancePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_default),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr(rName, "fabric_flush_interval", "1000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_flush_interval", "1000"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr(rName, "fabric_flush_interval", "450"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_flush_interval", "450"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},

			{

				Config: CreateAccQOSInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"dot1p-preserve"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "dot1p-preserve"),
				),
			},
			{

				Config: CreateAccQOSInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"dot1p-preserve", "none"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "dot1p-preserve"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "none"),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"none", "dot1p-preserve"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "none"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "dot1p-preserve"),
				),
			}, {
				Config: CreateAccQOSInstancePolicyUpdatedAttr(rName, "uburst_spine_queues", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "uburst_spine_queues", "100"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr(rName, "uburst_spine_queues", "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "uburst_spine_queues", "50"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr(rName, "uburst_tor_queues", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "uburst_tor_queues", "100"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},
			{
				Config: CreateAccQOSInstancePolicyUpdatedAttr(rName, "uburst_tor_queues", "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciQOSInstancePolicyExists(resourceName, &qos_instance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "uburst_tor_queues", "50"),
					testAccCheckAciQOSInstancePolicyIdEqual(&qos_instance_policy_default, &qos_instance_policy_updated),
				),
			},

			{
				Config: CreateAccQOSInstancePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciQOSInstancePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciQOSInstancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccQOSInstancePolicyConfig(rName),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "etrap_age_timer", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "etrap_bw_thresh", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "etrap_byte_ct", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "etrap_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "fabric_flush_interval", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "fabric_flush_interval", "99"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "fabric_flush_interval", "1001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "fabric_flush_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"dot1p-preserve", "dot1p-preserve"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			// TODO: add unspecified case for "ctrl" if applicable

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "uburst_spine_queues", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "uburst_spine_queues", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "uburst_spine_queues", "101"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "uburst_tor_queues", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "uburst_tor_queues", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, "uburst_tor_queues", "101"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccQOSInstancePolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccQOSInstancePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciQOSInstancePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciQOSInstancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccQOSInstancePolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciQOSInstancePolicyExists(name string, qos_instance_policy *models.QOSInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("QOS Instance Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No QOS Instance Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		qos_instance_policyFound := models.QOSInstancePolicyFromContainer(cont)
		if qos_instance_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("QOS Instance Policy %s not found", rs.Primary.ID)
		}
		*qos_instance_policy = *qos_instance_policyFound
		return nil
	}
}

func testAccCheckAciQOSInstancePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing qos_instance_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_qos_instance_policy" {
			cont, err := client.Get(rs.Primary.ID)
			qos_instance_policy := models.QOSInstancePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("QOS Instance Policy %s Still exists", qos_instance_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciQOSInstancePolicyIdEqual(m1, m2 *models.QOSInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("qos_instance_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciQOSInstancePolicyIdNotEqual(m1, m2 *models.QOSInstancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("qos_instance_policy DNs are equal")
		}
		return nil
	}
}

func CreateQOSInstancePolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing qos_instance_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_qos_instance_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccQOSInstancePolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing qos_instance_policy creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccQOSInstancePolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing qos_instance_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccQOSInstancePolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing qos_instance_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccQOSInstancePolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple qos_instance_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccQOSInstancePolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing qos_instance_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_qos_instance_policy"
		etrap_st = "yes"
		fabric_flush_interval = "101"
		fabric_flush_st = "yes"
		ctrl = ["dot1p-preserve"]
		uburst_spine_queues = "1"
		uburst_tor_queues = "1"
		
	}
	`, rName)

	return resource
}

func CreateAccQOSInstancePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing qos_instance_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_qos_instance_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_qos_instance_policy"
		etrap_st = "yes"
		fabric_flush_interval = "101"
		fabric_flush_st = "yes"
		ctrl = ["dot1p-preserve"]
		uburst_spine_queues = "1"
		uburst_tor_queues = "1"
		
	}
	`)

	return resource
}

func CreateAccQOSInstancePolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing qos_instance_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccQOSInstancePolicyUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing qos_instance_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
