package acctest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciApplicationEPGDataSource_Basic(t *testing.T) {
	var application_epg_default models.ApplicationEPG
	var application_epg_updated models.ApplicationEPG
	resourceName := "aci_application_epg.test"
	dataSourceName := "data.aci_application_epg.test"
	rName := acctest.RandString(5)
	randomParamter := acctest.RandString(10)
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationEPGDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccApplicationEPGDSWithoutApplicationProfile(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccApplicationEPGDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccApplicationEPGConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "exception_tag", resourceName, "exception_tag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "flood_on_encap", resourceName, "flood_on_encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fwd_ctrl", resourceName, "fwd_ctrl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "has_mcast_source", resourceName, "has_mcast_source"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_attr_based_epg", resourceName, "is_attr_based_epg"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_t", resourceName, "match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pc_enf_pref", resourceName, "pc_enf_pref"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref_gr_memb", resourceName, "pref_gr_memb"),
					resource.TestCheckResourceAttrPair(dataSourceName, "shutdown", resourceName, "shutdown"),
				),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedConfigDataSource(rName, randomParamter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccApplicationEPGUpdatedConfigDataSource(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated),
				),
			},
			{
				Config:      CreateAccApplicationEPGDSWithInvalidName(rName),                     // data source configuration with invalid application profile profile name
				ExpectError: regexp.MustCompile(`Error retriving Object: Object may not exists`), // test step expect error which should be match with defined regex
			},
		},
	})
}

func CreateAccApplicationEPGUpdatedConfigDataSource(rName, attribute, value string) string {
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
		%s = "%s"
	}

	data "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = aci_application_epg.test.name
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccApplicationEPGConfigDataSource(rName string) string {
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	data "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = aci_application_epg.test.name
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccApplicationEPGDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing Application EPG reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	data "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "${aci_application_epg.test.name}abc"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccApplicationEPGDSWithoutApplicationProfile(rName string) string {
	fmt.Println("=== STEP  Basic: testing Application EPG reading without giving application_profile_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	data "aci_application_epg" "test" {
		name = "%s"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccApplicationEPGDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing Application EPG reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	data "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
	}
	`, rName, rName, rName)
	return resource
}
