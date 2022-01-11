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

func TestAccAciAnnotation_Basic(t *testing.T) {
	var annotation_default models.Annotation
	var annotation_updated models.Annotation
	resourceName := "aci_annotation.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	key := makeTestVariable(acctest.RandString(5))
	keyUpdated := makeTestVariable(acctest.RandString(5))
	faultInstName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAnnotationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAnnotationWithoutRequired(faultInstName, key, "fault_inst_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAnnotationWithoutRequired(faultInstName, key, "key"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAnnotationConfig(faultInstName, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists(resourceName, &annotation_default),
					resource.TestCheckResourceAttr(resourceName, "fault_inst_dn", GetParentDn(annotation_default.DistinguishedName, fmt.Sprintf("/annotationKey-[%s]", key))),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "value", ""),
				),
			},
			{
				Config: CreateAccAnnotationConfigWithOptionalValues(faultInstName, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists(resourceName, &annotation_updated),
					resource.TestCheckResourceAttr(resourceName, "fault_inst_dn", GetParentDn(annotation_updated.DistinguishedName, fmt.Sprintf("/annotationKey-[%s]", key))),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_annotation"),

					resource.TestCheckResourceAttr(resourceName, "value", ""),

					testAccCheckAciAnnotationIdEqual(&annotation_default, &annotation_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccAnnotationRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAnnotationConfigWithRequiredParams(rNameUpdated, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists(resourceName, &annotation_updated),
					resource.TestCheckResourceAttr(resourceName, "fault_inst_dn", GetParentDn(annotation_updated.DistinguishedName, fmt.Sprintf("/annotationKey-[%s]", key))),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					testAccCheckAciAnnotationIdNotEqual(&annotation_default, &annotation_updated),
				),
			},
			{
				Config: CreateAccAnnotationConfig(faultInstName, key),
			},
			{
				Config: CreateAccAnnotationConfigWithRequiredParams(rName, keyUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnnotationExists(resourceName, &annotation_updated),
					resource.TestCheckResourceAttr(resourceName, "fault_inst_dn", GetParentDn(annotation_updated.DistinguishedName, fmt.Sprintf("/annotationKey-[%s]", key))),
					resource.TestCheckResourceAttr(resourceName, "key", keyUpdated),
					testAccCheckAciAnnotationIdNotEqual(&annotation_default, &annotation_updated),
				),
			},
		},
	})
}

func TestAccAciAnnotation_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	key := makeTestVariable(acctest.RandString(5))

	faultInstName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAnnotationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAnnotationConfig(faultInstName, key),
			},
			{
				Config:      CreateAccAnnotationWithInValidParentDn(rName, key),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAnnotationUpdatedAttr(faultInstName, key, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAnnotationUpdatedAttr(faultInstName, key, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAnnotationUpdatedAttr(faultInstName, key, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccAnnotationUpdatedAttr(faultInstName, key, "value", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccAnnotationUpdatedAttr(faultInstName, key, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAnnotationConfig(faultInstName, key),
			},
		},
	})
}

func TestAccAciAnnotation_MultipleCreateDelete(t *testing.T) {
	key := makeTestVariable(acctest.RandString(5))

	faultInstName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAnnotationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAnnotationConfigMultiple(faultInstName, key),
			},
		},
	})
}

func testAccCheckAciAnnotationExists(name string, annotation *models.Annotation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Annotation %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Annotation dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		annotationFound := models.AnnotationFromContainer(cont)
		if annotationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Annotation %s not found", rs.Primary.ID)
		}
		*annotation = *annotationFound
		return nil
	}
}

func testAccCheckAciAnnotationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing annotation destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_annotation" {
			cont, err := client.Get(rs.Primary.ID)
			annotation := models.AnnotationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Annotation %s Still exists", annotation.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAnnotationIdEqual(m1, m2 *models.Annotation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("annotation DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciAnnotationIdNotEqual(m1, m2 *models.Annotation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("annotation DNs are equal")
		}
		return nil
	}
}

func CreateAnnotationWithoutRequired(faultInstName, key, attrName string) string {
	fmt.Println("=== STEP  Basic: testing annotation creation without ", attrName)
	rBlock := `
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "fault_inst_dn":
		rBlock += `
	resource "aci_annotation" "test" {
	#	fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}
		`
	case "key":
		rBlock += `
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
	#	key  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, faultInstName, key)
}

func CreateAccAnnotationConfigWithRequiredParams(faultInstName, key string) string {
	fmt.Println("=== STEP  testing annotation creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}
	`, faultInstName, key)
	return resource
}

func CreateAccAnnotationConfig(faultInstName, key string) string {
	fmt.Println("=== STEP  testing annotation creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}
	`, faultInstName, key)
	return resource
}

func CreateAccAnnotationConfigMultiple(faultInstName, key string) string {
	fmt.Println("=== STEP  testing multiple annotation creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s_${count.index}"
		count = 5
	}
	`, faultInstName, key)
	return resource
}

func CreateAccAnnotationWithInValidParentDn(rName, key string) string {
	fmt.Println("=== STEP  Negative Case: testing annotation creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_tenant.test.id
		key  = "%s"	
	}
	`, rName, key)
	return resource
}

func CreateAccAnnotationConfigWithOptionalValues(faultInstName, key string) string {
	fmt.Println("=== STEP  Basic: testing annotation creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		fault_inst_dn  = "${aci_fault_inst.test.id}"
		key  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_annotation"
		value = ""
		
	}
	`, faultInstName, key)

	return resource
}

func CreateAccAnnotationRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing annotation updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_annotation" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_annotation"
		value = ""
		
	}
	`)

	return resource
}

func CreateAccAnnotationUpdatedAttr(faultInstName, key, attribute, value string) string {
	fmt.Printf("=== STEP  testing annotation attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
		%s = "%s"
	}
	`, faultInstName, key, attribute, value)
	return resource
}
