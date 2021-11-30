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

func TestAccAciApplicationProfile_Basic(t *testing.T) {
	var application_profile_default models.ApplicationProfile // variable of ApplicationProfile's model type would be useful to compare ids
	var application_profile_updated models.ApplicationProfile // variable of ApplicationProfile's model type would be useful to compare ids
	resourceName := "aci_application_profile.test"            // declared resource on which all operation would be performed
	rName := acctest.RandString(5)                            // randomly created string of 5 alphanumeric characters' for resource name
	rOther := acctest.RandString(5)                           // randomly created string of 5 alphanumeric characters' for another resource name
	prOther := acctest.RandString(5)                          // randomly created string of 5 alphanumeric characters' for another parent resource name
	longrName := acctest.RandString(65)                       // randomly created string of 65 alphanumeric characters' for negative resource name test case
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				// terraform will try to create application profile without required arguement tenant_dn
				Config:      CreateAccApplicationProfileWithoutTenant(rName), // configuration to check creation of application profile without tenant
				ExpectError: regexp.MustCompile(`Missing required argument`), // test step expect error which should be match with defined regex
			},
			{
				// terraform will try to create application profile without required arguement name
				Config:      CreateAccApplicationProfileWithoutName(rName), // configuration to check creation of application profile without tenant
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				// step terraform will create application profile with only required arguements i.e. name and tenant_dn
				Config: CreateAccApplicationProfileConfig(rName), // configuration to create application profile with required fields only
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default), // this function will check whether any resource is exist or not in state file with given resource name
					// now will compare value of all attributes with default for given resource
					resource.TestCheckResourceAttr(resourceName, "description", ""), // no default value for description so comparing with ""
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),  // no default value for name_alias so comparing with ""
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),       // comparing with default value of annotation
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),                        // comparing with default value of prio
					resource.TestCheckResourceAttr(resourceName, "name", rName),                                // comparing application profile name with randomly created name
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)), // comparing tenant_dn with its proper format
				),
			},
			{
				// this step will import state of particular resource and will test state file with configuration file
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// in this step all optional attribute expect realational attribute are given for the same resource and then compared
				Config: CreateAccApplicationProfileConfigWithOptionalValues(rName), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "from terraform"), // comparing description with value which is given in configuration
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ap"),         // comparing name_alias with value which is given in configuration
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag"), // comparing annotation with value which is given in configuration
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),    // comparing prio with value which is given in configuration
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccApplicationProfileConfigUpdatedName(rName, longrName), // passing invalid name for application profile
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of ap-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccApplicationProfileConfigWithParentAndName(rName, rOther), // creating resource with same parent name and different resource name
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rOther),                                            // comparing name attribute of applicaiton profile
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),              // comparing tenant_dn attribute of application profile
					testAccCheckAciApplicationProfileIdNotEqual(&application_profile_default, &application_profile_updated), // checking whether id or dn of both resource are different because name changed and terraform need to create another resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileConfig(rName), // creating resource with required parameters only
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileConfigWithParentAndName(prOther, rName), // creating resource with same name but different parent resource name
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", prOther)),
					testAccCheckAciApplicationProfileIdNotEqual(&application_profile_default, &application_profile_updated), // checking whether id or dn of both resource are different because tenant_dn changed and terraform need to create another resource
				),
			},
		},
	})
}

func TestAccApplicationProfile_Update(t *testing.T) {
	var application_profile_default models.ApplicationProfile
	var application_profile_updated models.ApplicationProfile
	resourceName := "aci_application_profile.test"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccApplicationProfileConfig(rName), // creating application profile with required arguements only
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default),
				),
			},
			{
				// this step will import state of particular resource and will test state file with configuration file
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// there are various value of prio parameter is possible so checking prio for each value
			{
				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level2"), // updating only prio parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"), // checking value updated value of prio parameter
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplicationProfile_NegativeCases(t *testing.T) {
	resourceName := "aci_application_profile.test"
	rName := acctest.RandString(5)
	longDescAnnotation := acctest.RandString(129)                                     // creating random string of 129 characters
	longNameAlias := acctest.RandString(64)                                           // creating random string of 64 characters
	randomPrio := acctest.RandString(6)                                               // creating random string of 6 characters
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz") // creating random string of 5 characters (to give as random parameter)
	randomValue := acctest.RandString(5)                                              // creating random string of 5 characters (to give as random value of random parameter)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccApplicationProfileConfig(rName), // creating application profile with required arguements only
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccApplicationProfileWithInValidTenantDn(rName),                       // checking application profile creation with invalid tenant_dn value
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvAp (.)+`), // test step expect error which should be match with defined regex
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "description", longDescAnnotation), // checking application profile creation with invalid description value
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "annotation", longDescAnnotation), // checking application profile creation with invalid annotation value
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "name_alias", longNameAlias), // checking application profile creation with invalid name_alias value
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "prio", randomPrio), // checking application profile creation with invalid prio value
				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, randomParameter, randomValue), // checking application profile creation with randomly created parameter and value
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccApplicationProfileConfig(rName), // creating application profile with required arguements only
			},
		},
	})
}

func TestAccApplicationProfile_reltionalParameters(t *testing.T) {
	var application_profile_default models.ApplicationProfile
	var application_profile_rel1 models.ApplicationProfile
	var application_profile_rel2 models.ApplicationProfile
	resourceName := "aci_application_profile.test"
	rName := acctest.RandString(5)
	monPolName1 := acctest.RandString(5) // randomly created name for relational resoruce
	monPolName2 := acctest.RandString(5) // randomly created name for relational resoruce
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccApplicationProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default), // creating application profile with required arguements only
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""),       // checking value of relation_fv_rs_ap_mon_pol parameter for given configuration
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileConfigInitial(rName, monPolName1), // creating application profile with relation_fv_rs_ap_mon_pol parameter for the first randomly generated name
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_rel1),                                                  // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, monPolName1)), // checking relation by comparing values
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_rel1),                                 // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				Config: CreateAccApplicationProfileConfigFinal(rName, monPolName2), // creating application profile with relation_fv_rs_ap_mon_pol parameter for the second randomly generated name (to verify update operation)
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_rel2),                                                  // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, monPolName2)), // checking relation by comparing values
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_rel2),
				),
			},
			{
				Config: CreateAccApplicationProfileConfig(rName), // this configuration will remove relation
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""), // checking removal of relation
				),
			},
		},
	})
}

// func TestAccApplicationProfile_MultipleCreateDelete(t *testing.T) {
// 	for i := 0; i < 20; i++ {
// 		resource.Test(t, resource.TestCase{
// 			PreCheck:     func() { testAccPreCheck(t) },
// 			Providers:    testAccProviders,
// 			CheckDestroy: testAccCheckAciApplicationProfileDestroy,
// 			Steps: []resource.TestStep{
// 				{
// 					Config: CreateAccApplicationProfileConfig(fmt.Sprintf("r%d", i)),
// 				},
// 			},
// 		})
// 	}
// }

func testAccCheckAciApplicationProfileExists(name string, application_profile *models.ApplicationProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Application Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Application Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		application_profileFound := models.ApplicationProfileFromContainer(cont)
		if application_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Application Profile %s not found", rs.Primary.ID)
		}
		*application_profile = *application_profileFound
		return nil
	}
}

func testAccCheckAciApplicationProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_application_profile" {
			cont, err := client.Get(rs.Primary.ID)
			application_profile := models.ApplicationProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Application Profile %s Still exists", application_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciApplicationProfileIdEqual(ap1, ap2 *models.ApplicationProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ap1.DistinguishedName != ap2.DistinguishedName {
			return fmt.Errorf("ApplicationProfile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciApplicationProfileIdNotEqual(ap1, ap2 *models.ApplicationProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ap1.DistinguishedName == ap2.DistinguishedName {
			return fmt.Errorf("ApplicationProfile DNs are equal")
		}
		return nil
	}
}

func CreateAccApplicationProfileWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile creation without creating tenant")
	resource := fmt.Sprintf(`
	resource "aci_application_profile" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}

func CreateAccApplicationProfileWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile creation without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccApplicationProfileConfigWithParentAndName(prName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing applicationProfile creation with tenant name %s name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, prName, rName)
	return resource
}

func CreateAccApplicationProfileConfig(rName string) string {
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccApplicationProfileWithInValidTenantDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing applicationProfile creation with invalid tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_vrf.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccApplicationProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		annotation = "tag"
		description = "from terraform"
		name_alias = "test_ap"
		prio = "level1"
	}
	`, rName, rName)
	return resource
}

func CreateAccApplicationProfileConfigInitial(rName, monPolName string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile creation with initial relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_monitoring_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		relation_fv_rs_ap_mon_pol = aci_monitoring_policy.test.id
	}
	`, rName, monPolName, rName)
	return resource
}

func CreateAccApplicationProfileConfigFinal(rName, monPolName string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile creation with final relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_monitoring_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		relation_fv_rs_ap_mon_pol = aci_monitoring_policy.test.id
	}
	`, rName, monPolName, rName)
	return resource
}

func CreateAccApplicationProfileConfigUpdatedName(rName, longrName string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile creation with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, longrName)
	return resource
}

func CreateAccApplicationProfileConfigWithChangedName(rName1, rName2 string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile creation with changed name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName1, rName2)
	return resource
}

func CreateAccApplicationProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}
