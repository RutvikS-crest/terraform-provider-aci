package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVsanPoolDataSource_Basic(t *testing.T) {
	resourceName := "aci_vsan_pool.test"
	dataSourceName := "data.aci_vsan_pool.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	allocMode := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVSANPoolDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateVsanPoolDSWithoutRequired(rName, allocMode, "alloc_mode"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateVsanPoolDSWithoutRequired(rName, allocMode, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVsanPoolConfigDataSource(rName, allocMode),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "alloc_mode", resourceName, "alloc_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccVsanPoolDataSourceUpdate(rName, allocMode, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccVsanPoolDSWithInvalidName(rName, allocMode),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccVsanPoolDataSourceUpdatedResource(rName, allocMode, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVsanPoolConfigDataSource(rName, allocMode string) string {
	fmt.Println("=== STEP  testing vsan_pool Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}

	data "aci_vsan_pool" "test" {
	
		name  = aci_vsan_pool.test.name
		alloc_mode  = aci_vsan_pool.test.allocMode
		depends_on = [ aci_vsan_pool.test ]
	}
	`, rName, allocMode)
	return resource
}

func CreateVsanPoolDSWithoutRequired(rName, allocMode, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vsan_pool Data Source without ", attrName)
	rBlock := `
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_vsan_pool" "test" {
	
	#	name  = aci_vsan_pool.test.name
		alloc_mode  = aci_vsan_pool.test.allocMode
		depends_on = [ aci_vsan_pool.test ]
	}
		`
	case "alloc_mode":
		rBlock += `
	data "aci_vsan_pool" "test" {
	
		name  = aci_vsan_pool.test.name
	#	alloc_mode  = aci_vsan_pool.test.allocMode
		depends_on = [ aci_vsan_pool.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, allocMode)
}

func CreateAccVsanPoolDSWithInvalidName(rName, allocMode string) string {
	fmt.Println("=== STEP  testing vsan_pool Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}

	data "aci_vsan_pool" "test" {
	
		name  = "${aci_vsan_pool.test.name}_invalid"
		alloc_mode  = aci_vsan_pool.test.allocMode
		depends_on = [ aci_vsan_pool.test ]
	}
	`, rName, allocMode)
	return resource
}

func CreateAccVsanPoolDataSourceUpdate(rName, allocMode, key, value string) string {
	fmt.Println("=== STEP  testing vsan_pool Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}

	data "aci_vsan_pool" "test" {
	
		name  = aci_vsan_pool.test.name
		alloc_mode  = aci_vsan_pool.test.allocMode
		%s = "%s"
		depends_on = [ aci_vsan_pool.test ]
	}
	`, rName, allocMode, key, value)
	return resource
}

func CreateAccVsanPoolDataSourceUpdatedResource(rName, allocMode, key, value string) string {
	fmt.Println("=== STEP  testing vsan_pool Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
		%s = "%s"
	}

	data "aci_vsan_pool" "test" {
	
		name  = aci_vsan_pool.test.name
		alloc_mode  = aci_vsan_pool.test.allocMode
		depends_on = [ aci_vsan_pool.test ]
	}
	`, rName, allocMode, key, value)
	return resource
}
