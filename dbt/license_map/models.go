package dbtlicensemap

type LicenseMap struct {
	Id                      int      `json:"id,omitempty"`
	AccountId               int      `json:"account_id"`
	LicenseType             string   `json:"license_type"`
	SsoLicenseMappingGroups []string `json:"sso_license_mapping_groups"`
	State                   int      `json:"state,omitempty"`
}

type GetLicenseMapsResponse struct {
	Data []LicenseMap `json:"data"`
}

type GetLicenseMapResponse struct {
	Data LicenseMap `json:"data"`
}
