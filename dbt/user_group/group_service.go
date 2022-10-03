package dbtusergroup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func CreateUserGroup(groupInput *UserGroup, serviceToken string) (*UserGroup, diag.Diagnostics) {
	url := fmt.Sprintf("https://cloud.getdbt.com/api/v3/accounts/%d/groups/", groupInput.AccountId)

	return CreateOrUpdateUserGroup(groupInput, serviceToken, url, http.StatusCreated)
}

func UpdateUserGroup(groupInput *UserGroup, serviceToken string) (*UserGroup, diag.Diagnostics) {
	url := fmt.Sprintf("https://cloud.getdbt.com/api/v3/accounts/%d/groups/%d/", groupInput.AccountId, groupInput.Id)

	return CreateOrUpdateUserGroup(groupInput, serviceToken, url, http.StatusOK)
}

func CreateOrUpdateUserGroup(groupInput *UserGroup, serviceToken string, url string, expectedStatusCode int) (*UserGroup, diag.Diagnostics) {
	var diags diag.Diagnostics

	response, err := PostAsJson(groupInput, url, serviceToken)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != expectedStatusCode {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "UserGroup PostAsJson return http error code in CreateOrUpdateUserGroup",
			Detail:   fmt.Sprintf("DBT returned StatusCode %d for (%s), message: %s", response.StatusCode, url, data),
		})
	}

	var groupResponse GetUserGroupResponse
	err = json.Unmarshal(data, &groupResponse)
	if err != nil {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not parse Dbt group response as json",
			Detail:   err.Error(),
		})
	}

	return &groupResponse.Data, nil
}

func ReadUserGroup(groupInput *UserGroup, serviceToken string) (*UserGroup, diag.Diagnostics) {
	var diags diag.Diagnostics

	url := fmt.Sprintf("https://cloud.getdbt.com/api/v3/accounts/%d/groups/%d/", groupInput.AccountId, groupInput.Id)

	response, err := GetRequest(url, serviceToken)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if response.StatusCode == 404 {
		return nil, nil
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "UserGroup GetRequest return http error code in ReadUserGroup",
			Detail:   fmt.Errorf("DBT returned StatusCode %d for (%s), message: %s", response.StatusCode, url, data).Error(),
		})
	}

	var groupResponse GetUserGroupResponse
	err = json.Unmarshal(data, &groupResponse)
	if err != nil {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not parse Dbt create group response as json",
			Detail:   err.Error(),
		})
	}

	return &groupResponse.Data, nil
}

func DeleteUserGroup(groupInput *UserGroup, serviceToken string) diag.Diagnostics {
	url := fmt.Sprintf("https://cloud.getdbt.com/api/v3/accounts/%d/groups/%d/", groupInput.AccountId, groupInput.Id)
	groupInput.State = 2

	response, err := PostAsJson(groupInput, url, serviceToken)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()

		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "UserGroup PostAsJson return http error code in CreateUserGroup",
				Detail:   fmt.Errorf("DBT returned StatusCode %d for (%s), message: %s", response.StatusCode, url, data).Error(),
			}}
	}

	return nil
}
