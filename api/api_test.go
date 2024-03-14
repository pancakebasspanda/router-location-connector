package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetRouterLocationData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		httpClient *retryablehttp.Client
		server     *httptest.Server
		want       *RouterLocationData
		wantErr    string
	}{
		{
			name:       "get router location data successfully",
			httpClient: retryablehttp.NewClient(),
			server: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !strings.Contains(r.URL.Path, "test") {
						t.Errorf("Expected to request 'test', got: %s", r.URL.Path)
					}

					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
  "routers": [
    {
      "id": 1,
      "name": "citadel-01",
      "location_id": 1,
      "router_links": [
        1
      ]
    },
    {
      "id": 2,
      "name": "citadel-02",
      "location_id": 1,
      "router_links": []
    },
    {
      "id": 3,
      "name": "core-07",
      "location_id": 7,
      "router_links": [
        15
      ]
    }
  ],
  "locations": [
    {
      "id": 1,
      "postcode": "BE12 2ND",
      "name": "Birmingham Motorcycle Museum"
    }
  ]
}`))
				}))
			}(),
			want: &RouterLocationData{
				Routers: []Router{
					{
						ID:          1,
						Name:        "citadel-01",
						LocationID:  1,
						RouterLinks: []int{1},
					},
					{
						ID:          2,
						Name:        "citadel-02",
						LocationID:  1,
						RouterLinks: []int{},
					},
					{
						ID:          3,
						Name:        "core-07",
						LocationID:  7,
						RouterLinks: []int{15},
					},
				},
				Locations: []Location{
					{
						ID:       1,
						Postcode: "BE12 2ND",
						Name:     "Birmingham Motorcycle Museum",
					},
				},
			},
			wantErr: "",
		},
		{
			name:       "returns an error when json is malformed",
			httpClient: retryablehttp.NewClient(),
			server: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !strings.Contains(r.URL.Path, "test") {
						t.Errorf("Expected to request 'test', got: %s", r.URL.Path)
					}

					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
  "routers": [
    {
}`))
				}))
			}(),
			want:    nil,
			wantErr: "unmarshal response body unexpected end of JSON input",
		},
		{
			name:       "retries on on server returning an error",
			httpClient: retryablehttp.NewClient(),
			server: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !strings.Contains(r.URL.Path, "test") {
						t.Errorf("Expected to request 'test', got: %s", r.URL.Path)
					}

					w.WriteHeader(http.StatusInternalServerError)
				}))
			}(),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(WithMaxRetries(2),
				WithBaseURL(tt.server.URL+"/test"),
				WithTimeout(10*time.Second))

			if tt.name == "retries on on server returning an error" {
				tt.wantErr = fmt.Sprintf("performing get router location data request: GET %s/test giving up after 3 attempt(s)", tt.server.URL)
			}

			got, err := c.GetRouterLocationData(ctx)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}
			assert.Equal(t, tt.want, got)

			tt.server.Close()
		})
	}
}

func TestClient_prepareGetRequest(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		requestURL string
		want       *retryablehttp.Request
		wantErr    string
	}{
		{
			name:       "returns request url",
			requestURL: "test/test",
			want: func() *retryablehttp.Request {
				req, err := retryablehttp.NewRequest("GET", "test/test", nil)
				if err != nil {
					t.Errorf("Expected to request url, got: %s", err)
				}

				return req
			}(),
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}

			got, err := c.prepareGetRequest(ctx, tt.requestURL)
			assert.Equal(t, tt.want.URL, got.URL)
			assert.Equal(t, tt.want.Body, got.Body)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}
