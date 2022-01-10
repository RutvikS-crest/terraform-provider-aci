package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLoginDomainProviderDataSource_Basic(t *testing.T) {
	resourceName := "aci_login_domain_provider.test"
	dataSourceName := "data.aci_login_domain_provider.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	aaaDuoProviderGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoginDomainProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLoginDomainProviderDSWithoutRequired(aaaDuoProviderGroupName, rName, "duo_provider_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLoginDomainProviderDSWithoutRequired(aaaDuoProviderGroupName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLoginDomainProviderConfigDataSource(aaaDuoProviderGroupName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "duo_provider_group_dn", resourceName, "duo_provider_group_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "order", resourceName, "order"),
				),
			},
			{
				Config:      CreateAccLoginDomainProviderDataSourceUpdate(aaaDuoProviderGroupName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLoginDomainProviderDSWithInvalidParentDn(aaaDuoProviderGroupName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccLoginDomainProviderDataSourceUpdatedResource(aaaDuoProviderGroupName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLoginDomainProviderConfigDataSource(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing login_domain_provider Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}

	data "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = aci_login_domain_provider.test.name
		depends_on = [ aci_login_domain_provider.test ]
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}

func CreateLoginDomainProviderDSWithoutRequired(aaaDuoProviderGroupName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing login_domain_provider Data Source without ", attrName)
	rBlock := `
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "duo_provider_group_dn":
		rBlock += `
	data "aci_login_domain_provider" "test" {
	#	duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = aci_login_domain_provider.test.name
		depends_on = [ aci_login_domain_provider.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
	#	name  = aci_login_domain_provider.test.name
		depends_on = [ aci_login_domain_provider.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaDuoProviderGroupName, rName)
}

func CreateAccLoginDomainProviderDSWithInvalidParentDn(aaaDuoProviderGroupName, rName string) string {
	fmt.Println("=== STEP  testing login_domain_provider Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}

	data "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "${aci_login_domain_provider.test.name}_invalid"
		depends_on = [ aci_login_domain_provider.test ]
	}
	`, aaaDuoProviderGroupName, rName)
	return resource
}

func CreateAccLoginDomainProviderDataSourceUpdate(aaaDuoProviderGroupName, rName, key, value string) string {
	fmt.Println("=== STEP  testing login_domain_provider Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s"
	}

	data "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = aci_login_domain_provider.test.name
		%s = "%s"
		depends_on = [ aci_login_domain_provider.test ]
	}
	`, aaaDuoProviderGroupName, rName, key, value)
	return resource
}

func CreateAccLoginDomainProviderDataSourceUpdatedResource(aaaDuoProviderGroupName, rName, key, value string) string {
	fmt.Println("=== STEP  testing login_domain_provider Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_login_domain_provider" "test" {
		duo_provider_group_dn  = aci_duo_provider_group.test.id
		name  = aci_login_domain_provider.test.name
		depends_on = [ aci_login_domain_provider.test ]
	}
	`, aaaDuoProviderGroupName, rName, key, value)
	return resource
}
