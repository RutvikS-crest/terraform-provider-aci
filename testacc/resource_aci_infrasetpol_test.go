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

func TestAccAciFabricWideSettingsPolicy_Basic(t *testing.T) {
	var fabric_wide_settings_default models.FabricWideSettingsPolicy
	var fabric_wide_settings_updated models.FabricWideSettingsPolicy
	resourceName := "aci_fabric_wide_settings.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricWideSettingsPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricWideSettingsPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_default),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "disable_ep_dampening", "no"),
					resource.TestCheckResourceAttr(resourceName, "domain_validation", "no"),
					resource.TestCheckResourceAttr(resourceName, "enable_mo_streaming", "no"),
					resource.TestCheckResourceAttr(resourceName, "enable_remote_leaf_direct", "no"),
					resource.TestCheckResourceAttr(resourceName, "enforce_subnet_check", "no"),
					resource.TestCheckResourceAttr(resourceName, "opflexp_authenticate_clients", "yes"),
					resource.TestCheckResourceAttr(resourceName, "opflexp_use_ssl", "yes"),
					resource.TestCheckResourceAttr(resourceName, "reallocate_gipo", "no"),
					resource.TestCheckResourceAttr(resourceName, "restrict_infra_vlan_traffic", "no"),
					resource.TestCheckResourceAttr(resourceName, "unicast_xr_ep_learn_disable", "no"),
					resource.TestCheckResourceAttr(resourceName, "validate_overlapping_vlans", "no"),
				),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricWideSettingsPolicyExists(resourceName, &fabric_wide_settings_updated),

					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fabric_wide_settings"),

					resource.TestCheckResourceAttr(resourceName, "disable_ep_dampening", "yes"),

					resource.TestCheckResourceAttr(resourceName, "domain_validation", "yes"),

					resource.TestCheckResourceAttr(resourceName, "enable_mo_streaming", "yes"),

					resource.TestCheckResourceAttr(resourceName, "enable_remote_leaf_direct", "yes"),

					resource.TestCheckResourceAttr(resourceName, "enforce_subnet_check", "yes"),

					resource.TestCheckResourceAttr(resourceName, "opflexp_authenticate_clients", "no"),

					resource.TestCheckResourceAttr(resourceName, "opflexp_use_ssl", "no"),

					resource.TestCheckResourceAttr(resourceName, "reallocate_gipo", "yes"),

					resource.TestCheckResourceAttr(resourceName, "restrict_infra_vlan_traffic", "yes"),

					resource.TestCheckResourceAttr(resourceName, "unicast_xr_ep_learn_disable", "yes"),

					resource.TestCheckResourceAttr(resourceName, "validate_overlapping_vlans", "yes"),

					testAccCheckAciFabricWideSettingsPolicyIdEqual(&fabric_wide_settings_default, &fabric_wide_settings_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFabricWideSettingsPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestAccAciFabricWideSettingsPolicy_Negative(t *testing.T) {

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricWideSettingsPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricWideSettingsPolicyConfig(),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("disable_ep_dampening", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("domain_validation", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("enable_mo_streaming", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("enable_remote_leaf_direct", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("enforce_subnet_check", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("opflexp_authenticate_clients", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("opflexp_use_ssl", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("reallocate_gipo", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("restrict_infra_vlan_traffic", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("unicast_xr_ep_learn_disable", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr("validate_overlapping_vlans", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricWideSettingsPolicyUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyConfig(),
			},
		},
	})
}

func testAccCheckAciFabricWideSettingsPolicyExists(name string, fabric_wide_settings *models.FabricWideSettingsPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric Wide Settings %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric Wide Settings dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_wide_settingsFound := models.FabricWideSettingsPolicyFromContainer(cont)
		if fabric_wide_settingsFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric Wide Settings %s not found", rs.Primary.ID)
		}
		*fabric_wide_settings = *fabric_wide_settingsFound
		return nil
	}
}

func testAccCheckAciFabricWideSettingsPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing fabric_wide_settings destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fabric_wide_settings" {
			cont, err := client.Get(rs.Primary.ID)
			fabric_wide_settings := models.FabricWideSettingsPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric Wide Settings %s Still exists", fabric_wide_settings.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFabricWideSettingsPolicyIdEqual(m1, m2 *models.FabricWideSettingsPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fabric_wide_settings DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFabricWideSettingsPolicyIdNotEqual(m1, m2 *models.FabricWideSettingsPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("fabric_wide_settings DNs are equal")
		}
		return nil
	}
}

func CreateAccFabricWideSettingsPolicyConfigWithRequiredParams() string {
	fmt.Println("=== STEP  testing fabric_wide_settings creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_wide_settings" "test" {
	
	}
	`)
	return resource
}
func CreateAccFabricWideSettingsPolicyConfig() string {
	fmt.Println("=== STEP  testing fabric_wide_settings creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_wide_settings" "test" {
	
	}
	`)
	return resource
}

func CreateAccFabricWideSettingsPolicyConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing fabric_wide_settings creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_wide_settings" "test" {
	
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_wide_settings"
		disable_ep_dampening = "yes"
		domain_validation = "yes"
		enable_mo_streaming = "yes"
		enable_remote_leaf_direct = "yes"
		enforce_subnet_check = "yes"
		opflexp_authenticate_clients = "no"
		opflexp_use_ssl = "no"
		reallocate_gipo = "yes"
		restrict_infra_vlan_traffic = "yes"
		unicast_xr_ep_learn_disable = "yes"
		validate_overlapping_vlans = "yes"
		
	}
	`)

	return resource
}

func CreateAccFabricWideSettingsPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing fabric_wide_settings updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_fabric_wide_settings" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_wide_settings"
		disable_ep_dampening = "yes"
		domain_validation = "yes"
		enable_mo_streaming = "yes"
		enable_remote_leaf_direct = "yes"
		enforce_subnet_check = "yes"
		opflexp_authenticate_clients = "no"
		opflexp_use_ssl = "no"
		reallocate_gipo = "yes"
		restrict_infra_vlan_traffic = "yes"
		unicast_xr_ep_learn_disable = "yes"
		validate_overlapping_vlans = "yes"
		
	}
	`)

	return resource
}

func CreateAccFabricWideSettingsPolicyUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing fabric_wide_settings attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_wide_settings" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
