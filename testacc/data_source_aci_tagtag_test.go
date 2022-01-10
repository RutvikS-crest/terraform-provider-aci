package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTagDataSource_Basic(t *testing.T) {
	resourceName := "aci_tag.test"
	dataSourceName := "data.aci_tag.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	key := makeTestVariable(acctest.RandString(5))
	faultInstName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTagDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTagDSWithoutRequired(faultInstName, key, "fault_inst_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTagDSWithoutRequired(faultInstName, key, "key"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTagConfigDataSource(faultInstName, key),
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
				Config:      CreateAccTagDataSourceUpdate(faultInstName, key, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTagDSWithInvalidParentDn(faultInstName, key),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccTagDataSourceUpdatedResource(faultInstName, key, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccTagConfigDataSource(faultInstName, key string) string {
	fmt.Println("=== STEP  testing tag Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}

	data "aci_tag" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = aci_tag.test.key
		depends_on = [ aci_tag.test ]
	}
	`, faultInstName, key)
	return resource
}

func CreateTagDSWithoutRequired(faultInstName, key, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tag Data Source without ", attrName)
	rBlock := `
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}
	`
	switch attrName {
	case "fault_inst_dn":
		rBlock += `
	data "aci_tag" "test" {
	#	fault_inst_dn  = aci_fault_inst.test.id
		key  = aci_tag.test.key
		depends_on = [ aci_tag.test ]
	}
		`
	case "key":
		rBlock += `
	data "aci_tag" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
	#	key  = aci_tag.test.key
		depends_on = [ aci_tag.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, faultInstName, key)
}

func CreateAccTagDSWithInvalidParentDn(faultInstName, key string) string {
	fmt.Println("=== STEP  testing tag Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}

	data "aci_tag" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "${aci_tag.test.key}_invalid"
		depends_on = [ aci_tag.test ]
	}
	`, faultInstName, key)
	return resource
}

func CreateAccTagDataSourceUpdate(faultInstName, key, attr, value string) string {
	fmt.Println("=== STEP  testing tag Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
	}

	data "aci_tag" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = aci_tag.test.key
		%s = "%s"
		depends_on = [ aci_tag.test ]
	}
	`, faultInstName, key, attr, value)
	return resource
}

func CreateAccTagDataSourceUpdatedResource(faultInstName, key, attr, value string) string {
	fmt.Println("=== STEP  testing tag Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_fault_inst" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = "%s"
		%s = "%s"
	}

	data "aci_tag" "test" {
		fault_inst_dn  = aci_fault_inst.test.id
		key  = aci_tag.test.key
		depends_on = [ aci_tag.test ]
	}
	`, faultInstName, key, attr, value)
	return resource
}
