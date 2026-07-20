// Package addons provides methods for managing addon classes and addon instances.
package addons

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// ListClasses lists all addon classes in the catalog.
func ListClasses(c *drycc.Client, results int) (api.AddonClasses, int, error) {
	u := "/v2/addon-classes/"
	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.AddonClass{}, -1, reqErr
	}
	var classes []api.AddonClass
	if err := json.Unmarshal([]byte(body), &classes); err != nil {
		return []api.AddonClass{}, -1, err
	}
	return classes, count, reqErr
}

// GetClass retrieves a single addon class by name.
func GetClass(c *drycc.Client, name string) (api.AddonClass, error) {
	u := fmt.Sprintf("/v2/addon-classes/%s/", name)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.AddonClass{}, reqErr
	}
	defer res.Body.Close()

	class := api.AddonClass{}
	if err := json.NewDecoder(res.Body).Decode(&class); err != nil {
		return api.AddonClass{}, err
	}
	return class, reqErr
}

// List lists addon instances attached to an app.
func List(c *drycc.Client, appID string, results int) (api.AddonInstances, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/addons/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.AddonInstance{}, -1, reqErr
	}
	var instances []api.AddonInstance
	if err := json.Unmarshal([]byte(body), &instances); err != nil {
		return []api.AddonInstance{}, -1, err
	}
	return instances, count, reqErr
}

// Get retrieves a single addon instance by name.
func Get(c *drycc.Client, appID string, name string) (api.AddonInstance, error) {
	u := fmt.Sprintf("/v2/apps/%s/addons/%s/", appID, name)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.AddonInstance{}, reqErr
	}
	defer res.Body.Close()

	instance := api.AddonInstance{}
	if err := json.NewDecoder(res.Body).Decode(&instance); err != nil {
		return api.AddonInstance{}, err
	}
	return instance, reqErr
}

// Connection retrieves the connection Secret data for an addon instance.
func Connection(c *drycc.Client, appID string, name string) (api.AddonConnection, error) {
	u := fmt.Sprintf("/v2/apps/%s/addons/%s/connection/", appID, name)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return nil, reqErr
	}
	defer res.Body.Close()

	connection := api.AddonConnection{}
	if err := json.NewDecoder(res.Body).Decode(&connection); err != nil {
		return nil, err
	}
	return connection, reqErr
}

// Upsert creates or updates an addon instance.
// The controller returns 201 on create and 200 on update; both carry the same
// response body, so callers do not need to distinguish between the two.
func Upsert(c *drycc.Client, appID string, name string, req api.AddonInstanceUpsertRequest) (api.AddonInstance, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return api.AddonInstance{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/addons/%s/", appID, name)
	res, reqErr := c.Request("PUT", u, body)
	if reqErr != nil {
		return api.AddonInstance{}, reqErr
	}
	defer res.Body.Close()

	instance := api.AddonInstance{}
	if err = json.NewDecoder(res.Body).Decode(&instance); err != nil {
		return api.AddonInstance{}, err
	}
	return instance, reqErr
}

// Delete removes an addon instance from an app.
func Delete(c *drycc.Client, appID string, name string) error {
	u := fmt.Sprintf("/v2/apps/%s/addons/%s/", appID, name)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
