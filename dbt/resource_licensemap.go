package dbt

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dbtlicensemap "terraform-provider-dbt/dbt/license_map"
	utils "terraform-provider-dbt/dbt/utils"
)

func resourceLicenseMap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLicenseMapCreate,
		ReadContext:   resourceLicenseMapRead,
		UpdateContext: resourceLicenseMapUpdate,
		DeleteContext: resourceLicenseMapDelete,
		Schema: map[string]*schema.Schema{
			"license_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
					value := i.(string)

					if value == "developer" || value == "read_only" {
						return diag.Diagnostics{}
					}

					return diag.Diagnostics{diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "License type is not valid",
						Detail:   fmt.Sprintf("%q is not a valid group permission. Must be either 'developer' or 'read_only'", value),
					}}
				},
			},
			"sso_license_mapping_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
				licenseType, ssoGroups, err := resourceServiceLicenseMapParseId(d.Id())

				if err != nil {
					return nil, err
				}

				d.Set("license_type", licenseType)
				d.Set("sso_license_mapping_groups", ssoGroups)

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceServiceLicenseMapParseId(id string) (string, []string, error) {
	parts := strings.SplitN(id, ":", -1)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" || !strings.HasPrefix(parts[1], "[") || !strings.HasSuffix(parts[1], "]") {
		return "", nil, fmt.Errorf("unexpected format of ID (%s), expected licenseType:[ssoGroup1 ssoGroup2]", id)
	}

	licenseType := parts[0]
	ssoGroups := strings.Split(parts[1][1:len(parts[1])-1], " ")

	return licenseType, ssoGroups, nil
}

func resourceLicenseMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	licenseType, mappingGroups, serviceToken, accountId := getInputData(d, m)

	licenseMap, diags := dbtlicensemap.CreateOrUpdateLicenseMap(
		accountId, licenseType, mappingGroups, nil, serviceToken)

	if diags != nil {
		return diags
	}

	setResourceData(d, licenseMap)

	return diags
}

func resourceLicenseMapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	licenseType, mappingGroups, serviceToken, accountId := getInputData(d, m)

	licenseMap, _ := dbtlicensemap.ReadLicenseMap(accountId, licenseType, mappingGroups, serviceToken)

	setResourceData(d, licenseMap)

	return nil
}

func resourceLicenseMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	licenseType, _, serviceToken, accountId := getInputData(d, m)

	old, new := d.GetChange("sso_license_mapping_groups")

	licenseMap, diags := dbtlicensemap.CreateOrUpdateLicenseMap(
		accountId, licenseType, utils.InterfaceToStringList(new), utils.InterfaceToStringList(old), serviceToken)

	if diags != nil {
		return diags
	}

	setResourceData(d, licenseMap)

	return diags
}

func resourceLicenseMapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	licenseType, mappingGroups, serviceToken, accountId := getInputData(d, m)

	_, diags := dbtlicensemap.CreateOrUpdateLicenseMap(
		accountId, licenseType, nil, mappingGroups, serviceToken)

	d.SetId("")

	return diags
}

func setResourceData(data *schema.ResourceData, licenseMap *dbtlicensemap.LicenseMap) {
	if licenseMap != nil {
		data.SetId(fmt.Sprintf("%s:%s", licenseMap.LicenseType, licenseMap.SsoLicenseMappingGroups))
		data.Set("license_type", licenseMap.LicenseType)
		data.Set("sso_license_mapping_groups", licenseMap.SsoLicenseMappingGroups)
	} else {
		data.SetId("")
	}
}

func getInputData(data *schema.ResourceData, m interface{}) (string, []string, string, int) {
	providerInput := m.(*DbtProviderInput)

	licenseType := data.Get("license_type").(string)
	ssoLicenseMappingGroups := utils.InterfaceToStringList(data.Get("sso_license_mapping_groups"))

	return licenseType, ssoLicenseMappingGroups, providerInput.ServiceToken, providerInput.AccountId
}
