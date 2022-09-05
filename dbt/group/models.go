package dbtgroup

type Group struct {
	Id               int                `json:"id,omitempty"`
	AccountId        int                `json:"account_id"`
	Name             string             `json:"name"`
	AssignByDefault  bool               `json:"assign_by_default"`
	SsoMappingGroups []string           `json:"sso_mapping_groups"`
	State            int                `json:"state,omitempty"`
	GroupPermissions *[]GroupPermission `json:"group_permissions,omitempty"`
}

type GetGroupResponse struct {
	Data Group `json:"data"`
}

type GroupPermission struct {
	GroupId       int    `json:"group_id"`
	AccountId     int    `json:"account_id"`
	PermissionSet string `json:"permission_set"`
	ProjectId     int    `json:"project_id"`
	AllProjects   bool   `json:"all_projects"`
}

type GroupPermissionsResponse struct {
	Data []GroupPermission `json:"data"`
}
