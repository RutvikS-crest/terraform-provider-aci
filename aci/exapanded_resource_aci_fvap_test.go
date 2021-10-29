package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciApplicationProfile_Basic(t *testing.T) {
	var application_profile models.ApplicationProfile
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
		},
	})
}

func TestAccApplicationProfile_description(t *testing.T) {
	var application_profile_default_description models.ApplicationProfile
	var application_profile_updated_description models.ApplicationProfile
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
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default_description),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				Config: CreateAccApplicationProfileUpdatedDescription(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated_description),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description for terraform test"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default_description, &application_profile_updated_description),
				),
			},
		},
	})
}

func TestAccApplicationProfile_annotation(t *testing.T) {
	var application_profile_default_annotation models.ApplicationProfile
	var application_profile_updated_annotation models.ApplicationProfile
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
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default_annotation),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
				),
			},
			{
				Config: CreateAccApplicationProfileUpdatedAnnotation(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated_annotation),
					resource.TestCheckResourceAttr(resourceName, "annotation", "updated_annotation_for_terraform_test"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default_annotation, &application_profile_updated_annotation),
				),
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

func TestAccApplicationProfile_prio(t *testing.T) {
	var application_profile_default_prio models.ApplicationProfile
	var application_profile_l1_prio models.ApplicationProfile
	var application_profile_l2_prio models.ApplicationProfile
	var application_profile_l3_prio models.ApplicationProfile
	var application_profile_l4_prio models.ApplicationProfile
	var application_profile_l5_prio models.ApplicationProfile
	var application_profile_l6_prio models.ApplicationProfile
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
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default_prio),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
				),
			},
			{
				Config: CreateAccApplicationProfileUpdatedPrio(rName, "level1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_l1_prio),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default_prio, &application_profile_l1_prio),
				),
			},
			{
				Config: CreateAccApplicationProfileUpdatedPrio(rName, "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_l2_prio),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default_prio, &application_profile_l2_prio),
				),
			},
			{
				Config: CreateAccApplicationProfileUpdatedPrio(rName, "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_l3_prio),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default_prio, &application_profile_l3_prio),
				),
			},
			{
				Config: CreateAccApplicationProfileUpdatedPrio(rName, "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_l4_prio),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default_prio, &application_profile_l4_prio),
				),
			},
			{
				Config: CreateAccApplicationProfileUpdatedPrio(rName, "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_l5_prio),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default_prio, &application_profile_l5_prio),
				),
			},
			{
				Config: CreateAccApplicationProfileUpdatedPrio(rName, "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_l6_prio),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default_prio, &application_profile_l6_prio),
				),
			},
		},
	})
}

func TestAccApplicationProfile_nameAlias(t *testing.T) {
	var application_profile_default_nameAlias models.ApplicationProfile
	var application_profile_updated_nameAlias models.ApplicationProfile
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
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_default_nameAlias),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccApplicationProfileUpdatedNameAlias(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists(resourceName, &application_profile_updated_nameAlias),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "updated_name_alias_for_terraform_test"),
					testAccCheckAciApplicationProfileIdEqual(&application_profile_default_nameAlias, &application_profile_updated_nameAlias),
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

func CreateAccApplicationProfileUpdatedPrio(rName, prio string) string {
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		prio = "%s"
	}
	`, rName, rName, prio)
	return resource
}

func CreateAccApplicationProfileUpdatedNameAlias(rName string) string {
	resource := fmt.Sprintf(`
 	resource "aci_tenant" "test" {
		name = "%s"
 	}
 
 	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		name_alias = "updated_name_alias_for_terraform_test"
 	}
 `, rName, rName)
	return resource
}

func CreateAccApplicationProfileUpdatedAnnotation(rName string) string {
	resource := fmt.Sprintf(`
 	resource "aci_tenant" "test" {
		name = "%s"
 	}
 
 	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		annotation = "updated_annotation_for_terraform_test"
 	}
 `, rName, rName)
	return resource
}

func CreateAccApplicationProfileUpdatedDescription(rName string) string {
	resource := fmt.Sprintf(`
 	resource "aci_tenant" "test" {
		name = "%s"
 	}
 
 	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		description = "updated description for terraform test"
 	}
 `, rName, rName)
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
