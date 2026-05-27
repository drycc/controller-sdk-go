package api

import (
	"sort"
	"testing"
)

func TestAppsSorted(t *testing.T) {
	apps := Apps{
		{Created: "2014-01-01T00:00:00UTC", ID: "Zulu", Workspace: "John", Updated: "2016-01-02", UUID: "d57be2ba-7ae2-4825-9ace-7c86cb893046"},
		{Created: "2014-01-01T00:00:00UTC", ID: "Alpha", Workspace: "John", Updated: "2016-01-02", UUID: "3d501190-1b8e-41ef-94c5-dd9a0bb707bb"},
		{Created: "2014-01-01T00:00:00UTC", ID: "Gamma", Workspace: "John", Updated: "2016-01-02", UUID: "41d95133-fd4d-4f4c-92a2-e454857371cc"},
		{Created: "2014-01-01T00:00:00UTC", ID: "Beta", Workspace: "John", Updated: "2016-01-02", UUID: "222ed1aa-e985-4bec-9966-a88215300661"},
	}

	sort.Sort(apps)
	expectedAppNames := []string{"Alpha", "Beta", "Gamma", "Zulu"}

	for i, app := range apps {
		if expectedAppNames[i] != app.ID {
			t.Errorf("Expected apps to be sorted %v, Got %v at index %v", expectedAppNames[i], app.ID, i)
		}
	}
}
