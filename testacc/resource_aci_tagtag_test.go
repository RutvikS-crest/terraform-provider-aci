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

func TestAccAciTag_Basic(t *testing.T) {
	var tag_default models.Tag
	var tag_updated models.Tag
	resourceName := "aci_tag.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	key := makeTestVariable(acctest.RandString(5))
	keyUpdated := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTagDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTagWithoutRequired(fvTenantName, key, "parent_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTagWithoutRequired(fvTenantName, key, "key"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTagConfig(fvTenantName, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists(resourceName, &tag_default),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", key)),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "value", ""),
				),
			},
			{
				Config: CreateAccTagConfigWithOptionalValues(fvTenantName, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists(resourceName, &tag_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", key)),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_tag"),

					resource.TestCheckResourceAttr(resourceName, "value", ""),

					testAccCheckAciTagIdEqual(&tag_default, &tag_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccTagRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTagConfigWithRequiredParams(rNameUpdated, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists(resourceName, &tag_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					testAccCheckAciTagIdNotEqual(&tag_default, &tag_updated),
				),
			},
			{
				Config: CreateAccTagConfig(fvTenantName, key),
			},
			{
				Config: CreateAccTagConfigWithRequiredParams(rName, keyUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists(resourceName, &tag_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "key", keyUpdated),
					testAccCheckAciTagIdNotEqual(&tag_default, &tag_updated),
				),
			},
		},
	})
}

func TestAccAciTag_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	key := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTagConfig(fvTenantName, key),
			},
			{
				Config:      CreateAccTagWithInValidParentDn(rName, key),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccTagUpdatedAttr(fvTenantName, key, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTagUpdatedAttr(fvTenantName, key, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTagUpdatedAttr(fvTenantName, key, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccTagUpdatedAttr(fvTenantName, key, "value", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccTagUpdatedAttr(fvTenantName, key, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccTagConfig(fvTenantName, key),
			},
		},
	})
}

func TestAccAciTag_MultipleCreateDelete(t *testing.T) {

	key := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTagConfigMultiple(fvTenantName, key),
			},
		},
	})
}

func testAccCheckAciTagExists(name string, tag *models.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Tag %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tag dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tagFound := models.TagFromContainer(cont)
		if tagFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Tag %s not found", rs.Primary.ID)
		}
		*tag = *tagFound
		return nil
	}
}

func testAccCheckAciTagDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing tag destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tag" {
			cont, err := client.Get(rs.Primary.ID)
			tag := models.TagFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Tag %s Still exists", tag.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTagIdEqual(m1, m2 *models.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("tag DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTagIdNotEqual(m1, m2 *models.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("tag DNs are equal")
		}
		return nil
	}
}

func CreateTagWithoutRequired(fvTenantName, key, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tag creation without ", attrName)
	rBlock := `
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "parent_dn":
		rBlock += `
	resource "aci_tag" "test" {
	#	parent_dn  = aci_fault_inst.test.id
		key  = "%s"
	}
		`
	case "key":
		rBlock += `
	resource "aci_tag" "test" {
		parent_dn  = aci_fault_inst.test.id
	#	key  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, key)
}

func CreateAccTagConfigWithRequiredParams(fvTenantName, key string) string {
	fmt.Println("=== STEP  testing tag creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_fault_inst.test.id
		key  = "%s"
	}
	`, fvTenantName, key)
	return resource
}

func CreateAccTagConfig(fvTenantName, key string) string {
	fmt.Println("=== STEP  testing tag creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_fault_inst.test.id
		key  = "%s"
	}
	`, fvTenantName, key)
	return resource
}

func CreateAccTagConfigMultiple(fvTenantName, key string) string {
	fmt.Println("=== STEP  testing multiple tag creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_fault_inst.test.id
		key  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, key)
	return resource
}

func CreateAccTagWithInValidParentDn(rName, key string) string {
	fmt.Println("=== STEP  Negative Case: testing tag creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"	
	}
	`, rName, key)
	return resource
}

func CreateAccTagConfigWithOptionalValues(fvTenantName, key string) string {
	fmt.Println("=== STEP  Basic: testing tag creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = "${aci_fault_inst.test.id}"
		key  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tag"
		value = ""
		
	}
	`, fvTenantName, key)

	return resource
}

func CreateAccTagRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing tag updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tag" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_tag"
		value = ""
		
	}
	`)

	return resource
}

func CreateAccTagUpdatedAttr(fvTenantName, key, attribute, value string) string {
	fmt.Printf("=== STEP  testing tag attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_fault_inst.test.id
		key  = "%s"
		%s = "%s"
	}
	`, fvTenantName, key, attribute, value)
	return resource
}

func CreateAccTagUpdatedAttrList(fvTenantName, key, attribute, value string) string {
	fmt.Printf("=== STEP  testing tag attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_fault_inst.test.id
		key  = "%s"
		%s = %s
	}
	`, fvTenantName, key, attribute, value)
	return resource
}
