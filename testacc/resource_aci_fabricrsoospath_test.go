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

func TestAccAciOutofServiceFabricPath_Basic(t *testing.T) {
	var interface_blacklist_default models.OutofServiceFabricPath
	var interface_blacklist_updated models.OutofServiceFabricPath
	resourceName := "aci_interface_blacklist.test"
	tDn := makeTestVariable(acctest.RandString(5))
	tDnUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOutofServiceFabricPathDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateOutofServiceFabricPathWithoutRequired(tDn, "t_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOutofServiceFabricPathConfig(tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofServiceFabricPathExists(resourceName, &interface_blacklist_default),

					resource.TestCheckResourceAttr(resourceName, "t_dn", tDn),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "lc", "in-service"),
				),
			},
			{
				Config: CreateAccOutofServiceFabricPathConfigWithOptionalValues(tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofServiceFabricPathExists(resourceName, &interface_blacklist_updated),

					resource.TestCheckResourceAttr(resourceName, "t_dn", tDn),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_interface_blacklist"),

					resource.TestCheckResourceAttr(resourceName, "lc", "blacklist"),

					testAccCheckAciOutofServiceFabricPathIdEqual(&interface_blacklist_default, &interface_blacklist_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccOutofServiceFabricPathRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccOutofServiceFabricPathConfigWithRequiredParams(tDnUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofServiceFabricPathExists(resourceName, &interface_blacklist_updated),

					resource.TestCheckResourceAttr(resourceName, "t_dn", tDnUpdated),
					testAccCheckAciOutofServiceFabricPathIdNotEqual(&interface_blacklist_default, &interface_blacklist_updated),
				),
			},
		},
	})
}

func TestAccAciOutofServiceFabricPath_Negative(t *testing.T) {

	tDn := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOutofServiceFabricPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOutofServiceFabricPathConfig(tDn),
			},

			{
				Config:      CreateAccOutofServiceFabricPathUpdatedAttr(tDn, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOutofServiceFabricPathUpdatedAttr(tDn, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOutofServiceFabricPathUpdatedAttr(tDn, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccOutofServiceFabricPathUpdatedAttr(tDn, "lc", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccOutofServiceFabricPathUpdatedAttr(tDn, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccOutofServiceFabricPathConfig(tDn),
			},
		},
	})
}

func TestAccAciOutofServiceFabricPath_MultipleCreateDelete(t *testing.T) {

	tDn := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOutofServiceFabricPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOutofServiceFabricPathConfigMultiple(tDn),
			},
		},
	})
}

func testAccCheckAciOutofServiceFabricPathExists(name string, interface_blacklist *models.OutofServiceFabricPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface Blacklist %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface Blacklist dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		interface_blacklistFound := models.OutofServiceFabricPathFromContainer(cont)
		if interface_blacklistFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface Blacklist %s not found", rs.Primary.ID)
		}
		*interface_blacklist = *interface_blacklistFound
		return nil
	}
}

func testAccCheckAciOutofServiceFabricPathDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing interface_blacklist destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_interface_blacklist" {
			cont, err := client.Get(rs.Primary.ID)
			interface_blacklist := models.OutofServiceFabricPathFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface Blacklist %s Still exists", interface_blacklist.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciOutofServiceFabricPathIdEqual(m1, m2 *models.OutofServiceFabricPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("interface_blacklist DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciOutofServiceFabricPathIdNotEqual(m1, m2 *models.OutofServiceFabricPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("interface_blacklist DNs are equal")
		}
		return nil
	}
}

func CreateOutofServiceFabricPathWithoutRequired(tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing interface_blacklist creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "t_dn":
		rBlock += `
	resource "aci_interface_blacklist" "test" {
	
	#	t_dn  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, tDn)
}

func CreateAccOutofServiceFabricPathConfigWithRequiredParams(tDn string) string {
	fmt.Println("=== STEP  testing interface_blacklist creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s"
	}
	`, tDn)
	return resource
}

func CreateAccOutofServiceFabricPathConfig(tDn string) string {
	fmt.Println("=== STEP  testing interface_blacklist creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s"
	}
	`, tDn)
	return resource
}

func CreateAccOutofServiceFabricPathConfigMultiple(tDn string) string {
	fmt.Println("=== STEP  testing multiple interface_blacklist creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s_${count.index}"
		count = 5
	}
	`, tDn)
	return resource
}

func CreateAccOutofServiceFabricPathConfigWithOptionalValues(tDn string) string {
	fmt.Println("=== STEP  Basic: testing interface_blacklist creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_interface_blacklist"
		lc = "blacklist"
		
	}
	`, tDn)

	return resource
}

func CreateAccOutofServiceFabricPathRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing interface_blacklist updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_interface_blacklist" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_interface_blacklist"
		lc = "blacklist"
		
	}
	`)

	return resource
}

func CreateAccOutofServiceFabricPathUpdatedAttr(tDn, attribute, value string) string {
	fmt.Printf("=== STEP  testing interface_blacklist attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s"
		%s = "%s"
	}
	`, tDn, attribute, value)
	return resource
}

func CreateAccOutofServiceFabricPathUpdatedAttrList(tDn, attribute, value string) string {
	fmt.Printf("=== STEP  testing interface_blacklist attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s"
		%s = %s
	}
	`, tDn, attribute, value)
	return resource
}
