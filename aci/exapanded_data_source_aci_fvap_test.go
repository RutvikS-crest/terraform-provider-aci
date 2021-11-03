package aci

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciApplicationProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_application_profile.test"        // defining name of resource
	dataSourceName := "data.aci_application_profile.test" // defining name of data source
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccApplicationProfileDSWithoutTenant(rName), // creating data source for application profile without required arguement tenant_dn
				ExpectError: regexp.MustCompile(`Missing required argument`),   // test step expect error which should be match with defined regex
			},
			{
				Config:      CreateAccApplicationProfileDSWithoutName(rName), // creating data source for application profile without required arguement name
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccApplicationProfileConfigDataSource(rName), // creating data source with required arguements from the resource
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"), // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),   // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),   // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),               // comparing value of parameter description in data source and resoruce
				),
			},
			{
				Config:      CreateAccApplicationProfileDSWithInvalidName(rName),                 // data source configuration with invalid application profile profile name
				ExpectError: regexp.MustCompile(`Error retriving Object: Object may not exists`), // test step expect error which should be match with defined regex
			},
		},
	})
}

func CreateAccApplicationProfileConfigDataSource(rName string) string {
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = aci_application_profile.test.name
	}
	`, rName, rName)
	return resource
}

func CreateAccApplicationProfileDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "${aci_application_profile.test.name}xyz"
	}
	`, rName, rName)
	return resource
}

func CreateAccApplicationProfileDSWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile reading without giving tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_application_profile" "test" {
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccApplicationProfileDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName, rName)
	return resource
}