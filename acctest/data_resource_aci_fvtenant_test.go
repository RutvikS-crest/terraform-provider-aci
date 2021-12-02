package acctest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTenantDataSource_Basic(t *testing.T) {
	resourceName := "aci_tenant.test"        
	dataSourceName := "data.aci_tenant.test" 
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccTenantDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTenantDataSource(rName), 
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"), // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),   // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),   // comparing value of parameter description in data source and resoruce
				),
			},
			{
				Config:      CreateAccTenantDSWithInvalidName(rName),                             // data source configuration with invalid application profile profile name
				ExpectError: regexp.MustCompile(`Error retriving Object: Object may not exists`), // test step expect error which should be match with defined regex
			},
		},
	})
}

func CreateAccTenantDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing tenant data source reading with giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	data "aci_tenant" "test" {
		name = "${aci_tenant.test.name}"
	}
	`, rName)
	return resource
}

func CreateAccTenantDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing tenant data source reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	data "aci_tenant" "test" {
		name = "${aci_tenant.test.name}xyz"
	}
	`, rName)
	return resource
}

func CreateAccTenantDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing tenant data source reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	data "aci_tenant" "test" {
	}
	`, rName)
	return resource
}
