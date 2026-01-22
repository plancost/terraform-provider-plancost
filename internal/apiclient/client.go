/*
Copyright (c) 2026 Plancost.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/version"
	"github.com/tidwall/gjson"
)

var sharedTransport *http.Transport

func init() {
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		sharedTransport = t.Clone()
		sharedTransport.MaxIdleConns = 100
		sharedTransport.MaxIdleConnsPerHost = 100
	} else {
		sharedTransport = &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		}
	}
}

type APIClient struct {
	httpClient *http.Client
	endpoint   string
	apiKey     string
	uuid       uuid.UUID
}

type GraphQLQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type Project struct {
	ID           string `json:"id"`
	FriendlyName string `json:"friendlyName"`
}

// NewAPIClient creates a new APIClient with the given endpoint and API key.
func NewAPIClient(endpoint, apiKey string) *APIClient {
	client := retryablehttp.NewClient()
	client.HTTPClient.Transport = sharedTransport
	client.HTTPClient.Timeout = time.Second * 30
	client.Logger = &LeveledLogger{Logger: logging.Logger}

	return &APIClient{
		httpClient: client.StandardClient(),
		endpoint:   endpoint,
		apiKey:     apiKey,
		uuid:       uuid.New(),
	}
}

func (c *APIClient) CreateProject(friendlyName string) (*Project, error) {
	reqBody := map[string]string{
		"friendlyName": friendlyName,
	}
	respBody, err := c.doRequest("POST", "/projects", reqBody)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.Unmarshal(respBody, &project); err != nil {
		return nil, errors.Wrap(err, "Error parsing project response")
	}
	return &project, nil
}

func (c *APIClient) GetProject(id string) (*Project, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("/projects/%s", id), nil)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.Unmarshal(respBody, &project); err != nil {
		return nil, errors.Wrap(err, "Error parsing project response")
	}
	return &project, nil
}

func (c *APIClient) DeleteProject(id string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/projects/%s", id), nil)
	return err
}

type ProjectQuota struct {
	Limit     int  `json:"limit"`
	Used      int  `json:"used"`
	Remaining int  `json:"remaining"`
	IsPaid    bool `json:"isPaid"`
}

func (c *APIClient) GetProjectQuota() (*ProjectQuota, error) {
	respBody, err := c.doRequest("GET", "/projects/quota", nil)
	if err != nil {
		return nil, err
	}
	var quotaResp ProjectQuota
	if err := json.Unmarshal(respBody, &quotaResp); err != nil {
		return nil, errors.Wrap(err, "Error parsing project quota response")
	}
	return &quotaResp, nil
}

func (c *APIClient) DoQueries(queries []GraphQLQuery) ([]gjson.Result, error) {
	if len(queries) == 0 {
		logging.Logger.Debug().Msg("Skipping GraphQL request as no queries have been specified")
		return []gjson.Result{}, nil
	}

	respBody, err := c.doRequest("POST", "/graphql", queries)
	return gjson.ParseBytes(respBody).Array(), err
}

func (c *APIClient) doRequest(method string, path string, d interface{}) ([]byte, error) {
	logging.Logger.Debug().Msgf("'%s' request to '%s' using trace_id: '%s'", method, path, c.uuid.String())

	var bodyReader io.Reader
	if d != nil {
		reqBody, err := json.Marshal(d)
		if err != nil {
			return []byte{}, errors.Wrap(err, "Error generating request body")
		}
		bodyReader = bytes.NewBuffer(reqBody)
	}

	req, err := http.NewRequest(method, c.endpoint+path, bodyReader)
	if err != nil {
		return []byte{}, errors.Wrap(err, "Error generating request")
	}

	c.AddAuthHeaders(req)

	client := c.httpClient
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, errors.Wrap(err, "Error sending API request")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("invalid API response %s %s", method, path)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return []byte{}, fmt.Errorf("API request failed with status %d, response: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *APIClient) AddDefaultHeaders(req *http.Request) {
	req.Header.Set("content-type", "application/json")
	req.Header.Set("User-Agent", userAgent())
}

func (c *APIClient) AddAuthHeaders(req *http.Request) {
	c.AddDefaultHeaders(req)
	if strings.HasPrefix(c.apiKey, "ics") {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	} else {
		req.Header.Set("X-Api-Key", c.apiKey)
	}
}

func userAgent() string {
	providerVersion := version.Version
	if providerVersion == "" {
		providerVersion = "dev"
	}

	return fmt.Sprintf("terraform-provider-plancost/%s", providerVersion)
}
