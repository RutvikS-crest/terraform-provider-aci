package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSAMLCertificateDataSource_Basic(t *testing.T) {
	resourceName := "aci_saml_certificate.test"
	dataSourceName := "data.aci_saml_certificate.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{

			{
				Config:      CreateSAMLCertificateDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSAMLCertificateConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "regenerate", resourceName, "regenerate"),
				),
			},
			{
				Config:      CreateAccSAMLCertificateDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSAMLCertificateDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccSAMLCertificateDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSAMLCertificateConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing saml_certificate Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_certificate" "test" {
	
		name  = "%s"
	}

	data "aci_saml_certificate" "test" {
	
		name  = aci_saml_certificate.test.name
		depends_on = [ aci_saml_certificate.test ]
	}
	`, rName)
	return resource
}

func CreateSAMLCertificateDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing saml_certificate Data Source without ", attrName)
	rBlock := `
	
	resource "aci_saml_certificate" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_saml_certificate" "test" {
	
	#	name  = aci_saml_certificate.test.name
		depends_on = [ aci_saml_certificate.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSAMLCertificateDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing saml_certificate Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_certificate" "test" {
	
		name  = "%s"
	}

	data "aci_saml_certificate" "test" {
	
		name  = "${aci_saml_certificate.test.name}_invalid"
		name  = aci_saml_certificate.test.name
		depends_on = [ aci_saml_certificate.test ]
	}
	`, rName)
	return resource
}

func CreateAccSAMLCertificateDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing saml_certificate Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_certificate" "test" {
	
		name  = "%s"
	}

	data "aci_saml_certificate" "test" {
	
		name  = aci_saml_certificate.test.name
		%s = "%s"
		depends_on = [ aci_saml_certificate.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccSAMLCertificateDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing saml_certificate Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_certificate" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_saml_certificate" "test" {
	
		name  = aci_saml_certificate.test.name
		depends_on = [ aci_saml_certificate.test ]
	}
	`, rName, key, value)
	return resource
}
