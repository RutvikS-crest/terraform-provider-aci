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

func TestAccAciSubnet_Basic(t *testing.T) {
	var subnet_default models.Subnet
	var subnet_updated models.Subnet
	resourceName := "aci_subnet.test"
	rName := acctest.RandString(5)
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
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
				Config: CreateAccSubnetConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "nd"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "preferred", "no"),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "private"),
					resource.TestCheckResourceAttr(resourceName, "virtual", "no"),
					// resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_nd_pfx_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/BD-%s", rName, rName)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetConfigWithOptionalValues(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag_subnet"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "nd"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "querier"),
					resource.TestCheckResourceAttr(resourceName, "description", "subnet"),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "alias_subnet"),
					resource.TestCheckResourceAttr(resourceName, "preferred", "no"),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "private"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared"),
					resource.TestCheckResourceAttr(resourceName, "virtual", "yes"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_nd_pfx_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/BD-%s", rName, rName)),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSubnetWithInavalidIP(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvSubnet (.)+`),
			},
		},
	})
}

// func ExpandList(attr, rName string, num int) []resource.TestCheckFunc {
// 	list := make([]resource.TestCheckFunc, 0, 1)
// 	list = append(list, resource.TestCheckResourceAttr(rName, fmt.Sprintf("%s.#", attr), strconv.Itoa(num)))
// 	for i := 0; i < num; i++ {
// 		list = append(list, resource.TestCheckResourceAttr(rName, fmt.Sprintf("%s.%d", attr, i), "querier"))
// 	}
// 	return list
// }

func TestAccSubnet_Update(t *testing.T) {
	var subnet_default models.Subnet
	var subnet_updated models.Subnet
	resourceName := "aci_subnet.test"
	rName := acctest.RandString(5)
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSubnetConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_default),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttr(rName, ip, "description", "updated description for terraform test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description for terraform test"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttr(rName, ip, "annotation", "updated_annotation_for_terraform_test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "updated_annotation_for_terraform_test"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttr(rName, ip, "preferred", "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "preferred", "yes"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttr(rName, ip, "virtual", "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "virtual", "yes"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"unspecified"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "unspecified"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"querier"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "querier"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"no-default-gateway"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "no-default-gateway"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"nd", "no-default-gateway"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "nd"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "no-default-gateway"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"nd", "querier"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "nd"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "querier"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"no-default-gateway", "querier"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "no-default-gateway"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "querier"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{"nd", "no-default-gateway", "querier"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "nd"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "no-default-gateway"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "querier"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"public"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "public"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"shared"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "shared"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"private", "public"})),
				ExpectError: regexp.MustCompile(`Invalid Configuration : Subnet scope cannot be both private and public`),
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"private", "shared"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "private"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"public", "shared"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scope.0", "public"),
					resource.TestCheckResourceAttr(resourceName, "scope.1", "shared"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{"private", "public", "shared"})),
				ExpectError: regexp.MustCompile(`Invalid Configuration : Subnet scope cannot be both private and public`),
			},
			{
				Config: CreateAccSubnetUpdatedAttr(rName, ip, "name_alias", "updated_name_alias_for_terraform_test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_updated),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "updated_name_alias_for_terraform_test"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_updated),
				),
			},
		},
	})
}

func TestAccSubnet_NegativeCases(t *testing.T) {
	rName := acctest.RandString(5)
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSubnetConfig(rName, ip),
			},
			{
				Config:      CreateAccSubnetWithInValidParentDn(rName, ip),
				ExpectError: regexp.MustCompile(`configured object (.)+ not found (.)+,`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, "virtual", randomValue),
				ExpectError: regexp.MustCompile(`expected virtual to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, "preferred", randomValue),
				ExpectError: regexp.MustCompile(`expected preferred to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "scope", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected scope.0 to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttrList(rName, ip, "ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected ctrl.0 to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccSubnetUpdatedAttr(rName, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSubnetConfig(rName, ip),
			},
		},
	})
}

func TestAccSubnet_reltionalParameters(t *testing.T) {
	var subnet_default models.Subnet
	var subnet_rel1 models.Subnet
	var subnet_rel2 models.Subnet
	resourceName := "aci_subnet.test"
	rName := acctest.RandString(5)
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	bdSubnetToProfileName1 := acctest.RandString(5)
	bdSubnetToProfileName2 := acctest.RandString(5)
	bdSubnetToOutName1 := acctest.RandString(5)
	bdSubnetToOutName2 := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSubnetConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_default),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_nd_pfx_pol", ""),
					// resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccSubnetUpdatedbdSubnetIntial(rName, ip, bdSubnetToProfileName1, bdSubnetToOutName1, "aci_bgp_route_control_profile.test.id", StringListtoStringWithoutQuoted([]string{"aci_l3_outside.test1.id"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_rel1),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", fmt.Sprintf("uni/tn-%s/prof-%s", rName, bdSubnetToProfileName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.0", fmt.Sprintf("uni/tn-%s/out-%s", rName, bdSubnetToOutName1)),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_rel1),
				),
			},
			{
				Config: CreateAccSubnetUpdatedbdSubnetFinal(rName, ip, bdSubnetToProfileName2, bdSubnetToOutName1, bdSubnetToOutName2, "aci_bgp_route_control_profile.test.id", StringListtoStringWithoutQuoted([]string{"aci_l3_outside.test1.id", "aci_l3_outside.test2.id"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_rel2),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", fmt.Sprintf("uni/tn-%s/prof-%s", rName, bdSubnetToProfileName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.#", "2"),
					testAccCheckAciSubnetIdEqual(&subnet_default, &subnet_rel2),
				),
			},
			{
				Config: CreateAccSubnetConfig(rName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists(resourceName, &subnet_default),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_nd_pfx_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_subnet_to_out.#", "0"),
				),
			},
		},
	})
}

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
	fmt.Println("=== STEP  Basic: testing subnet creation without giving ip")
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

func CreateAccSubnetConfig(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with required arguements only")
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

func CreateAccSubnetWithInValidParentDn(rName, ip string) string {
	fmt.Println("=== STEP  Negative Case: testing subnet creation with invalid parent_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_subnet" "test" {
		parent_dn = "${aci_bridge_domain.test.id}xyz"
		ip = "%s"
	}
	`, rName, rName, ip)
	return resource
}

func CreateAccSubnetConfigWithOptionalValues(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with optional parameters")
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
		description = "subnet"
        annotation = "tag_subnet"
        ctrl = ["nd", "querier"]
        name_alias = "alias_subnet"
        preferred = "no"
        scope = ["private", "shared"]
        virtual = "yes"
	}
	`, rName, rName, ip)
	return resource
}

func CreateAccSubnetUpdatedbdSubnetIntial(rName, ip, bdSubnetToProfileName, bdSubnetToOutName, bdSubnetToProfileRef, bdSubnetToOutRef string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with initial relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bgp_route_control_profile" "test" {
		parent_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_l3_outside" "test1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
		relation_fv_rs_bd_subnet_to_profile = %s
		relation_fv_rs_bd_subnet_to_out = %s
	}
	`, rName, bdSubnetToProfileName, bdSubnetToOutName, rName, ip, bdSubnetToProfileRef, bdSubnetToOutRef)
	return resource
}

func CreateAccSubnetUpdatedbdSubnetFinal(rName, ip, bdSubnetToProfileName, bdSubnetToOutName1, bdSubnetToOutName2, bdSubnetToProfileRef, bdSubnetToOutRef string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with final relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bgp_route_control_profile" "test" {
		parent_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_l3_outside" "test1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_l3_outside" "test2" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_subnet" "test" {
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
		relation_fv_rs_bd_subnet_to_profile = %s
		relation_fv_rs_bd_subnet_to_out = %s
	}
	`, rName, bdSubnetToProfileName, bdSubnetToOutName1, bdSubnetToOutName2, rName, ip, bdSubnetToProfileRef, bdSubnetToOutRef)
	return resource
}

func CreateAccSubnetWithInavalidIP(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing subnet creation with invalid IP")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s0"
	}
	`, rName, rName, ip)
	return resource
}

func CreateAccSubnetUpdatedAttr(rName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
		%s = "%s"
	}
	`, rName, rName, ip, attribute, value)
	return resource
}

func CreateAccSubnetUpdatedAttrList(rName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_subnet" "test"{
		parent_dn = aci_bridge_domain.test.id
		ip = "%s"
		%s = %s
	}
	`, rName, rName, ip, attribute, value)
	return resource
}
