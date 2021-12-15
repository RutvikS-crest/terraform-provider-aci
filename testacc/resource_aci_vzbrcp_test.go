package acctest

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

func TestAccAciContract_Basic(t *testing.T) {
	var contract_default models.Contract
	var contract_updated models.Contract
	resourceName := "aci_contract.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	prOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateAccContractWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccContractWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccContractConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "scope", "context"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
				),
			},
			{
				Config: CreateAccContractConfigOptional(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_alias"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "scope", "tenant"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS0"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccContractRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccContractConfigWithFilterResources(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.annotation", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_name", rName),
					resource.TestCheckResourceAttr(resourceName, "filter.0.name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.apply_to_frag", "no"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.arp_opc", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_from_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_to_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_annotation", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_description", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.filter_entry_name", rName),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv4_t", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv6_t", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.match_dscp", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_from_port", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_to_port", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.stateful", "no"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.tcp_rules", "unspecified"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			// {
			// 	Config: CreateAccContractConfigWithFilterResourcesOptional(rName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.annotation", "filter_annotation"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.description", "filter_description"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_name", rName),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.name_alias", "filter_name_alias"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.apply_to_frag", "no"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.arp_opc", "unspecified"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_from_port", "20"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_to_port", "20"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_annotation", ""),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_description", ""),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_name_alias", ""),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "unspecified"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.filter_entry_name", rName),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv4_t", "unspecified"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv6_t", "unspecified"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.match_dscp", "unspecified"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "unspecified"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_from_port", "unspecified"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_to_port", "unspecified"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.stateful", "no"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.tcp_rules", ""),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			{
				Config:      CreateAccContractConfigWithParentAndName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of brc-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccContractConfigWithParentAndName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rOther),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciContrctIdNotEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractConfig(rName),
			},
			{
				Config: CreateAccContractConfigWithParentAndName(prOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", prOther)),
					testAccCheckAciContrctIdNotEqual(&contract_default, &contract_updated),
				),
			},
		},
	})
}

func TestAccAciContract_Update(t *testing.T) {
	var contract_default models.Contract
	var contract_updated models.Contract
	resourceName := "aci_contract.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccContractConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_default),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "scope", "global"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "scope", "global"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "scope", "application-profile"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "scope", "application-profile"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "CS1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS1"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF11"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF11"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF12"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF12"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF13"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF13"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "CS2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS2"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF21"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF21"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF22"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF22"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF23"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF23"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "CS3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS3"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF31"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF31"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF32"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF33"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF33"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "CS4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS4"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF41"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF41"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF42"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF42"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF43"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF43"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "CS5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS5"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "VA"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "VA"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "EF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "EF"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "CS6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS6"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "CS7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS7"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
		},
	})
}

func TestAccAciContract_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longAnnotationDesc := acctest.RandString(129)
	longNameAlias := acctest.RandString(65)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccContractConfig(rName),
			},
			{
				Config:      CreateAccContractWithInValidTenantDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class vzBrCP (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "description", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "annotation", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "target_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected target_dscp to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "scope", randomValue),
				ExpectError: regexp.MustCompile(`expected scope to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccContractUpdatedFilterAttr(rName, "description", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedFilterAttr(rName, "annotation", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedFilterAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedFilterAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			{
				Config:      CreateAccContractFilterWithInvalidName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of flt-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccContractConfig(rName),
			},
		},
	})
}

func TestAccContract_RelationParameters(t *testing.T) {
	var contract_default models.Contract
	var contract_rel1 models.Contract
	var contract_rel2 models.Contract
	resourceName := "aci_contract.test"
	rName := makeTestVariable(acctest.RandString(5))
	relRes1 := makeTestVariable(acctest.RandString(5))
	relRes2 := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccContractConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_default),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_graph_att", ""),
				),
			},
			{
				Config: CreateAccContractRelations(rName, relRes1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_rel1),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_graph_att", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", rName, relRes1)),
					testAccCheckAciContractdEqual(&contract_default, &contract_rel1),
				),
			},
			{
				Config: CreateAccContractRelations(rName, relRes2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_rel2),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_graph_att", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", rName, relRes2)),
					testAccCheckAciContractdEqual(&contract_default, &contract_rel2),
				),
			},
			{
				Config: CreateAccContractConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_graph_att", ""),
				),
			},
		},
	})
}

func CreateAccContractRelations(rName, relName string) string {
	fmt.Printf("=== STEP  testing vrf creation with resource name %s and relation resource name %s\n", rName, relName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_l4_l7_service_graph_template" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		relation_vz_rs_graph_att = aci_l4_l7_service_graph_template.test.id
	}
	`, rName, relName, rName)
	return resource
}

func CreateAccContractFilterWithInvalidName(rName, longrName string) string {
	fmt.Printf("=== STEP  testing contract's filter creation with name = %s\n", longrName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		filter{
			filter_name = "%s"
		}
	}

	`, rName, rName, longrName)
	return resource
}

func CreateAccContractUpdatedFilterAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing contract's filter with %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		filter{
			filter_name = "%s"
			%s = "%s"
		}
	}

	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccContractWithInValidTenantDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing contract creation with invalid tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_application_profile.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccContractUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing contract with %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}

	`, rName, rName, attribute, value)
	return resource
}

// func CreateAccContractConfigWithFilterResourcesOptional(rName string) string {
// 	fmt.Println("=== STEP  testing contract creation with optional parameters of filter resources")
// 	resource := fmt.Sprintf(`
// 	resource "aci_tenant" "test"{
// 		name = "%s"
// 	}

// 	resource "aci_contract" "test" {
// 		tenant_dn = aci_tenant.test.id
// 		name = "%s"
// 		filter {
// 		  filter_name = "%s"
// 		  annotation = "filter_annotation"
// 		  description = "filter_description"
// 		  name_alias = "filter_name_alias"
// 		  filter_entry {
// 			filter_entry_name = "%s"
// 			apply_to_frag = "no"
// 			arp_opc = "unspecified"
// 			d_from_port = "ftpData"
// 			d_to_port = "ftpData"
// 			entry_annotation = "entry_annotation"
// 			entry_description = "entry_description"
// 			entry_name_alias = "entry_name_alias"
// 			ether_t = "ipv4"
// 			icmpv4_t = "echo-rep"
// 			icmpv6_t = "dst-unreach"
// 			match_dscp = "CS0"
// 			prot = "tcp"
// 			stateful = "yes"
// 			tcp_rules = "est"
// 		  }
// 		}
// 	}
// 	`, rName, rName, rName, rName)
// 	return resource
// }

func CreateAccContractConfigWithParentAndName(prName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing contract creation with tenant name %s name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, prName, rName)
	return resource
}

func CreateAccContractRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing contract updation without required fields")
	resource := fmt.Sprintln(`
	resource "aci_contract" "test" {
		annotation = "test_annotation"
		description = "test_description"
		name_alias = "test_alias"
		prio = "level1"
		scope = "tenant"
		target_dscp = "CS0"
	}
	`)
	return resource
}

func CreateAccContractConfigWithFilterResources(rName string) string {
	fmt.Println("=== STEP  testing contract creation with filter resources")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		filter {
			filter_name = "%s"
			filter_entry {
			  filter_entry_name = "%s"
			}
		}
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccContractConfigOptional(rName string) string {
	fmt.Println("=== STEP  testing contract creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		annotation = "test_annotation"
		description = "test_description"
		name_alias = "test_alias"
		prio = "level1"
		scope = "tenant"
		target_dscp = "CS0"
	}
	`, rName, rName)
	return resource
}

func CreateAccContractConfig(rName string) string {
	fmt.Println("=== STEP  testing contract creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccContractWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract creation without name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccContractWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract creation without creating tenant")
	resource := fmt.Sprintf(`
	resource "aci_contract" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}

func testAccCheckAciContractExists(name string, contract *models.Contract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Contract %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Contract dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		contractFound := models.ContractFromContainer(cont)
		if contractFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Contract %s not found", rs.Primary.ID)
		}
		*contract = *contractFound
		return nil
	}
}

func testAccCheckAciContractDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing contract destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_contract" {
			cont, err := client.Get(rs.Primary.ID)
			contract := models.ContractFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Contract %s Still exists", contract.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciContrctIdNotEqual(c1, c2 *models.Contract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if c1.DistinguishedName == c2.DistinguishedName {
			return fmt.Errorf("Contract DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciContractdEqual(c1, c2 *models.Contract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if c1.DistinguishedName != c2.DistinguishedName {
			return fmt.Errorf("Contract DNs are not equal")
		}
		return nil
	}
}
