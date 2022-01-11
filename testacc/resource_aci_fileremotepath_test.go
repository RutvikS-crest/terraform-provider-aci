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

func TestAccAciRemotePathofaFile_Basic(t *testing.T) {
	var file_remote_path_default models.RemotePathofaFile
	var file_remote_path_updated models.RemotePathofaFile
	resourceName := "aci_file_remote_path.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRemotePathofaFileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateRemotePathofaFileWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRemotePathofaFileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRemotePathofaFileExists(resourceName, &file_remote_path_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "usePassword"),
					resource.TestCheckResourceAttr(resourceName, "host", ""),
					resource.TestCheckResourceAttr(resourceName, "identity_private_key_contents", ""),
					resource.TestCheckResourceAttr(resourceName, "identity_private_key_passphrase", ""),
					resource.TestCheckResourceAttr(resourceName, "identity_public_key_contents", ""),
					resource.TestCheckResourceAttr(resourceName, "protocol", "sftp"),
					resource.TestCheckResourceAttr(resourceName, "remote_path", ""),
					resource.TestCheckResourceAttr(resourceName, "remote_port", ""),
					resource.TestCheckResourceAttr(resourceName, "user_name", ""),
					resource.TestCheckResourceAttr(resourceName, "user_passwd", ""),
				),
			},
			{
				Config: CreateAccRemotePathofaFileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRemotePathofaFileExists(resourceName, &file_remote_path_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_file_remote_path"),

					resource.TestCheckResourceAttr(resourceName, "auth_type", "useSshKeyContents"),

					resource.TestCheckResourceAttr(resourceName, "host", ""),

					resource.TestCheckResourceAttr(resourceName, "identity_private_key_contents", ""),

					resource.TestCheckResourceAttr(resourceName, "identity_private_key_passphrase", ""),

					resource.TestCheckResourceAttr(resourceName, "identity_public_key_contents", ""),

					resource.TestCheckResourceAttr(resourceName, "protocol", "ftp"),

					resource.TestCheckResourceAttr(resourceName, "remote_path", ""),

					resource.TestCheckResourceAttr(resourceName, "user_name", ""),

					resource.TestCheckResourceAttr(resourceName, "user_passwd", ""),

					testAccCheckAciRemotePathofaFileIdEqual(&file_remote_path_default, &file_remote_path_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccRemotePathofaFileConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccRemotePathofaFileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccRemotePathofaFileConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRemotePathofaFileExists(resourceName, &file_remote_path_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciRemotePathofaFileIdNotEqual(&file_remote_path_default, &file_remote_path_updated),
				),
			},
		},
	})
}

func TestAccAciRemotePathofaFile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRemotePathofaFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRemotePathofaFileConfig(rName),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "auth_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "host", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "identity_private_key_contents", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "identity_private_key_passphrase", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "identity_public_key_contents", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "protocol", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "remote_path", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "remote_port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "user_name", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, "user_passwd", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRemotePathofaFileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccRemotePathofaFileConfig(rName),
			},
		},
	})
}

func TestAccAciRemotePathofaFile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRemotePathofaFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRemotePathofaFileConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciRemotePathofaFileExists(name string, file_remote_path *models.RemotePathofaFile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("File Remote Path %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No File Remote Path dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		file_remote_pathFound := models.RemotePathofaFileFromContainer(cont)
		if file_remote_pathFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("File Remote Path %s not found", rs.Primary.ID)
		}
		*file_remote_path = *file_remote_pathFound
		return nil
	}
}

func testAccCheckAciRemotePathofaFileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing file_remote_path destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_file_remote_path" {
			cont, err := client.Get(rs.Primary.ID)
			file_remote_path := models.RemotePathofaFileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("File Remote Path %s Still exists", file_remote_path.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRemotePathofaFileIdEqual(m1, m2 *models.RemotePathofaFile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("file_remote_path DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciRemotePathofaFileIdNotEqual(m1, m2 *models.RemotePathofaFile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("file_remote_path DNs are equal")
		}
		return nil
	}
}

func CreateRemotePathofaFileWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing file_remote_path creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_file_remote_path" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccRemotePathofaFileConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing file_remote_path creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccRemotePathofaFileConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing file_remote_path creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccRemotePathofaFileConfig(rName string) string {
	fmt.Println("=== STEP  testing file_remote_path creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccRemotePathofaFileConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple file_remote_path creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccRemotePathofaFileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing file_remote_path creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_file_remote_path" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_file_remote_path"
		auth_type = "useSshKeyContents"
		host = ""
		identity_private_key_contents = ""
		identity_private_key_passphrase = ""
		identity_public_key_contents = ""
		protocol = "ftp"
		remote_path = ""
		user_name = ""
		user_passwd = ""
		
	}
	`, rName)

	return resource
}

func CreateAccRemotePathofaFileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing file_remote_path updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_file_remote_path" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_file_remote_path"
		auth_type = "useSshKeyContents"
		host = ""
		identity_private_key_contents = ""
		identity_private_key_passphrase = ""
		identity_public_key_contents = ""
		protocol = "ftp"
		remote_path = ""
		user_name = ""
		user_passwd = ""
		
	}
	`)

	return resource
}
