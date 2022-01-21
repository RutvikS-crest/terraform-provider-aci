package testacc

// import (
// 	"fmt"
// 	"regexp"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// )




  

// func TestAccAciLDAPGroupMaprulerefDataSource_Basic(t *testing.T) {
// 	resourceName := "aci_ldap_group_mapruleref.test"
// 	dataSourceName := "data.aci_ldap_group_mapruleref.test"
// 	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
// 	randomValue := acctest.RandString(10)
// 	rName := makeTestVariable(acctest.RandString(5))
	
// 	aaaLdapGroupMapName := makeTestVariable(acctest.RandString(5))
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:	  func(){ testAccPreCheck(t) },
// 		ProviderFactories:    testAccProviders,
// 		CheckDestroy: testAccCheckAciLDAPGroupMaprulerefDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config:      CreateLDAPGroupMaprulerefDSWithoutRequired(aaaLdapGroupMapName, rName,"ldap_group_map_dn"),
// 				ExpectError: regexp.MustCompile(`Missing required argument`),
// 			},
// 			{
// 				Config:      CreateLDAPGroupMaprulerefDSWithoutRequired(aaaLdapGroupMapName, rName, "name"),
// 				ExpectError: regexp.MustCompile(`Missing required argument`),
// 			},
// 			{
// 				Config: CreateAccLDAPGroupMaprulerefConfigDataSource(aaaLdapGroupMapName, rName),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttrPair(dataSourceName, "ldap_group_map_dn", resourceName, "ldap_group_map_dn",),
// 					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
// 					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
// 					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
// 					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					
// 				),
// 			},
// 			{
// 				Config:      CreateAccLDAPGroupMaprulerefDataSourceUpdate(aaaLdapGroupMapName, rName, randomParameter, randomValue),
// 				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
// 			},
			
// 			{
// 				Config:      CreateAccLDAPGroupMaprulerefDSWithInvalidParentDn(aaaLdapGroupMapName, rName),
// 				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
// 			},
			
// 			{
// 				Config: CreateAccLDAPGroupMaprulerefDataSourceUpdatedResource(aaaLdapGroupMapName, rName, "annotation", "orchestrator:terraform-testacc"),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
// 				),
// 			},
// 		},
// 	})
// }


// func CreateAccLDAPGroupMaprulerefConfigDataSource(aaaLdapGroupMapName, rName string) string {
// 	fmt.Println("=== STEP  testing ldap_group_mapruleref Data Source with required arguments only")
// 	resource := fmt.Sprintf(`
	
// 	resource "aci_ldap_group_map" "test" {
// 		name 		= "%s"
	
// 	}
	
// 	resource "aci_ldap_group_mapruleref" "test" {
// 		ldap_group_map_dn  = aci_ldap_group_map.test.id
// 		name  = "%s"
// 	}

// 	data "aci_ldap_group_mapruleref" "test" {
// 		ldap_group_map_dn  = aci_ldap_group_map.test.id
// 		name  = aci_ldap_group_mapruleref.test.name
// 		depends_on = [ aci_ldap_group_mapruleref.test ]
// 	}
// 	`, aaaLdapGroupMapName, rName)
// 	return resource
// }

// func CreateLDAPGroupMaprulerefDSWithoutRequired(aaaLdapGroupMapName, rName, attrName string) string {
// 	fmt.Println("=== STEP  Basic: testing ldap_group_mapruleref Data Source without ",attrName)
// 	rBlock := `
	
// 	resource "aci_ldap_group_map" "test" {
// 		name 		= "%s"
	
// 	}
	
// 	resource "aci_ldap_group_mapruleref" "test" {
// 		ldap_group_map_dn  = aci_ldap_group_map.test.id
// 		name  = "%s"
// 	}
// 	`
// 	switch attrName {
// 	case "ldap_group_map_dn":
// 		rBlock += `
// 	data "aci_ldap_group_mapruleref" "test" {
// 	#	ldap_group_map_dn  = aci_ldap_group_map.test.id
// 		name  = aci_ldap_group_mapruleref.test.name
// 		depends_on = [ aci_ldap_group_mapruleref.test ]
// 	}
// 		`
// 	case "name":
// 		rBlock += `
// 	data "aci_ldap_group_mapruleref" "test" {
// 		ldap_group_map_dn  = aci_ldap_group_map.test.id
// 	#	name  = aci_ldap_group_mapruleref.test.name
// 		depends_on = [ aci_ldap_group_mapruleref.test ]
// 	}
// 		`
// 	}
// 	return fmt.Sprintf(rBlock,aaaLdapGroupMapName, rName)
// }

// func CreateAccLDAPGroupMaprulerefDSWithInvalidParentDn(aaaLdapGroupMapName, rName string) string {
// 	fmt.Println("=== STEP  testing ldap_group_mapruleref Data Source with Invalid Parent Dn")
// 	resource := fmt.Sprintf(`
	
// 	resource "aci_ldap_group_map" "test" {
// 		name 		= "%s"
	
// 	}
	
// 	resource "aci_ldap_group_mapruleref" "test" {
// 		ldap_group_map_dn  = aci_ldap_group_map.test.id
// 		name  = "%s"
// 	}

// 	data "aci_ldap_group_mapruleref" "test" {
// 		ldap_group_map_dn  = aci_ldap_group_map.test.id
// 		name  = "${aci_ldap_group_mapruleref.test.name}_invalid"
// 		depends_on = [ aci_ldap_group_mapruleref.test ]
// 	}
// 	`, aaaLdapGroupMapName, rName)
// 	return resource
// }

// func CreateAccLDAPGroupMaprulerefDataSourceUpdate(aaaLdapGroupMapName, rName, key, value string) string {
// 	fmt.Println("=== STEP  testing ldap_group_mapruleref Data Source with random attribute")
// 	resource := fmt.Sprintf(`
	
// 	resource "aci_ldap_group_map" "test" {
// 		name 		= "%s"
	
// 	}
	
// 	resource "aci_ldap_group_mapruleref" "test" {
// 		ldap_group_map_dn  = aci_ldap_group_map.test.id
// 		name  = "%s"
// 	}

// 	data "aci_ldap_group_mapruleref" "test" {
// 		ldap_group_map_dn  = aci_ldap_group_map.test.id
// 		name  = aci_ldap_group_mapruleref.test.name
// 		%s = "%s"
// 		depends_on = [ aci_ldap_group_mapruleref.test ]
// 	}
// 	`, aaaLdapGroupMapName, rName,key,value)
// 	return resource
// }

// func CreateAccLDAPGroupMaprulerefDataSourceUpdatedResource(aaaLdapGroupMapName, rName, key, value string) string {
// 	fmt.Println("=== STEP  testing ldap_group_mapruleref Data Source with updated resource")
// 	resource := fmt.Sprintf(`
	
// 	resource "aci_ldap_group_map" "test" {
// 		name 		= "%s"
	
// 	}
	
// 	resource "aci_ldap_group_mapruleref" "test" {
// 		ldap_group_map_dn  = aci_ldap_group_map.test.id
// 		name  = "%s"
// 		%s = "%s"
// 	}

// 	data "aci_ldap_group_mapruleref" "test" {
// 		ldap_group_map_dn  = aci_ldap_group_map.test.id
// 		name  = aci_ldap_group_mapruleref.test.name
// 		depends_on = [ aci_ldap_group_mapruleref.test ]
// 	}
// 	`, aaaLdapGroupMapName, rName,key,value)
// 	return resource
// }