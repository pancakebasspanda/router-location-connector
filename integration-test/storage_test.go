package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"router-location-connecter/api"
	"router-location-connecter/storage"
)

const (
	_redisAddress  = "localhost:6379"
	_redisPassword = ""
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	redisClient, err := storage.New(nil, _redisAddress, _redisPassword)
	if err != nil {
		panic(err)
	}
	// Run tests
	exitVal := m.Run()

	// Teardown
	if err := redisClient.FlushAll(ctx); err != nil {
		panic(err)
	}

	if err := redisClient.Close(); err != nil {
		panic(err)
	}

	// Exit
	os.Exit(exitVal)
}

func TestStorage_Add_Retrieve_Router(t *testing.T) {
	redisHandler, err := storage.New(nil, _redisAddress, _redisPassword)
	if err != nil {
		assert.NoError(t, err)
	}

	tests := []struct {
		name    string
		Rh      storage.Storage
		router  *api.Router
		wantErr string
	}{
		{
			name: "add router to storage",
			Rh:   redisHandler,
			router: &api.Router{
				ID:          1,
				Name:        "Router 1",
				LocationID:  1,
				RouterLinks: []int{2},
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// add router
			if err := tt.Rh.AddRouter(tt.router); err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}

			// retrieve added router router
			router, err := tt.Rh.GetRouter(tt.router.ID)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}

			assert.Equal(t, tt.router, router)

		})
	}
}

func TestStorage_Add_Retrieve_Location(t *testing.T) {
	redisHandler, err := storage.New(nil, _redisAddress, _redisPassword)
	if err != nil {
		assert.NoError(t, err)
	}

	tests := []struct {
		name     string
		Rh       storage.Storage
		location *api.Location
		wantErr  string
	}{
		{
			name: "add router to storage",
			Rh:   redisHandler,
			location: &api.Location{
				ID:       1,
				Postcode: "AB1 2CD",
				Name:     "Location A",
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// add router
			if err := tt.Rh.AddLocation(tt.location); err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}

			// retrieve added router router
			lcoation, err := tt.Rh.GetLocation(tt.location.ID)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}

			assert.Equal(t, tt.location, lcoation)

		})
	}

}

func TestStorage_Add_Retrieve_RouterLocationLinks(t *testing.T) {
	redisHandler, err := storage.New(nil, _redisAddress, _redisPassword)
	if err != nil {
		assert.NoError(t, err)
	}

	tests := []struct {
		name       string
		Rh         storage.Storage
		routerLink *api.RouterLocationLink
		wantErr    string
	}{
		{
			name: "add router to storage",
			Rh:   redisHandler,
			routerLink: &api.RouterLocationLink{
				UniqueID:   "Location A: Location B",
				Connection: fmt.Sprintf("[%s] <-> [%s]", "Location A", "Location B"),
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// add router
			if err := tt.Rh.AddRouterLocationLink(tt.routerLink); err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}

			// retrieve added router router
			routerLocationLink, err := tt.Rh.GetRouterLocationLink(tt.routerLink.UniqueID)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}

			assert.Equal(t, tt.routerLink, routerLocationLink)
		})
	}
}
