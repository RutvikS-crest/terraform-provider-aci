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

func resourceAciAAAAuthentication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciAAAAuthenticationCreate,
		UpdateContext: resourceAciAAAAuthenticationUpdate,
		ReadContext:   resourceAciAAAAuthenticationRead,
		DeleteContext: resourceAciAAAAuthenticationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAAAAuthenticationImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"def_role_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"assign-default-role",
					"no-login",
				}, false),
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteAAAAuthentication(client *client.Client, dn string) (*models.AAAAuthentication, error) {
	aaaAuthRealmCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaAuthRealm := models.AAAAuthenticationFromContainer(aaaAuthRealmCont)
	if aaaAuthRealm.DistinguishedName == "" {
		return nil, fmt.Errorf("AAAAuthentication %s not found", aaaAuthRealm.DistinguishedName)
	}
	return aaaAuthRealm, nil
}

func setAAAAuthenticationAttributes(aaaAuthRealm *models.AAAAuthentication, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaAuthRealm.DistinguishedName)
	d.Set("description", aaaAuthRealm.Description)
	aaaAuthRealmMap, err := aaaAuthRealm.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaAuthRealmMap["annotation"])
	d.Set("def_role_policy", aaaAuthRealmMap["defRolePolicy"])
	d.Set("name", aaaAuthRealmMap["name"])
	d.Set("name_alias", aaaAuthRealmMap["nameAlias"])
	return d, nil
}

func resourceAciAAAAuthenticationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaAuthRealm, err := getRemoteAAAAuthentication(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setAAAAuthenticationAttributes(aaaAuthRealm, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAAAAuthenticationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AAAAuthentication: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	aaaAuthRealmAttr := models.AAAAuthenticationAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaAuthRealmAttr.Annotation = Annotation.(string)
	} else {
		aaaAuthRealmAttr.Annotation = "{}"
	}

	if DefRolePolicy, ok := d.GetOk("def_role_policy"); ok {
		aaaAuthRealmAttr.DefRolePolicy = DefRolePolicy.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaAuthRealmAttr.Name = Name.(string)
	}
	aaaAuthRealm := models.NewAAAAuthentication(fmt.Sprintf("userext/authrealm"), "uni", desc, nameAlias, aaaAuthRealmAttr)
	aaaAuthRealm.Status = "modified"
	err := aciClient.Save(aaaAuthRealm)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaAuthRealm.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciAAAAuthenticationRead(ctx, d, m)
}

func resourceAciAAAAuthenticationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AAAAuthentication: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	aaaAuthRealmAttr := models.AAAAuthenticationAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaAuthRealmAttr.Annotation = Annotation.(string)
	} else {
		aaaAuthRealmAttr.Annotation = "{}"
	}

	if DefRolePolicy, ok := d.GetOk("def_role_policy"); ok {
		aaaAuthRealmAttr.DefRolePolicy = DefRolePolicy.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaAuthRealmAttr.Name = Name.(string)
	}
	aaaAuthRealm := models.NewAAAAuthentication(fmt.Sprintf("userext/authrealm"), "uni", desc, nameAlias, aaaAuthRealmAttr)
	aaaAuthRealm.Status = "modified"
	err := aciClient.Save(aaaAuthRealm)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaAuthRealm.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciAAAAuthenticationRead(ctx, d, m)
}

func resourceAciAAAAuthenticationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaAuthRealm, err := getRemoteAAAAuthentication(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setAAAAuthenticationAttributes(aaaAuthRealm, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciAAAAuthenticationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name aaaAuthRealm cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}
