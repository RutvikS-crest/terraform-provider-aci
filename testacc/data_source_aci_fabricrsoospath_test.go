package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciInterfaceBlacklistDataSource_Basic(t *testing.T) {
	resourceName := "aci_interface_blacklist.test"
	dataSourceName := "data.aci_interface_blacklist.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	tDn := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOutofServiceFabricPathDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateInterfaceBlacklistDSWithoutRequired(tDn, "t_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccInterfaceBlacklistConfigDataSource(tDn),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "t_dn", resourceName, "t_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lc", resourceName, "lc"),
				),
			},
			{
				Config:      CreateAccInterfaceBlacklistDataSourceUpdate(tDn, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccInterfaceBlacklistDSWithInvalidName(tDn),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccInterfaceBlacklistDataSourceUpdatedResource(tDn, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccInterfaceBlacklistConfigDataSource(tDn string) string {
	fmt.Println("=== STEP  testing interface_blacklist Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s"
	}

	data "aci_interface_blacklist" "test" {
	
		t_dn  = aci_interface_blacklist.test.tDn
		depends_on = [ aci_interface_blacklist.test ]
	}
	`, tDn)
	return resource
}

func CreateInterfaceBlacklistDSWithoutRequired(tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing interface_blacklist Data Source without ", attrName)
	rBlock := `
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s"
	}
	`
	switch attrName {
	case "t_dn":
		rBlock += `
	data "aci_interface_blacklist" "test" {
	
	#	t_dn  = aci_interface_blacklist.test.tDn
		depends_on = [ aci_interface_blacklist.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, tDn)
}

func CreateAccInterfaceBlacklistDSWithInvalidName(tDn string) string {
	fmt.Println("=== STEP  testing interface_blacklist Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s"
	}

	data "aci_interface_blacklist" "test" {
	
		t_dn  = aci_interface_blacklist.test.tDn
		depends_on = [ aci_interface_blacklist.test ]
	}
	`, tDn)
	return resource
}

func CreateAccInterfaceBlacklistDataSourceUpdate(tDn, key, value string) string {
	fmt.Println("=== STEP  testing interface_blacklist Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s"
	}

	data "aci_interface_blacklist" "test" {
	
		t_dn  = aci_interface_blacklist.test.tDn
		%s = "%s"
		depends_on = [ aci_interface_blacklist.test ]
	}
	`, tDn, key, value)
	return resource
}

func CreateAccInterfaceBlacklistDataSourceUpdatedResource(tDn, key, value string) string {
	fmt.Println("=== STEP  testing interface_blacklist Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		t_dn  = "%s"
		%s = "%s"
	}

	data "aci_interface_blacklist" "test" {
	
		t_dn  = aci_interface_blacklist.test.tDn
		depends_on = [ aci_interface_blacklist.test ]
	}
	`, tDn, key, value)
	return resource
}
