package dbtlicensemap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"terraform-provider-dbt/dbt/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var lock sync.Mutex

func ReadLicenseMap(accountId int, licenseType string, ssoLicenseMappingGroups []string, serviceToken string) (*LicenseMap, diag.Diagnostics) {
	licenseMap, diags := readLicenseMapFromLicenseType(accountId, serviceToken, licenseType)

	if diags != nil {
		return nil, diags
	}

	if licenseMap == nil {
		return nil, nil
	}

	var groups []string
	for _, val := range ssoLicenseMappingGroups {
		if utils.Contains(licenseMap.SsoLicenseMappingGroups, val) {
			groups = append(groups, val)
		}
	}

	licenseMap.SsoLicenseMappingGroups = groups

	return licenseMap, nil
}

func readLicenseMapFromLicenseType(accountId int, serviceToken string, licenseType string) (*LicenseMap, diag.Diagnostics) {
	url := fmt.Sprintf("https://cloud.getdbt.com/api/v3/accounts/%d/license-maps/", accountId)

	getLicenseMapsResponse, err := utils.GetAsObject[GetLicenseMapsResponse](url, serviceToken)

	if err != nil {
		return nil, diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading Licence maps",
			Detail:   err.Error(),
		}}
	}

	for _, val := range getLicenseMapsResponse.Data {
		if val.LicenseType == licenseType {
			return &val, nil
		}
	}

	return nil, nil
}

func CreateOrUpdateLicenseMap(accountId int, licenseType string, mappingsToAdd []string, mappingsToRemove []string, serviceToken string) (*LicenseMap, diag.Diagnostics) {
	lock.Lock()
	defer lock.Unlock()
	existingLicenceMap, diags := readLicenseMapFromLicenseType(accountId, serviceToken, licenseType)

	if diags != nil {
		return nil, diags
	}

	var url string
	var expectedStatusCode int

	request := LicenseMap{
		AccountId:   accountId,
		LicenseType: licenseType,
	}

	if existingLicenceMap == nil {
		url = fmt.Sprintf("https://cloud.getdbt.com/api/v3/accounts/%d/license-maps/", accountId)
		expectedStatusCode = http.StatusCreated

		request.SsoLicenseMappingGroups = mappingsToAdd
	} else {
		url = fmt.Sprintf("https://cloud.getdbt.com/api/v3/accounts/%d/license-maps/%d/", accountId, existingLicenceMap.Id)
		expectedStatusCode = http.StatusOK

		request.Id = existingLicenceMap.Id
		request.SsoLicenseMappingGroups = utils.RemoveFromList(append(existingLicenceMap.SsoLicenseMappingGroups, mappingsToAdd...), mappingsToRemove)

		if len(request.SsoLicenseMappingGroups) == 0 {
			request.State = 2
		}
	}

	response, err := utils.PostAsJson(request, url, serviceToken)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != expectedStatusCode {
		body, _ := json.Marshal(request)
		return nil, diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "CreateOrUpdateLicenseMap returned error code",
			Detail:   fmt.Sprintf("DBT returned StatusCode %d for (%s), message: %s, body: %s", response.StatusCode, url, data, body),
		}}
	}

	var getLicenseMapResponse GetLicenseMapResponse
	err = json.Unmarshal(data, &getLicenseMapResponse)

	if err != nil {
		return nil, diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not parse Dbt licenseMap response as json",
			Detail:   err.Error(),
		}}
	}

	return &getLicenseMapResponse.Data, nil
}
