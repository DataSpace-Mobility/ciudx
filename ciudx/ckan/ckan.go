package ckan

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	CKAN_ACTION_API = "/api/3/action"
)

type CkanResponseResult struct {
	Records []map[string]interface{} `json:"records,omitempty"`
}

type CkanResponse struct {
	Result CkanResponseResult `json:"result,omitempty"`
}

type Ckan struct {
	baseURL string
}

// InitCkan initiates CKAN link.
func InitCkan(baseURL string) *Ckan {
	if baseURL == "" {
		baseURL = "https://dataspace.mobi"
	}
	return &Ckan{
		baseURL: baseURL,
	}
}

func (ckan *Ckan) Search(resourceId string) ([]map[string]interface{}, error) {
	moduleURL := "/datastore_search"
	urlParams := "?resource_id=" + (strings.Split(resourceId, ":")[1])

	res, err := http.Get(ckan.baseURL + CKAN_ACTION_API + moduleURL + urlParams)
	if err != nil {
		return nil, err
	}

	results, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var jsonResponse CkanResponse
	err = json.Unmarshal(results, &jsonResponse)
	if err != nil {
		return nil, nil
	}

	return jsonResponse.Result.Records, nil
}
