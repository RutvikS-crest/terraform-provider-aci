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

func TestAccAciX509Certificate_Basic(t *testing.T) {
	var x509_certificate_default models.X509Certificate
	var x509_certificate_updated models.X509Certificate
	resourceName := "aci_x509_certificate.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	aaaUserName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciX509CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateX509CertificateWithoutRequired(aaaUserName, rName, "local_user_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateX509CertificateWithoutRequired(aaaUserName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccX509CertificateConfig(aaaUserName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciX509CertificateExists(resourceName, &x509_certificate_default),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", aaaUserName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "data", ""),
				),
			},
			{
				Config: CreateAccX509CertificateConfigWithOptionalValues(aaaUserName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciX509CertificateExists(resourceName, &x509_certificate_updated),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", aaaUserName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_x509_certificate"),

					resource.TestCheckResourceAttr(resourceName, "data", ""),

					testAccCheckAciX509CertificateIdEqual(&x509_certificate_default, &x509_certificate_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccX509CertificateConfigUpdatedName(aaaUserName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccX509CertificateRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccX509CertificateConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciX509CertificateExists(resourceName, &x509_certificate_updated),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciX509CertificateIdNotEqual(&x509_certificate_default, &x509_certificate_updated),
				),
			},
			{
				Config: CreateAccX509CertificateConfig(aaaUserName, rName),
			},
			{
				Config: CreateAccX509CertificateConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciX509CertificateExists(resourceName, &x509_certificate_updated),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciX509CertificateIdNotEqual(&x509_certificate_default, &x509_certificate_updated),
				),
			},
		},
	})
}

func TestAccAciX509Certificate_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciX509CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccX509CertificateConfig(aaaUserName, rName),
			},
			{
				Config:      CreateAccX509CertificateWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccX509CertificateUpdatedAttr(aaaUserName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccX509CertificateUpdatedAttr(aaaUserName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccX509CertificateUpdatedAttr(aaaUserName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccX509CertificateUpdatedAttr(aaaUserName, rName, "data", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccX509CertificateUpdatedAttr(aaaUserName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccX509CertificateConfig(aaaUserName, rName),
			},
		},
	})
}

func TestAccAciX509Certificate_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciX509CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccX509CertificateConfigMultiple(aaaUserName, rName),
			},
		},
	})
}

func testAccCheckAciX509CertificateExists(name string, x509_certificate *models.X509Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("X509 Certificate %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No X509 Certificate dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		x509_certificateFound := models.X509CertificateFromContainer(cont)
		if x509_certificateFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("X509 Certificate %s not found", rs.Primary.ID)
		}
		*x509_certificate = *x509_certificateFound
		return nil
	}
}

func testAccCheckAciX509CertificateDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing x509_certificate destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_x509_certificate" {
			cont, err := client.Get(rs.Primary.ID)
			x509_certificate := models.X509CertificateFromContainer(cont)
			if err == nil {
				return fmt.Errorf("X509 Certificate %s Still exists", x509_certificate.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciX509CertificateIdEqual(m1, m2 *models.X509Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("x509_certificate DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciX509CertificateIdNotEqual(m1, m2 *models.X509Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("x509_certificate DNs are equal")
		}
		return nil
	}
}

func CreateX509CertificateWithoutRequired(aaaUserName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing x509_certificate creation without ", attrName)
	rBlock := `
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "local_user_dn":
		rBlock += `
	resource "aci_x509_certificate" "test" {
	#	local_user_dn  = aci_local_user.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaUserName, rName)
}

func CreateAccX509CertificateConfigWithRequiredParams(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate creation with updated name")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
	}
	`, aaaUserName, rName)
	return resource
}
func CreateAccX509CertificateConfigUpdatedName(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
	}
	`, aaaUserName, rName)
	return resource
}

func CreateAccX509CertificateConfig(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
	}
	`, aaaUserName, rName)
	return resource
}

func CreateAccX509CertificateConfigMultiple(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing multiple x509_certificate creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, aaaUserName, rName)
	return resource
}

func CreateAccX509CertificateWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing x509_certificate creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccX509CertificateConfigWithOptionalValues(aaaUserName, rName string) string {
	fmt.Println("=== STEP  Basic: testing x509_certificate creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = "${aci_local_user.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_x509_certificate"
		data = ""
		
	}
	`, aaaUserName, rName)

	return resource
}

func CreateAccX509CertificateRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing x509_certificate updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_x509_certificate" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_x509_certificate"
		data = ""
		
	}
	`)

	return resource
}

func CreateAccX509CertificateUpdatedAttr(aaaUserName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing x509_certificate attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
		%s = "%s"
	}
	`, aaaUserName, rName, attribute, value)
	return resource
}

func CreateAccX509CertificateUpdatedAttrList(aaaUserName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing x509_certificate attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
		%s = %s
	}
	`, aaaUserName, rName, attribute, value)
	return resource
}
