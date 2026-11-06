package services

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

const servicesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "created": "2023-04-19T00:00:00UTC",
            "updated": "2023-04-19T00:00:00UTC",
            "domain": "example-go-worker.example-go.svc",
            "ptype": "worker",
            "ports": [
                {
                    "name": "tcp-5000",
                    "port": 5000,
                    "protocol": "TCP",
                    "target_port": 5000
                }
            ]
        }
    ]
}`

const (
	serviceApplyExpected string = `{"ptype":"worker","ports":[{"port":5000,"protocol":"TCP","target_port":5000},{"port":80,"protocol":"UDP","target_port":8080}]}`
	serviceInfoResponse  string = `{"app":"example-go","domain":"example-go-worker.example-go.svc","ptype":"worker","ports":[{"name":"tcp-5000","port":5000,"protocol":"TCP","target_port":5000},{"name":"udp-80","port":80,"protocol":"UDP","target_port":8080}]}`
)

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/services/" && req.Method == "GET" {
		res.Write([]byte(servicesFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/services/worker/" && req.Method == "PUT" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if string(body) != serviceApplyExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", serviceApplyExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte(serviceInfoResponse))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/services/worker/" && req.Method == "GET" {
		res.Write([]byte(serviceInfoResponse))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/services/worker/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestServicesList(t *testing.T) {
	t.Parallel()

	expected := api.Services{
		{
			App:     "example-go",
			Created: "2023-04-19T00:00:00UTC",
			Updated: "2023-04-19T00:00:00UTC",
			Domain:  "example-go-worker.example-go.svc",
			Ptype:   "worker",
			Ports: []api.ServicePort{
				{Port: 5000, Protocol: "TCP", TargetPort: 5000},
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

	actual, count, err := List(drycc, "example-go", 100)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected count 1, got %d", count)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestServicesApply(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	req := api.ServiceUpdateRequest{
		Ptype: "worker",
		Ports: []api.ServicePort{
			{Port: 5000, Protocol: "TCP", TargetPort: 5000},
			{Port: 80, Protocol: "UDP", TargetPort: 8080},
		},
	}

	info, err := Apply(drycc, "example-go", req)
	if err != nil {
		t.Fatal(err)
	}
	if info.Domain != "example-go-worker.example-go.svc" {
		t.Fatalf("Unexpected domain %s", info.Domain)
	}
	if len(info.Ports) != 2 {
		t.Fatalf("Expected 2 ports, got %d", len(info.Ports))
	}
}

func TestServicesInfo(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	info, err := Info(drycc, "example-go", "worker")
	if err != nil {
		t.Fatal(err)
	}
	if info.Ptype != "worker" {
		t.Fatalf("Expected worker, got %s", info.Ptype)
	}
	if len(info.Ports) != 2 {
		t.Fatalf("Expected 2 ports, got %d", len(info.Ports))
	}
}

func TestServicesRemove(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", "worker"); err != nil {
		t.Fatal(err)
	}
}
