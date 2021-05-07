package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciTriggerScheduler() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciTriggerSchedulerRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciTriggerSchedulerRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/schedp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	trigSchedP, err := getRemoteTriggerScheduler(aciClient, dn)

	if err != nil {
		return err
	}
	setTriggerSchedulerAttributes(trigSchedP, d)
	return nil
}
