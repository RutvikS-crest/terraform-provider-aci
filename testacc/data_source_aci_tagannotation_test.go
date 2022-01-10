package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciAnnotationDataSource_Basic(t *testing.T) {
	resourceName := "aci_annotation.test"
	dataSourceName := "data.aci_annotation.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	key := makeTestVariable(acctest.RandString(5))
	faultInstName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAnnotationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAnnotationDSWithoutRequired(faultInstName, key, "fault_inst_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAnnotationDSWithoutRequired(faultInstName, key, "key"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAnnotationConfigDataSource(faultInstName, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "fault_inst_dn", resourceName, "fault_inst_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "key", resourceName, "key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "value", resourceName, "value"),
				),
			},
			{
				Config:      CreateAccAnnotationDataSourceUpdate(faultInstName, key, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccAnnotationDSWithInvalidParentDn(faultInstName, key),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccAnnotationDataSourceUpdatedResource(faultInstName, key, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccAnnotationConfigDataSource(faultInstName, key string) string {
	fmt.Println("=== STEP  testing annotation Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}

	data "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = aci_annotation.test.key
		depends_on = [ aci_annotation.test ]
	}
	`, faultInstName, key)
	return resource
}

func CreateAnnotationDSWithoutRequired(faultInstName, key, attrName string) string {
	fmt.Println("=== STEP  Basic: testing annotation Data Source without ", attrName)
	rBlock := `
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}
	`
	switch attrName {
	case "fault_inst_dn":
		rBlock += `
	data "aci_annotation" "test" {
	#	fault_inst_dn  = aci_fault_inst.test.id
		key  = aci_annotation.test.key
		depends_on = [ aci_annotation.test ]
	}
		`
	case "key":
		rBlock += `
	data "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
	#	key  = aci_annotation.test.key
		depends_on = [ aci_annotation.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, faultInstName, key)
}

func CreateAccAnnotationDSWithInvalidParentDn(faultInstName, key string) string {
	fmt.Println("=== STEP  testing annotation Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}

	data "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "${aci_annotation.test.key}_invalid"
		depends_on = [ aci_annotation.test ]
	}
	`, faultInstName, key)
	return resource
}

func CreateAccAnnotationDataSourceUpdate(faultInstName, key, attr, value string) string {
	fmt.Println("=== STEP  testing annotation Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}

	data "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = aci_annotation.test.key
		%s = "%s"
		depends_on = [ aci_annotation.test ]
	}
	`, faultInstName, key, attr, value)
	return resource
}

func CreateAccAnnotationDataSourceUpdatedResource(faultInstName, key, attr, value string) string {
	fmt.Println("=== STEP  testing annotation Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
		%s = "%s"
	}

	data "aci_annotation" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = aci_annotation.test.key
		depends_on = [ aci_annotation.test ]
	}
	`, faultInstName, key, attr, value)
	return resource
}
