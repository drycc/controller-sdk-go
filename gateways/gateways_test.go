package gateways

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

const gatewaysFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "name": "example-go",
            "created": "2023-04-19T00:00:00UTC",
            "updated": "2023-04-19T00:00:00UTC",
			"ports": [
                {
                    "port": 80,
					"protocol": "HTTP"
                },
                {
                    "port": 443,
					"protocol": "HTTPS"
                }
            ],
            "addresses": [
                {
                    "type": "IPAddress",
                    "value": "172.22.108.207"
                }
            ]
        }
    ]
}`

const (
	gatewayApplyExpected string = `{"app":"example-go","name":"example-go","ports":[{"port":80,"protocol":"HTTP"},{"port":443,"protocol":"HTTPS"}]}`
	gatewayInfoResponse  string = `{"name":"example-go","ports":[{"port":80,"protocol":"HTTP"},{"port":443,"protocol":"HTTPS"}]}`
)

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/gateways/" && req.Method == "GET" {
		res.Write([]byte(gatewaysFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/gateways/example-go/" && req.Method == "PUT" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if string(body) != gatewayApplyExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", gatewayApplyExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte(gatewayInfoResponse))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/gateways/example-go/" && req.Method == "GET" {
		res.Write([]byte(gatewayInfoResponse))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/gateways/example-go/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestGatewaysList(t *testing.T) {
	t.Parallel()

	expected := api.Gateways{
		{
			App:     "example-go",
			Created: "2023-04-19T00:00:00UTC",
			Name:    "example-go",
			Updated: "2023-04-19T00:00:00UTC",
			Ports: []api.GatewayPort{
				{
					Port:     80,
					Protocol: "HTTP",
				},
				{
					Port:     443,
					Protocol: "HTTPS",
				},
			},
			Addresses: []api.Address{
				{
					Type:  "IPAddress",
					Value: "172.22.108.207",
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

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestGatewaysApply(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	req := api.GatewayUpdateRequest{
		Name: "example-go",
		Ports: []api.GatewayPort{
			{Port: 80, Protocol: "HTTP"},
			{Port: 443, Protocol: "HTTPS"},
		},
	}

	info, err := Apply(drycc, "example-go", req)
	if err != nil {
		t.Fatal(err)
	}

	if len(info.Ports) != 2 {
		t.Fatalf("Expected 2 ports, got %d", len(info.Ports))
	}
}

func TestGatewaysInfo(t *testing.T) {
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

	if len(info.Ports) != 2 {
		t.Fatalf("Expected 2 ports, got %d", len(info.Ports))
	}
}

func TestGatewaysRemove(t *testing.T) {
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
