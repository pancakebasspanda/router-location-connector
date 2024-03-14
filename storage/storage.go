package storage

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
	goredis "github.com/redis/go-redis/v9"
	"router-location-connecter/api"
)

// Storage is the interface for storage operations
type Storage interface {
	AddRouter(router *api.Router) error
	GetRouter(id int) (*api.Router, error)
	AddLocation(location *api.Location) error
	GetLocation(id int) (*api.Location, error)
	AddRouterLocationLink(links *api.RouterLocationLink) error
	GetRouterLocationLink(uniqueID string) (*api.RouterLocationLink, error)
	FlushAll(ctx context.Context) error
	Close() error
}

// Redis is the implementation of Storage interface
type Redis struct {
	Rh     *rejson.Handler
	Client *goredis.Client
}

func (r *Redis) FlushAll(ctx context.Context) error {
	// Flush all data from the selected database
	if err := r.Client.FlushAll(ctx).Err(); err != nil {
		return err
	}

	return nil
}

func (r *Redis) Close() error {
	if err := r.Client.Close(); err != nil {
		return err
	}

	return nil
}

func (r *Redis) AddRouterLocationLink(link *api.RouterLocationLink) error {
	res, err := r.Rh.JSONSet(link.UniqueID, ".", link)
	if err != nil {
		return err
	}

	if res.(string) != "OK" {
		return err
	}

	return nil
}

func (r *Redis) GetRouterLocationLink(uniqueID string) (*api.RouterLocationLink, error) {
	value, err := redis.Bytes(r.Rh.JSONGet(uniqueID, "."))
	if err != nil {
		return nil, err
	}

	link := api.RouterLocationLink{}
	if err = json.Unmarshal(value, &link); err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *Redis) AddRouter(router *api.Router) error {
	res, err := r.Rh.JSONSet("router_id_"+strconv.Itoa(router.ID), ".", router)
	if err != nil {
		return err
	}

	if res.(string) != "OK" {
		return err
	}

	return nil
}

func (r *Redis) GetRouter(id int) (*api.Router, error) {
	value, err := redis.Bytes(r.Rh.JSONGet("router_id_"+strconv.Itoa(id), "."))
	if err != nil {
		return nil, err
	}

	router := api.Router{}
	if err = json.Unmarshal(value, &router); err != nil {
		return nil, err
	}

	return &router, nil
}

func (r *Redis) AddLocation(location *api.Location) error {
	res, err := r.Rh.JSONSet("location_id_"+strconv.Itoa(location.ID), ".", location)
	if err != nil {
		return err
	}

	if res.(string) != "OK" {
		return err
	}

	return nil
}

func (r *Redis) GetLocation(id int) (*api.Location, error) {
	value, err := redis.Bytes(r.Rh.JSONGet("location_id_"+strconv.Itoa(id), "."))
	if err != nil {
		return nil, err
	}

	location := api.Location{}
	if err = json.Unmarshal(value, &location); err != nil {
		return nil, err
	}

	return &location, nil
}

func New(ctx context.Context, address, password string) (Storage, error) {
	reJsonHandler := rejson.NewReJSONHandler()

	client := goredis.NewClient(&goredis.Options{
		Addr:     address, // Assuming Redis is running on localhost
		Password: password,
	})

	reJsonHandler.SetGoRedisClientWithContext(ctx, client)

	return &Redis{
		Rh:     reJsonHandler,
		Client: client,
	}, nil
}

//func newPool(address, password string) *redis.Pool {
//	return &redis.Pool{
//		Wait:        true,
//		MaxIdle:     30,
//		IdleTimeout: 240 * time.Second,
//		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", address, redis.DialPassword(password)) },
//	}
//}

//func (r *Redis) getReJSONHandler() rejson.Handler {
//	reJsonHandler := rejson.NewReJSONHandler()
//
//
//	reJsonHandler.SetRedigoClient(r.pool.Get())
//
//	return reJsonHandler
//}
