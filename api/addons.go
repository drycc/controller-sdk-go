package api

// AddonPlan is a pricing/spec tier within an AddonClass.
type AddonPlan struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Defaults    map[string]any `json:"defaults,omitempty"`
	Overrides   map[string]any `json:"overrides,omitempty"`
	AllowCreate []string       `json:"allow_create,omitempty"`
	AllowUpdate []string       `json:"allow_update,omitempty"`
}

// AddonClass is a catalog entry returned by GET /v2/addon-classes.
type AddonClass struct {
	Name         string      `json:"name"`
	Description  string      `json:"description,omitempty"`
	Kind         string      `json:"kind,omitempty"`
	StorageModel string      `json:"storage_model,omitempty"`
	Plans        []AddonPlan `json:"plans,omitempty"`
}

// AddonClasses is a collection of AddonClass.
type AddonClasses []AddonClass

func (a AddonClasses) Len() int           { return len(a) }
func (a AddonClasses) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AddonClasses) Less(i, j int) bool { return a[i].Name < a[j].Name }

// AddonInstance is an addon bound to an app.
type AddonInstance struct {
	UUID       string         `json:"uuid,omitempty"`
	Name       string         `json:"name,omitempty"`
	App        string         `json:"app,omitempty"`
	Plan       string         `json:"plan,omitempty"`
	Kind       string         `json:"kind,omitempty"`
	Multiplier int            `json:"multiplier,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty"`
	Created    string         `json:"created,omitempty"`
	Updated    string         `json:"updated,omitempty"`
}

// AddonInstances is a collection of AddonInstance.
type AddonInstances []AddonInstance

func (a AddonInstances) Len() int           { return len(a) }
func (a AddonInstances) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AddonInstances) Less(i, j int) bool { return a[i].Name < a[j].Name }

// AddonInstanceUpsertRequest is the body of PUT /v2/apps/<app id>/addons/<name>/.
type AddonInstanceUpsertRequest struct {
	Kind       string         `json:"kind"`
	Plan       string         `json:"plan"`
	Parameters map[string]any `json:"parameters"`
}

// AddonConnection is the decoded connection Secret data for an addon instance,
// returned by GET /v2/apps/<app id>/addons/<name>/connection/.
type AddonConnection map[string]any
