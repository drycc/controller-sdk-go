// Package config provides methods for managing configuration of apps.
package volumes

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List list an app's volumes.
func List(c *drycc.Client, appID string, results int) (api.Volumes, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Volume{}, -1, reqErr
	}
	var volumes []api.Volume
	if err := json.Unmarshal([]byte(body), &volumes); err != nil {
		return []api.Volume{}, -1, err
	}
	return volumes, count, reqErr
}

// Get an app's volume.
func Get(c *drycc.Client, appID string, name string) (api.Volume, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/", appID, name)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Volume{}, reqErr
	}
	defer res.Body.Close()

	volume := api.Volume{}
	if err := json.NewDecoder(res.Body).Decode(&volume); err != nil {
		return volume, err
	}

	return volume, nil
}

// Create create an app's Volume.
func Create(c *drycc.Client, appID string, volume api.Volume) (api.Volume, error) {
	body, err := json.Marshal(volume)
	if err != nil {
		return api.Volume{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/volumes/", appID)
	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil {
		return api.Volume{}, reqErr
	}
	defer res.Body.Close()
	newVolume := api.Volume{}
	if err = json.NewDecoder(res.Body).Decode(&newVolume); err != nil {
		return api.Volume{}, err
	}
	return newVolume, reqErr
}

// Expand create an app's Volume.
func Expand(c *drycc.Client, appID string, volume api.Volume) (api.Volume, error) {
	body, err := json.Marshal(volume)
	if err != nil {
		return api.Volume{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/", appID, volume.Name)
	res, reqErr := c.Request("PATCH", u, body)
	if reqErr != nil {
		return api.Volume{}, reqErr
	}
	defer res.Body.Close()
	newVolume := api.Volume{}
	if err = json.NewDecoder(res.Body).Decode(&newVolume); err != nil {
		return api.Volume{}, err
	}
	return newVolume, reqErr
}

// Delete delete an app's Volume.
func Delete(c *drycc.Client, appID string, name string) error {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/", appID, name)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Mount mount an app's volume and creates a new release.
// This is a patching operation, which means when you call Mount() with an api.Volumes:
//
//   - If the variable does not exist, it will be set.
//   - If the variable exists, it will be overwritten.
//   - If the variable is set to nil, it will be unmount.
//   - If the variable was ignored in the api.Volumes, it will remain unchanged.
//
// Calling Mount() with an empty api.Volume will return a drycc.ErrConflict.
// Trying to Unmount a key that does not exist returns a drycc.ErrUnprocessable.
func Mount(c *drycc.Client, appID string, name string, volume api.Volume) (api.Volume, error) {
	body, err := json.Marshal(volume)
	if err != nil {
		return api.Volume{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/path/", appID, name)
	res, reqErr := c.Request("PATCH", u, body)
	if reqErr != nil {
		return api.Volume{}, reqErr
	}
	defer res.Body.Close()
	newVolume := api.Volume{}
	if err = json.NewDecoder(res.Body).Decode(&newVolume); err != nil {
		return api.Volume{}, err
	}
	return newVolume, reqErr
}
