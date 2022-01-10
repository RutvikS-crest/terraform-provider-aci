package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciX509CertificateDataSource_Basic(t *testing.T) {
	resourceName := "aci_x509_certificate.test"
	dataSourceName := "data.aci_x509_certificate.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciX509CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateX509CertificateDSWithoutRequired(rName, rName, "local_user_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateX509CertificateDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccX509CertificateConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "local_user_dn", resourceName, "local_user_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "data", resourceName, "data"),
				),
			},
			{
				Config:      CreateAccX509CertificateDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccX509CertificateDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccX509CertificateDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccX509CertificateConfigDataSource(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
	}

	data "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = aci_x509_certificate.test.name
		depends_on = [ aci_x509_certificate.test ]
	}
	`, aaaUserName, rName)
	return resource
}

func CreateX509CertificateDSWithoutRequired(aaaUserName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing x509_certificate creation without ", attrName)
	rBlock := `
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "local_user_dn":
		rBlock += `
	data "aci_x509_certificate" "test" {
	#	local_user_dn  = aci_local_user.test.id
		name  = aci_x509_certificate.test.name
		depends_on = [ aci_x509_certificate.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
	#	name  = aci_x509_certificate.test.name
		depends_on = [ aci_x509_certificate.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaUserName, rName)
}

func CreateAccX509CertificateDSWithInvalidParentDn(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
	}

	data "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "${aci_x509_certificate.test.name}_invalid"
		depends_on = [ aci_x509_certificate.test ]
	}
	`, aaaUserName, rName)
	return resource
}

func CreateAccX509CertificateDataSourceUpdate(aaaUserName, rName, key, value string) string {
	fmt.Println("=== STEP  testing x509_certificate Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
	}

	data "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = aci_x509_certificate.test.name
		%s = "%s"
		depends_on = [ aci_x509_certificate.test ]
	}
	`, aaaUserName, rName, key, value)
	return resource
}

func CreateAccX509CertificateDataSourceUpdatedResource(aaaUserName, rName, key, value string) string {
	fmt.Println("=== STEP  testing x509_certificate Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = aci_x509_certificate.test.name
		depends_on = [ aci_x509_certificate.test ]
	}
	`, aaaUserName, rName, key, value)
	return resource
}
