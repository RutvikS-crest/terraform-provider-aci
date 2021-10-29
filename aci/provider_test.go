package aci

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//TODO: check password is not showing in state file

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"aci": testAccProvider,
	}
}
func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	// We will use this function later on to make sure our test environment is valid.
	// For example, you can make sure here that some environment variables are set.
	if v := os.Getenv("ACI_USERNAME"); v == "" {
		t.Fatal("ACI_USERNAME env variable must be set for acceptance tests")
	}
	if v := os.Getenv("ACI_PASSWORD"); v == "" {
		privateKey := os.Getenv("ACI_PRIVATE_KEY")
		certName := os.Getenv("ACI_CERT_NAME")
		if privateKey == "" && certName == "" {
			t.Fatal("Either of ACI_PASSWORD or ACI_PRIVATE_KEY/ACI_CERT_NAME env variables must be set for acceptance tests")
		}
	}
	if v := os.Getenv("ACI_URL"); v == "" {
		t.Fatal("ACI_URL env variable must be set for acceptance tests")
	}
}

func TestAccAciProviderWithInvalidCredentials(t *testing.T) {
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTenantWithWrongPassword(rName),
				ExpectError: regexp.MustCompile(`User credential is incorrect - FAILED local authentication`),
			},
		},
	})
}

func CreateTenantWithWrongPassword(rName string) string {
	fmt.Println("=== STEP  testing creation with invalid credentials")
	resource := fmt.Sprintf(`
	provider "aci" {
		username = "%s"
		password = "%sxyz"
		url      = "%s"
		insecure = true
	}

	resource "aci_tenant" "test" {
		name = "%s"
	}`, os.Getenv("ACI_USERNAME"), os.Getenv("ACI_PASSWORD"), os.Getenv("ACI_URL"), rName)
	return resource
}
