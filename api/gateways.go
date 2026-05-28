package api

// Gateway is the structure of an app's gateways.
type Gateway struct {
	// App is the app name. It cannot be updated at all right now.
	App string `json:"app,omitempty"`
	// Created is the time that the application settings was created and cannot be updated.
	Created string `json:"created,omitempty"`
	// Updated is the last time the application settings was changed and cannot be updated.
	Updated string `json:"updated,omitempty"`
	// UUID is a unique string reflecting the application settings in its current state.
	// It changes every time the application settings is changed and cannot be updated.
	UUID      string        `json:"uuid,omitempty"`
	Name      string        `json:"name,omitempty"`
	Ports     []GatewayPort `json:"ports,omitempty"`
	Addresses []Address     `json:"addresses,omitempty"`
}

// Address represents a gateway address configuration.
type Address struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

// Gateways defines a collection of gateway objects.
type Gateways []Gateway

// GatewayUpdateRequest is the structure of PUT /v2/apps/<app id>/gateways/<name>/.
type GatewayUpdateRequest struct {
	App   string        `json:"app"`
	Name  string        `json:"name"`
	Ports []GatewayPort `json:"ports"`
}

// GatewayPort represents a port in a gateway apply request.
type GatewayPort struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

// GatewayInfo is the structure returned by GET /v2/apps/<app id>/gateways/<name>/.
type GatewayInfo struct {
	App       string        `json:"app,omitempty"`
	Name      string        `json:"name,omitempty"`
	Ports     []GatewayPort `json:"ports"`
	Addresses []Address     `json:"addresses,omitempty"`
}
