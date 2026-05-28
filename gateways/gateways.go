// Package gateways provides methods for managing an app's gateways.
package gateways

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List gateways registered with an app.
func List(c *drycc.Client, appID string, results int) (api.Gateways, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/gateways/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Gateway{}, -1, reqErr
	}

	var gateways []api.Gateway
	if err := json.Unmarshal([]byte(body), &gateways); err != nil {
		return []api.Gateway{}, -1, err
	}

	return gateways, count, reqErr
}

// Apply creates or updates a gateway for an app.
func Apply(c *drycc.Client, appID string, req api.GatewayUpdateRequest) (api.GatewayInfo, error) {
	name := req.Name
	req.App = appID
	u := fmt.Sprintf("/v2/apps/%s/gateways/%s/", appID, name)

	body, err := json.Marshal(req)
	if err != nil {
		return api.GatewayInfo{}, err
	}

	res, reqErr := c.Request("PUT", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.GatewayInfo{}, reqErr
	}
	defer res.Body.Close()

	var info api.GatewayInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return api.GatewayInfo{}, err
	}

	return info, reqErr
}

// Info retrieves information about a gateway.
func Info(c *drycc.Client, appID string, name string) (api.GatewayInfo, error) {
	u := fmt.Sprintf("/v2/apps/%s/gateways/%s/", appID, name)

	res, err := c.Request("GET", u, nil)
	if err != nil {
		return api.GatewayInfo{}, err
	}
	defer res.Body.Close()

	var info api.GatewayInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return api.GatewayInfo{}, err
	}

	return info, nil
}

// Delete removes a gateway from an app.
func Delete(c *drycc.Client, appID string, name string) error {
	u := fmt.Sprintf("/v2/apps/%s/gateways/%s/", appID, name)

	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
