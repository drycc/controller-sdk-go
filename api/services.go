package api

// Service is the structure of an app's service.
type Service struct {
	App     string        `json:"app,omitempty"`
	Created string        `json:"created,omitempty"`
	Updated string        `json:"updated,omitempty"`
	Domain  string        `json:"domain,omitempty"`
	Name    string        `json:"name,omitempty"`
	Ptype   string        `json:"ptype,omitempty"`
	Ports   []ServicePort `json:"ports,omitempty"`
}

// ServicePort represents a port in a service apply request.
type ServicePort struct {
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	TargetPort int    `json:"target_port"`
}

// Services defines a collection of service objects.
type Services []Service

func (s Services) Len() int           { return len(s) }
func (s Services) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Services) Less(i, j int) bool { return s[i].Name < s[j].Name }

// ServiceUpdateRequest is the structure of PUT /v2/apps/<app id>/services/<name>/.
type ServiceUpdateRequest struct {
	Name  string        `json:"name"`
	Ptype string        `json:"ptype"`
	Ports []ServicePort `json:"ports"`
}

// ServiceInfo is the structure returned by GET /v2/apps/<app id>/services/<name>/.
type ServiceInfo struct {
	App    string        `json:"app,omitempty"`
	Domain string        `json:"domain,omitempty"`
	Name   string        `json:"name,omitempty"`
	Ptype  string        `json:"ptype"`
	Ports  []ServicePort `json:"ports"`
}
