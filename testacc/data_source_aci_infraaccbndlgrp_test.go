package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLeafAccessBundlePolicyGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_leaf_access_bundle_policy_group.test"
	dataSourceName := "data.aci_leaf_access_bundle_policy_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafAccessBundlePolicyGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLeafAccessBundlePolicyGroupDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafAccessBundlePolicyGroupConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "lag_t", resourceName, "lag_t"),
				),
			},
			{
				Config:      CreateAccLeafAccessBundlePolicyGroupDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLeafAccessBundlePolicyGroupDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccLeafAccessBundlePolicyGroupDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLeafAccessBundlePolicyGroupConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing leaf_access_bundle_policy_group Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s"
	}

	data "aci_leaf_access_bundle_policy_group" "test" {
		name  = aci_leaf_access_bundle_policy_group.test.name
		depends_on = [ aci_leaf_access_bundle_policy_group.test ]
	}
	`, rName)
	return resource
}

func CreateAccLeafAccessBundlePolicyGroupDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing leaf_access_bundle_policy_group Data Source with Invalid Name")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s"
	}

	data "aci_leaf_access_bundle_policy_group" "test" {
		name  = "${aci_leaf_access_bundle_policy_group.test.name}_invalid"
		depends_on = [ aci_leaf_access_bundle_policy_group.test ]
	}
	`, rName)
	return resource
}

func CreateLeafAccessBundlePolicyGroupDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_access_bundle_policy_group Data Source without ", attrName)
	rBlock := `
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_leaf_access_bundle_policy_group" "test" {
	#	name  = "%s"
		depends_on = [ aci_leaf_access_bundle_policy_group.test ]
	}
	`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLeafAccessBundlePolicyGroupDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing leaf_access_bundle_policy_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s"
	}

	data "aci_leaf_access_bundle_policy_group" "test" {
		name  = aci_leaf_access_bundle_policy_group.test.name
		%s = "%s"
		depends_on = [ aci_leaf_access_bundle_policy_group.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccLeafAccessBundlePolicyGroupDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing leaf_access_bundle_policy_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s"
		%s = "%s"
	}

	data "aci_leaf_access_bundle_policy_group" "test" {
		name  = aci_leaf_access_bundle_policy_group.test.name
		depends_on = [ aci_leaf_access_bundle_policy_group.test ]
	}
	`, rName, key, value)
	return resource
}
