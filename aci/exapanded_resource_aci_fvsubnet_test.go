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

func TestAccAciApplicationProfile_Basic(t *testing.T) {
	var subnet_initial models.Subnet
	var subnet_updated models.Subnet
	var subnet_using_tenant models.Subnet
	var subnet_using_tenant_updated models.Subnet
	resourceName := "aci_subnet.test"
	rName := acctest.RandString(5)
	prefixLen := acctest.RandIntRange(0, 31)
	ip, _ := acctest.RandIpAddress(string(prefixLen))
	longerName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSubnetWithoutParentDn(ip),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSubnetWithoutIP(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSubnetConfigParentEpg(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_initial),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetConfigParentBridgeDomain(rName,ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					testAccCheckAciSubnetIdNotEqual(&subnet_initial, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetConfigParentL2OutExtEpg(rName,ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					testAccCheckAciSubnetIdNotEqual(&subnet_initial, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetConfig(rName,ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_using_tenant),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "preferred", "no"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s", rName)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// in this step all optional attribute expect realational attribute are given for the same resource and then compared
				Config: CreateAccSubnetConfigWithOptionalValues(rName), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &application_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "from terraform"), // comparing description with value which is given in configuration
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ap"),         // comparing name_alias with value which is given in configuration
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag"), // comparing annotation with value which is given in configuration
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),    // comparing prio with value which is given in configuration
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSubnetConfigUpdatedName(rName, longerName), // passing invalid name for application profile
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of ap-%s failed validation for value '%s'", longerName, longerName)),
			},
		},
	})
}

// func TestAccApplicationProfile_Update(t *testing.T) {
// 	var application_profile_default models.ApplicationProfile
// 	var application_profile_updated models.ApplicationProfile
// 	resourceName := "aci_application_profile.test"
// 	rName := acctest.RandString(5)
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckAciSubnetDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: CreateAccSubnetConfig(rName), // creating application profile with required arguements only
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_default),
// 				),
// 			},
// 			{
// 				// this step will import state of particular resource and will test state file with configuration file
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: CreateAccApplicationProfileUpdatedAttr(rName, "description", "updated description for terraform test"), // updating only description parameter
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_updated),                               // checking whether resource is exist or not in state file
// 					resource.TestCheckResourceAttr(resourceName, "description", "updated description for terraform test"), // checking value updated value of description parameter
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_updated),  // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: CreateAccApplicationProfileUpdatedAttr(rName, "annotation", "updated_annotation_for_terraform_test"), // updating only description parameter
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_updated),
// 					resource.TestCheckResourceAttr(resourceName, "annotation", "updated_annotation_for_terraform_test"), // checking value updated value of description parameter
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_updated),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			// there are various value of prio parameter is possible so checking prio for each value
// 			{
// 				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level1"), // updating only prio parameter
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_updated),
// 					resource.TestCheckResourceAttr(resourceName, "prio", "level1"), // checking value updated value of prio parameter
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_updated),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level2"), // updating only prio parameter
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_updated),
// 					resource.TestCheckResourceAttr(resourceName, "prio", "level2"), // checking value updated value of prio parameter
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_updated),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level3"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_updated),
// 					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_updated),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level4"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_updated),
// 					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_updated),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level5"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_updated),
// 					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_updated),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: CreateAccApplicationProfileUpdatedAttr(rName, "prio", "level6"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_updated),
// 					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_updated),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: CreateAccApplicationProfileUpdatedAttr(rName, "name_alias", "updated_name_alias_for_terraform_test"), // updating only name_alias parameter
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_updated),
// 					resource.TestCheckResourceAttr(resourceName, "name_alias", "updated_name_alias_for_terraform_test"), // checking value updated value of prio parameter
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_updated),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccApplicationProfile_NegativeCases(t *testing.T) {
// 	var application_profile_default models.ApplicationProfile
// 	resourceName := "aci_application_profile.test"
// 	rName := acctest.RandString(5)
// 	longDescAnnotation := acctest.RandString(129) // creating random string of 129 characters
// 	longNameAlias := acctest.RandString(64)       // creating random string of 64 characters
// 	randomPrio := acctest.RandString(6)           // creating random string of 6 characters
// 	randomParameter := acctest.RandString(5)      // creating random string of 5 characters (to give as random parameter)
// 	randomValue := acctest.RandString(5)          // creating random string of 5 characters (to give as random value of random parameter)
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckAciSubnetDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: CreateAccSubnetConfig(rName), // creating application profile with required arguements only
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_default),
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config:      CreateAccApplicationProfileWithInValidTenantDn(rName),                       // checking application profile creation with invalid tenant_dn value
// 				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvAp (.)+`), // test step expect error which should be match with defined regex
// 			},
// 			{
// 				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "description", longDescAnnotation), // checking application profile creation with invalid description value
// 				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
// 			},
// 			{
// 				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "annotation", longDescAnnotation), // checking application profile creation with invalid annotation value
// 				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
// 			},
// 			{
// 				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "name_alias", longNameAlias), // checking application profile creation with invalid name_alias value
// 				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
// 			},
// 			{
// 				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "prio", randomPrio), // checking application profile creation with invalid prio value
// 				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
// 			},
// 			{
// 				Config:      CreateAccApplicationProfileUpdatedAttr(rName, randomParameter, randomValue), // checking application profile creation with randomly created parameter and value
// 				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
// 			},
// 			{
// 				Config: CreateAccSubnetConfig(rName), // creating application profile with required arguements only
// 			},
// 		},
// 	})
// }

// func TestAccApplicationProfile_relMonPol(t *testing.T) {
// 	var application_profile_default models.ApplicationProfile
// 	var application_profile_relMonPol1 models.ApplicationProfile
// 	var application_profile_relMonPol2 models.ApplicationProfile
// 	resourceName := "aci_application_profile.test"
// 	rName := acctest.RandString(5)
// 	monPolName1 := acctest.RandString(5) // randomly created name for relational resoruce
// 	monPolName2 := acctest.RandString(5) // randomly created name for relational resoruce
// 	//TODO: Invalid relation check
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckAciSubnetDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: CreateAccSubnetConfig(rName),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_default),       // creating application profile with required arguements only
// 					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""), // checking value of relation_fv_rs_ap_mon_pol parameter for given configuration
// 				),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: CreateAccApplicationProfileUpdatedMonPol(rName, monPolName1), // creating application profile with relation_fv_rs_ap_mon_pol parameter for the first randomly generated name
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_relMonPol1),                                                        // checking whether resource is exist or not in state file
// 					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, monPolName1)), // checking relation by comparing values
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_relMonPol1),                           // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
// 				),
// 			},
// 			{
// 				Config: CreateAccApplicationProfileUpdatedMonPol(rName, monPolName2), // creating application profile with relation_fv_rs_ap_mon_pol parameter for the second randomly generated name (to verify update operation)
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_relMonPol2),                                                        // checking whether resource is exist or not in state file
// 					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, monPolName2)), // checking relation by comparing values
// 					testAccCheckAciSubnetIdEqual(&application_profile_default, &application_profile_relMonPol2),
// 				),
// 			},
// 			{
// 				Config: CreateAccSubnetConfig(rName), // this configuration will remove relation
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckAciSubnetExists(resourceName, &application_profile_default),
// 					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ap_mon_pol", ""), // checking removal of relation
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccApplicationProfile_MultipleCreateDelete(t *testing.T) {
// 	for i := 0; i < 20; i++ {
// 		resource.Test(t, resource.TestCase{
// 			PreCheck:     func() { testAccPreCheck(t) },
// 			Providers:    testAccProviders,
// 			CheckDestroy: testAccCheckAciSubnetDestroy,
// 			Steps: []resource.TestStep{
// 				{
// 					Config: CreateAccSubnetConfig(fmt.Sprintf("r%d", i)),
// 				},
// 			},
// 		})
// 	}
// }

func testAccCheckAciSubnetExists(name string, subnet *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		subnetFound := models.SubnetFromContainer(cont)
		if subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Subnet %s not found", rs.Primary.ID)
		}
		*subnet = *subnetFound
		return nil
	}
}

func testAccCheckAciSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			subnet := models.SubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Subnet %s Still exists", subnet.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSubnetIdEqual(sn1, sn2 *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if sn1.DistinguishedName != sn2.DistinguishedName {
			return fmt.Errorf("Subnet DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSubnetIdNotEqual(sn1, sn2 *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if sn1.DistinguishedName == sn2.DistinguishedName {
			return fmt.Errorf("Subnet DNs are equal")
		}
		return nil
	}
}

func CreateSubnetWithoutParentDn(ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation without creating parent resource")
	resource := fmt.Sprintf(`
	resource "aci_subnet" "test" {
		ip = "%s"
	}
	`, ip)
	return resource
}

func CreateSubnetWithoutIP(rName string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_subnet" "test" {
		parent_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccSubnetConfigParentEpg(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with applicationEpg as parent resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_application_epg" "test"{
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_subnet" "test" {
		parent_dn = aci_application_epg.test.id
		ip = "%s"
	}
	`, rName, rName, rName, ip)
	return resource
}

func CreateAccSubnetConfigParentBridgeDomain(rName,ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with bridgeDomain as parent resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_bridge_domain" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
	}
	`, rName, rName, ip)
	return resource
}

func CreateAccSubnetConfigParentL2OutExtEpg(rName,ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with l2out_extepg as parent resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_l2_outside" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_l2out_extepg" "test"{
		name = "%s"
		l2_outside_dn = aci_l2_outside.test.id
	}

	resource "aci_subnet" "test" {
		parent_dn = aci_l2out_extepg.test.id
		ip = "%s"
	}
	`, rName, rName, rName, ip)
	return resource
}

func CreateAccSubnetConfig(rName,ip string) string {
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_subnet" "test" {
		parent_dn = aci_tenant.test.id
		ip = "%s"
	}
	`, rName, ip)
	return resource
}

func CreateAccApplicationProfileWithInValidTenantDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing applicationProfile creation with invalid tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_vrf.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccSubnetConfigWithOptionalValues(rName string) string {
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

func CreateAccSubnetConfigUpdatedName(rName, longerName string) string {
	fmt.Println("=== STEP  Basic: testing applicationProfile creation with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, longerName)
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
