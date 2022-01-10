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

func TestAccAciPortTracking_Basic(t *testing.T) {
	var port_tracking_default models.PortTracking
	var port_tracking_updated models.PortTracking
	resourceName := "aci_port_tracking.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPortTrackingDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreatePortTrackingWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPortTrackingConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "off"),
					resource.TestCheckResourceAttr(resourceName, "delay", "120"),
					resource.TestCheckResourceAttr(resourceName, "include_apic_ports", "no"),
					resource.TestCheckResourceAttr(resourceName, "minlinks", "0"),
				),
			},
			{
				Config: CreateAccPortTrackingConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_port_tracking"),

					resource.TestCheckResourceAttr(resourceName, "admin_st", "on"),
					resource.TestCheckResourceAttr(resourceName, "delay", "2"),

					resource.TestCheckResourceAttr(resourceName, "include_apic_ports", "yes"),
					resource.TestCheckResourceAttr(resourceName, "minlinks", "1"),

					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccPortTrackingConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccPortTrackingRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccPortTrackingConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciPortTrackingIdNotEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
		},
	})
}

func TestAccAciPortTracking_Update(t *testing.T) {
	var port_tracking_default models.PortTracking
	var port_tracking_updated models.PortTracking
	resourceName := "aci_port_tracking.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPortTrackingDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPortTrackingConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_default),
				),
			},
			{
				Config: CreateAccPortTrackingUpdatedAttr(rName, "delay", "300"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "delay", "300"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				Config: CreateAccPortTrackingUpdatedAttr(rName, "delay", "149"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "delay", "149"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				Config: CreateAccPortTrackingUpdatedAttr(rName, "minlinks", "48"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "minlinks", "48"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},
			{
				Config: CreateAccPortTrackingUpdatedAttr(rName, "minlinks", "24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortTrackingExists(resourceName, &port_tracking_updated),
					resource.TestCheckResourceAttr(resourceName, "minlinks", "24"),
					testAccCheckAciPortTrackingIdEqual(&port_tracking_default, &port_tracking_updated),
				),
			},

			{
				Config: CreateAccPortTrackingConfig(rName),
			},
		},
	})
}

func TestAccAciPortTracking_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPortTrackingDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPortTrackingConfig(rName),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "delay", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "delay", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "delay", "301"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "include_apic_ports", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "minlinks", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "minlinks", "-1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, "minlinks", "49"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccPortTrackingUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccPortTrackingConfig(rName),
			},
		},
	})
}

func TestAccAciPortTracking_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPortTrackingDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPortTrackingConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciPortTrackingExists(name string, port_tracking *models.PortTracking) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Port Tracking %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Port Tracking dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		port_trackingFound := models.PortTrackingFromContainer(cont)
		if port_trackingFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Port Tracking %s not found", rs.Primary.ID)
		}
		*port_tracking = *port_trackingFound
		return nil
	}
}

func testAccCheckAciPortTrackingDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing port_tracking destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_port_tracking" {
			cont, err := client.Get(rs.Primary.ID)
			port_tracking := models.PortTrackingFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Port Tracking %s Still exists", port_tracking.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciPortTrackingIdEqual(m1, m2 *models.PortTracking) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("port_tracking DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciPortTrackingIdNotEqual(m1, m2 *models.PortTracking) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("port_tracking DNs are equal")
		}
		return nil
	}
}

func CreatePortTrackingWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing port_tracking creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_port_tracking" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccPortTrackingConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing port_tracking creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccPortTrackingConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing port_tracking creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccPortTrackingConfig(rName string) string {
	fmt.Println("=== STEP  testing port_tracking creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccPortTrackingConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple port_tracking creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccPortTrackingConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing port_tracking creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_port_tracking"
		admin_st = "on"
		delay = "2"
		include_apic_ports = "yes"
		minlinks = "1"
		
	}
	`, rName)

	return resource
}

func CreateAccPortTrackingRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing port_tracking updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_port_tracking" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_port_tracking"
		admin_st = "on"
		delay = "2"
		include_apic_ports = "yes"
		minlinks = "1"
		
	}
	`)

	return resource
}

func CreateAccPortTrackingUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing port_tracking attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccPortTrackingUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing port_tracking attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
