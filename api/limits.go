package api

// LimitSpec is the definition of GET /v2/limits/specs/
type LimitSpec struct {
	ID       string         `json:"id"`
	CPU      map[string]any `json:"cpu"`
	Memory   map[string]any `json:"memory"`
	Features map[string]any `json:"features"`
	Keywords []string       `json:"keywords"`
	Disabled bool           `json:"disabled"`
}

// LimitPlan is the definition of GET /v2/limits/plans/
type LimitPlan struct {
	ID       string         `json:"id"`
	Spec     LimitSpec      `json:"spec"`
	CPU      int            `json:"cpu"`
	Memory   int            `json:"memory"`
	Features map[string]any `json:"features"`
	Disabled bool           `json:"disabled"`
}
