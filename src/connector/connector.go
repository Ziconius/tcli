package connector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/tines/go-sdk/tines"
)

type TinesAPI struct {
	API *APIClient
	SDK *tines.Client
}

type APIClient struct {
	tenantURL string
	apiKey    string
	// userAgent  string  // TODO: Implement.
	httpClient *http.Client // TODO: Resolve later
}

/*

Currently this module has a workaround as the Tines SDK does not support all of the required
	endpoints
*/

// Used to interact with the Tines API for pulling tcli enabled stories.
func NewTinesAPI(tenant, apiKey string) (TinesAPI, error) {

	tinesAPI := TinesAPI{}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	tinesAPI.API = &APIClient{
		tenantURL:  tenant,
		apiKey:     apiKey,
		httpClient: &client,
	}

	sdk, err := tines.NewClient(
		tines.SetTenantUrl(tenant),
		tines.SetApiKey(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	tinesAPI.SDK = sdk

	return tinesAPI, nil
}

func (api *APIClient) GetNotes() (Notes, error) {
	ep := fmt.Sprintf("%vapi/v1/notes", api.tenantURL)
	slog.Info("Getting notes", "ep", ep)

	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		slog.Error("fail new request", "error", err)

		return Notes{}, err
	}
	key := "Bearer " + api.apiKey

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("User-Agent", "tcliDev")
	req.Header.Set("Authorization", key)

	fmt.Printf("Headers %#v\n", req.Header)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		slog.Error("Failed DO", "error", err)

		return Notes{}, err
	}
	slog.Info("Status code info", "sc", resp.StatusCode)
	slog.Info("Data", "body", resp.Body)

	notes := Notes{}
	byteData, _ := io.ReadAll(resp.Body)
	fmt.Printf("DATA: %v\n", string(byteData))
	err = json.Unmarshal(byteData, &notes)
	if err != nil {
		return notes, err
	}

	return notes, nil
}

type Notes struct {
	Annotations []struct {
		ID        int    `json:"id"`
		Content   string `json:"content"`
		StoryMode string `json:"story_mode"`
		DraftName string `json:"draft_name,omitempty"`
	} `json:"annotations"`
}

// "annotations": [
//     {
//       "id": 21405,
//       "story_id": 813,
//       "story_mode": "LIVE",
//       "group_id": null,
//       "position": {
//         "x": 0,
//         "y": 0
//       },
//       "content": "# Example note",
//       "draft_id": 12345,
//       "draft_name": "My draft"
//     }
//   ]

// TODO: Iterate pages - current tenant has a single page of resources.
func (api *APIClient) ListResources() (ResourseList, error) {
	ep := fmt.Sprintf("%vapi/v1/global_resources", api.tenantURL)
	slog.Info("Getting resource list", "ep", ep)

	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		slog.Error("fail new request", "error", err)

		return ResourseList{}, err
	}
	key := "Bearer " + api.apiKey

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("User-Agent", "tcliDev")
	req.Header.Set("Authorization", key)

	fmt.Printf("Headers %#v\n", req.Header)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		slog.Error("Failed DO", "error", err)

		return ResourseList{}, err
	}
	slog.Info("Status code info", "sc", resp.StatusCode)
	slog.Info("Data", "body", resp.Body)

	notes := ResourseList{}
	byteData, _ := io.ReadAll(resp.Body)
	fmt.Printf("DATA: %v\n", string(byteData))
	err = json.Unmarshal(byteData, &notes)
	if err != nil {
		return notes, err
	}

	return notes, nil
}

type ResourseList struct {
	GlobalResources []Resource `json:"global_resources"`
}

type Resource struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Value string `json:"value"`
}

/*
{
  "global_resources": [
    {
      "id": 1,
      "name": "tines_jira_url",
      "value": "tinesio.atlassian.net",
      "team_id": 2,
      "folder_id": 1,
      "user_id": 1,
      "read_access": "TEAM",
      "shared_team_slugs": [],
      "slug": "tines_jira_url",
      "created_at": "2019-12-12T21:34:16.540Z",
      "updated_at": "2019-12-12T21:34:16.540Z",
      "description": "",
      "test_resource_enabled": true,
      "test_resource": {
        "value": "test resource",
        "created_at": "2019-12-12T22:34:16.540Z",
        "updated_at": "2019-12-12T22:34:16.540Z"
      },
      "referencing_action_ids": [431]
    }
  ],
  "meta": {
    "current_page": "https://ancient-resonance-5800.tines.com/api/v1/global_resources?per_page=20&page=1",
    "previous_page": null,
    "next_page": null,
    "next_page_number": null,
    "per_page": 20,
    "pages": 1,
    "count": 1
  }
}
*/
