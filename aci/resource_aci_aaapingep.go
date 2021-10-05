package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciDefaultRadiusAuthenticationSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciDefaultRadiusAuthenticationSettingsCreate,
		UpdateContext: resourceAciDefaultRadiusAuthenticationSettingsUpdate,
		ReadContext:   resourceAciDefaultRadiusAuthenticationSettingsRead,
		DeleteContext: resourceAciDefaultRadiusAuthenticationSettingsDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDefaultRadiusAuthenticationSettingsImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ping_check": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"false",
					"true",
				}, false),
			},
			"retries": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteDefaultRadiusAuthenticationSettings(client *client.Client, dn string) (*models.DefaultRadiusAuthenticationSettings, error) {
	aaaPingEpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaPingEp := models.DefaultRadiusAuthenticationSettingsFromContainer(aaaPingEpCont)
	if aaaPingEp.DistinguishedName == "" {
		return nil, fmt.Errorf("DefaultRadiusAuthenticationSettings %s not found", aaaPingEp.DistinguishedName)
	}
	return aaaPingEp, nil
}

func setDefaultRadiusAuthenticationSettingsAttributes(aaaPingEp *models.DefaultRadiusAuthenticationSettings, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaPingEp.DistinguishedName)
	d.Set("description", aaaPingEp.Description)
	aaaPingEpMap, err := aaaPingEp.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaPingEpMap["annotation"])
	d.Set("name", aaaPingEpMap["name"])
	d.Set("ping_check", aaaPingEpMap["pingCheck"])
	d.Set("retries", aaaPingEpMap["retries"])
	d.Set("timeout", aaaPingEpMap["timeout"])
	d.Set("name_alias", aaaPingEpMap["nameAlias"])
	return d, nil
}

func resourceAciDefaultRadiusAuthenticationSettingsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaPingEp, err := getRemoteDefaultRadiusAuthenticationSettings(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setDefaultRadiusAuthenticationSettingsAttributes(aaaPingEp, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDefaultRadiusAuthenticationSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DefaultRadiusAuthenticationSettings: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	aaaPingEpAttr := models.DefaultRadiusAuthenticationSettingsAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaPingEpAttr.Annotation = Annotation.(string)
	} else {
		aaaPingEpAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaPingEpAttr.Name = Name.(string)
	}

	if PingCheck, ok := d.GetOk("ping_check"); ok {
		aaaPingEpAttr.PingCheck = PingCheck.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaPingEpAttr.Retries = Retries.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaPingEpAttr.Timeout = Timeout.(string)
	}
	aaaPingEp := models.NewDefaultRadiusAuthenticationSettings(fmt.Sprintf("userext/pingext"), "uni", desc, nameAlias, aaaPingEpAttr)
	err := aciClient.Save(aaaPingEp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaPingEp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciDefaultRadiusAuthenticationSettingsRead(ctx, d, m)
}

func resourceAciDefaultRadiusAuthenticationSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DefaultRadiusAuthenticationSettings: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	aaaPingEpAttr := models.DefaultRadiusAuthenticationSettingsAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaPingEpAttr.Annotation = Annotation.(string)
	} else {
		aaaPingEpAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaPingEpAttr.Name = Name.(string)
	}

	if PingCheck, ok := d.GetOk("ping_check"); ok {
		aaaPingEpAttr.PingCheck = PingCheck.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaPingEpAttr.Retries = Retries.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaPingEpAttr.Timeout = Timeout.(string)
	}
	aaaPingEp := models.NewDefaultRadiusAuthenticationSettings(fmt.Sprintf("userext/pingext"), "uni", desc, nameAlias, aaaPingEpAttr)
	aaaPingEp.Status = "modified"
	err := aciClient.Save(aaaPingEp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaPingEp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciDefaultRadiusAuthenticationSettingsRead(ctx, d, m)
}

func resourceAciDefaultRadiusAuthenticationSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaPingEp, err := getRemoteDefaultRadiusAuthenticationSettings(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setDefaultRadiusAuthenticationSettingsAttributes(aaaPingEp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciDefaultRadiusAuthenticationSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaPingEp")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
