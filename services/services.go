// Package services provides methods for managing an app's services.
package services

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List services registered with an app.
func List(c *drycc.Client, appID string, results int) (api.Services, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/services/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Service{}, -1, reqErr
	}

	var services []api.Service
	if err := json.Unmarshal([]byte(body), &services); err != nil {
		return []api.Service{}, -1, err
	}

	return services, count, reqErr
}

// Apply creates or updates a service for an app.
func Apply(c *drycc.Client, appID string, req api.ServiceUpdateRequest) (api.ServiceInfo, error) {
	ptype := req.Ptype
	u := fmt.Sprintf("/v2/apps/%s/services/%s/", appID, ptype)

	body, err := json.Marshal(req)
	if err != nil {
		return api.ServiceInfo{}, err
	}

	res, reqErr := c.Request("PUT", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.ServiceInfo{}, reqErr
	}
	defer res.Body.Close()

	var info api.ServiceInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return api.ServiceInfo{}, err
	}

	return info, reqErr
}

// Info retrieves information about a service.
func Info(c *drycc.Client, appID string, ptype string) (api.ServiceInfo, error) {
	u := fmt.Sprintf("/v2/apps/%s/services/%s/", appID, ptype)

	res, err := c.Request("GET", u, nil)
	if err != nil {
		return api.ServiceInfo{}, err
	}
	defer res.Body.Close()

	var info api.ServiceInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return api.ServiceInfo{}, err
	}

	return info, nil
}

// Delete removes a service from an app.
func Delete(c *drycc.Client, appID string, ptype string) error {
	u := fmt.Sprintf("/v2/apps/%s/services/%s/", appID, ptype)

	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
