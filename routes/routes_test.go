package routes

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

const routesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "created": "2023-04-19T00:00:00UTC",
            "updated": "2023-04-19T00:00:00UTC",
            "name": "example-go",
            "kind": "HTTPRoute",
            "parent_refs": [
                {
                    "name": "example-go",
                    "port": 80
                }
            ],
            "rules": [{
				"backend_refs": [{
                    "kind": "Service",
                    "name": "example-go",
                    "port": 5000,
                    "weight": 100
                }]
            }]
        }
    ]
}`

const routeApplyExpected string = `{"app":"example-go","name":"example-go","kind":"HTTPRoute","parent_refs":[{"name":"example-go","port":80}],"rules":[{"backend_refs":[{"kind":"Service","name":"example-go","port":5000,"weight":100}]}]}`

const routeInfoResponse string = `{"name":"example-go","kind":"HTTPRoute","parent_refs":[{"name":"example-go","port":80}],"rules":[{"backend_refs":[{"kind":"Service","name":"example-go","port":5000,"weight":100}]}]}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/routes/" && req.Method == "GET" {
		res.Write([]byte(routesFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/routes/example-go/" && req.Method == "PUT" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if string(body) != routeApplyExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", routeApplyExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte(routeInfoResponse))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/routes/example-go/" && req.Method == "GET" {
		res.Write([]byte(routeInfoResponse))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/routes/example-go/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestRoutesList(t *testing.T) {
	t.Parallel()

	expected := api.Routes{
		{
			App:     "example-go",
			Created: "2023-04-19T00:00:00UTC",
			Name:    "example-go",
			Updated: "2023-04-19T00:00:00UTC",
			Kind:    "HTTPRoute",
			ParentRefs: []api.ParentRef{
				{
					Name: "example-go",
					Port: 80,
				},
			},
			Rules: []api.RouteRule{
				{
					"backend_refs": []map[string]any{{
						"kind":   "Service",
						"name":   "example-go",
						"port":   5000,
						"weight": 100,
					}},
				},
			},
		},
	}
	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := List(drycc, "example-go", 100)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%v\n", actual)
	fmt.Printf("%v\n", expected)
	fmt.Printf("%v\n", fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", expected))

	if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestRoutesApply(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	req := api.RouteUpdateRequest{
		Name: "example-go",
		Kind: "HTTPRoute",
		ParentRefs: []api.RouteParentRef{
			{Name: "example-go", Port: 80},
		},
		Rules: []api.RouteRule{
			{
				"backend_refs": []map[string]any{{
					"kind":   "Service",
					"name":   "example-go",
					"port":   5000,
					"weight": 100,
				}},
			},
		},
	}

	info, err := Apply(drycc, "example-go", req)
	if err != nil {
		t.Fatal(err)
	}

	if info.Kind != "HTTPRoute" {
		t.Fatalf("Expected HTTPRoute, got %s", info.Kind)
	}
}

func TestRoutesInfo(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	info, err := Info(drycc, "example-go", "example-go")
	if err != nil {
		t.Fatal(err)
	}

	if info.Kind != "HTTPRoute" {
		t.Fatalf("Expected HTTPRoute, got %s", info.Kind)
	}
}

func TestRoutesRemove(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", "example-go"); err != nil {
		t.Fatal(err)
	}
}
