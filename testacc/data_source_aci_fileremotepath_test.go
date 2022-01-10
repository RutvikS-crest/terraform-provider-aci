package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFileRemotePathDataSource_Basic(t *testing.T) {
	resourceName := "aci_file_remote_path.test"
	dataSourceName := "data.aci_file_remote_path.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRemotePathofaFileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFileRemotePathDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFileRemotePathConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_type", resourceName, "auth_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "host", resourceName, "host"),
					resource.TestCheckResourceAttrPair(dataSourceName, "identity_private_key_contents", resourceName, "identity_private_key_contents"),
					resource.TestCheckResourceAttrPair(dataSourceName, "identity_private_key_passphrase", resourceName, "identity_private_key_passphrase"),
					resource.TestCheckResourceAttrPair(dataSourceName, "identity_public_key_contents", resourceName, "identity_public_key_contents"),
					resource.TestCheckResourceAttrPair(dataSourceName, "protocol", resourceName, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourceName, "remote_path", resourceName, "remote_path"),
					resource.TestCheckResourceAttrPair(dataSourceName, "remote_port", resourceName, "remote_port"),
					resource.TestCheckResourceAttrPair(dataSourceName, "user_name", resourceName, "user_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "user_passwd", resourceName, "user_passwd"),
				),
			},
			{
				Config:      CreateAccFileRemotePathDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFileRemotePathDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccFileRemotePathDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFileRemotePathConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing file_remote_path Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
	}

	data "aci_file_remote_path" "test" {
	
		name  = aci_file_remote_path.test.name
		depends_on = [ aci_file_remote_path.test ]
	}
	`, rName)
	return resource
}

func CreateFileRemotePathDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing file_remote_path Data Source without ", attrName)
	rBlock := `
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_file_remote_path" "test" {
	
	#	name  = aci_file_remote_path.test.name
		depends_on = [ aci_file_remote_path.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFileRemotePathDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing file_remote_path Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
	}

	data "aci_file_remote_path" "test" {
	
		name  = "${aci_file_remote_path.test.name}_invalid"
		name  = aci_file_remote_path.test.name
		depends_on = [ aci_file_remote_path.test ]
	}
	`, rName)
	return resource
}

func CreateAccFileRemotePathDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing file_remote_path Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
	}

	data "aci_file_remote_path" "test" {
	
		name  = aci_file_remote_path.test.name
		%s = "%s"
		depends_on = [ aci_file_remote_path.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccFileRemotePathDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing file_remote_path Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_file_remote_path" "test" {
	
		name  = aci_file_remote_path.test.name
		depends_on = [ aci_file_remote_path.test ]
	}
	`, rName, key, value)
	return resource
}
