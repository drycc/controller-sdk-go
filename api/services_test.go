package api

import (
	"sort"
	"testing"
)

func TestServicesSortByNameWithDuplicatePtype(t *testing.T) {
	t.Parallel()

	services := Services{
		{Name: "worker-b", Ptype: "worker"},
		{Name: "worker-a", Ptype: "worker"},
	}
	sort.Sort(services)

	if services[0].Name != "worker-a" || services[1].Name != "worker-b" {
		t.Fatalf("expected services sorted by name, got %+v", services)
	}
}
