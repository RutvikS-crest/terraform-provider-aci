package acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers("aci_tenant",
		&resource.Sweeper{
			Name: "aci_tenant",
			F:    aciTenantSweeper,
		})
}
