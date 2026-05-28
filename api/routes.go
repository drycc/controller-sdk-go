package api

// Route is the structure of an app's route.
type Route struct {
	// App is the app name. It cannot be updated at all right now.
	App string `json:"app,omitempty"`
	// Created is the time that the application settings was created and cannot be updated.
	Created string `json:"created,omitempty"`
	// Updated is the last time the application settings was changed and cannot be updated.
	Updated string `json:"updated,omitempty"`
	// UUID is a unique string reflecting the application settings in its current state.
	// It changes every time the application settings is changed and cannot be updated.
	UUID       string      `json:"uuid,omitempty"`
	Name       string      `json:"name,omitempty"`
	Kind       string      `json:"kind,omitempty"`
	ParentRefs []ParentRef `json:"parent_refs,omitempty"`
	Rules      []RouteRule `json:"rules,omitempty"`
}

// ParentRef represents a reference to a parent gateway.
type ParentRef struct {
	Name string `json:"name,omitempty"`
	Port int    `json:"port,omitempty"`
}

// RouteRule represents a rule in a route configuration.
type RouteRule map[string]any

// Routes defines a collection of Route objects.
type Routes []Route

// RouteUpdateRequest is the structure of PUT /v2/apps/<app_id>/routes/<name>/.
type RouteUpdateRequest struct {
	App        string           `json:"app"`
	Name       string           `json:"name"`
	Kind       string           `json:"kind"`
	ParentRefs []RouteParentRef `json:"parent_refs"`
	Rules      []RouteRule      `json:"rules"`
}

// RouteParentRef represents a reference to a parent gateway in apply request.
type RouteParentRef struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

// RouteInfo is the structure returned by GET /v2/apps/<app_id>/routes/<name>/.
type RouteInfo struct {
	App        string           `json:"app,omitempty"`
	Name       string           `json:"name,omitempty"`
	Kind       string           `json:"kind"`
	ParentRefs []RouteParentRef `json:"parent_refs"`
	Rules      []RouteRule      `json:"rules"`
	Routable   *bool            `json:"routable,omitempty"`
}
