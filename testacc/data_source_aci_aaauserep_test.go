package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciGlobalSecurityDataSource_Basic(t *testing.T) {
	resourceName := "aci_global_security.test"
	dataSourceName := "data.aci_global_security.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciGlobalSecurityDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccGlobalSecurityConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pwd_strength_check", resourceName, "pwd_strength_check"),
				),
			},
			{
				Config:      CreateAccGlobalSecurityDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccGlobalSecurityDSWithInvalidName(),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccGlobalSecurityDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccGlobalSecurityConfigDataSource() string {
	fmt.Println("=== STEP  testing global_security Data Source with required arguments only")
	resource := `
	
	resource "aci_global_security" "test" {
	
	}

	data "aci_global_security" "test" {
	
		depends_on = [ aci_global_security.test ]
	}
	`
	return resource
}

func CreateAccGlobalSecurityDSWithInvalidName() string {
	fmt.Println("=== STEP  testing global_security Data Source with required arguments only")
	resource := `
	
	resource "aci_global_security" "test" {
	
	}

	data "aci_global_security" "test" {
	
		depends_on = [ aci_global_security.test ]
	}
	`
	return resource
}

func CreateAccGlobalSecurityDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing global_security Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {
	
	}

	data "aci_global_security" "test" {
	
		%s = "%s"
		depends_on = [ aci_global_security.test ]
	}
	`, key, value)
	return resource
}

func CreateAccGlobalSecurityDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing global_security Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {
	
		%s = "%s"
	}

	data "aci_global_security" "test" {
	
		depends_on = [ aci_global_security.test ]
	}
	`, key, value)
	return resource
}
