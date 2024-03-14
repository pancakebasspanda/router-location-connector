package app

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"router-location-connecter/api"
	mock_api "router-location-connecter/api/mocks"
	"router-location-connecter/storage"
	mock_storage "router-location-connecter/storage/mocks"
	"testing"
)

func Test_sortLocations(t *testing.T) {
	tests := []struct {
		name string
		loc1 string
		loc2 string
		want string
	}{
		{
			name: "sorts different locations alphabetically and returns them in a concatenated string",
			loc1: "Birmingham Motorcycle Museum",
			loc2: "Winterbourne House",
			want: "Birmingham Motorcycle Museum:Winterbourne House",
		},
		{
			name: "sorts same locations alphabetically and returns them in a concatenated string",
			loc1: "Birmingham Motorcycle Museum",
			loc2: "Birmingham Motorcycle Museum",
			want: "Birmingham Motorcycle Museum:Birmingham Motorcycle Museum",
		},
		{
			name: "sorts same locations alphabetically when locations start with same letter",
			loc1: "Birmingham Motorcycle Museum1",
			loc2: "Birmingham Motorcycle Museum2",
			want: "Birmingham Motorcycle Museum1:Birmingham Motorcycle Museum2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortLocations(tt.loc1, tt.loc2)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_app_CalculateLink(t *testing.T) {
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	mockController := gomock.NewController(t)

	storageMock := mock_storage.NewMockStorage(mockController)

	defer mockController.Finish()

	tests := []struct {
		name                string
		log                 zerolog.Logger
		storageMockOutcomes func(storageMock *mock_storage.MockStorage)
		srcLocationID       int
		destLocationID      int
		wantErr             string
	}{
		{
			name: "calculates location link between src and destination routers",
			log:  log,
			storageMockOutcomes: func(storageMock *mock_storage.MockStorage) {
				storageMock.EXPECT().
					GetLocation(1).
					Times(1).
					Return(&api.Location{
						ID:       1,
						Postcode: "BE13 1EQ",
						Name:     "Winterbourne House",
					}, nil)
				storageMock.EXPECT().
					GetLocation(2).
					Times(1).
					Return(&api.Location{
						ID:       2,
						Postcode: "BE12 2ND",
						Name:     "Birmingham Hippodrome",
					}, nil)
				storageMock.EXPECT().
					GetRouterLocationLink(sortLocations("Winterbourne House", "Birmingham Hippodrome")).
					Times(1).
					Return(nil, redis.Nil)
				storageMock.EXPECT().
					AddRouterLocationLink(&api.RouterLocationLink{
						UniqueID:   sortLocations("Winterbourne House", "Birmingham Hippodrome"),
						Connection: fmt.Sprintf("[%s] <-> [%s]", "Winterbourne House", "Birmingham Hippodrome"),
					}).Times(1).
					Return(nil)
			},
			srcLocationID:  1,
			destLocationID: 2,
			wantErr:        "",
		},
		{
			name: "calculates location link between src and destination routers, but doesnt save it when src and destination from the previous example are switched",
			log:  log,
			storageMockOutcomes: func(storageMock *mock_storage.MockStorage) {
				storageMock.EXPECT().
					GetLocation(1).
					Times(1).
					Return(&api.Location{
						ID:       1,
						Postcode: "BE13 1EQ",
						Name:     "Winterbourne House",
					}, nil)
				storageMock.EXPECT().
					GetLocation(2).
					Times(1).
					Return(&api.Location{
						ID:       2,
						Postcode: "BE12 2ND",
						Name:     "Birmingham Hippodrome",
					}, nil)
				storageMock.EXPECT().
					GetRouterLocationLink(sortLocations("Winterbourne House", "Birmingham Hippodrome")).
					Times(1).
					Return(&api.RouterLocationLink{
						UniqueID:   sortLocations("Winterbourne House", "Birmingham Hippodrome"),
						Connection: fmt.Sprintf("[%s] <-> [%s]", "Winterbourne House", "Birmingham Hippodrome"),
					}, nil)
			},
			srcLocationID:  2,
			destLocationID: 1,
			wantErr:        "",
		},
		{
			name: "when src and destination are the same",
			log:  log,
			storageMockOutcomes: func(storageMock *mock_storage.MockStorage) {
				storageMock.EXPECT().
					GetLocation(1).
					Times(2).
					Return(&api.Location{
						ID:       1,
						Postcode: "BE13 1EQ",
						Name:     "Winterbourne House",
					}, nil)
				storageMock.EXPECT().
					GetRouterLocationLink(sortLocations("Winterbourne House", "Winterbourne House")).
					Times(1).
					Return(nil, redis.Nil)
				storageMock.EXPECT().
					AddRouterLocationLink(&api.RouterLocationLink{
						UniqueID:   sortLocations("Winterbourne House", "Winterbourne House"),
						Connection: fmt.Sprintf("[%s] <-> [%s]", "Winterbourne House", "Winterbourne House"),
					}).Times(1).
					Return(nil)
			},
			srcLocationID:  1,
			destLocationID: 1,
			wantErr:        "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storageMockOutcomes(storageMock)
			a := &app{
				storage: storageMock,
				log:     tt.log,
			}
			if err := a.CalculateLink(tt.srcLocationID, tt.destLocationID); err != nil {
				assert.Equal(t, err, tt.wantErr)
			}
		})
	}
}

func Test_app_Process(t *testing.T) {
	ctx := context.Background()

	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	mockController := gomock.NewController(t)
	storageMock := mock_storage.NewMockStorage(mockController)
	apiMock := mock_api.NewMockAPI(mockController)

	defer mockController.Finish()

	tests := []struct {
		name                string
		api                 api.API
		storage             storage.Storage
		apiMockOutcomes     func(apiMock *mock_api.MockAPI)
		storageMockOutcomes func(storageMock *mock_storage.MockStorage)
		log                 zerolog.Logger
	}{
		{
			name: "process the API data, store router and location data and print the linked router locations",
			api:  apiMock,
			apiMockOutcomes: func(apiMock *mock_api.MockAPI) {
				apiMock.EXPECT().
					GetRouterLocationData(ctx).
					Times(1).
					Return(&api.RouterLocationData{
						Routers: []api.Router{
							{
								ID:          1,
								Name:        "Router A",
								LocationID:  1,
								RouterLinks: []int{2},
							},
							{
								ID:          2,
								Name:        "Router B",
								LocationID:  2,
								RouterLinks: []int{1, 3},
							},
							{
								ID:          3,
								Name:        "Router C",
								LocationID:  3,
								RouterLinks: []int{2},
							},
							{
								ID:          4,
								Name:        "Router D",
								LocationID:  1,
								RouterLinks: []int{1}, // at location 1, linked to 1
							},
						},
						Locations: []api.Location{
							{
								ID:       1,
								Postcode: "A",
								Name:     "Location A",
							},
							{
								ID:       2,
								Postcode: "B",
								Name:     "Location B",
							},
							{
								ID:       3,
								Postcode: "C",
								Name:     "Location C",
							},
						},
					}, nil)
			},
			storage: storageMock,
			storageMockOutcomes: func(storageMock *mock_storage.MockStorage) {
				storageMock.EXPECT().AddLocation(&api.Location{
					ID:       1,
					Postcode: "A",
					Name:     "Location A",
				}).Times(1).Return(nil)
				storageMock.EXPECT().AddLocation(&api.Location{
					ID:       2,
					Postcode: "B",
					Name:     "Location B",
				}).Times(1).Return(nil)
				storageMock.EXPECT().AddLocation(&api.Location{
					ID:       3,
					Postcode: "C",
					Name:     "Location C",
				}).Times(1).Return(nil)
				storageMock.EXPECT().
					GetLocation(1).
					Times(1).
					Return(&api.Location{
						ID:       1,
						Postcode: "A",
						Name:     "Location A",
					}, nil)
				storageMock.EXPECT().
					GetLocation(2).
					Times(2).
					Return(&api.Location{
						ID:       2,
						Postcode: "B",
						Name:     "Location B",
					}, nil)
				storageMock.EXPECT().
					GetLocation(3).
					Times(1).
					Return(&api.Location{
						ID:       3,
						Postcode: "C",
						Name:     "Location C",
					}, nil)
				storageMock.EXPECT().AddRouter(&api.Router{
					ID:          1,
					Name:        "Router A",
					LocationID:  1,
					RouterLinks: []int{2},
				}).Times(1).Return(nil)
				storageMock.EXPECT().AddRouter(&api.Router{
					ID:          2,
					Name:        "Router B",
					LocationID:  2,
					RouterLinks: []int{1, 3},
				}).Times(1).Return(nil)
				storageMock.EXPECT().AddRouter(&api.Router{
					ID:          3,
					Name:        "Router C",
					LocationID:  3,
					RouterLinks: []int{2},
				}).Times(1).Return(nil)
				storageMock.EXPECT().AddRouter(&api.Router{
					ID:          4,
					Name:        "Router D",
					LocationID:  1,
					RouterLinks: []int{1}, // at location 1, linked to 1
				}).Times(1).Return(nil)
				storageMock.EXPECT().GetRouter(2).Times(2).Return(&api.Router{
					ID:          2,
					Name:        "Router B",
					LocationID:  2,
					RouterLinks: []int{1, 3},
				}, nil)
				storageMock.EXPECT().GetRouter(3).Times(1).Return(&api.Router{
					ID:          3,
					Name:        "Router C",
					LocationID:  3,
					RouterLinks: []int{2},
				}, nil)
				storageMock.EXPECT().GetRouter(1).Times(1).Return(&api.Router{
					ID:          1,
					Name:        "Router A",
					LocationID:  1,
					RouterLinks: []int{1},
				}, nil)
				storageMock.EXPECT().
					GetRouterLocationLink("Location A:Location B").
					Times(1).
					Return(nil, redis.Nil)
				storageMock.EXPECT().AddRouterLocationLink(&api.RouterLocationLink{
					UniqueID:   "Location A:Location B",
					Connection: fmt.Sprintf("[%s] <-> [%s]", "Location A", "Location B"),
				}).Times(1).
					Return(nil)
				storageMock.EXPECT().
					GetRouterLocationLink("Location B:Location C").
					Times(1).
					Return(nil, redis.Nil)
				storageMock.EXPECT().AddRouterLocationLink(&api.RouterLocationLink{
					UniqueID:   sortLocations("Location B", "Location C"),
					Connection: fmt.Sprintf("[%s] <-> [%s]", "Location B", "Location C"),
				}).Times(1).
					Return(nil)
			},
			log: log,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storageMockOutcomes(storageMock)
			tt.apiMockOutcomes(apiMock)

			a := &app{
				apiClient: tt.api,
				storage:   tt.storage,
				log:       tt.log,
			}

			a.Process(ctx)
		})
	}
}

func Test_app_processLinkedRouter(t *testing.T) {
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	mockController := gomock.NewController(t)
	storageMock := mock_storage.NewMockStorage(mockController)
	defer mockController.Finish()

	tests := []struct {
		name                string
		parentRouter        *api.Router
		linkedRouter        *api.Router
		storage             storage.Storage
		storageMockOutcomes func(storageMock *mock_storage.MockStorage)
		log                 zerolog.Logger
		processedRouters    map[int]struct{}
	}{
		{
			name: "processes linked router connections as well as all other links in the linked router",
			parentRouter: &api.Router{
				ID:          2,
				Name:        "Router B",
				LocationID:  2,
				RouterLinks: []int{1, 3},
			},
			linkedRouter: &api.Router{
				ID:          3,
				Name:        "Router C",
				LocationID:  3,
				RouterLinks: []int{2},
			},
			storage: storageMock,
			storageMockOutcomes: func(storageMock *mock_storage.MockStorage) {
				storageMock.EXPECT().
					GetLocation(2).
					Times(1).
					Return(&api.Location{
						ID:       2,
						Postcode: "B",
						Name:     "Location B",
					}, nil)
				storageMock.EXPECT().
					GetLocation(3).
					Times(1).
					Return(&api.Location{
						ID:       3,
						Postcode: "C",
						Name:     "Location C",
					}, nil)
				storageMock.EXPECT().
					GetRouterLocationLink("Location B:Location C").
					Times(1).
					Return(nil, redis.Nil)
				storageMock.EXPECT().AddRouterLocationLink(&api.RouterLocationLink{
					UniqueID:   sortLocations("Location B", "Location C"),
					Connection: fmt.Sprintf("[%s] <-> [%s]", "Location B", "Location C"),
				}).Times(1).
					Return(nil)
			},
			log:              log,
			processedRouters: map[int]struct{}{},
		},
		{
			name: "processes linked router connections when parent already exists in processed routers",
			parentRouter: &api.Router{
				ID:          2,
				Name:        "Router B",
				LocationID:  2,
				RouterLinks: []int{1, 3},
			},
			linkedRouter: &api.Router{
				ID:          3,
				Name:        "Router C",
				LocationID:  3,
				RouterLinks: []int{2},
			},
			storage: storageMock,
			storageMockOutcomes: func(storageMock *mock_storage.MockStorage) {
				storageMock.EXPECT().
					GetLocation(2).
					Times(1).
					Return(&api.Location{
						ID:       2,
						Postcode: "B",
						Name:     "Location B",
					}, nil)
				storageMock.EXPECT().
					GetLocation(3).
					Times(1).
					Return(&api.Location{
						ID:       3,
						Postcode: "C",
						Name:     "Location C",
					}, nil)
				storageMock.EXPECT().
					GetRouterLocationLink("Location B:Location C").
					Times(1).
					Return(nil, redis.Nil)
				storageMock.EXPECT().AddRouterLocationLink(&api.RouterLocationLink{
					UniqueID:   sortLocations("Location B", "Location C"),
					Connection: fmt.Sprintf("[%s] <-> [%s]", "Location B", "Location C"),
				}).Times(1).
					Return(nil)
			},
			log:              log,
			processedRouters: map[int]struct{}{2: {}},
		},
		{
			name: "processes linked router connections when linked router already exists in processed routers",
			parentRouter: &api.Router{
				ID:          2,
				Name:        "Router B",
				LocationID:  2,
				RouterLinks: []int{1, 3},
			},
			linkedRouter: &api.Router{
				ID:          3,
				Name:        "Router C",
				LocationID:  3,
				RouterLinks: []int{2},
			},
			storage: storageMock,
			storageMockOutcomes: func(storageMock *mock_storage.MockStorage) {
			},
			log:              log,
			processedRouters: map[int]struct{}{3: {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storageMockOutcomes(storageMock)
			a := &app{
				storage: tt.storage,
				log:     tt.log,
			}
			a.processLinkedRouter(tt.parentRouter, tt.linkedRouter, tt.processedRouters)
		})
	}
}
