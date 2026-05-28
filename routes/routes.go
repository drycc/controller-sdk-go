// Package routes provides methods for managing an app's routes.
package routes

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List routes registered with an app.
func List(c *drycc.Client, appID string, results int) (api.Routes, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/routes/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Route{}, -1, reqErr
	}

	var routes []api.Route
	if err := json.Unmarshal([]byte(body), &routes); err != nil {
		return []api.Route{}, -1, err
	}

	return routes, count, reqErr
}

// Apply creates or updates a route for an app.
func Apply(c *drycc.Client, appID string, req api.RouteUpdateRequest) (api.RouteInfo, error) {
	name := req.Name
	req.App = appID
	u := fmt.Sprintf("/v2/apps/%s/routes/%s/", appID, name)

	body, err := json.Marshal(req)
	if err != nil {
		return api.RouteInfo{}, err
	}

	res, reqErr := c.Request("PUT", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.RouteInfo{}, reqErr
	}
	defer res.Body.Close()

	var info api.RouteInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return api.RouteInfo{}, err
	}

	return info, reqErr
}

// Info retrieves information about a route.
func Info(c *drycc.Client, appID string, name string) (api.RouteInfo, error) {
	u := fmt.Sprintf("/v2/apps/%s/routes/%s/", appID, name)

	res, err := c.Request("GET", u, nil)
	if err != nil {
		return api.RouteInfo{}, err
	}
	defer res.Body.Close()

	var info api.RouteInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return api.RouteInfo{}, err
	}

	return info, nil
}

// Delete removes a route from an app.
func Delete(c *drycc.Client, appID string, name string) error {
	u := fmt.Sprintf("/v2/apps/%s/routes/%s/", appID, name)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
