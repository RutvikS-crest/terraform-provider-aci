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

func TestAccAciRadiusProviderGroup_Basic(t *testing.T) {
	var radius_provider_group_default models.RADIUSProviderGroup
	var radius_provider_group_updated models.RADIUSProviderGroup
	resourceName := "aci_radius_provider_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRadiusProviderGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateRadiusProviderGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRadiusProviderGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderGroupExists(resourceName, &radius_provider_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccRadiusProviderGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderGroupExists(resourceName, &radius_provider_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_radius_provider_group"),

					testAccCheckAciRadiusProviderGroupIdEqual(&radius_provider_group_default, &radius_provider_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccRadiusProviderGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccRadiusProviderGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccRadiusProviderGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRadiusProviderGroupExists(resourceName, &radius_provider_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciRadiusProviderGroupIdNotEqual(&radius_provider_group_default, &radius_provider_group_updated),
				),
			},
		},
	})
}

func TestAccAciRadiusProviderGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRadiusProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRadiusProviderGroupConfig(rName),
			},

			{
				Config:      CreateAccRadiusProviderGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRadiusProviderGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRadiusProviderGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccRadiusProviderGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccRadiusProviderGroupConfig(rName),
			},
		},
	})
}

func TestAccAciRadiusProviderGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRadiusProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRadiusProviderGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciRadiusProviderGroupExists(name string, radius_provider_group *models.RADIUSProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Radius Provider Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Radius Provider Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		radius_provider_groupFound := models.RADIUSProviderGroupFromContainer(cont)
		if radius_provider_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Radius Provider Group %s not found", rs.Primary.ID)
		}
		*radius_provider_group = *radius_provider_groupFound
		return nil
	}
}

func testAccCheckAciRadiusProviderGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing radius_provider_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_radius_provider_group" {
			cont, err := client.Get(rs.Primary.ID)
			radius_provider_group := models.RADIUSProviderGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Radius Provider Group %s Still exists", radius_provider_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRadiusProviderGroupIdEqual(m1, m2 *models.RADIUSProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("radius_provider_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciRadiusProviderGroupIdNotEqual(m1, m2 *models.RADIUSProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("radius_provider_group DNs are equal")
		}
		return nil
	}
}

func CreateRadiusProviderGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing radius_provider_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_radius_provider_group" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccRadiusProviderGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing radius_provider_group creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccRadiusProviderGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing radius_provider_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccRadiusProviderGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing radius_provider_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccRadiusProviderGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple radius_provider_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider_group" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccRadiusProviderGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing radius_provider_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_radius_provider_group"
		
	}
	`, rName)

	return resource
}

func CreateAccRadiusProviderGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing radius_provider_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_radius_provider_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_radius_provider_group"
		
	}
	`)

	return resource
}

func CreateAccRadiusProviderGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing radius_provider_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccRadiusProviderGroupUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing radius_provider_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_radius_provider_group" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
