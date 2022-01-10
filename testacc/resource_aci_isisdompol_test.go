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

func TestAccAciISISDomainPolicy_Basic(t *testing.T) {
	var isis_domain_policy_default models.ISISDomainPolicy
	var isis_domain_policy_updated models.ISISDomainPolicy
	resourceName := "aci_isis_domain_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciISISDomainPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateISISDomainPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccISISDomainPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "mtu", "1492"),
					resource.TestCheckResourceAttr(resourceName, "redistrib_metric", "63"),
				),
			},
			{
				Config: CreateAccISISDomainPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_isis_domain_policy"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "257"),
					resource.TestCheckResourceAttr(resourceName, "redistrib_metric", "2"),

					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccISISDomainPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccISISDomainPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccISISDomainPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciISISDomainPolicyIdNotEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
		},
	})
}

func TestAccAciISISDomainPolicy_Update(t *testing.T) {
	var isis_domain_policy_default models.ISISDomainPolicy
	var isis_domain_policy_updated models.ISISDomainPolicy
	resourceName := "aci_isis_domain_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciISISDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccISISDomainPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_default),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr(rName, "mtu", "4352"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "mtu", "4352"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr(rName, "mtu", "2048"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "mtu", "2048"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr(rName, "redistrib_metric", "63"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "redistrib_metric", "63"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},
			{
				Config: CreateAccISISDomainPolicyUpdatedAttr(rName, "redistrib_metric", "31"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciISISDomainPolicyExists(resourceName, &isis_domain_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "redistrib_metric", "31"),
					testAccCheckAciISISDomainPolicyIdEqual(&isis_domain_policy_default, &isis_domain_policy_updated),
				),
			},

			{
				Config: CreateAccISISDomainPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciISISDomainPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciISISDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccISISDomainPolicyConfig(rName),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(rName, "mtu", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(rName, "mtu", "255"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(rName, "mtu", "4353"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(rName, "redistrib_metric", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(rName, "redistrib_metric", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(rName, "redistrib_metric", "64"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccISISDomainPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccISISDomainPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciISISDomainPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciISISDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccISISDomainPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciISISDomainPolicyExists(name string, isis_domain_policy *models.ISISDomainPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("ISIS Domain Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ISIS Domain Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		isis_domain_policyFound := models.ISISDomainPolicyFromContainer(cont)
		if isis_domain_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("ISIS Domain Policy %s not found", rs.Primary.ID)
		}
		*isis_domain_policy = *isis_domain_policyFound
		return nil
	}
}

func testAccCheckAciISISDomainPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing isis_domain_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_isis_domain_policy" {
			cont, err := client.Get(rs.Primary.ID)
			isis_domain_policy := models.ISISDomainPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("ISIS Domain Policy %s Still exists", isis_domain_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciISISDomainPolicyIdEqual(m1, m2 *models.ISISDomainPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("isis_domain_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciISISDomainPolicyIdNotEqual(m1, m2 *models.ISISDomainPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("isis_domain_policy DNs are equal")
		}
		return nil
	}
}

func CreateISISDomainPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing isis_domain_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_isis_domain_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccISISDomainPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing isis_domain_policy creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccISISDomainPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing isis_domain_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccISISDomainPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing isis_domain_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccISISDomainPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple isis_domain_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccISISDomainPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing isis_domain_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_isis_domain_policy"
		mtu = "257"
		redistrib_metric = "2"
		
	}
	`, rName)

	return resource
}

func CreateAccISISDomainPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing isis_domain_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_isis_domain_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_isis_domain_policy"
		mtu = "257"
		redistrib_metric = "2"
		
	}
	`)

	return resource
}

func CreateAccISISDomainPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing isis_domain_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccISISDomainPolicyUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing isis_domain_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_isis_domain_policy" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
