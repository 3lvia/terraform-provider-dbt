package dbtusergroup

type UserGroup struct {
	Id                   int                    `json:"id,omitempty"`
	AccountId            int                    `json:"account_id"`
	Name                 string                 `json:"name"`
	AssignByDefault      bool                   `json:"assign_by_default"`
	SsoMappingUserGroups []string               `json:"sso_mapping_groups"`
	State                int                    `json:"state,omitempty"`
	UserGroupPermissions *[]UserGroupPermission `json:"group_permissions,omitempty"`
}

type GetUserGroupResponse struct {
	Data UserGroup `json:"data"`
}

type UserGroupPermission struct {
	UserGroupId   int    `json:"group_id"`
	AccountId     int    `json:"account_id"`
	PermissionSet string `json:"permission_set"`
	ProjectId     int    `json:"project_id"`
	AllProjects   bool   `json:"all_projects"`
}

type UserGroupPermissionsResponse struct {
	Data []UserGroupPermission `json:"data"`
}
