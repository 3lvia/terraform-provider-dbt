package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

func PostAsJson[T any](requestBody T, url string, serviceToken string) (*http.Response, error) {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceToken))

	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

func GetRequest(url string, serviceToken string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceToken))

	return client.Do(req)
}

func GetAsObject[T any](url string, serviceToken string) (*T, error) {

	response, err := GetRequest(url, serviceToken)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == 404 {
		return nil, nil
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Get request to: '%s' returned status code %d. message: %s", url, response.StatusCode, data)
	}

	var object T
	err = json.Unmarshal(data, &object)

	if err != nil {
		return nil, err
	}

	return &object, nil
}
