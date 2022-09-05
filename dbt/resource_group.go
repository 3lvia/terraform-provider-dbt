package dbt

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dbtgroup "terraform-provider-dbt/dbt/group"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"assign_by_default": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"sso_mapping_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_set": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
								value := i.(string)
								var diags diag.Diagnostics

								validGroupPermissions := []string{
									"owner",
									"member",
									"account_admin",
									"admin",
									"database_admin",
									"git_admin",
									"team_admin",
									"job_admin",
									"job_viewer",
									"analyst",
									"developer",
									"stakeholder",
									"readonly",
									"project_creator",
									"account_viewer",
									"metadata_only",
									"webhooks_only",
								}

								for _, val := range validGroupPermissions {
									if val == value {
										return diags
									}
								}

								return append(diags, diag.Diagnostic{
									Severity: diag.Warning,
									Summary:  "Group Permission not valid",
									Detail:   fmt.Sprintf("%q is not a valid group permission. Must be one of: [%s]", value, strings.Join(validGroupPermissions, ", ")),
								})
							},
						},
						"project_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"all_projects": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*DbtProviderInput)
	groupInput := readGroupFromResourceData(d, providerInput.AccountId)

	group, diags := dbtgroup.CreateGroup(groupInput, providerInput.ServiceToken)
	if diags != nil {
		return diags
	}
	groupPermisisonsInput := readGroupPermissionsFromResourceData(d, group.Id, group.AccountId)

	if group != nil {
		setStateFromGroup(d, group)
	}

	groupPermissions, diags := dbtgroup.CreateOrUpdateGroupPermissions(groupPermisisonsInput, group.Id, group.AccountId, providerInput.ServiceToken)
	if diags != nil {
		return diags
	}

	d.Set("group_permissions", flattenGroupPermissions(groupPermissions))

	return diags
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*DbtProviderInput)
	groupInput := readGroupFromResourceData(d, providerInput.AccountId)

	group, diags := dbtgroup.ReadGroup(groupInput, providerInput.ServiceToken)

	if diags != nil {
		return diags
	}

	if group != nil {
		setStateFromGroup(d, group)
	} else {
		d.SetId("")
	}

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	providerInput := m.(*DbtProviderInput)
	groupInput := readGroupFromResourceData(d, providerInput.AccountId)
	groupPermisisonsInput := readGroupPermissionsFromResourceData(d, groupInput.Id, groupInput.AccountId)

	if groupHasChange(d) {
		group, diags := dbtgroup.UpdateGroup(groupInput, providerInput.ServiceToken)
		if group != nil {
			setStateFromGroup(d, group)
		}
		if diags != nil {
			return diags
		}
	}

	if d.HasChange("group_permissions") {
		groupPermissions, diags := dbtgroup.CreateOrUpdateGroupPermissions(groupPermisisonsInput, groupInput.Id, groupInput.AccountId, providerInput.ServiceToken)
		d.Set("group_permissions", flattenGroupPermissions(groupPermissions))
		return diags
	}

	return diags
}

func groupHasChange(d *schema.ResourceData) bool {
	hasChange := d.HasChange("assign_by_default") ||
		d.HasChange("name") ||
		d.HasChange("account_id") ||
		d.HasChange("sso_mapping_groups")

	return hasChange
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*DbtProviderInput)
	groupInput := readGroupFromResourceData(d, providerInput.AccountId)

	diags := dbtgroup.DeleteGroup(groupInput, providerInput.ServiceToken)

	d.SetId("")

	return diags
}

func setStateFromGroup(d *schema.ResourceData, group *dbtgroup.Group) {
	d.SetId(strconv.Itoa(group.Id))
	d.Set("assign_by_default", group.AssignByDefault)
	d.Set("name", group.Name)
	d.Set("account_id", group.AccountId)
	d.Set("sso_mapping_groups", group.SsoMappingGroups)
	d.Set("group_permissions", flattenGroupPermissions(group.GroupPermissions))
}

func readGroupFromResourceData(data *schema.ResourceData, accountId int) *dbtgroup.Group {
	id, _ := strconv.Atoi(data.Id())

	group := &dbtgroup.Group{
		Id:               id,
		AccountId:        accountId,
		Name:             data.Get("name").(string),
		AssignByDefault:  data.Get("assign_by_default").(bool),
		SsoMappingGroups: getStringArrayFromResourceSet(data, "sso_mapping_groups"),
	}

	return group
}

func readGroupPermissionsFromResourceData(data *schema.ResourceData, groupId int, accountId int) *[]dbtgroup.GroupPermission {
	rawGroupPermissions := data.Get("group_permissions").(*schema.Set).List()
	groupPermissions := []dbtgroup.GroupPermission{}
	for _, item := range rawGroupPermissions {
		gp := item.(map[string]interface{})
		permission := dbtgroup.GroupPermission{
			GroupId:       groupId,
			AccountId:     accountId,
			PermissionSet: gp["permission_set"].(string),
			ProjectId:     gp["project_id"].(int),
			AllProjects:   gp["all_projects"].(bool),
		}

		groupPermissions = append(groupPermissions, permission)
	}
	return &groupPermissions
}

func flattenGroupPermissions(groupPermissions *[]dbtgroup.GroupPermission) []interface{} {
	if groupPermissions == nil {
		return make([]interface{}, 0)
	}

	permissions := make([]interface{}, len(*groupPermissions))
	for i, permission := range *groupPermissions {
		p := make(map[string]interface{})

		p["permission_set"] = permission.PermissionSet
		p["project_id"] = permission.ProjectId
		p["all_projects"] = permission.AllProjects

		permissions[i] = p
	}
	return permissions
}

func getStringArrayFromResourceSet(d *schema.ResourceData, name string) []string {
	rawList := d.Get(name).(*schema.Set).List()
	stringList := make([]string, len(rawList))
	for i, v := range rawList {
		stringList[i] = v.(string)
	}
	return stringList
}
