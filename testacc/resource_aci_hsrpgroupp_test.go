package testacc

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


func TestAccAciL3outHsrpInterfaceGroup_Basic(t *testing.T) {
	var l3out_hsrp_interface_group_default models.L3outHsrpInterfaceGroup
	var l3out_hsrp_interface_group_updated models.L3outHsrpInterfaceGroup
	resourceName := "aci_l3out_hsrp_interface_group.testacc"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHsrpInterfaceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outHsrpInterfaceGroupWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName,"logical_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outHsrpInterfaceGroupWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_default),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", GetParentDn(l3out_hsrp_interface_group_default.DistinguishedName, fmt.Sprintf("/hsrpGroupP-%s", name))),
					resource.TestCheckResourceAttr(resourceName, "name",""),
					resource.TestCheckResourceAttr(resourceName, "annotation","orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description",""),
					resource.TestCheckResourceAttr(resourceName, "name_alias",""),
					resource.TestCheckResourceAttr(resourceName, "config_issues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config_issues.0", "none"),
					resource.TestCheckResourceAttr(resourceName, "group_af", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "group_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "group_name", ""),
					resource.TestCheckResourceAttr(resourceName, "ip", ""),
					resource.TestCheckResourceAttr(resourceName, "ip_obtain_mode", "admin"),
					resource.TestCheckResourceAttr(resourceName, "mac", ""),
					
				),
			},
			{
				// in this step all optional attribute expect realational attribute are given for the same resource and then compared
				Config: CreateAccL3outHsrpInterfaceGroupConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_hsrp_interface_group"),
					resource.TestCheckResourceAttr(resourceName, "config_issues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config_issues.0", "GroupMac-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "group_af", "ipv6"),
					resource.TestCheckResourceAttr(resourceName, "group_id", "1"),resource.TestCheckResourceAttr(resourceName, "group_id", ""),
					resource.TestCheckResourceAttr(resourceName, "group_name", ""),
					resource.TestCheckResourceAttr(resourceName, "ip", ""),
					resource.TestCheckResourceAttr(resourceName, "ip_obtain_mode", "auto"),
					resource.TestCheckResourceAttr(resourceName, "mac", ""),
					
				),
			},  
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)* failed validation`),
			},
			
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfigWithRequiredParams(rNameUpdated,rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", GetParentDn(l3out_hsrp_interface_group_default.DistinguishedName, fmt.Sprintf("/hsrpGroupP-%s", name))),
					resource.TestCheckResourceAttr(resourceName, "name",rName),
					testAccCheckAciL3outHsrpInterfaceGroupIdNotEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfigWithRequiredParams(rName,rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", GetParentDn(l3out_hsrp_interface_group_default.DistinguishedName, fmt.Sprintf("/hsrpGroupP-%s", name))),
					resource.TestCheckResourceAttr(resourceName, "name",rNameUpdated),
					testAccCheckAciL3outHsrpInterfaceGroupIdNotEqual(&l3out_hsrp_interface_group_default, &l3out_hsrp_interface_group_updated),
				),
			},
		},
	})
}

func TestAccAciL3outHsrpInterfaceGroup_Update(t *testing.T) {
	var l3out_hsrp_interface_group_default models.L3outHsrpInterfaceGroup
	var l3out_hsrp_interface_group_updated models.L3outHsrpInterfaceGroup
	resourceName := "aci_l3out_hsrp_interface_group.testacc"
	rName := makeTestVariable(acctest.RandString(5))
	
	fvTenantName := makeTestVariable(acctest.RandString(5))
	l3extOutName := makeTestVariable(acctest.RandString(5))
	l3extLNodePName := makeTestVariable(acctest.RandString(5))
	l3extLIfPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHsrpInterfaceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_default),
				),
			},
			
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupMac-Conflicts-Other-Group"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupMac-Conflicts-Other-Group"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupMac-Conflicts-Other-Group","GroupName-Conflicts-Other-Group"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupMac-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupName-Conflicts-Other-Group"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupName-Conflicts-Other-Group"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupName-Conflicts-Other-Group"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupMac-Conflicts-Other-Group","GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupMac-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "GroupVIP-Conflicts-Other-Group"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupVIP-Conflicts-Other-Group"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupVIP-Conflicts-Other-Group"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupVIP-Conflicts-Other-Group"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupMac-Conflicts-Other-Group","GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupMac-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Multiple-Version-On-Interface"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Multiple-Version-On-Interface"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Multiple-Version-On-Interface"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Multiple-Version-On-Interface"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Multiple-Version-On-Interface"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupMac-Conflicts-Other-Group","GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupMac-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "Secondary-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Secondary-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Secondary-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Secondary-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Secondary-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Secondary-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupMac-Conflicts-Other-Group","GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupMac-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.5", "Secondary-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "Secondary-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Secondary-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Secondary-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Secondary-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Secondary-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Secondary-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupMac-Conflicts-Other-Group","GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "7"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupMac-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.5", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.6", "group-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.5", "group-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "group-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "group-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "group-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "group-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"group-vip-conflicts-if-ip"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "group-vip-conflicts-if-ip"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupMac-Conflicts-Other-Group","GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "8"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupMac-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.5", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.6", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.7", "group-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "7"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.5", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.6", "group-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.5", "group-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "group-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "group-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "group-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"group-vip-conflicts-if-ip","group-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "group-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"group-vip-subnet-mismatch"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "group-vip-subnet-mismatch"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupMac-Conflicts-Other-Group","GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch","none"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "9"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupMac-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.5", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.6", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.7", "group-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.8", "none"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupName-Conflicts-Other-Group","GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch","none"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "8"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupName-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.5", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.6", "group-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.7", "none"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"GroupVIP-Conflicts-Other-Group","Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch","none"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "7"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "GroupVIP-Conflicts-Other-Group"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.5", "group-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.6", "none"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Multiple-Version-On-Interface","Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch","none"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Multiple-Version-On-Interface"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "group-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.5", "none"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Secondary-vip-conflicts-if-ip","Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch","none"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Secondary-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "group-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.4", "none"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"Secondary-vip-subnet-mismatch","group-vip-conflicts-if-ip","group-vip-subnet-mismatch","none"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "Secondary-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "group-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.3", "none"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"group-vip-conflicts-if-ip","group-vip-subnet-mismatch","none"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "group-vip-conflicts-if-ip"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "group-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.2", "none"),
				),
			},
			{
				
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"group-vip-subnet-mismatch","none"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHsrpInterfaceGroupExists(resourceName, &l3out_hsrp_interface_group_updated),
					resource.TestCheckResourceAttr(resourceName, "configIssues.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.0", "group-vip-subnet-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "configIssues.1", "none"),
				),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "configIssues", StringListtoString([]string{"none"})),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName),
			},
		},
	})
}

func TestAccAciL3outHsrpInterfaceGroup_Negative(t *testing.T) {
	var l3out_hsrp_interface_group_default models.L3outHsrpInterfaceGroup
	var l3out_hsrp_interface_group_updated models.L3outHsrpInterfaceGroup
	resourceName := "aci_l3out_hsrp_interface_group.testacc"
	rName := makeTestVariable(acctest.RandString(5))
	
	fvTenantName := makeTestVariable(acctest.RandString(5))
	l3extOutName := makeTestVariable(acctest.RandString(5))
	l3extLNodePName := makeTestVariable(acctest.RandString(5))
	l3extLIfPName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHsrpInterfaceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupWithInValidParentDn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName),
				ExpectError: regexp.MustCompile(`configured object (.)+ not found (.)+,`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "config_issues", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected(.*)to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "config_issues", StringListtoString([]string{"GroupMac-Conflicts-Other-Group", "GroupMac-Conflicts-Other-Group"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			// TODO: add unspecified case for "config_issues" if applicable
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "group_af", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "group_id", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "group_name", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "ip", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "ip_obtain_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, "mac", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			
			{
				Config:      CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccL3outHsrpInterfaceGroupConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName),
			},
		},
	})
}

func testAccCheckAciL3outHsrpInterfaceGroupExists(name string, l3out_hsrp_interface_group *models.L3outHsrpInterfaceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Hsrp Interface Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Hsrp Interface Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_hsrp_interface_groupFound := models.L3outHsrpInterfaceGroupFromContainer(cont)
		if l3out_hsrp_interface_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Hsrp Interface Group %s not found", rs.Primary.ID)
		}
		*l3out_hsrp_interface_group = *l3out_hsrp_interface_groupFound
		return nil
	}
}

func testAccCheckAciL3outHsrpInterfaceGroupDestroy(s *terraform.State) error {	
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		 if rs.Type == "aci_l3out_hsrp_interface_group" {
			cont,err := client.Get(rs.Primary.ID)
			l3out_hsrp_interface_group := models.L3outHsrpInterfaceGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Hsrp Interface Group %s Still exists",l3out_hsrp_interface_group.DistinguishedName)
			}
		}else{
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outHsrpInterfaceGroupIdEqual(m1, m2 *models.L3outHsrpInterfaceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_hsrp_interface_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outHsrpInterfaceGroupIdNotEqual(m1, m2 *models.L3outHsrpInterfaceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_hsrp_interface_group DNs are equal")
		}
		return nil
	}
}

func CreateL3outHsrpInterfaceGroupWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_group creation without ",attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
		
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	`
	switch attrName {
	case "logical_interface_profile_dn":
		rBlock += `
	resource "aci_l3out_hsrp_interface_group" "test" {
	#	logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		name  = "%s"
		description = "created while acceptance testing"
	}
		`
	case "name":
		rBlock += `
	resource "aci_l3out_hsrp_interface_group" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	#	name  = "%s"
		description = "created while acceptance testing"
	}
		`
	}
	return fmt.Sprintf(rBlock,fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
}

func CreateAccL3outHsrpInterfaceGroupConfigWithRequiredParams(rName,rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		name  = "%s"
	}
	`, rName,rName)
	return resource
}


func CreateAccL3outHsrpInterfaceGroupConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		name  = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
	return resource
}

func CreateAccL3outHsrpInterfaceGroupWithInValidParentDn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_hsrp_interface_group creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}invalid"
		name  = "%s"	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
	return resource
}


func CreateAccL3outHsrpInterfaceGroupConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_hsrp_interface_group"
		config_issues = ["GroupMac-Conflicts-Other-Group"]
		group_af = "ipv6"group_id = "1"group_id = ""group_name = ""ip = ""ip_obtain_mode = "auto"mac = ""
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)

	return resource
}

func CreateAccL3outHsrpInterfaceGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_hsrp_interface_group creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_l3out_hsrp_interface_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_hsrp_interface_group"
		config_issues = ["GroupMac-Conflicts-Other-Group"]
		group_af = "ipv6"group_id = "1"group_id = ""group_name = ""ip = ""ip_obtain_mode = "auto"mac = ""
	}
	`)

	return resource
}

func CreateAccL3outHsrpInterfaceGroupUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName,attribute,value string) string {
	fmt.Printf("=== STEP  testing l3out_hsrp_interface_group attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_hsrp_interface_group" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, rName)
	return resource
}