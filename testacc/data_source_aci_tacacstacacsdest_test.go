package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTACACSAccountingDestinationDataSource_Basic(t *testing.T) {
	resourceName := "aci_tacacs_accounting_destination.test"
	dataSourceName := "data.aci_tacacs_accounting_destination.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	host := makeTestVariable(acctest.RandString(5))

	port := makeTestVariable(acctest.RandString(5))
	tacacsGroupName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSAccountingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTACACSAccountingDestinationDSWithoutRequired(tacacsGroupName, host, port, "tacacs_monitoring_destination_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTACACSAccountingDestinationDSWithoutRequired(tacacsGroupName, host, port, "host"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateTACACSAccountingDestinationDSWithoutRequired(tacacsGroupName, host, port, "port"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSAccountingDestinationConfigDataSource(tacacsGroupName, host, port),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tacacs_monitoring_destination_group_dn", resourceName, "tacacs_monitoring_destination_group_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "host", resourceName, "host"),
					resource.TestCheckResourceAttrPair(dataSourceName, "port", resourceName, "port"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_protocol", resourceName, "auth_protocol"),
					resource.TestCheckResourceAttrPair(dataSourceName, "key", resourceName, "key"),
				),
			},
			{
				Config:      CreateAccTACACSAccountingDestinationDataSourceUpdate(tacacsGroupName, host, port, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTACACSAccountingDestinationDSWithInvalidParentDn(tacacsGroupName, host, port),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccTACACSAccountingDestinationDataSourceUpdatedResource(tacacsGroupName, host, port, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccTACACSAccountingDestinationConfigDataSource(tacacsGroupName, host, port string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_monitoring_destination_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = "%s"
		port  = "%s"
	}

	data "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = aci_tacacs_accounting_destination.test.host
		port  = aci_tacacs_accounting_destination.test.port
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
	`, tacacsGroupName, host, port)
	return resource
}

func CreateTACACSAccountingDestinationDSWithoutRequired(tacacsGroupName, host, port, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting_destination Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tacacs_monitoring_destination_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = "%s"
		port  = "%s"
	}
	`
	switch attrName {
	case "tacacs_monitoring_destination_group_dn":
		rBlock += `
	data "aci_tacacs_accounting_destination" "test" {
	#	tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = aci_tacacs_accounting_destination.test.host	port  = aci_tacacs_accounting_destination.test.port
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
		`
	case "host":
		rBlock += `
	data "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
	#	host  = aci_tacacs_accounting_destination.test.host
		port  = aci_tacacs_accounting_destination.test.port
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
		`
	case "port":
		rBlock += `
	data "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = aci_tacacs_accounting_destination.test.host
	#	port  = aci_tacacs_accounting_destination.test.port
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, tacacsGroupName, host, port)
}

func CreateAccTACACSAccountingDestinationDSWithInvalidParentDn(tacacsGroupName, host, port string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_monitoring_destination_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = "%s"
		port  = "%s"
	}

	data "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = "${aci_tacacs_accounting_destination.test.host}_invalid"
		port  = "${aci_tacacs_accounting_destination.test.port}_invalid"
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
	`, tacacsGroupName, host, port)
	return resource
}

func CreateAccTACACSAccountingDestinationDataSourceUpdate(tacacsGroupName, host, port, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_monitoring_destination_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = "%s"
		port  = "%s"
	}

	data "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = aci_tacacs_accounting_destination.test.host
		port  = aci_tacacs_accounting_destination.test.port
		%s = "%s"
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
	`, tacacsGroupName, host, port, key, value)
	return resource
}

func CreateAccTACACSAccountingDestinationDataSourceUpdatedResource(tacacsGroupName, host, port, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_accounting_destination Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_monitoring_destination_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = "%s"
		port  = "%s"
		%s = "%s"
	}

	data "aci_tacacs_accounting_destination" "test" {
		tacacs_monitoring_destination_group_dn  = aci_tacacs_monitoring_destination_group.test.id
		host  = aci_tacacs_accounting_destination.test.host
		port  = aci_tacacs_accounting_destination.test.port
		depends_on = [ aci_tacacs_accounting_destination.test ]
	}
	`, tacacsGroupName, host, port, key, value)
	return resource
}
