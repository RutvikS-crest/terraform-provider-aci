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

func TestAccAciEncryptionKey_Basic(t *testing.T) {
	var encryption_key_default models.AESEncryptionPassphraseandKeysforConfigExportImport
	var encryption_key_updated models.AESEncryptionPassphraseandKeysforConfigExportImport
	resourceName := "aci_encryption_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEncryptionKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEncryptionKeyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEncryptionKeyExists(resourceName, &encryption_key_default),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "clear_encryption_key", "no"),
					resource.TestCheckResourceAttr(resourceName, "passphrase", ""),
					resource.TestCheckResourceAttr(resourceName, "passphrase_key_derivation_version", "0"),
					resource.TestCheckResourceAttr(resourceName, "strong_encryption_enabled", "no"),
				),
			},
			{
				Config: CreateAccEncryptionKeyConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEncryptionKeyExists(resourceName, &encryption_key_updated),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_encryption_key"),

					resource.TestCheckResourceAttr(resourceName, "clear_encryption_key", "yes"),

					resource.TestCheckResourceAttr(resourceName, "passphrase", ""),

					resource.TestCheckResourceAttr(resourceName, "passphrase_key_derivation_version", "v1"),

					resource.TestCheckResourceAttr(resourceName, "strong_encryption_enabled", "yes"),

					testAccCheckAciEncryptionKeyIdEqual(&encryption_key_default, &encryption_key_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccEncryptionKeyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestAccAciEncryptionKey_Negative(t *testing.T) {

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEncryptionKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEncryptionKeyConfig(),
			},

			{
				Config:      CreateAccEncryptionKeyUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEncryptionKeyUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEncryptionKeyUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccEncryptionKeyUpdatedAttr("clear_encryption_key", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccEncryptionKeyUpdatedAttr("passphrase", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccEncryptionKeyUpdatedAttr("passphrase_key_derivation_version", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccEncryptionKeyUpdatedAttr("strong_encryption_enabled", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccEncryptionKeyUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEncryptionKeyConfig(),
			},
		},
	})
}

func testAccCheckAciEncryptionKeyExists(name string, encryption_key *models.AESEncryptionPassphraseandKeysforConfigExportImport) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Encryption Key %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Encryption Key dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		encryption_keyFound := models.AESEncryptionPassphraseandKeysforConfigExportImportFromContainer(cont)
		if encryption_keyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Encryption Key %s not found", rs.Primary.ID)
		}
		*encryption_key = *encryption_keyFound
		return nil
	}
}

func testAccCheckAciEncryptionKeyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing encryption_key destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_encryption_key" {
			cont, err := client.Get(rs.Primary.ID)
			encryption_key := models.AESEncryptionPassphraseandKeysforConfigExportImportFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Encryption Key %s Still exists", encryption_key.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciEncryptionKeyIdEqual(m1, m2 *models.AESEncryptionPassphraseandKeysforConfigExportImport) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("encryption_key DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciEncryptionKeyIdNotEqual(m1, m2 *models.AESEncryptionPassphraseandKeysforConfigExportImport) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("encryption_key DNs are equal")
		}
		return nil
	}
}

func CreateAccEncryptionKeyConfigWithRequiredParams() string {
	fmt.Println("=== STEP  testing encryption_key creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
	}
	`)
	return resource
}

func CreateAccEncryptionKeyConfig() string {
	fmt.Println("=== STEP  testing encryption_key creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
	}
	`)
	return resource
}

func CreateAccEncryptionKeyConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing encryption_key creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_encryption_key"
		clear_encryption_key = "yes"
		passphrase = ""
		passphrase_key_derivation_version = "v1"
		strong_encryption_enabled = "yes"
		
	}
	`)

	return resource
}

func CreateAccEncryptionKeyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing encryption_key updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_encryption_key" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_encryption_key"
		clear_encryption_key = "yes"
		passphrase = ""
		passphrase_key_derivation_version = "v1"
		strong_encryption_enabled = "yes"
		
	}
	`)

	return resource
}

func CreateAccEncryptionKeyUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing encryption_key attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}

func CreateAccEncryptionKeyUpdatedAttrList(attribute, value string) string {
	fmt.Printf("=== STEP  testing encryption_key attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
		%s = %s
	}
	`, attribute, value)
	return resource
}
