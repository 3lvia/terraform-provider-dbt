package dbtgroup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func CreateOrUpdateGroupPermissions(groupPermissionsInput *[]GroupPermission, groupId int, accountId int, serviceToken string) (*[]GroupPermission, diag.Diagnostics) {
	url := fmt.Sprintf("https://cloud.getdbt.com/api/v3/accounts/%d/group-permissions/%d/", accountId, groupId)
	response, err := PostAsJson(groupPermissionsInput, url, serviceToken)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Dbt returned an http error code in CreateOrUpdatePermissions",
			Detail:   fmt.Errorf("Returned statusCode %d for %s, message: %s", response.StatusCode, url, data).Error(),
		}}
	}

	var groupPermissionsResponse GroupPermissionsResponse
	err = json.Unmarshal(data, &groupPermissionsResponse)
	if err != nil {
		return nil, diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not parse Dbt groupPermissions response",
			Detail:   err.Error(),
		}}
	}

	return &groupPermissionsResponse.Data, nil
}
