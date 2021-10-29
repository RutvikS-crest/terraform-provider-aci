package aci

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

//TODO: Multiple resource

func TestAccAciApplicationProfile_Basic(t *testing.T) {
	var application_profile models.ApplicationProfile
	resourceName := "aci_application_profile.test"
	rName := acctest.RandString(5)
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccApplicationProfileWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccApplicationProfileWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccApplicationProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile),
					resource.TestCheckResourceAttr(resourceName, "description", "from terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ap"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccApplicationProfileConfigUpdatedName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of ap-%s failed validation for value '%s'", longrName, longrName)),
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
				Config: CreateAccApplicationProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileUpdatedAttr(rName, "description", "updated description for terraform test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description for terraform test"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileUpdatedAttr(rName, "annotation", "updated_annotation_for_terraform_test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "updated_annotation_for_terraform_test"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
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
			{
				Config: CreateAccApplicationProfileUpdatedAttr(rName, "name_alias", "updated_name_alias_for_terraform_test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "updated_name_alias_for_terraform_test"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_updated),
				),
			},
		},
	})
}

func TestAccApplicationProfile_NegativeCases(t *testing.T) {
	var application_profile_default models.ApplicationProfile
	resourceName := "aci_application_profile.test"
	rName := acctest.RandString(5)
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomPrio := acctest.RandString(6)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccApplicationProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			//TODO: Invalid parent dn
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "prio", randomPrio),
				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
			},
			//TODO: randomly generate parameter and value
			{
				Config: CreateAccApplicationProfileConfig(rName),
			},
		},
	})
}

func TestAccApplicationProfile_relMonPol(t *testing.T) {
	var application_profile_default models.ApplicationProfile
	var application_profile_relMonPol1 models.ApplicationProfile
	var application_profile_relMonPol2 models.ApplicationProfile
	resourceName := "aci_application_profile.test"
	rName := acctest.RandString(5)
	monPolName1 := acctest.RandString(5)
	monPolName2 := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccApplicationProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccApplicationProfileUpdatedMonPol(rName, monPolName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_relMonPol1),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, monPolName1)),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_relMonPol1),
				),
			},
			{
				Config: CreateAccApplicationProfileUpdatedMonPol(rName, monPolName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_relMonPol2),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, monPolName2)),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default, &application_profile_relMonPol2),
				),
			},
			// TODO: add dn which is not allowed
			{
				Config: CreateAccApplicationProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""),
				),
			},
		},
	})
}

func testAccCheckAciApplicationProfileIdEqual(ap1, ap2 *models.ApplicationProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ap1.DistinguishedName != ap2.DistinguishedName {
			return fmt.Errorf("ApplicationProfile DNs are not equal")
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

func CreateAccApplicationProfileUpdatedMonPol(rName, monPolName string) string {
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
