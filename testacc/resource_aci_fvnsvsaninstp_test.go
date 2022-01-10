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

func TestAccAciVSANPool_Basic(t *testing.T) {
	var vsan_pool_default models.VSANPool
	var vsan_pool_updated models.VSANPool
	resourceName := "aci_vsan_pool.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	allocMode := "dynamic"
	allocModeUpdated := "static"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVSANPoolDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateVSANPoolWithoutRequired(rName, allocMode, "alloc_mode"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateVSANPoolWithoutRequired(rName, allocMode, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVSANPoolConfig(rName, allocMode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSANPoolExists(resourceName, &vsan_pool_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "alloc_mode", allocMode),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccVSANPoolConfigWithOptionalValues(rName, allocMode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSANPoolExists(resourceName, &vsan_pool_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "alloc_mode", allocMode),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vsan_pool"),

					testAccCheckAciVSANPoolIdEqual(&vsan_pool_default, &vsan_pool_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVSANPoolConfigUpdatedName(acctest.RandString(65), "dynamic"),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccVSANPoolRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccVSANPoolConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSANPoolExists(resourceName, &vsan_pool_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciVSANPoolIdNotEqual(&vsan_pool_default, &vsan_pool_updated),
				),
			},
			{
				Config: CreateAccVSANPoolConfigWithRequiredParams(allocMode, allocModeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSANPoolExists(resourceName, &vsan_pool_updated),
					resource.TestCheckResourceAttr(resourceName, "alloc_mode", allocModeUpdated),
					testAccCheckAciVSANPoolIdNotEqual(&vsan_pool_default, &vsan_pool_updated),
				),
			},
		},
	})
}

func TestAccAciVSANPool_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	allocMode := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVSANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVSANPoolConfig(rName, allocMode),
			},

			{
				Config:      CreateAccVSANPoolUpdatedAttr(rName, allocMode, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVSANPoolUpdatedAttr(rName, allocMode, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVSANPoolUpdatedAttr(rName, allocMode, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccVSANPoolUpdatedAttr(rName, allocMode, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVSANPoolConfig(rName, allocMode),
			},
		},
	})
}

func TestAccAciVSANPool_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	allocMode := "dynamic"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVSANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVSANPoolConfigMultiple(rName, allocMode),
			},
		},
	})
}

func testAccCheckAciVSANPoolExists(name string, vsan_pool *models.VSANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Vsan Pool %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vsan Pool dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vsan_poolFound := models.VSANPoolFromContainer(cont)
		if vsan_poolFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Vsan Pool %s not found", rs.Primary.ID)
		}
		*vsan_pool = *vsan_poolFound
		return nil
	}
}

func testAccCheckAciVSANPoolDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vsan_pool destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vsan_pool" {
			cont, err := client.Get(rs.Primary.ID)
			vsan_pool := models.VSANPoolFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Vsan Pool %s Still exists", vsan_pool.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVSANPoolIdEqual(m1, m2 *models.VSANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vsan_pool DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVSANPoolIdNotEqual(m1, m2 *models.VSANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vsan_pool DNs are equal")
		}
		return nil
	}
}

func CreateVSANPoolWithoutRequired(rName, allocMode, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vsan_pool creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_vsan_pool" "test" {
	
	#	name  = "%s"
		alloc_mode  = "%s"
	}
		`
	case "alloc_mode":
		rBlock += `
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
	#	alloc_mode  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, allocMode)
}

func CreateAccVSANPoolConfigWithRequiredParams(rName, allocMode string) string {
	fmt.Println("=== STEP  testing vsan_pool creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}
	`, rName, allocMode)
	return resource
}
func CreateAccVSANPoolConfigUpdatedName(rName, allocMode string) string {
	fmt.Println("=== STEP  testing vsan_pool creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}
	`, rName, allocMode)
	return resource
}

func CreateAccVSANPoolConfig(rName, allocMode string) string {
	fmt.Println("=== STEP  testing vsan_pool creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
	}
	`, rName, allocMode)
	return resource
}

func CreateAccVSANPoolConfigMultiple(rName, allocMode string) string {
	fmt.Println("=== STEP  testing multiple vsan_pool creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s_${count.index}"
		alloc_mode  = "%s_${count.index}"
		count = 5
	}
	`, rName, allocMode)
	return resource
}

func CreateAccVSANPoolConfigWithOptionalValues(rName, allocMode string) string {
	fmt.Println("=== STEP  Basic: testing vsan_pool creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vsan_pool"
		
	}
	`, rName, allocMode)

	return resource
}

func CreateAccVSANPoolRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vsan_pool updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vsan_pool" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vsan_pool"
		
	}
	`)

	return resource
}

func CreateAccVSANPoolUpdatedAttr(rName, allocMode, attribute, value string) string {
	fmt.Printf("=== STEP  testing vsan_pool attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
		%s = "%s"
	}
	`, rName, allocMode, attribute, value)
	return resource
}

func CreateAccVSANPoolUpdatedAttrList(rName, allocMode, attribute, value string) string {
	fmt.Printf("=== STEP  testing vsan_pool attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_vsan_pool" "test" {
	
		name  = "%s"
		alloc_mode  = "%s"
		%s = %s
	}
	`, rName, allocMode, attribute, value)
	return resource
}
