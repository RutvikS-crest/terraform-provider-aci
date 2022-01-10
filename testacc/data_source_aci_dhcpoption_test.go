package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciDHCPOptionDataSource_Basic(t *testing.T) {
	resourceName := "aci_dhcp_option.test"
	dataSourceName := "data.aci_dhcp_option.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	dhcpOptionPolName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateDHCPOptionDSWithoutRequired(fvTenantName, dhcpOptionPolName, rName, "dhcp_option_policy_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateDHCPOptionDSWithoutRequired(fvTenantName, dhcpOptionPolName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDHCPOptionConfigDataSource(fvTenantName, dhcpOptionPolName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "dhcp_option_policy_dn", resourceName, "dhcp_option_policy_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "data", resourceName, "data"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dhcp_option_id", resourceName, "dhcp_option_id"),
				),
			},
			{
				Config:      CreateAccDHCPOptionDataSourceUpdate(fvTenantName, dhcpOptionPolName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccDHCPOptionDSWithInvalidParentDn(fvTenantName, dhcpOptionPolName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccDHCPOptionDataSourceUpdatedResource(fvTenantName, dhcpOptionPolName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccDHCPOptionConfigDataSource(fvTenantName, dhcpOptionPolName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_option Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = "%s"
	}

	data "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = aci_dhcp_option.test.name
		depends_on = [ aci_dhcp_option.test ]
	}
	`, fvTenantName, dhcpOptionPolName, rName)
	return resource
}

func CreateDHCPOptionDSWithoutRequired(fvTenantName, dhcpOptionPolName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing dhcp_option Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "dhcp_option_policy_dn":
		rBlock += `
	data "aci_dhcp_option" "test" {
	#	dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = aci_dhcp_option.test.name
		depends_on = [ aci_dhcp_option.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
	#	name  = aci_dhcp_option.test.name
		depends_on = [ aci_dhcp_option.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, dhcpOptionPolName, rName)
}

func CreateAccDHCPOptionDSWithInvalidParentDn(fvTenantName, dhcpOptionPolName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_option Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = "%s"
	}

	data "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = "${aci_dhcp_option.test.name}_invalid"
		depends_on = [ aci_dhcp_option.test ]
	}
	`, fvTenantName, dhcpOptionPolName, rName)
	return resource
}

func CreateAccDHCPOptionDataSourceUpdate(fvTenantName, dhcpOptionPolName, rName, key, value string) string {
	fmt.Println("=== STEP  testing dhcp_option Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = "%s"
	}

	data "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = aci_dhcp_option.test.name
		%s = "%s"
		depends_on = [ aci_dhcp_option.test ]
	}
	`, fvTenantName, dhcpOptionPolName, rName, key, value)
	return resource
}

func CreateAccDHCPOptionDataSourceUpdatedResource(fvTenantName, dhcpOptionPolName, rName, key, value string) string {
	fmt.Println("=== STEP  testing dhcp_option Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_dhcp_option" "test" {
		dhcp_option_policy_dn  = aci_dhcp_option_policy.test.id
		name  = aci_dhcp_option.test.name
		depends_on = [ aci_dhcp_option.test ]
	}
	`, fvTenantName, dhcpOptionPolName, rName, key, value)
	return resource
}
