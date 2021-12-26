package acctest

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

func TestAccAciEPGsUsingFunction_Basic(t *testing.T) {
	var epgs_using_function_default models.EPGsUsingFunction
	var epgs_using_function_updated models.EPGsUsingFunction
	resourceName := "aci_epgs_using_function.testacc"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	tDn := makeTestVariable(acctest.RandString(5))
	tDnUpdated := makeTestVariable(acctest.RandString(5))
	infraAttEntityPName := makeTestVariable(acctest.RandString(5))
	infraGenericName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPGsUsingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateEPGsUsingFunctionWithoutRequired(infraAttEntityPName, infraGenericName, tDn, "access_generic_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateEPGsUsingFunctionWithoutRequired(infraAttEntityPName, infraGenericName, tDn, "tDn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfig(infraAttEntityPName, infraGenericName, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_default),
					resource.TestCheckResourceAttr(resourceName, "access_generic_dn", GetParentDn(epgs_using_function_default.DistinguishedName, fmt.Sprintf("/rsfuncToEpg-[%s]", tDn))),
					resource.TestCheckResourceAttr(resourceName, "tDn", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "encap", ""),
					resource.TestCheckResourceAttr(resourceName, "instr_imedcy", "lazy"),
					resource.TestCheckResourceAttr(resourceName, "mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "primary_encap", ""),
				),
			},
			{
				// in this step all optional attribute expect realational attribute are given for the same resource and then compared
				Config: CreateAccEPGsUsingFunctionConfigWithOptionalValues(infraAttEntityPName, infraGenericName, tDn), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_epgs_using_function"),
					resource.TestCheckResourceAttr(resourceName, "encap", ""),
					resource.TestCheckResourceAttr(resourceName, "instr_imedcy", "immediate"),
					resource.TestCheckResourceAttr(resourceName, "mode", "native"),
					resource.TestCheckResourceAttr(resourceName, "primary_encap", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: CreateAccEPGsUsingFunctionConfigWithRequiredParams(rNameUpdated, tDn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_updated),
					resource.TestCheckResourceAttr(resourceName, "access_generic_dn", GetParentDn(epgs_using_function_default.DistinguishedName, fmt.Sprintf("/rsfuncToEpg-[%s]", tDn))),
					resource.TestCheckResourceAttr(resourceName, "tDn", tDn),
					testAccCheckAciEPGsUsingFunctionIdNotEqual(&epgs_using_function_default, &epgs_using_function_updated),
				),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfig(infraAttEntityPName, infraGenericName, tDn),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfigWithRequiredParams(rName, tDnUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists(resourceName, &epgs_using_function_updated),
					resource.TestCheckResourceAttr(resourceName, "access_generic_dn", GetParentDn(epgs_using_function_default.DistinguishedName, fmt.Sprintf("/rsfuncToEpg-[%s]", tDn))),
					resource.TestCheckResourceAttr(resourceName, "tDn", tDnUpdated),
					testAccCheckAciEPGsUsingFunctionIdNotEqual(&epgs_using_function_default, &epgs_using_function_updated),
				),
			},
		},
	})
}

func TestAccAciEPGsUsingFunction_Negative(t *testing.T) {

	// rName := makeTestVariable(acctest.RandString(5))

	tDn := makeTestVariable(acctest.RandString(5))
	// tDnUpdated := makeTestVariable(acctest.RandString(5))
	infraAttEntityPName := makeTestVariable(acctest.RandString(5))
	infraGenericName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEPGsUsingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEPGsUsingFunctionConfig(infraAttEntityPName, infraGenericName, tDn),
			},
			{
				Config:      CreateAccEPGsUsingFunctionWithInValidParentDn(infraAttEntityPName, infraGenericName, tDn),
				ExpectError: regexp.MustCompile(`configured object (.)+ not found (.)+,`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, infraGenericName, tDn, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, infraGenericName, tDn, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, infraGenericName, tDn, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, infraGenericName, tDn, "encap", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, infraGenericName, tDn, "instr_imedcy", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, infraGenericName, tDn, "mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, infraGenericName, tDn, "primary_encap", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},

			{
				Config:      CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, infraGenericName, tDn, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccEPGsUsingFunctionConfig(infraAttEntityPName, infraGenericName, tDn),
			},
		},
	})
}

func testAccCheckAciEPGsUsingFunctionExists(name string, epgs_using_function *models.EPGsUsingFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("EPGs Using Function %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EPGs Using Function dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		epgs_using_functionFound := models.EPGsUsingFunctionFromContainer(cont)
		if epgs_using_functionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("EPGs Using Function %s not found", rs.Primary.ID)
		}
		*epgs_using_function = *epgs_using_functionFound
		return nil
	}
}

func testAccCheckAciEPGsUsingFunctionDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing epgs_using_function destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_epgs_using_function" {
			cont, err := client.Get(rs.Primary.ID)
			epgs_using_function := models.EPGsUsingFunctionFromContainer(cont)
			if err == nil {
				return fmt.Errorf("EPGs Using Function %s Still exists", epgs_using_function.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciEPGsUsingFunctionIdEqual(m1, m2 *models.EPGsUsingFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("epgs_using_function DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciEPGsUsingFunctionIdNotEqual(m1, m2 *models.EPGsUsingFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("epgs_using_function DNs are equal")
		}
		return nil
	}
}

func CreateEPGsUsingFunctionWithoutRequired(infraAttEntityPName, infraGenericName, tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing epgs_using_function creation without ", attrName)
	rBlock := `
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
		description = "attachable_access_entity_profile created while acceptance testing"
		
	}
	
	resource "aci_access_generic" "test" {
		name 		= "%s"
		description = "access_generic created while acceptance testing"
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
	}
	
	`
	switch attrName {
	case "access_generic_dn":
		rBlock += `
	resource "aci_epgs_using_function" "test" {
	#	access_generic_dn  = aci_access_generic.test.id
		tDn  = "%s"
		description = "created while acceptance testing"
	}
		`
	case "tDn":
		rBlock += `
	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
	#	tDn  = "%s"
		description = "created while acceptance testing"
	}
		`
	}
	return fmt.Sprintf(rBlock, infraAttEntityPName, infraGenericName, tDn)
}

func CreateAccEPGsUsingFunctionConfigWithRequiredParams(rName, tDn string) string {
	fmt.Println("=== STEP  testing epgs_using_function creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
		description = "attachable_access_entity_profile created while acceptance testing"
	
	}
	
	resource "aci_access_generic" "test" {
		name 		= "%s"
		description = "access_generic created while acceptance testing"
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
	}
	
	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tDn  = "%s"
	}
	`, rName, rName, tDn)
	return resource
}

func CreateAccEPGsUsingFunctionConfig(infraAttEntityPName, infraGenericName, tDn string) string {
	fmt.Println("=== STEP  testing epgs_using_function creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
		description = "attachable_access_entity_profile created while acceptance testing"
	
	}
	
	resource "aci_access_generic" "test" {
		name 		= "%s"
		description = "access_generic created while acceptance testing"
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
	}
	
	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tDn  = "%s"
	}
	`, infraAttEntityPName, infraGenericName, tDn)
	return resource
}

func CreateAccEPGsUsingFunctionWithInValidParentDn(infraAttEntityPName, infraGenericName, tDn string) string {
	fmt.Println("=== STEP  Negative Case: testing epgs_using_function creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
		description = "attachable_access_entity_profile created while acceptance testing"
	
	}
	
	resource "aci_access_generic" "test" {
		name 		= "%s"
		description = "access_generic created while acceptance testing"
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
	}
	
	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = "${aci_access_generic.test.id}invalid"
		tDn  = "%s"	}
	`, infraAttEntityPName, infraGenericName, tDn)
	return resource
}

func CreateAccEPGsUsingFunctionConfigWithOptionalValues(infraAttEntityPName, infraGenericName, tDn string) string {
	fmt.Println("=== STEP  Basic: testing epgs_using_function creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
		description = "attachable_access_entity_profile created while acceptance testing"
	
	}
	
	resource "aci_access_generic" "test" {
		name 		= "%s"
		description = "access_generic created while acceptance testing"
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
	}
	
	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = "${aci_access_generic.test.id}"
		tDn  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_epgs_using_function"
		encap = ""instr_imedcy = "immediate"mode = "native"primary_encap = ""
	}
	`, infraAttEntityPName, infraGenericName, tDn)

	return resource
}

func CreateAccEPGsUsingFunctionRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing epgs_using_function creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_epgs_using_function" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_epgs_using_function"
		encap = ""instr_imedcy = "immediate"mode = "native"primary_encap = ""
	}
	`)

	return resource
}

func CreateAccEPGsUsingFunctionUpdatedAttr(infraAttEntityPName, infraGenericName, tDn, attribute, value string) string {
	fmt.Printf("=== STEP  testing epgs_using_function attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
		description = "attachable_access_entity_profile created while acceptance testing"
	
	}
	
	resource "aci_access_generic" "test" {
		name 		= "%s"
		description = "access_generic created while acceptance testing"
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.test.id
	}
	
	resource "aci_epgs_using_function" "test" {
		access_generic_dn  = aci_access_generic.test.id
		tDn  = "%s"
		%s = "%s"
	}
	`, infraAttEntityPName, infraGenericName, tDn, attribute, value)
	return resource
}
