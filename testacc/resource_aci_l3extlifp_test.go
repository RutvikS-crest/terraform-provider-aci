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

func TestAccAciLogicalInterfaceProfile_Basic(t *testing.T) {
	var logicalInterfaceProfile_default models.LogicalInterfaceProfile
	var logicalInterfaceProfile_updated models.LogicalInterfaceProfile
	resourceName := "aci_logical_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	prOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccLogicalInterfaceProfileWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileWithoutLogicalNodeProfileDn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "tag", "yellow-green"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l_if_p_to_netflow_monitor_pol.#", "0"),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag_prof"),
					resource.TestCheckResourceAttr(resourceName, "description", "Sample logical interface profile"),
					resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "alias_prof"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "tag", "navy"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_l_if_p_to_netflow_monitor_pol.#", "0"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfigWithAnotherName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					testAccCheckAciLogicalInterfaceIdNotEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
			},
			{
				Config: 	CreateAccl3outsideConfigWithAnotherLogicalNodeProfileDn(prOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					testAccCheckAciLogicalInterfaceIdNotEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileConfigUpdateWithoutRequiredAttri(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileConfigUpdateWithInvalidName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of lifp-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
			},
		},
	})
}

func TestAccAciLogicalInterfaceProfile_Update(t *testing.T) {
	var logicalInterfaceProfile_default models.LogicalInterfaceProfile
	var logicalInterfaceProfile_updated models.LogicalInterfaceProfile
	resourceName := "aci_logical_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_default),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "medium-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "medium-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "teal"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "teal"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-cyan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-cyan"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "deep-sky-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "deep-sky-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-turquoise"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-turquoise"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "medium-spring-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "medium-spring-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "lime"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "lime"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "spring-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "spring-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "aqua"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "aqua"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			// {
			// 	Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "cyan"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "tag", "cyan"),
			// 		testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
			// 	),
			// },
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "midnight-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "midnight-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dodger-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dodger-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-sea-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-sea-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "forest-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "forest-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "sea-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "sea-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-slate-gray"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-slate-gray"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "lime-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "lime-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "medium-sea-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "medium-sea-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "turquoise"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "turquoise"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "royal-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "royal-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "steel-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "steel-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-slate-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-slate-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "medium-turquoise"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "medium-turquoise"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "indigo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "indigo"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-olive-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-olive-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "cadet-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "cadet-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "cornflower-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "cornflower-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "medium-aquamarine"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "medium-aquamarine"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dim-gray"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dim-gray"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "slate-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "slate-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "olive-drab"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "olive-drab"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "slate-gray"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "slate-gray"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-slate-gray"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-slate-gray"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "medium-slate-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "medium-slate-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "lawn-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "lawn-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "chartreuse"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "chartreuse"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "aquamarine"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "aquamarine"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "maroon"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "maroon"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "purple"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "purple"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "olive"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "olive"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "gray"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "gray"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "sky-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "sky-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-sky-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-sky-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "blue-violet"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "blue-violet"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-red"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-red"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-magenta"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-magenta"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "saddle-brown"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "saddle-brown"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-sea-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-sea-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "medium-purple"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "medium-purple"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-violet"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-violet"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "pale-green"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "pale-green"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-orchid"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-orchid"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "black"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "black"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "sienna"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "sienna"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "brown"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "brown"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-gray"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-gray"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "green-yellow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "green-yellow"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "pale-turquoise"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "pale-turquoise"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-steel-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-steel-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "powder-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "powder-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "fire-brick"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "fire-brick"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-goldenrod"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-goldenrod"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "medium-orchid"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "medium-orchid"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "rosy-brown"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "rosy-brown"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-khaki"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-khaki"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "silver"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "silver"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "medium-violet-red"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "medium-violet-red"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "indian-red"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "indian-red"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "peru"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "peru"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "chocolate"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "chocolate"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "tan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "tan"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-gray"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-gray"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "thistle"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "thistle"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "orchid"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "orchid"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "goldenrod"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "goldenrod"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "pale-violet-red"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "pale-violet-red"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "crimson"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "crimson"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "gainsboro"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "gainsboro"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "plum"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "plum"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "burlywood"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "burlywood"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-cyan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-cyan"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "lavender"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "lavender"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-salmon"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-salmon"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "violet"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "violet"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "pale-goldenrod"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "pale-goldenrod"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-coral"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-coral"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "khaki"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "khaki"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "alice-blue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "alice-blue"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "honeydew"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "honeydew"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "azure"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "azure"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "sandy-brown"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "sandy-brown"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "wheat"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "wheat"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "beige"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "beige"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "white-smoke"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "white-smoke"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "mint-cream"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "mint-cream"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "ghost-white"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "ghost-white"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "salmon"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "salmon"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "antique-white"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "antique-white"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "linen"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "linen"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-goldenrod-yellow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-goldenrod-yellow"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "old-lace"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "old-lace"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "red"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "red"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "fuchsia"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "fuchsia"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			// {
			// 	Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "magenta"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "tag", "magenta"),
			// 		testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
			// 	),
			// },
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "deep-pink"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "deep-pink"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "orange-red"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "orange-red"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "tomato"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "tomato"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "hot-pink"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "hot-pink"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "coral"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "coral"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "dark-orange"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "dark-orange"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-salmon"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-salmon"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "orange"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "orange"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-pink"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-pink"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "pink"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "pink"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "gold"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "gold"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "peachpuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "peachpuff"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "navajo-white"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "navajo-white"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "moccasin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "moccasin"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "bisque"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "bisque"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "misty-rose"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "misty-rose"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "blanched-almond"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "blanched-almond"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "papaya-whip"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "papaya-whip"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "lavender-blush"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "lavender-blush"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			 {
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "seashell"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "seashell"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "cornsilk"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "cornsilk"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "lemon-chiffon"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "lemon-chiffon"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "floral-white"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "floral-white"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "snow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "snow"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "yellow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "yellow"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "light-yellow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "light-yellow"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "ivory"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "ivory"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
			 {
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", "white"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogicalInterfaceProfileExists(resourceName, &logicalInterfaceProfile_updated),
					resource.TestCheckResourceAttr(resourceName, "tag", "white"),
					testAccCheckAciLogicalInterfaceProfileIdEqual(&logicalInterfaceProfile_default, &logicalInterfaceProfile_updated),
				),
			},
		},
	})
}
func TestAccAciLogicalInterfaceProfile_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
			},
			{
				Config:      CreateAccl3outsideConfigWithInvalidLogicalNodeProfiledn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class l3extLIfP (.)+`),
			},
			{
				Config: CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name prio, class l3extLIfP (.)+`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, "tag", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name tag, class l3extLIfP (.)+`),
			},
			{
				Config:      CreateAccLogicalInterfaceProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLogicalInterfaceProfileConfig(rName),
			},
		},
	})
}

func TestAccAciLogicalInterfaceProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalInterfaceProfileConfigMultiple(rName),
			},
		},
	})
}

func CreateAccLogicalInterfaceProfileConfigUpdateWithInvalidName(parentName, rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile update with invalid Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
  resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

   resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	   resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
  }
	`, parentName, parentName, parentName, rName)
	return resource
}

func testAccCheckLogicalInterfaceProfileExists(name string, logicalInterfaceProfile *models.LogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LogicalInterfaceProfile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LogicalInterfaceProfile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		logicalInterfaceProfileFound := models.LogicalInterfaceProfileFromContainer(cont)
		if logicalInterfaceProfileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LogicalInterfaceProfile %s not found", rs.Primary.ID)
		}
		*logicalInterfaceProfile = *logicalInterfaceProfileFound
		return nil
	}
}

func testAccCheckAciLogicalInterfaceProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing LogicalInterfaceProfile destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_logical_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			aci := models.LogicalInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LogicalInterfaceProfile %s Still exists", aci.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLogicalInterfaceProfileIdEqual(logicalInterfaceProfile1, logicalInterfaceProfile2 *models.LogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if logicalInterfaceProfile1.DistinguishedName != logicalInterfaceProfile2.DistinguishedName {
			return fmt.Errorf("LogicalInterfaceProfile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLogicalInterfaceIdNotEqual(logicalInterfaceProfile1, logicalInterfaceProfile2 *models.LogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if logicalInterfaceProfile1.DistinguishedName == logicalInterfaceProfile2.DistinguishedName {
			return fmt.Errorf("LogicalInterfaceProfile DNs are equal")
		}
		return nil
	}
}

func CreateAccLogicalInterfaceProfileWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile creation without giving Name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
  resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

   resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	   resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileWithoutLogicalNodeProfileDn(rName string) string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile creation without giving LogicalNodeProfile dn")
	resource := fmt.Sprintf(`
	resource "aci_logical_interface_profile" "test" {
        name                    = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileConfigUpdateWithoutRequiredAttri() string {
	fmt.Println("=== STEP  Basic: testing LogicalInterfaceProfile update without giving required Attributes")
	resource := fmt.Sprintf(`
    resource "aci_logical_interface_profile" "test" {
		description             = "Sample logical interface profile"
        annotation              = "tag_prof"
        name_alias              = "alias_prof"
        prio                    = "level1"
        tag                     = "navy"
  }
	`)
	return resource
}



func CreateAccLogicalInterfaceProfileConfigWithAnotherName(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing LogicalInterfaceProfile creation with different LogicalInterfaceProfile name %s \n", rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
  resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

   resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	   resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
  }
	`, parentName, parentName, parentName, rName)
	return resource
}

func CreateAccl3outsideConfigWithAnotherLogicalNodeProfileDn(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing LogicalInterfaceProfile creation with different parent %s \n", parentName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
  resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

   resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
	   resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
  }
	`, parentName, parentName, parentName, rName)
	return resource
}

func CreateAccl3outsideConfigWithInvalidLogicalNodeProfiledn(rName string) string {
	fmt.Printf("=== STEP  Basic: testing LogicalInterfaceProfile creation with invalid Logical Node Profile dn \n")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
    resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_tenant.test.id
        name                    = "%s"
		description             = "Sample logical interface profile"
        annotation              = "tag_prof"
        name_alias              = "alias_prof"
        prio                    = "level1"
        tag                     = "navy"
  }
	`, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileConfig(rName string) string {
	fmt.Println("=== STEP testing LogicalInterfaceProfile creation with required attributes")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
    resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

    resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
    resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
  }
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing logicalInterfaceProfile creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
    resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

    resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
    resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
		description             = "Sample logical interface profile"
        annotation              = "tag_prof"
        name_alias              = "alias_prof"
        prio                    = "level1"
        tag                     = "navy"
  }
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccLogicalInterfaceProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
    resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

    resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
    resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name           = "%s"
		%s = "%s"
  }
	`, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccLogicalInterfaceProfileConfigMultiple(rName string) string {
	fmt.Println("=== STEP  Creating Multiple LogicalInterfaceProfile")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
    resource "aci_l3_outside" "test" {
        tenant_dn      = aci_tenant.test.id
        description    = "from terraform"
        name           = "%s"
    }

    resource "aci_logical_node_profile" "test" {
        l3_outside_dn = aci_l3_outside.test.id
        name          = "%s"
      }
    resource "aci_logical_interface_profile" "test" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
    }
	resource "aci_logical_interface_profile" "test1" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
    }
	resource "aci_logical_interface_profile" "test2" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
    }
	resource "aci_logical_interface_profile" "test3" {
        logical_node_profile_dn = aci_logical_node_profile.test.id
        name                    = "%s"
    }
	`, rName, rName, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource
}