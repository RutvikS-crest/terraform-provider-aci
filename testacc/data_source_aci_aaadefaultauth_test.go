package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciDefaultAuthenticationDataSource_Basic(t *testing.T) {
	resourceName := "aci_default_authentication.test"
	dataSourceName := "data.aci_default_authentication.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDefaultAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDefaultAuthenticationConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fallback_check", resourceName, "fallback_check"),
					resource.TestCheckResourceAttrPair(dataSourceName, "provider_group", resourceName, "provider_group"),
					resource.TestCheckResourceAttrPair(dataSourceName, "realm", resourceName, "realm"),
					resource.TestCheckResourceAttrPair(dataSourceName, "realm_sub_type", resourceName, "realm_sub_type"),
				),
			},
			{
				Config:      CreateAccDefaultAuthenticationDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccDefaultAuthenticationDSWithInvalidName(),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccDefaultAuthenticationDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccDefaultAuthenticationConfigDataSource() string {
	fmt.Println("=== STEP  testing default_authentication Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
	}

	data "aci_default_authentication" "test" {
	
		depends_on = [ aci_default_authentication.test ]
	}
	`)
	return resource
}

func CreateDefaultAuthenticationDSWithoutRequired(attrName string) string {
	fmt.Println("=== STEP  Basic: testing default_authentication Data Source without ", attrName)
	rBlock := `
	
	resource "aci_default_authentication" "test" {
	
	}
	data "aci_default_authentication" "test" {
	
		depends_on = [ aci_default_authentication.test ]
	}
	`
	return fmt.Sprintf(rBlock)
}

func CreateAccDefaultAuthenticationDSWithInvalidName() string {
	fmt.Println("=== STEP  testing default_authentication Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
	}

	data "aci_default_authentication" "test" {
	
		depends_on = [ aci_default_authentication.test ]
	}
	`)
	return resource
}

func CreateAccDefaultAuthenticationDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing default_authentication Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
	}

	data "aci_default_authentication" "test" {
	
		%s = "%s"
		depends_on = [ aci_default_authentication.test ]
	}
	`, key, value)
	return resource
}

func CreateAccDefaultAuthenticationDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing default_authentication Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_default_authentication" "test" {
	
		%s = "%s"
	}

	data "aci_default_authentication" "test" {
	
		depends_on = [ aci_default_authentication.test ]
	}
	`, key, value)
	return resource
}
