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

func TestAccAciFabricIFPol_Basic(t *testing.T) {
	var fabric_if_pol_default models.LinkLevelPolicy
	var fabric_if_pol_updated models.LinkLevelPolicy
	resourceName := "aci_fabric_if_pol.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIFPolDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFabricIFPolWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricIFPolConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auto_neg", "on"),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "inherit"),
					resource.TestCheckResourceAttr(resourceName, "link_debounce", "100"),
					resource.TestCheckResourceAttr(resourceName, "speed", "inherit"),
				),
			},
			{
				// in this step all optional attribute expect realational attribute are given for the same resource and then compared
				Config: CreateAccFabricIFPolConfigWithOptionalValues(rName), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fabric_if_pol"),
					resource.TestCheckResourceAttr(resourceName, "auto_neg", "off"),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "auto-fec"),
					resource.TestCheckResourceAttr(resourceName, "link_debounce", "2"),
					resource.TestCheckResourceAttr(resourceName, "speed", "100G"),

					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFabricIFPolConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccFabricIFPolRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccFabricIFPolConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFabricIFPolIdNotEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
		},
	})
}

func TestAccAciFabricIFPol_Update(t *testing.T) {
	var fabric_if_pol_default models.LinkLevelPolicy
	var fabric_if_pol_updated models.LinkLevelPolicy
	resourceName := "aci_fabric_if_pol.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIFPolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricIFPolConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_default),
				),
			},

			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "fec_mode", "cl74-fc-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "cl74-fc-fec"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "fec_mode", "cl91-rs-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "cl91-rs-fec"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "fec_mode", "cons16-rs-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "cons16-rs-fec"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "fec_mode", "disable-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "disable-fec"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "fec_mode", "ieee-rs-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "ieee-rs-fec"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "fec_mode", "kp-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "kp-fec"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "speed", "100M"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "100M"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "speed", "10G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "10G"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "speed", "1G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "1G"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "speed", "200G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "200G"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "speed", "25G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "25G"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "speed", "400G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "400G"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "speed", "40G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "40G"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "speed", "50G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "50G"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolUpdatedAttr(rName, "speed", "unknown"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIFPolExists(resourceName, &fabric_if_pol_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "unknown"),
					testAccCheckAciFabricIFPolIdEqual(&fabric_if_pol_default, &fabric_if_pol_updated),
				),
			},
			{
				Config: CreateAccFabricIFPolConfig(rName),
			},
		},
	})
}

func TestAccAciFabricIFPol_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIFPolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricIFPolConfig(rName),
			},

			{
				Config:      CreateAccFabricIFPolUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricIFPolUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricIFPolUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricIFPolUpdatedAttr(rName, "auto_neg", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccFabricIFPolUpdatedAttr(rName, "fec_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccFabricIFPolUpdatedAttr(rName, "link_debounce", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFabricIFPolUpdatedAttr(rName, "speed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricIFPolUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFabricIFPolConfig(rName),
			},
		},
	})
}

func testAccCheckAciFabricIFPolExists(name string, fabric_if_pol *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric IF Pol %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric IF Pol dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_if_polFound := models.LinkLevelPolicyFromContainer(cont)
		if fabric_if_polFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric IF Pol %s not found", rs.Primary.ID)
		}
		*fabric_if_pol = *fabric_if_polFound
		return nil
	}
}

func testAccCheckAciFabricIFPolDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing fabric_if_pol destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fabric_if_pol" {
			cont, err := client.Get(rs.Primary.ID)
			fabric_if_pol := models.LinkLevelPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric IF Pol %s Still exists", fabric_if_pol.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFabricIFPolIdEqual(m1, m2 *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fabric_if_pol DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFabricIFPolIdNotEqual(m1, m2 *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("fabric_if_pol DNs are equal")
		}
		return nil
	}
}

func CreateFabricIFPolWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_if_pol creation without ", attrName)
	rBlock := ``
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

func CreateAccFabricIFPolConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFabricIFPolConfig(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFabricIFPolConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with Updated Name")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFabricIFPolConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_if_pol creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_if_pol"
		auto_neg = "off"
		fec_mode = "auto-fec"
		link_debounce = "2"
		speed = "100G"
	}
	`, rName)

	return resource
}

func CreateAccFabricIFPolRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing fabric_if_pol creation with optional parameters")
	resource := `
	resource "aci_fabric_if_pol" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_if_pol"
		auto_neg = "off"
		fec_mode = "auto-fec"
		link_debounce = "1"
		speed = "100G"
	}
	`

	return resource
}

func CreateAccFabricIFPolUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing fabric_if_pol attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
