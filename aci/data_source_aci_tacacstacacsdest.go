package aci




import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciTACACSDestination() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAciTACACSDestinationRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tacacs_monitoring_destination_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				
			},
            "auth_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				
			},
            "host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				
			},
            "key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				
			},
            "name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				
			},
            "port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				
			},
            	    
            })),
    }
}

func dataSourceAciTACACSDestinationRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	host := d.Get("host").(string)
	port := d.Get("port").(string)
    TACACSMonitoringDestinationGroupDn := d.Get("tacacs_monitoring_destination_group_dn").(string)
	rn := fmt.Sprintf("tacacsdest-%s-port-%s", host,port,)
    dn := fmt.Sprintf("%s/%s",TACACSMonitoringDestinationGroupDn, rn)
	tacacsTacacsDest, err := getRemoteTACACSDestination(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setTACACSDestinationAttributes(tacacsTacacsDest, d)
	return nil
}