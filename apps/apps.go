// Package apps provides methods for managing drycc apps.
package apps

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List lists apps on a Drycc controller.
func List(c *drycc.Client, workspace string, results int) (api.Apps, int, error) {
	body, count, reqErr := c.LimitedRequest(fmt.Sprintf("/v2/apps/?workspace=%s", workspace), results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.App{}, -1, reqErr
	}

	var apps []api.App
	if err := json.Unmarshal([]byte(body), &apps); err != nil {
		return []api.App{}, -1, err
	}

	return apps, count, reqErr
}

// New creates a new app with the given appID and workspace.
// Passing an empty appID will result in a randomized app name.
//
// If the app name already exists, the error drycc.ErrDuplicateApp will be returned.
func New(c *drycc.Client, appID string, workspace string) (api.App, error) {
	req := api.AppCreateRequest{ID: appID, Workspace: workspace}
	body, err := json.Marshal(req)
	if err != nil {
		return api.App{}, err
	}

	res, reqErr := c.Request("POST", "/v2/apps/", body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.App{}, reqErr
	}
	defer res.Body.Close()

	app := api.App{}
	if err := json.NewDecoder(res.Body).Decode(&app); err != nil {
		return api.App{}, err
	}

	return app, reqErr
}

// Get app details from a controller.
func Get(c *drycc.Client, appID string) (api.App, error) {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.App{}, reqErr
	}
	defer res.Body.Close()

	app := api.App{}

	if err := json.NewDecoder(res.Body).Decode(&app); err != nil {
		return api.App{}, err
	}

	return app, reqErr
}

// Quota gets an app's resource quota and current usage.
func Quota(c *drycc.Client, appID string) (api.AppQuota, error) {
	u := fmt.Sprintf("/v2/apps/%s/quota", appID)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.AppQuota{}, reqErr
	}
	defer res.Body.Close()

	quota := api.AppQuota{}
	if err := json.NewDecoder(res.Body).Decode(&quota); err != nil {
		return api.AppQuota{}, err
	}

	return quota, reqErr
}

// Run a one-time command in your app. This will start a kubernetes job with the
// same container image and environment as the rest of the app.
func Run(c *drycc.Client, appID string, command string, volumes map[string]any, timeout, expires uint32) error {
	req := api.AppRunRequest{
		Command: command,
		Volumes: volumes,
		Timeout: timeout,
		Expires: expires,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	u := fmt.Sprintf("/v2/apps/%s/run", appID)

	res, reqErr := c.Request("POST", u, body)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return reqErr
	}
	defer res.Body.Close()
	return reqErr
}

// Archive archives an app.
func Archive(c *drycc.Client, appID string) error {
	u := fmt.Sprintf("/v2/apps/%s/archive", appID)

	res, reqErr := c.Request("POST", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return reqErr
	}
	defer res.Body.Close()
	return reqErr
}

// Unarchive unarchives an app.
func Unarchive(c *drycc.Client, appID string) error {
	u := fmt.Sprintf("/v2/apps/%s/unarchive", appID)

	res, reqErr := c.Request("POST", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return reqErr
	}
	defer res.Body.Close()
	return reqErr
}

// Delete an app.
func Delete(c *drycc.Client, appID string) error {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Transfer moves an app to another workspace.
func Transfer(c *drycc.Client, appID string, workspace string) error {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	req := api.AppUpdateRequest{Workspace: workspace}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	res, err := c.Request("PATCH", u, body)
	if err == nil {
		res.Body.Close()
	}
	return err
}
