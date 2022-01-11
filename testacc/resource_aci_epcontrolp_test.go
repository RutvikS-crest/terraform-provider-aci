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

func TestAccAciEndpointControlPolicy_Basic(t *testing.T) {
	var endpoint_controls_default models.EndpointControlPolicy
	var endpoint_controls_updated models.EndpointControlPolicy
	resourceName := "aci_endpoint_controls.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEndpointControlPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateEndpointControlPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndpointControlPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlPolicyExists(resourceName, &endpoint_controls_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "1800"),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_intvl", "60"),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_mult", "4"),
				),
			},
			{
				Config: CreateAccEndpointControlPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlPolicyExists(resourceName, &endpoint_controls_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_endpoint_controls"),

					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "301"),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_intvl", "31"),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_mult", "3"),

					testAccCheckAciEndpointControlPolicyIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccEndpointControlPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccEndpointControlPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccEndpointControlPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlPolicyExists(resourceName, &endpoint_controls_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciEndpointControlPolicyIdNotEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
		},
	})
}

func TestAccAciEndpointControlPolicy_Update(t *testing.T) {
	var endpoint_controls_default models.EndpointControlPolicy
	var endpoint_controls_updated models.EndpointControlPolicy
	resourceName := "aci_endpoint_controls.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEndpointControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointControlPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlPolicyExists(resourceName, &endpoint_controls_default),
				),
			},
			{
				Config: CreateAccEndpointControlPolicyUpdatedAttr(rName, "hold_intvl", "3600"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlPolicyExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "3600"),
					testAccCheckAciEndpointControlPolicyIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlPolicyUpdatedAttr(rName, "hold_intvl", "1650"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlPolicyExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "1650"),
					testAccCheckAciEndpointControlPolicyIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlPolicyUpdatedAttr(rName, "rogue_ep_detect_intvl", "3600"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlPolicyExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_intvl", "3600"),
					testAccCheckAciEndpointControlPolicyIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlPolicyUpdatedAttr(rName, "rogue_ep_detect_intvl", "1785"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlPolicyExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_intvl", "1785"),
					testAccCheckAciEndpointControlPolicyIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlPolicyUpdatedAttr(rName, "rogue_ep_detect_mult", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlPolicyExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_mult", "65535"),
					testAccCheckAciEndpointControlPolicyIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},
			{
				Config: CreateAccEndpointControlPolicyUpdatedAttr(rName, "rogue_ep_detect_mult", "32766"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointControlPolicyExists(resourceName, &endpoint_controls_updated),
					resource.TestCheckResourceAttr(resourceName, "rogue_ep_detect_mult", "32766"),
					testAccCheckAciEndpointControlPolicyIdEqual(&endpoint_controls_default, &endpoint_controls_updated),
				),
			},

			{
				Config: CreateAccEndpointControlPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciEndpointControlPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEndpointControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointControlPolicyConfig(rName),
			},

			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "hold_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "hold_intvl", "299"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "hold_intvl", "3601"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "rogue_ep_detect_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "rogue_ep_detect_intvl", "29"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "rogue_ep_detect_intvl", "3601"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "rogue_ep_detect_mult", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "rogue_ep_detect_mult", "1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, "rogue_ep_detect_mult", "65536"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccEndpointControlPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEndpointControlPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciEndpointControlPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEndpointControlPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointControlPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciEndpointControlPolicyExists(name string, endpoint_controls *models.EndpointControlPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Controls %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Controls dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_controlsFound := models.EndpointControlPolicyFromContainer(cont)
		if endpoint_controlsFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Controls %s not found", rs.Primary.ID)
		}
		*endpoint_controls = *endpoint_controlsFound
		return nil
	}
}

func testAccCheckAciEndpointControlPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing endpoint_controls destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_endpoint_controls" {
			cont, err := client.Get(rs.Primary.ID)
			endpoint_controls := models.EndpointControlPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Endpoint Controls %s Still exists", endpoint_controls.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciEndpointControlPolicyIdEqual(m1, m2 *models.EndpointControlPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("endpoint_controls DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciEndpointControlPolicyIdNotEqual(m1, m2 *models.EndpointControlPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("endpoint_controls DNs are equal")
		}
		return nil
	}
}

func CreateEndpointControlPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_controls creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_endpoint_controls" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccEndpointControlPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing endpoint_controls creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccEndpointControlPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing endpoint_controls creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccEndpointControlPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing endpoint_controls creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccEndpointControlPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple endpoint_controls creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccEndpointControlPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_controls creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_endpoint_controls"
		admin_st = "enabled"
		hold_intvl = "301"
		rogue_ep_detect_intvl = "31"
		rogue_ep_detect_mult = "3"
		
	}
	`, rName)

	return resource
}

func CreateAccEndpointControlPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing endpoint_controls updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_endpoint_controls" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_endpoint_controls"
		admin_st = "enabled"
		hold_intvl = "301"
		rogue_ep_detect_intvl = "31"
		rogue_ep_detect_mult = "3"
		
	}
	`)

	return resource
}

func CreateAccEndpointControlPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing endpoint_controls attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
