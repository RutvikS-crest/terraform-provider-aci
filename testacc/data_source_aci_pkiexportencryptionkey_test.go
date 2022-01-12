package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciEncryptionKeyDataSource_Basic(t *testing.T) {
	resourceName := "aci_encryption_key.test"
	dataSourceName := "data.aci_encryption_key.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEncryptionKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEncryptionKeyConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "clear_encryption_key", resourceName, "clear_encryption_key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "passphrase", resourceName, "passphrase"),
					resource.TestCheckResourceAttrPair(dataSourceName, "passphrase_key_derivation_version", resourceName, "passphrase_key_derivation_version"),
					resource.TestCheckResourceAttrPair(dataSourceName, "strong_encryption_enabled", resourceName, "strong_encryption_enabled"),
				),
			},
			{
				Config:      CreateAccEncryptionKeyDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEncryptionKeyDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccEncryptionKeyConfigDataSource() string {
	fmt.Println("=== STEP  testing encryption_key Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
	}

	data "aci_encryption_key" "test" {
	
		depends_on = [ aci_encryption_key.test ]
	}
	`)
	return resource
}

func CreateAccEncryptionKeyDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing encryption_key Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
	}

	data "aci_encryption_key" "test" {
	
		%s = "%s"
		depends_on = [ aci_encryption_key.test ]
	}
	`, key, value)
	return resource
}

func CreateAccEncryptionKeyDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing encryption_key Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
		%s = "%s"
	}

	data "aci_encryption_key" "test" {
	
		depends_on = [ aci_encryption_key.test ]
	}
	`, key, value)
	return resource
}
