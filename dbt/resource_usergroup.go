package dbt

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dbtusergroup "terraform-provider-dbt/dbt/group"
)

func resourceUserUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,
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

								validUserGroupPermissions := []string{
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

								for _, val := range validUserGroupPermissions {
									if val == value {
										return diags
									}
								}

								return append(diags, diag.Diagnostic{
									Severity: diag.Warning,
									Summary:  "UserGroup Permission not valid",
									Detail:   fmt.Sprintf("%q is not a valid group permission. Must be one of: [%s]", value, strings.Join(validUserGroupPermissions, ", ")),
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

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*DbtProviderInput)
	groupInput := readUserGroupFromResourceData(d, providerInput.AccountId)

	group, diags := dbtusergroup.CreateUserGroup(groupInput, providerInput.ServiceToken)
	if diags != nil {
		return diags
	}
	groupPermisisonsInput := readUserGroupPermissionsFromResourceData(d, group.Id, group.AccountId)

	if group != nil {
		setStateFromUserGroup(d, group)
	}

	groupPermissions, diags := dbtusergroup.CreateOrUpdateUserGroupPermissions(groupPermisisonsInput, group.Id, group.AccountId, providerInput.ServiceToken)
	if diags != nil {
		return diags
	}

	d.Set("group_permissions", flattenUserGroupPermissions(groupPermissions))

	return diags
}

func resourceUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*DbtProviderInput)
	groupInput := readUserGroupFromResourceData(d, providerInput.AccountId)

	group, diags := dbtusergroup.ReadUserGroup(groupInput, providerInput.ServiceToken)

	if diags != nil {
		return diags
	}

	if group != nil {
		setStateFromUserGroup(d, group)
	} else {
		d.SetId("")
	}

	return diags
}

func resourceUserGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	providerInput := m.(*DbtProviderInput)
	groupInput := readUserGroupFromResourceData(d, providerInput.AccountId)
	groupPermisisonsInput := readUserGroupPermissionsFromResourceData(d, groupInput.Id, groupInput.AccountId)

	if groupHasChange(d) {
		group, diags := dbtusergroup.UpdateUserGroup(groupInput, providerInput.ServiceToken)
		if group != nil {
			setStateFromUserGroup(d, group)
		}
		if diags != nil {
			return diags
		}
	}

	if d.HasChange("group_permissions") {
		groupPermissions, diags := dbtusergroup.CreateOrUpdateUserGroupPermissions(groupPermisisonsInput, groupInput.Id, groupInput.AccountId, providerInput.ServiceToken)
		d.Set("group_permissions", flattenUserGroupPermissions(groupPermissions))
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

func resourceUserGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*DbtProviderInput)
	groupInput := readUserGroupFromResourceData(d, providerInput.AccountId)

	diags := dbtusergroup.DeleteUserGroup(groupInput, providerInput.ServiceToken)

	d.SetId("")

	return diags
}

func setStateFromUserGroup(d *schema.ResourceData, group *dbtusergroup.UserGroup) {
	d.SetId(strconv.Itoa(group.Id))
	d.Set("assign_by_default", group.AssignByDefault)
	d.Set("name", group.Name)
	d.Set("account_id", group.AccountId)
	d.Set("sso_mapping_groups", group.SsoMappingUserGroups)
	d.Set("group_permissions", flattenUserGroupPermissions(group.UserGroupPermissions))
}

func readUserGroupFromResourceData(data *schema.ResourceData, accountId int) *dbtusergroup.UserGroup {
	id, _ := strconv.Atoi(data.Id())

	group := &dbtusergroup.UserGroup{
		Id:                   id,
		AccountId:            accountId,
		Name:                 data.Get("name").(string),
		AssignByDefault:      data.Get("assign_by_default").(bool),
		SsoMappingUserGroups: getStringArrayFromResourceSet(data, "sso_mapping_groups"),
	}

	return group
}

func readUserGroupPermissionsFromResourceData(data *schema.ResourceData, groupId int, accountId int) *[]dbtusergroup.UserGroupPermission {
	rawUserGroupPermissions := data.Get("group_permissions").(*schema.Set).List()
	groupPermissions := []dbtusergroup.UserGroupPermission{}
	for _, item := range rawUserGroupPermissions {
		gp := item.(map[string]interface{})
		permission := dbtusergroup.UserGroupPermission{
			UserGroupId:   groupId,
			AccountId:     accountId,
			PermissionSet: gp["permission_set"].(string),
			ProjectId:     gp["project_id"].(int),
			AllProjects:   gp["all_projects"].(bool),
		}

		groupPermissions = append(groupPermissions, permission)
	}
	return &groupPermissions
}

func flattenUserGroupPermissions(groupPermissions *[]dbtusergroup.UserGroupPermission) []interface{} {
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
