package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciEndpointIpAgingProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_endpoint_ip_aging_profile.test"
	dataSourceName := "data.aci_endpoint_ip_aging_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciIPAgingPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateEndpointIpAgingProfileDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndpointIpAgingProfileConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
				),
			},
			{
				Config:      CreateAccEndpointIpAgingProfileDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccEndpointIpAgingProfileDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccEndpointIpAgingProfileDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccEndpointIpAgingProfileConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s"
	}

	data "aci_endpoint_ip_aging_profile" "test" {
	
		name  = aci_endpoint_ip_aging_profile.test.name
		depends_on = [ aci_endpoint_ip_aging_profile.test ]
	}
	`, rName)
	return resource
}

func CreateEndpointIpAgingProfileDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_ip_aging_profile Data Source without ", attrName)
	rBlock := `
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_endpoint_ip_aging_profile" "test" {
	
	#	name  = aci_endpoint_ip_aging_profile.test.name
		depends_on = [ aci_endpoint_ip_aging_profile.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccEndpointIpAgingProfileDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s"
	}

	data "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "${aci_endpoint_ip_aging_profile.test.name}_invalid"
		name  = aci_endpoint_ip_aging_profile.test.name
		depends_on = [ aci_endpoint_ip_aging_profile.test ]
	}
	`, rName)
	return resource
}

func CreateAccEndpointIpAgingProfileDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s"
	}

	data "aci_endpoint_ip_aging_profile" "test" {
	
		name  = aci_endpoint_ip_aging_profile.test.name
		%s = "%s"
		depends_on = [ aci_endpoint_ip_aging_profile.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccEndpointIpAgingProfileDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_endpoint_ip_aging_profile" "test" {
	
		name  = aci_endpoint_ip_aging_profile.test.name
		depends_on = [ aci_endpoint_ip_aging_profile.test ]
	}
	`, rName, key, value)
	return resource
}
