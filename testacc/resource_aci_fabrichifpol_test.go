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

func TestAccAciFabricIfPol_Basic(t *testing.T) {
	var fabric_if_pol_default models.LinkLevelPolicy
	var fabric_if_pol_updated models.LinkLevelPolicy
	resourceName := "aci_fabric_if_pol.testacc"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIfPolDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFabricIfPolWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricIfPolConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_default),

					resource.TestCheckResourceAttr(resourceName, "name", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auto_neg", "on"),
					resource.TestCheckResourceAttr(resourceName, "dfe_delay_ms", "0"),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "inherit"),
					resource.TestCheckResourceAttr(resourceName, "link_debounce", "1"),
					resource.TestCheckResourceAttr(resourceName, "speed", "inherit"),
				),
			},
			{
				Config: CreateAccFabricIfPolConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fabric_if_pol"),
					resource.TestCheckResourceAttr(resourceName, "auto_neg", "off"),
					resource.TestCheckResourceAttr(resourceName, "dfe_delay_ms", "1"), resource.TestCheckResourceAttr(resourceName, "dfe_delay_ms", ""),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "auto-fec"),
					resource.TestCheckResourceAttr(resourceName, "link_debounce", "1"), resource.TestCheckResourceAttr(resourceName, "link_debounce", ""),
					resource.TestCheckResourceAttr(resourceName, "speed", "100G"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFabricIfPolConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)* failed validation`),
			},

			{
				Config: CreateAccFabricIfPolConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFabricIfPolIdNotEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
		},
	})
}

func TestAccAciFabricIfPol_Update(t *testing.T) {
	var fabric_if_pol_default models.LinkLevelPolicy
	var fabric_if_pol_updated models.LinkLevelPolicy
	resourceName := "aci_fabric_if_pol.testacc"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIfPolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricIfPolConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_default),
				),
			},

			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "fecMode", "cl74-fc-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fecMode", "cl74-fc-fec"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "fecMode", "cl91-rs-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fecMode", "cl91-rs-fec"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "fecMode", "cons16-rs-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fecMode", "cons16-rs-fec"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "fecMode", "disable-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fecMode", "disable-fec"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "fecMode", "ieee-rs-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fecMode", "ieee-rs-fec"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "fecMode", "kp-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fecMode", "kp-fec"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "speed", "100M"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "100M"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "speed", "10G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "10G"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "speed", "1G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "1G"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "speed", "200G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "200G"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "speed", "25G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "25G"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "speed", "400G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "400G"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "speed", "40G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "40G"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "speed", "50G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "50G"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolUpdatedAttr(rName, "speed", "unknown"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "unknown"),
					testAccCheckAciFabricIfPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolConfig(rName),
			},
		},
	})
}

func TestAccAciFabricIfPol_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIfPolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricIfPolConfig(rName),
			},

			{
				Config:      CreateAccFabricIfPolUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricIfPolUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricIfPolUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricIfPolUpdatedAttr(rName, "auto_neg", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccFabricIfPolUpdatedAttr(rName, "dfe_delay_ms", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccFabricIfPolUpdatedAttr(rName, "fec_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccFabricIfPolUpdatedAttr(rName, "link_debounce", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccFabricIfPolUpdatedAttr(rName, "speed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},

			{
				Config:      CreateAccFabricIfPolUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccFabricIfPolConfig(rName),
			},
		},
	})
}

func testAccCheckAciFabricIfPolExists(name string, fabric_if_pol *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric If Pol %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric If Pol dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_if_polFound := models.LinkLevelPolicyFromContainer(cont)
		if fabric_if_polFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric If Pol %s not found", rs.Primary.ID)
		}
		*fabric_if_pol = *fabric_if_polFound
		return nil
	}
}

func testAccCheckAciFabricIfPolDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing fabric_if_pol destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fabric_if_pol" {
			cont, err := client.Get(rs.Primary.ID)
			fabric_if_pol := models.LinkLevelPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric If Pol %s Still exists", fabric_if_pol.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFabricIfPolIdEqual(m1, m2 *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fabric_if_pol DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFabricIfPolIdNotEqual(m1, m2 *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("fabric_if_pol DNs are equal")
		}
		return nil
	}
}

func CreateFabricIfPolWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_if_pol creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_fabric_if_pol" "test" {
	
	#	name  = "%s"
		description = "created while acceptance testing"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFabricIfPolConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFabricIfPolConfig(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFabricIfPolConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFabricIfPolConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_if_pol creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_if_pol"
		auto_neg = "off"dfe_delay_ms = "1"dfe_delay_ms = ""fec_mode = "auto-fec"link_debounce = "1"link_debounce = ""speed = "100G"
	}
	`, rName)

	return resource
}

func CreateAccFabricIfPolRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing fabric_if_pol creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_fabric_if_pol" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_if_pol"
		auto_neg = "off"dfe_delay_ms = "1"dfe_delay_ms = ""fec_mode = "auto-fec"link_debounce = "1"link_debounce = ""speed = "100G"
	}
	`)

	return resource
}

func CreateAccFabricIfPolUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing fabric_if_pol attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
