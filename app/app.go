package app

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	
	"router-location-connecter/api"
	"router-location-connecter/storage"
)

type app struct {
	apiClient api.API
	storage   storage.Storage
	log       zerolog.Logger
}

func NewApp(client api.API, redisClient storage.Storage, l zerolog.Logger) app {
	return app{
		apiClient: client,
		storage:   redisClient,
		log:       l,
	}
}

// Process runs the logic of coordinating the retrieval of data and processing it
func (a *app) Process(ctx context.Context) {
	// request api data
	rLocData, err := a.apiClient.GetRouterLocationData(ctx)
	if err != nil {
		log.Error().Err(err).Msg("get router location data")

		return
	}

	// keep track of routers already processed as we process based on links
	processedRouters := make(map[int]struct{}, 0)

	a.SaveRouterData(rLocData.Routers)
	a.SaveLocationData(rLocData.Locations)

	// output list of connections between locations
	for _, router := range rLocData.Routers {
		if _, ok := processedRouters[router.ID]; ok {
			// already processed router entry
			continue
		}

		for _, rLinkID := range router.RouterLinks {

			if router.ID == rLinkID {
				continue
			}

			linkedRouter, err := a.storage.GetRouter(rLinkID)
			if err != nil {
				if err != redis.Nil {
					log.Error().Err(err).Msg("get router data")
				}
				continue
			}

			a.processLinkedRouter(&router, linkedRouter, processedRouters)
		}
	}
}

// processLinkedRouter is a recursive function responsible for crawling through router links and calculating connections
func (a *app) processLinkedRouter(parentRouter, linkedRouter *api.Router, processedRouters map[int]struct{}) {
	if _, ok := processedRouters[linkedRouter.ID]; ok {
		// already processed linked router entry and since bidirectional its ok to skip
		return
	}

	for _, link := range linkedRouter.RouterLinks {
		//  check if link is bidirectional and follow links within the linked routers link
		// 2 routers connected at same location
		if parentRouter.LocationID == linkedRouter.LocationID {
			continue
		}

		if link == parentRouter.ID {
			processedRouters[parentRouter.ID] = struct{}{}
			if err := a.CalculateLink(parentRouter.LocationID, linkedRouter.LocationID); err != nil {
				if err != redis.Nil {
					log.Error().Err(err).Msg("calculate link")
				}
			}
		} else {
			// get routers for other the linked routers within the linked router
			addLinkedRouter, err := a.storage.GetRouter(link)
			if err != nil {
				if err != redis.Nil {
					log.Error().Err(err).Msg("get router data")
				}
				continue
			}

			a.processLinkedRouter(linkedRouter, addLinkedRouter, processedRouters)

		}
	}
}

// SaveRouterData makes a call to storage to save router data
func (a *app) SaveRouterData(routers []api.Router) {
	for _, router := range routers {
		if err := a.storage.AddRouter(&router); err != nil {
			log.Error().Err(err).Msg("store router data")
		}
	}
}

// SaveLocationData makes a call to storage to save location data
func (a *app) SaveLocationData(locations []api.Location) {
	for _, location := range locations {
		if err := a.storage.AddLocation(&location); err != nil {
			log.Error().Err(err).Msg("store location data")
		}
	}
}

// CalculateLink calculates links between router locations and prints to stdout
func (a *app) CalculateLink(srcLocationID, destLocationID int) error {
	// link is bi-directional, get link details
	srcLocation, err := a.storage.GetLocation(srcLocationID)
	if err != nil {
		return err
	}

	destLocation, err := a.storage.GetLocation(destLocationID)
	if err != nil {
		return err
	}

	// check if link(in any direction) already exists
	// we generate a unique alphabetically sorted ID
	linkUniqueID := sortLocations(srcLocation.Name, destLocation.Name)
	exists, err := a.storage.GetRouterLocationLink(linkUniqueID)
	if err != redis.Nil {
		return err
	}

	if exists != nil {
		return nil
	}

	fmt.Printf("[%s] <-> [%s]\n", srcLocation.Name, destLocation.Name)
	// store locationLink
	a.storage.AddRouterLocationLink(&api.RouterLocationLink{
		UniqueID:   linkUniqueID,
		Connection: fmt.Sprintf("[%s] <-> [%s]", srcLocation.Name, destLocation.Name),
	})

	return nil
}

// Function to sort two locations alphabetically and concatenates them
func sortLocations(loc1, loc2 string) string {
	if loc1 < loc2 {
		return loc1 + ":" + loc2
	}
	return loc2 + ":" + loc1
}
