package addons

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

const addonClassGetFixture string = `
{
    "name": "valkey",
    "description": "Valkey in-memory data store",
    "kind": "Valkey",
            "storage_model": "bundle",
    "plans": [
        {
            "name": "micro",
            "description": "1 GB, Standalone",
            "defaults": {},
            "overrides": {
                "shards": 1,
                "replicas": 0
            },
            "allow_create": ["users", "config", "tls.enabled"],
            "allow_update": ["users", "config", "tls.enabled"]
        },
        {
            "name": "small",
            "description": "2 GB, 1 Master + 1 Replica",
            "defaults": {},
            "overrides": {
                "shards": 1,
                "replicas": 1
            },
            "allow_create": ["users", "config", "tls.enabled"],
            "allow_update": ["users", "config", "tls.enabled"]
        }
    ]
}
`

const addonClassesListFixture string = `
{
    "count": 1,
    "results": [
        {
            "name": "valkey",
            "description": "Valkey in-memory data store",
            "kind": "Valkey",
    "storage_model": "bundle",
            "plans": [
                {
                    "name": "micro",
                    "description": "1 GB, Standalone",
                    "defaults": {},
                    "overrides": {
                        "shards": 1,
                        "replicas": 0
                    },
                    "allow_create": ["users", "config", "tls.enabled"],
                    "allow_update": ["users", "config", "tls.enabled"]
                }
            ]
        }
    ]
}
`

const addonInstanceGetFixture string = `
{
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
    "name": "my-addon",
    "app": "example-go",
    "plan": "micro",
    "kind": "Valkey",
    "multiplier": 1,
    "parameters": {},
    "created": "2025-07-20T00:00:00UTC",
    "updated": "2025-07-20T00:00:00UTC"
}
`

const addonInstancesListFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
            "name": "my-addon",
            "app": "example-go",
            "plan": "micro",
            "kind": "Valkey",
            "multiplier": 1,
            "parameters": {},
            "created": "2025-07-20T00:00:00UTC",
            "updated": "2025-07-20T00:00:00UTC"
        }
    ]
}
`

const addonInstanceUpsertFixture string = `
{
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
    "name": "my-addon",
    "app": "example-go",
    "plan": "micro",
    "kind": "Valkey",
    "multiplier": 1,
    "parameters": {},
    "created": "2025-07-20T00:00:00UTC",
    "updated": "2025-07-20T00:00:00UTC"
}
`

const addonInstanceUpsertExpected string = `{"kind":"Valkey","plan":"micro","parameters":{}}`

const addonInstanceConnectionFixture string = `{"host":"localhost","port":"6379"}`

type fakeHTTPServer struct{}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	// Addon classes - list
	if req.URL.Path == "/v2/addon-classes/" && req.Method == "GET" {
		res.Write([]byte(addonClassesListFixture))
		return
	}

	// Addon classes - get
	if req.URL.Path == "/v2/addon-classes/valkey/" && req.Method == "GET" {
		res.Write([]byte(addonClassGetFixture))
		return
	}

	// Addon instances - list
	if req.URL.Path == "/v2/apps/example-go/addons/" && req.Method == "GET" {
		res.Write([]byte(addonInstancesListFixture))
		return
	}

	// Addon instances - get
	if req.URL.Path == "/v2/apps/example-go/addons/my-addon/" && req.Method == "GET" {
		res.Write([]byte(addonInstanceGetFixture))
		return
	}

	// Addon instances - connection
	if req.URL.Path == "/v2/apps/example-go/addons/my-addon/connection/" && req.Method == "GET" {
		res.Write([]byte(addonInstanceConnectionFixture))
		return
	}

	// Addon instances - upsert (PUT)
	if req.URL.Path == "/v2/apps/example-go/addons/my-addon/" && req.Method == "PUT" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != addonInstanceUpsertExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", addonInstanceUpsertExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte(addonInstanceUpsertFixture))
		return
	}

	// Addon instances - delete
	if req.URL.Path == "/v2/apps/example-go/addons/my-addon/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestListClasses(t *testing.T) {
	t.Parallel()

	expected := api.AddonClasses{
		{
			Name:         "valkey",
			Description:  "Valkey in-memory data store",
			Kind:         "Valkey",
			StorageModel: "bundle",
			Plans: []api.AddonPlan{
				{
					Name:        "micro",
					Description: "1 GB, Standalone",
					Defaults:    map[string]any{},
					Overrides: map[string]any{
						"shards":   float64(1),
						"replicas": float64(0),
					},
					AllowCreate: []string{"users", "config", "tls.enabled"},
					AllowUpdate: []string{"users", "config", "tls.enabled"},
				},
			},
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, count, err := ListClasses(drycc, 100)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, Got %d", count)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestGetClass(t *testing.T) {
	t.Parallel()

	expected := api.AddonClass{
		Name:         "valkey",
		Description:  "Valkey in-memory data store",
		Kind:         "Valkey",
		StorageModel: "bundle",
		Plans: []api.AddonPlan{
			{
				Name:        "micro",
				Description: "1 GB, Standalone",
				Defaults:    map[string]any{},
				Overrides: map[string]any{
					"shards":   float64(1),
					"replicas": float64(0),
				},
				AllowCreate: []string{"users", "config", "tls.enabled"},
				AllowUpdate: []string{"users", "config", "tls.enabled"},
			},
			{
				Name:        "small",
				Description: "2 GB, 1 Master + 1 Replica",
				Defaults:    map[string]any{},
				Overrides: map[string]any{
					"shards":   float64(1),
					"replicas": float64(1),
				},
				AllowCreate: []string{"users", "config", "tls.enabled"},
				AllowUpdate: []string{"users", "config", "tls.enabled"},
			},
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := GetClass(drycc, "valkey")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestList(t *testing.T) {
	t.Parallel()

	expected := api.AddonInstances{
		{
			UUID:       "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			Name:       "my-addon",
			App:        "example-go",
			Plan:       "micro",
			Kind:       "Valkey",
			Multiplier: 1,
			Parameters: map[string]any{},
			Created:    "2025-07-20T00:00:00UTC",
			Updated:    "2025-07-20T00:00:00UTC",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, count, err := List(drycc, "example-go", 100)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, Got %d", count)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	expected := api.AddonInstance{
		UUID:       "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Name:       "my-addon",
		App:        "example-go",
		Plan:       "micro",
		Kind:       "Valkey",
		Multiplier: 1,
		Parameters: map[string]any{},
		Created:    "2025-07-20T00:00:00UTC",
		Updated:    "2025-07-20T00:00:00UTC",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Get(drycc, "example-go", "my-addon")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestUpsert(t *testing.T) {
	t.Parallel()

	expected := api.AddonInstance{
		UUID:       "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Name:       "my-addon",
		App:        "example-go",
		Plan:       "micro",
		Kind:       "Valkey",
		Multiplier: 1,
		Parameters: map[string]any{},
		Created:    "2025-07-20T00:00:00UTC",
		Updated:    "2025-07-20T00:00:00UTC",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	req := api.AddonInstanceUpsertRequest{
		Kind:       "Valkey",
		Plan:       "micro",
		Parameters: map[string]any{},
	}
	actual, err := Upsert(drycc, "example-go", "my-addon", req)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", "my-addon"); err != nil {
		t.Fatal(err)
	}
}

func TestConnection(t *testing.T) {
	t.Parallel()

	expected := api.AddonConnection{
		"host": "localhost",
		"port": "6379",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Connection(drycc, "example-go", "my-addon")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}
