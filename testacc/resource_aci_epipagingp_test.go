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

func TestAccAciIPAgingPolicy_Basic(t *testing.T) {
	var endpoint_ip_aging_profile_default models.IPAgingPolicy
	var endpoint_ip_aging_profile_updated models.IPAgingPolicy
	resourceName := "aci_endpoint_ip_aging_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciIPAgingPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateIPAgingPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccIPAgingPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIPAgingPolicyExists(resourceName, &endpoint_ip_aging_profile_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),
				),
			},
			{
				Config: CreateAccIPAgingPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIPAgingPolicyExists(resourceName, &endpoint_ip_aging_profile_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_endpoint_ip_aging_profile"),

					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),

					testAccCheckAciIPAgingPolicyIdEqual(&endpoint_ip_aging_profile_default, &endpoint_ip_aging_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccIPAgingPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccIPAgingPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccIPAgingPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIPAgingPolicyExists(resourceName, &endpoint_ip_aging_profile_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciIPAgingPolicyIdNotEqual(&endpoint_ip_aging_profile_default, &endpoint_ip_aging_profile_updated),
				),
			},
		},
	})
}

func TestAccAciIPAgingPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciIPAgingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccIPAgingPolicyConfig(rName),
			},

			{
				Config:      CreateAccIPAgingPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccIPAgingPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccIPAgingPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccIPAgingPolicyUpdatedAttr(rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccIPAgingPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccIPAgingPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciIPAgingPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciIPAgingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccIPAgingPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciIPAgingPolicyExists(name string, endpoint_ip_aging_profile *models.IPAgingPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Ip Aging Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Ip Aging Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_ip_aging_profileFound := models.IPAgingPolicyFromContainer(cont)
		if endpoint_ip_aging_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Ip Aging Profile %s not found", rs.Primary.ID)
		}
		*endpoint_ip_aging_profile = *endpoint_ip_aging_profileFound
		return nil
	}
}

func testAccCheckAciIPAgingPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_endpoint_ip_aging_profile" {
			cont, err := client.Get(rs.Primary.ID)
			endpoint_ip_aging_profile := models.IPAgingPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Endpoint Ip Aging Profile %s Still exists", endpoint_ip_aging_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciIPAgingPolicyIdEqual(m1, m2 *models.IPAgingPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("endpoint_ip_aging_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciIPAgingPolicyIdNotEqual(m1, m2 *models.IPAgingPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("endpoint_ip_aging_profile DNs are equal")
		}
		return nil
	}
}

func CreateIPAgingPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_ip_aging_profile creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_endpoint_ip_aging_profile" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccIPAgingPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccIPAgingPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccIPAgingPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccIPAgingPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple endpoint_ip_aging_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccIPAgingPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing endpoint_ip_aging_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_endpoint_ip_aging_profile"
		admin_st = "enabled"
		
	}
	`, rName)

	return resource
}

func CreateAccIPAgingPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing endpoint_ip_aging_profile updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_endpoint_ip_aging_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_endpoint_ip_aging_profile"
		admin_st = "enabled"
		
	}
	`)

	return resource
}

func CreateAccIPAgingPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing endpoint_ip_aging_profile attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
